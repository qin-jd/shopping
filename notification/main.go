package main

import (
	"github.com/micro/go-config"
	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
	//"github.com/opentracing/opentracing-go"
	//"os"

	//"github.com/micro/go-micro/broker"
	//"github.com/micro/go-plugins/broker/rabbitmq"

	"shopping/notification/subscriber"
	//wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing"
	//zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
)

func main() {
	//加载配置项
	err := config.LoadFile("./config.json")
	if err != nil {
		log.Fatalf("Could not load config file: %s", err.Error())
		return
	}

	b := rabbitmq.NewBroker(
		broker.Addrs(config.Get("rabbitmq_addr").String("")),
	)

	b.Init()
	b.Connect()

	// boot trace
	//TraceBoot()

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.notification"),
		micro.Version("latest"),
		//micro.WrapSubscriber(wrapperTrace.NewSubscriberWrapper()),
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

//func TraceBoot() {
//	apiURL := "http://192.168.0.111:9411/api/v1/spans"
//	hostPort,_ := os.Hostname()
//	serviceName := "go.micro.srv.notification"
//
//	collector, err := zipkin.NewHTTPCollector(apiURL)
//	if err != nil {
//		log.Fatalf("unable to create Zipkin HTTP collector: %v", err)
//		return
//	}
//	tracer, err := zipkin.NewTracer(
//		zipkin.NewRecorder(collector, false, hostPort, serviceName),
//	)
//	if err != nil {
//		log.Fatalf("unable to create Zipkin tracer: %v", err)
//		return
//	}
//	opentracing.InitGlobalTracer(tracer)
//	return
//}