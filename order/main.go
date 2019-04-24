package main

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"

	"shopping/order/handler"
	"shopping/order/model"
	"shopping/order/repository"

	"github.com/micro/go-plugins/broker/rabbitmq"
	order "shopping/order/proto/order"
	product "shopping/product/proto/product"
)

func main() {

	db,err := CreateConnection()
	defer db.Close()

	db.AutoMigrate(&model.Order{})

	if err != nil {
		log.Fatalf("connection error : %v \n" , err)
	}

	repo := &repository.Order{db}

	b := rabbitmq.NewBroker(
		broker.Addrs("amqp://用户名:密码@主机host:端口port"),
	)

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.order"),
		micro.Version("latest"),
		micro.Broker(b),
	)

	// Initialise service
	service.Init()

	//创建消息发布者
	publisher := micro.NewPublisher("notification.submit" , service.Client())

	//product-srv client
	productCli := product.NewProductService("go.micro.srv.product" , service.Client())

	// Register Handler
	order.RegisterOrderServiceHandler(service.Server(), &handler.Order{repo , productCli , publisher})

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.order", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.order", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
