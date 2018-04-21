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
	"strconv"
	"sync"
)

type OrderSubscriber interface {
	Subscribe(regionId int32)
	Init()
}

//20 regions at ~50,000 items equates to ~83m of memory required... not too shabby
type itemPrice struct {
	itemId, readings int32
	min, max, avg    float64
}

type movingAverages struct {
	Region map[int32] /*regionId*/ map[bool] /*buy|sell*/ map[int64] /*itemId*/ itemPrice
}

type orderSubscriber struct {
	s           *mgo.Session
	persistChan chan int32

	buMu      sync.Mutex
	bulkS     map[int32]*mgo.Session
	bulk      map[int32]*mgo.Bulk
	batchSize map[int32]int

	mu       sync.Mutex
	averages movingAverages
}

// Subscribe puts a new Handler function on the broker's subscription whatever.
// This function will produce 4 such Handler's:
//   1. regionId.state.begin listener
//   2. buy.regionId listener
//   3. sell.regionId listener
//   4. regionId.state.end listener
func (sub *orderSubscriber) Subscribe(regionId int32) {
	log, err := zap.NewProduction()
	slog := log.Sugar()

	if err = broker.Init(); err != nil {
		slog.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		slog.Fatalf("Broker Connect error: %v", err)
	}

	_, err = broker.Subscribe("buy."+strconv.Itoa(int(regionId)), func(p broker.Publication) error {
		return sub.orderPublication(p)
	})

	_, err = broker.Subscribe("sell."+strconv.Itoa(int(regionId)), func(p broker.Publication) error {
		return sub.orderPublication(p)
	})

	_, err = broker.Subscribe(strconv.Itoa(int(regionId))+".state.begin", func(p broker.Publication) error {
		return sub.beginPublication(p)
	})

	_, err = broker.Subscribe(strconv.Itoa(int(regionId))+".state.end", func(p broker.Publication) error {
		return sub.endPublication(p)
	})

	if err != nil {
		slog.Fatal(err.Error())
	}
}

func (sub *orderSubscriber) orderPublication(p broker.Publication) error {
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

	sub.mu.Lock()
	region := sub.averages.Region[12345]

	myAverages := region[op.IsBuyOrder][op.OrderId]

	if op.Price < myAverages.min {
		myAverages.min = op.Price
	}

	if op.Price > myAverages.max {
		myAverages.max = op.Price
	}

	oldAvgDividend := myAverages.avg * float64(myAverages.readings)
	myAverages.readings += 1
	myAverages.avg = (op.Price + oldAvgDividend) / float64(myAverages.readings)
	sub.mu.Unlock()

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

	//TODO: The below being implemented is a golive blocker
	//It's implementation will grab all orders for this region that don't have this FetchRequestId
	//and deletes them adjusting the sub's averages accordingly
	//cleanupOldOrders(op.FetchRequestId)

	return nil
}

func (sub *orderSubscriber) beginPublication(p broker.Publication) error {
	return nil
}

func (sub *orderSubscriber) endPublication(p broker.Publication) error {

	return nil
}

func (sub *orderSubscriber) Init() {
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

	go sub.persister()
}

func (sub *orderSubscriber) persister() {
	for {
		select {
		case regionId := <-sub.persistChan:
			fmt.Printf("Persisting %v\n", regionId)
			_, err := sub.bulk[regionId].Run()
			if err != nil {
				fmt.Printf("Had an error running the bulk: %s\n", err.Error())
			}

			sub.bulk[regionId] = nil
			sub.bulkS[regionId] = nil
			sub.buMu.Unlock()
		}
	}
}

func NewOrderSubscriber(sess *mgo.Session) OrderSubscriber {
	return &orderSubscriber{
		s:           sess,
		persistChan: make(chan int32),
		bulkS:       make(map[int32]*mgo.Session),
		bulk:        make(map[int32]*mgo.Bulk),
		batchSize:   make(map[int32]int),
	}
}
