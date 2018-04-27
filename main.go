package main

import (
	"fmt"
	"github.com/abaeve/pricing-service/subscriber"
	"github.com/micro/go-micro/cmd"
	"gopkg.in/mgo.v2"
)

func main() {
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
