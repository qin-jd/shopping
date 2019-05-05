### 服务监控
可能读者会问，链路追踪和服务监控不是一个东西么？以下知识普及来源于网上。
我这里从网上找到一张图 很形象。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190426162807213.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)
Logging - 用于记录离散的事件。例如，应用程序的调试信息或错误信息。它是我们诊断问题的依据。
Metrics - 用于记录可聚合的数据。例如，队列的当前深度可被定义为一个度量值，在元素入队或出队时被更新；HTTP 请求个数可被定义为一个计数器，新请求到来时进行累加。
Tracing - 用于记录请求范围内的信息。例如，一次远程方法调用的执行过程和耗时。它是我们排查系统性能问题的利器。

通过上述信息，我们可以对已有系统进行分类。例如，Zipkin 专注于 tracing 领域；Prometheus 开始专注于 metrics，随着时间推移可能会集成更多的 tracing 功能，但不太可能深入 logging 领域； ELK，阿里云日志服务这样的系统开始专注于 logging 领域，但同时也不断地集成其他领域的特性到系统中来，正向上图中的圆心靠近。

我们上一节讲了trace。这节来讲metrics。这一领域目前备受推崇的是Prometheus。云原生应用时代的比较流行的监控系统。
本次实验教程是用Prometheus+grafana来做

### 安装prometheus
在/home/www/下创建go-micro目录，里面创建prometheus.yml
```
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'
    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s
    static_configs:
      - targets: ['192.168.0.110:8085']
```
`docker run -d -p --network=host 9090:9090 -v /home/www/go-micro/prometheus.yml:/etc/prometheus/prometheus.yml --name prometheus prom/prometheus`
此时打开host:9090可以看到prometheus的界面

### 安装grafana
`docker run -d -p 3000:3000 grafana/grafana`
此时打开host:3000可以看到grafana的界面，默认用户名密码是admin

### 1.go-micro集成prometheus
go-micro里的插件里已经有prometheus中间件了。在go-plugins/wrapper/monitoring/prometheus/里，我们可以直接使用。

#### 1.1在order服务里的main.go里引入插件
```
import(
	...
	wrapperPrometheus "github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	...
)
```

#### 1.2在创建服务里注册prometheus中间件
```
service := grpc.NewService(
		micro.Name("go.micro.srv.order"),
		micro.Version("latest"),
		micro.Broker(b),
		micro.WrapHandler(wrapperTrace.NewHandlerWrapper()),//trace中间件
		micro.WrapClient(wrapperTrace.NewClientWrapper()),//trace中间件
		micro.WrapHandler(wrapperPrometheus.NewHandlerWrapper()),//这里是prometheus中间件
	)
```
prometheus中间件导出了三个用于监控的指标，可以查看插件里prometheus.go源码
request_total请求的总数量，upstream_latency_microseconds请求延迟，request_duration_seconds请求耗时
如果我们需要关心特殊的指标，我们可以自行参照prometheus.go源码来扩展我们自己关心的指标中间件。

#### 1.3为prometheus提供metrics接口
在main.go里写上以下方法
```
func PrometheusBoot(){
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("192.168.0.110:8085", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
```
然后在main func里调用该方法即可。

注意：这里的metrics接口是监听的是8085端口。在docker里启动时配置的targets就是该metrics接口监听的地址和端口。docker使用了host的网络模式，可以直接访问宿主机的网络。

### 2.grafana里可视化prometheus数据
#### 2.1配置数据源
在data source里找到prometheus，填写prometheus服务的地址。也就是host：9090端口那个
据说在选址Access那里选择Browser，目前还没深入了解这两个选项什么意思。

#### 2.2配置dashboard
按+即可创建一个dashboard，然后在添加panel那点击Add Query
然后在输入框那输入一些指标的name即可。
比如：micro_upstream_latency_microseconds,micro_request_duration_seconds_bucket,micro_request_total,go_threads,go_gc_duration_seconds

创建完成以后即可看到
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190426161855728.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)

### 3.grafana里的alerting告警机制
grafana从4.0开始支持alerting告警了。加上告警这就是完整的监控系统了。

#### 3.1 先配置grafana里的smtp邮件设置。
grafana配置文件在docker容器里的/etc/grafana/grafana.ini

开启smtp,找到smtp配置。把xxx替换成自己的smtp配置
```
enabled = true
host = smtp.qq.com:465
user = xxx@qq.com
# If the password contains # or ; you have to wrap it with trippel quotes. Ex """#password;"""
password = xxx
;cert_file =
;key_file =
skip_verify = true
from_address = xxx@qq.com
from_name = Grafana
# EHLO identity in SMTP dialog (defaults to instance_name)
;ehlo_identity = dashboard.example.com

```

#### 3.2在grafana里创建通知通道
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190426161923601.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)


#### 3.3在监控的指标下面配置触发告警的规则
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190426161944161.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)
#### 3.4然后test rules
触发告警规则，我的邮箱里收到了通知邮件
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190426162015793.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)
### 总结
1.基于prometheus+grafana组合，可以很容易创建高大上的监控告警系统。
2.我们对于go-micro微服务，基于中间件机制，我们可以跟容易为我们的微服务加上监控
3.对于监控的指标，go-plugins的prometheus中间件里的指标，感觉可能对于请求的延迟时间，我们关心的更多一些。
4.后期做k8s的docker集群监控时还会再讲prometheus的使用。