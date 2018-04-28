package main

import (
	"fmt"
	"github.com/abaeve/pricing-service/handler"
	"github.com/abaeve/pricing-service/proto"
	"github.com/abaeve/pricing-service/subscriber"
	"github.com/chremoas/services-common/config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/prometheus/common/log"
	"gopkg.in/mgo.v2"
)

var Version = "1.0.0"
var service micro.Service
var sub *subscriber.OrderSubImpl
var mongoConfig = "localhost"
var confFile string

var serviceType = "srv"
var serviceName = "pricing"
var fullName string

func main() {
	cmd.Init()

	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	session, err := mgo.Dial(mongoConfig)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	sub = subscriber.NewOrderSubscriber(session)
	sub.Init()
	sub.Subscribe(10000002)
	sub.Subscribe(10000043)
	sub.Subscribe(10000032)
	sub.Subscribe(10000030)

	conf := config.Configuration{}
	conf.Namespace = "com.abaeve"
	fullName = conf.LookupService(serviceType, serviceName)

	service = micro.NewService(
		micro.Version(Version),
		micro.Name(fullName),
	)

	service.Init()

	pricing.RegisterPricesHandler(service.Server(), handler.NewPricesHandler(sub))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

//// This function is a callback from the config.NewService function.  Read those docs
//func initialize(config *config.Configuration) error {
//	mongoConfig = config.Extensions["mongodb"].(map[interface{}]interface{})["dial"].(string)
//
//	pricing.RegisterPricesHandler(service.Server(), handler.NewPricesHandler(sub))
//
//	return nil
//}

func oldMain() {
	cmd.Init()

	forever := make(chan struct{})

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	sub := subscriber.NewOrderSubscriber(session)
	sub.Init()
	sub.Subscribe(10000002)
	sub.Subscribe(10000043)
	sub.Subscribe(10000032)
	sub.Subscribe(10000030)

	fmt.Println("Waiting...")

	<-forever
}
