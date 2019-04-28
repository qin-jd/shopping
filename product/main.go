package main

import (
	"github.com/micro/go-config"
	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	//"github.com/opentracing/opentracing-go"
	//"os"
	"shopping/product/handler"
	"shopping/product/model"
	"shopping/product/repository"

	product "shopping/product/proto/product"
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
	conf := config.Map()

	//db
	db, err := CreateConnection(conf["mysql"].(map[string]interface{}))
	defer db.Close()

	db.AutoMigrate(&model.Product{})

	if err != nil {
		log.Fatalf("connection error : %v \n" , err)
	}

	repo := &repository.Product{db}

	// boot trace
	//TraceBoot()

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.product"),
		micro.Version("latest"),
		//micro.WrapHandler(wrapperTrace.NewHandlerWrapper()),
		//micro.WrapClient(wrapperTrace.NewClientWrapper()),
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

//func TraceBoot() {
//	apiURL := "http://192.168.0.111:9411/api/v1/spans"
//	hostPort,_ := os.Hostname()
//	serviceName := "go.micro.srv.product"
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


