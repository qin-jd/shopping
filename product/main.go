package main

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"shopping/product/handler"
	"shopping/product/model"
	"shopping/product/repository"

	product "shopping/product/proto/product"
)

func main() {

	db,err := CreateConnection()
	defer db.Close()

	db.AutoMigrate(&model.Product{})

	if err != nil {
		log.Fatalf("connection error : %v \n" , err)
	}

	repo := &repository.Product{db}

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.product"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	product.RegisterProductServiceHandler(service.Server(), &handler.Product{repo})

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.product", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.product", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
