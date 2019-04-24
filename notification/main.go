package main

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"

	"shopping/notification/subscriber"
)

func main() {

	b := rabbitmq.NewBroker(
		broker.Addrs("amqp://用户名:密码@主机host:端口port"),
	)

	b.Init()
	b.Connect()

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.notification"),
		micro.Version("latest"),
		micro.Broker(b),
	)

	// Initialise service
	service.Init()

	// Register Handler
	//example.RegisterExampleHandler(service.Server(), new(handler.Example))

	//defer sub.Unsubscribe()
	micro.RegisterSubscriber("notification.submit", service.Server(), new(subscriber.Notification))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.notification", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
