// Package subscriber holds the bits and bobs that listen to the pricing messages being firehosed onto the queue
// and does stuff to these messages... like... stores them and creates stats and stuff.
package subscriber

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/abaeve/pricing-service/model"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"strconv"
	"sync"
	"time"
)

type OrderSubscriber interface {
	Subscribe(regionId int32)
	Init()
}

type stats struct {
	region map[int32] /*regionId*/ map[int32] /*typeId*/ *model.TypeStat
}

type OrderSubImpl struct {
	s           *mgo.Session
	persistChan chan int32

	buMu      sync.Mutex
	bulkS     map[int32]*mgo.Session
	bulk      map[int32]*mgo.Bulk
	batchSize map[int32]int

	statMu sync.Mutex
	stats  stats

	bench map[int32]map[string] /*start|stop*/ time.Time
}

// Subscribe puts a new Handler function on the broker's subscription whatever.
// This function will produce 4 such Handler's:
//   1. regionId.state.begin listener
//   2. buy.regionId listener
//   3. sell.regionId listener
//   4. regionId.state.end listener
func (sub *OrderSubImpl) Subscribe(regionId int32) {
	log, err := zap.NewProduction()
	slog := log.Sugar()

	_, err = broker.Subscribe("buy."+strconv.Itoa(int(regionId)), func(p broker.Publication) error {
		return sub.orderPublication(p)
	})

	_, err = broker.Subscribe("sell."+strconv.Itoa(int(regionId)), func(p broker.Publication) error {
		return sub.orderPublication(p)
	})

	_, err = broker.Subscribe(strconv.Itoa(int(regionId))+".state.begin", func(p broker.Publication) error {
		return sub.beginPublication(regionId, p)
	})

	_, err = broker.Subscribe(strconv.Itoa(int(regionId))+".state.end", func(p broker.Publication) error {
		go sub.endPublication(regionId, p)
		return nil
	})

	if err != nil {
		slog.Fatal(err.Error())
	}
}

func (sub *OrderSubImpl) Init() {
	session := sub.s.Copy()
	defer session.Close()

	c := session.DB("market").C("orders")

	index := mgo.Index{
		Key:        []string{"request_id", "order_id", "location_id", "type_id"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	c = session.DB("market").C("stats")

	index = mgo.Index{
		Key:        []string{"region_id", "type_id"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	go sub.persister()
}

func (sub *OrderSubImpl) Stats(regionId, typeId int32) *model.TypeStat {
	fmt.Printf("RegionId: %d, TypeId: %d\n", regionId, typeId)
	fmt.Println(strconv.Itoa(int(regionId)) + "-" + strconv.Itoa(int(typeId)))
	result := model.TypeStat{}

	session := sub.s.Copy()
	defer session.Close()

	c := session.DB("market").C("stats")

	query := c.Find(bson.M{"_id": strconv.Itoa(int(regionId)) + "-" + strconv.Itoa(int(typeId))}).Iter()

	if ok := query.Next(&result); !ok {
		fmt.Printf("%v", result)
		if err := query.Err(); err != nil {
			fmt.Errorf("errored fetching stats for region %d typeid %d with message %s", regionId, typeId, err.Error())
		}
	}

	return &result
}

func (sub *OrderSubImpl) orderPublication(p broker.Publication) error {
	op := &model.OrderPayload{}

	err := json.Unmarshal(p.Message().Body, op)
	if err != nil {
		return err
	}

	//BGN Validation
	if op.OrderId == 0 && op.TypeId == 0 {
		return errors.New("order_id and type_id cannot be 0")
	}
	if op.TypeId == 0 {
		return errors.New("type_id cannot be 0")
	}
	//END Validation

	sub.statMu.Lock()
	//fmt.Print(".")
	if sub.stats.region[op.RegionId] == nil {
		//fmt.Print("!")
		sub.stats.region[op.RegionId] = make(map[int32]*model.TypeStat)
	}

	if sub.stats.region[op.RegionId] == nil {
		//fmt.Print("@")
		sub.stats.region[op.RegionId] = make(map[int32]*model.TypeStat)
	}

	if sub.stats.region[op.RegionId][op.TypeId] == nil {
		//fmt.Print("#")
		sub.stats.region[op.RegionId][op.TypeId] = &model.TypeStat{
			RegionId: op.RegionId,
			TypeId:   op.TypeId,
		}

		sub.stats.region[op.RegionId][op.TypeId].Buy.Min = math.MaxFloat64
		sub.stats.region[op.RegionId][op.TypeId].Sell.Min = math.MaxFloat64
	}
	//fmt.Print("\n")

	if op.IsBuyOrder {
		if op.Price < sub.stats.region[op.RegionId][op.TypeId].Buy.Min {
			sub.stats.region[op.RegionId][op.TypeId].Buy.Min = op.Price
		}

		if op.Price > sub.stats.region[op.RegionId][op.TypeId].Buy.Max {
			sub.stats.region[op.RegionId][op.TypeId].Buy.Max = op.Price
		}

		oldAvgDividend := sub.stats.region[op.RegionId][op.TypeId].Buy.Avg * float64(sub.stats.region[op.RegionId][op.TypeId].Buy.Ord)
		sub.stats.region[op.RegionId][op.TypeId].Buy.Ord += 1
		sub.stats.region[op.RegionId][op.TypeId].Buy.Vol += int64(op.VolumeRemain)
		sub.stats.region[op.RegionId][op.TypeId].Buy.Avg = (op.Price + oldAvgDividend) / float64(sub.stats.region[op.RegionId][op.TypeId].Buy.Ord)
	} else {
		if op.Price < sub.stats.region[op.RegionId][op.TypeId].Sell.Min {
			sub.stats.region[op.RegionId][op.TypeId].Sell.Min = op.Price
		}

		if op.Price > sub.stats.region[op.RegionId][op.TypeId].Sell.Max {
			sub.stats.region[op.RegionId][op.TypeId].Sell.Max = op.Price
		}

		oldAvgDividend := sub.stats.region[op.RegionId][op.TypeId].Sell.Avg * float64(sub.stats.region[op.RegionId][op.TypeId].Sell.Ord)
		sub.stats.region[op.RegionId][op.TypeId].Sell.Ord += 1
		sub.stats.region[op.RegionId][op.TypeId].Sell.Vol += int64(op.VolumeRemain)
		sub.stats.region[op.RegionId][op.TypeId].Sell.Avg = (op.Price + oldAvgDividend) / float64(sub.stats.region[op.RegionId][op.TypeId].Sell.Ord)
	}
	sub.statMu.Unlock()

	sub.buMu.Lock()
	if sub.bulk[op.RegionId] == nil {
		sub.bulkS[op.RegionId] = sub.s.Copy()
		sub.bulk[op.RegionId] = sub.bulkS[op.RegionId].DB("market").C("orders").Bulk()
	}
	sub.bulk[op.RegionId].Upsert(bson.M{"_id": op.OrderId}, op)
	sub.batchSize[op.RegionId] += 1

	//We're going to continue if we call for a bulk run but we'll pass control of the lock to the persister
	if sub.batchSize[op.RegionId] >= 999 {
		sub.batchSize[op.RegionId] = 0
		sub.persistChan <- op.RegionId
	} else {
		sub.buMu.Unlock()
	}

	return nil
}

func (sub *OrderSubImpl) beginPublication(regionId int32, p broker.Publication) error {
	fmt.Printf("%s BGN FetchRequestId: %s\n", time.Now().Format("2006-01-02T15:04:05.999999-07:00"), string(p.Message().Body))
	go func(regionId int32, fetchRequestId string) {
		sub.statMu.Lock()

		sub.statMu.Unlock()
	}(regionId, string(p.Message().Body))
	return nil
}

func (sub *OrderSubImpl) endPublication(regionId int32, p broker.Publication) error {
	fmt.Printf("%s END FetchRequestId: %s\n", time.Now().Format("2006-01-02T15:04:05.999999-07:00"), string(p.Message().Body))
	fetchRequestId := string(p.Message().Body)

	time.Sleep(time.Second * 3)
	sub.buMu.Lock()
	sub.persistChan <- regionId
	time.Sleep(time.Second)

	//fmt.Println("About to iterate over stats...")

	sub.statMu.Lock()
	session := sub.s.Copy()

	c := session.DB("market").C("stats").Bulk()

	documentCount := 0

	for typeId, typeStat := range sub.stats.region[regionId] {
		//fmt.Println("Looping over stats...")

		savingTypeStat := &model.TypeStat{}
		savingTypeStat.RegionId = regionId
		savingTypeStat.TypeId = typeId

		savingTypeStat.Buy.Min = typeStat.Buy.Min
		savingTypeStat.Buy.Avg = typeStat.Buy.Avg
		savingTypeStat.Buy.Max = typeStat.Buy.Max
		savingTypeStat.Buy.Vol = typeStat.Buy.Vol
		savingTypeStat.Buy.Ord = typeStat.Buy.Ord

		savingTypeStat.Sell.Min = typeStat.Sell.Min
		savingTypeStat.Sell.Avg = typeStat.Sell.Avg
		savingTypeStat.Sell.Max = typeStat.Sell.Max
		savingTypeStat.Sell.Vol = typeStat.Sell.Vol
		savingTypeStat.Sell.Ord = typeStat.Sell.Ord

		//c.Upsert(bson.M{"_id": strconv.Itoa(int(time.Now().UnixNano())) + "-" + strconv.Itoa(int(savingTypeStat.RegionId)) + "-" + strconv.Itoa(int(savingTypeStat.TypeId))}, savingTypeStat)
		c.Upsert(bson.M{"_id": strconv.Itoa(int(savingTypeStat.RegionId)) + "-" + strconv.Itoa(int(savingTypeStat.TypeId))}, savingTypeStat)

		documentCount += 1
		//fmt.Printf("Document count %d\n", documentCount)
		if documentCount >= 999 {
			//fmt.Println("Persisting stats...")
			_, err := c.Run()
			if err != nil {
				fmt.Errorf("errored saving stats: %s\n", err.Error())
			}

			documentCount = 0
			session.Close()
			session = sub.s.Copy()
			c = session.DB("market").C("stats").Bulk()
		}
	}
	sub.statMu.Unlock()

	_, err := c.Run()
	if err != nil {
		fmt.Errorf("errored saving stats: %s\n", err.Error())
	}
	session.Close()

	sub.cleanupStats(regionId, fetchRequestId)
	sub.cleanupForFetchRequestId(regionId, fetchRequestId)

	return nil
}

func (sub *OrderSubImpl) persister() {
	for {
		select {
		case regionId := <-sub.persistChan:
			//fmt.Printf("Persisting %v\n", regionId)
			/*bulkResult*/
			if sub.bulk[regionId] != nil {
				_, err := sub.bulk[regionId].Run()
				if err != nil {
					fmt.Errorf("had an error running the bulk: %s\n", err.Error())
				}

				//if bulkResult != nil {
				//	fmt.Printf("Persisted: %d documents\n", bulkResult.Matched)
				//}

				sub.bulkS[regionId].Close()
			}
			sub.bulk[regionId] = nil
			sub.bulkS[regionId] = nil
			sub.buMu.Unlock()
		}
	}
}

func (sub *OrderSubImpl) cleanupStats(regionId int32, fetchRequestId string) {
	sub.statMu.Lock()
	sub.stats.region[regionId] = nil
	sub.statMu.Unlock()
	//I'd rather do something else with this information such as publish to a queue for others to observe.
	//Our job here however is maintain the correct current state
	//session := sub.s.Copy()
	//defer session.Close()
	//
	//c := session.DB("market").C("orders")
	//
	//result := c.Find(bson.M{"region_id": regionId, "fetch_request_id": bson.M{"$ne": fetchRequestId}}).Iter()
	//
	//var payload *model.OrderPayload
	//
	//ok := result.Next(payload)
	//for ok {
	//	sub.statMu.Lock()
	//
	//	if sub.stats.region[payload.RegionId][payload.TypeId].Buy.Min == payload.Price {
	//		sub.stats.region[payload.RegionId][payload.TypeId].Buy.Min = math.MaxFloat64
	//	}
	//
	//	if sub.stats.region[payload.RegionId][payload.TypeId].Buy.Max == payload.Price {
	//		sub.stats.region[payload.RegionId][payload.TypeId].Buy.Max = 0
	//	}
	//
	//	if sub.stats.region[payload.RegionId][payload.TypeId].Sell.Min == payload.Price {
	//		sub.stats.region[payload.RegionId][payload.TypeId].Sell.Min = math.MaxFloat64
	//	}
	//
	//	if sub.stats.region[payload.RegionId][payload.TypeId].Sell.Max == payload.Price {
	//		sub.stats.region[payload.RegionId][payload.TypeId].Sell.Max = 0
	//	}
	//
	//	oldAvgDividend := sub.stats.region[payload.RegionId][payload.TypeId].Buy.Avg * float64(sub.stats.region[payload.RegionId][payload.TypeId].Buy.Ord)
	//	sub.stats.region[payload.RegionId][payload.TypeId].Buy.Ord -= 1
	//	sub.stats.region[payload.RegionId][payload.TypeId].Buy.Avg = (oldAvgDividend - payload.Price) / float64(sub.stats.region[payload.RegionId][payload.TypeId].Buy.Ord)
	//
	//	oldAvgDividend = sub.stats.region[payload.RegionId][payload.TypeId].Sell.Avg * float64(sub.stats.region[payload.RegionId][payload.TypeId].Sell.Ord)
	//	sub.stats.region[payload.RegionId][payload.TypeId].Sell.Ord -= 1
	//	sub.stats.region[payload.RegionId][payload.TypeId].Sell.Avg = (oldAvgDividend - payload.Price) / float64(sub.stats.region[payload.RegionId][payload.TypeId].Sell.Ord)
	//
	//	if ok = result.Next(payload); !ok {
	//		err := result.Err()
	//		if err != nil {
	//			fmt.Errorf("had an error while cleaning up stats: %s", err.Error())
	//		}
	//	}
	//	sub.statMu.Unlock()
	//}
}

func (sub *OrderSubImpl) cleanupForFetchRequestId(regionId int32, fetchRequestId string) {
	session := sub.s.Copy()
	defer session.Close()

	c := session.DB("market").C("orders")

	info, err := c.RemoveAll(bson.M{"region_id": regionId, "fetch_request_id": bson.M{"$ne": fetchRequestId}})
	if err != nil {
		fmt.Errorf("errored out with: %s\n", err.Error())
	}

	if info != nil {
		fmt.Printf("Removed: %d for region: %d, and request_id: %s\n", info.Removed, regionId, fetchRequestId)
	}
}

func NewOrderSubscriber(sess *mgo.Session) *OrderSubImpl {
	result := &OrderSubImpl{
		s:           sess,
		persistChan: make(chan int32),
		bulkS:       make(map[int32]*mgo.Session),
		bulk:        make(map[int32]*mgo.Bulk),
		batchSize:   make(map[int32]int),

		stats: stats{
			region: make(map[int32]map[int32]*model.TypeStat),
		},
	}
	return result
}
