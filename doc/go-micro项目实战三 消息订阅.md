### 通知服务
提供发送通知服务

### 新建服务
`micro new shopping/notification`
该服务比较简单，只实验了消息发布和订阅的功能，未提供真正通知的逻辑。
实现的功能是：订单提交成功后，通知用户订单已经提交。

### 消息代理
基于go-micro强大的插件机制。go-plugins内置了诸如grpc,rabbitmq，nats,redis。几乎可以在这些代理之间无缝切换。
本次实验采用rabbitmq代理来实现基本功能。

### 修改main.go
```
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

```
### 创建事件的接口约束
shopping/notification/proto/notification/notification.proto
```
syntax = "proto3";

import "user/proto/user/user.proto";
import "product/proto/product/product.proto";

package go.micro.srv.notification;

service ProductService {
	rpc NotifyOrder (NotifyOrderRequest) returns (Response){}
}

message NotifyOrderRequest{
	go.micro.srv.user.User user = 1;
	go.micro.srv.product.Product product = 2;
}

message SubmitRequest {
	uint32 productId = 1;
	uint32 count = 2;
	uint32 uid = 3;
}

message Response{
	string code = 1;
	string msg = 2;
}
```

### 修改事件的handler
shopping/notification/subscriber/notification.go
```
package subscriber

import (
	"context"
	"fmt"
	"github.com/micro/go-log"

	notification "shopping/notification/proto/notification"
)

type Notification struct{}

func (e *Notification) Handle(ctx context.Context, req *notification.SubmitRequest) error {
	log.Log(fmt.Sprintf("Handler Received message: %v 购买了商品ID为：%v 的物品" , req.Uid , req.ProductId))
	//执行通知的逻辑
	
	return nil
}

```
### 启动事件订阅者服务
`go run main.go plugin.go`

### 修改订单服务，提交订单的过程中发布消息
shopping/order/main.go，引入rabbitmq的消息代理插件，然后创建rabbitmq的broker，提供给micro
```
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
	
	//创建消息发布者
	publisher := micro.NewPublisher("notification.submit" , service.Client())
	
	//在注册订单handler里传进去publisher
	order.RegisterOrderServiceHandler(service.Server(), &handler.Order{repo , productCli , publisher})
	
```
shopping/order/handler/order.go里增加发布消息的操作
```
	//异步发送通知给用户订单信息
	if err := h.Publisher.Publish(ctx , req);err != nil {
		return errors.BadRequest("notification" , err.Error())
	}
```
### 运行订单服务
`go run main.go database.go plugin.go --broker=rabbitmq`

### 测试
请求接口
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190424141624953.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)
消息订阅者那里
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190424141635966.png)

### 总结
go-micro的插件机制很强大，我们可以随意更换消息的代理实现，从rabbitmq切换到别的实现几乎不用动代码。
使用异步消息用于实现功能解耦。基于go-micro我们不用考虑消息的发送和协议转换操作。专注于消息本身即可。
注意：
笔者在go-plugins v1.1.0版本，在使用grpc Broker时有bug。
https://github.com/micro/go-micro/issues/459
暂时无法正常使用grpc Broker。坐等新版本升级后再重新尝试