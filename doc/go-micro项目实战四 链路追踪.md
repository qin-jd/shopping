
### 链路追踪
微服务架构下，所有的服务都分散在不同的地方，一旦某个服务出现问题，排查起来很费时费力。所以在微服务的演进下，需要一个全链路追踪系统来分析服务的运行状态。

### go-micro的trace插件
Micro通过Wrapper实现了三种trace接口，aswxray,opencensus,opentracing。第一个是亚马逊AWS的。没有尝试。
opentracing是一个开源的标准。提供对厂商中立的 API，用来向应用程序添加追踪功能并将追踪数据发送到分布式的追踪系统。已经快成功行业的标准了。
opencensus是谷歌开源的数据收集和分布式跟踪框架。OpenCensus也是实现了opentracing标准。OpenCensus 不仅提供规范，还提供开发语言的实现，和连接协议，而且它不仅只做追踪，还引入了额外的度量指标，这些一般不在分布式追踪系统的职责范围。opencensus也支持把数据导出到别的系统做分析。比如zipkin和Prometheus等

我们本次实验通过两种组合来实现go-micro场景下的链路追踪。
opencensus+zipkin
opentracing+zipkin
opentracing+Jaeger

注：zipkin 是 twitter 开源的分布式跟踪系统，并且具有UI界面来显示每个跟踪请求的状态。
Prometheus是当前热门的监控方案，也是用go+Grafana开发的，据说是未来云原生应用下的监控解决方案。我们下一节详细讲这个。

### 运行zipkin
`docker run -d -p 9411:9411 openzipkin/zipkin`
然后浏览器访问host:9411端口，即可看到zipkin的UI界面


### 1.opencensus+zipkin
我们使用opencensus的trace功能。也是实现了opentracing的标准。具体实现是由opencensus的trace来做的，然后通过zipkin的exporter把trace收集到的数据丢给zipkin。

##### 1.1 需要引入的包
```
import (
	...
	"go.opencensus.io/trace"
	"go.opencensus.io/exporter/zipkin"
	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opencensus"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	...
)
```
#### 1.2修改order微服务下的main.go
创建TraceBoot方法
```
func TraceBoot() {
	apiURL := "http://192.168.0.111:9411/api/v2/spans"
	hostPort,_ := os.Hostname()
	serviceName := "go.micro.srv.order"

	localEndpoint, err := openzipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		log.Fatalf("Failed to create the local zipkinEndpoint: %v", err)
	}
	reporter := zipkinHTTP.NewReporter(apiURL)
	ze := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(ze)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return
}
```
在func main里引用
```
	// boot trace
	TraceBoot()

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.order"),
		micro.Version("latest"),
		micro.Broker(b),
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper()),
		micro.WrapClient(wrapperTrace.NewClientWrapper()),
	)
```
#### 1.3修改order微服务下的handler
将原来handler里的上下文传进去。因为为了识别是同一个请求，需要把相同的traceId传过去。这个是在上下文环境里的。具体操作就是将原来代码里的context.TODO()替换成ctx即可。

### 1.4启动服务，执行请求，查看链路追踪分析、
![](http://doc.dev.tanikawa.com/Public/Uploads/2019-04-25/5cc162495dce9.png)

### 2.opentracing+zipkin
我们使用opentracing定义的标准，具体的实现我们使用opentracing-go（opentracing标准的具体GO语言实现包），然后把tracer收集到的数据丢给zipkin。

#### 2.1需要引入的包
```
import (
	...
	"github.com/opentracing/opentracing-go"
	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing"
	...
)
```
#### 2.2修改order微服务下的main.go
创建TraceBoot方法
```
func TraceBoot() {
	apiURL := "http://192.168.0.111:9411/api/v1/spans"
	hostPort,_ := os.Hostname()
	serviceName := "go.micro.srv.order"

	collector, err := zipkin.NewHTTPCollector(apiURL)
	if err != nil {
		log.Fatalf("unable to create Zipkin HTTP collector: %v", err)
		return
	}
	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, hostPort, serviceName),
	)
	if err != nil {
		log.Fatalf("unable to create Zipkin tracer: %v", err)
		return
	}
	opentracing.InitGlobalTracer(tracer)
	return
}
```
除了import，其他部分不用修改，即可很方便地从opencensus+zipkin切换到opentracing+zipkin

### opentracing+jaeger
```
func TraceBoot() {
	serviceName := "go.micro.srv.order"
	cfg := jaegercfg.Configuration{
		ServiceName : serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			//CollectorEndpoint :"http://192.168.0.111:14268/api/traces",
			LogSpans: true,
			LocalAgentHostPort:  "192.168.0.111:6831",
		},
	}

	metricsFactory := jprom.New().Namespace(metrics.NSOptions{Name: serviceName, Tags: nil})
	jaegerLogger := jaegerLoggerAdapter{log.With("serviceName" , serviceName)}

	tracer, _, err := cfg.NewTracer(
		jaegercfg.Logger(jaegerLogger),
		jaegercfg.Metrics(metricsFactory),
		jaegercfg.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		log.Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}

	opentracing.SetGlobalTracer(tracer)
}
```

### jaeger截图
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190430095755762.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)
### 总结
统一标准的opentracing确实带来了很大的方便，几乎不费力就可以切换zipkin到opencensus或者是Jaeger。
个人更喜欢jaeger的界面。而且jaeger是拿go写的。





