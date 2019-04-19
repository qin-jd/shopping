### 项目简介
该项目是我自学go-micro的练手项目。作为一个用来学习的小型的电商项目。该项目会包括以下微服务：
user-srv  用户服务，提供注册、登录、修改密码等操作
product-srv 商品服务，提供商品搜索，商品详情，减少库存等操作
order-srv 订单服务，提供提交订单、查看订单详情等操作

代码比较糙，因为我完全没有系统地学习过go语言。只是花了两个小时看了一下go的语法，来不及细致地学这门语言了，就上手做项目来顺便熟悉语法。

项目版本库地址
https://github.com/qin-jd/shopping

### 用到的包
go-micro、protobuf、gorm

### 学习计划
- go-micro服务的熟练使用
- 服务间的相互调用
- 服务的全链路监控、
- 实验go-config
- 单元测试
- 实验服务的熔断，降级等机制
- go-micro微服务的容器化部署
- go-micro微服务最佳实践

### 每个微服务的代码组织
在使用go-micro开发时，关于如何组织自己的代码结构时，我查阅了很多资料：
https://mp.weixin.qq.com/s/yh2WtuK2m-bDE6oVflXM-Q 
https://github.com/golang-standards/project-layout

方式一：单一package，类似于gorm这种，所有代码都在package gorm下面。这种提倡的就是简洁高效。符合go的特点，我不是很喜欢这种。我是那种一看到一个包里的文件很多很多的时候就头疼，宁愿放到多个文件夹。

方式二：MVC这种，经典。之前做PHP的,都了解优雅的框架laravel，类似这种的。

方式三：按模块划分

我能想到我们的应用场景，虽然是拿go-micro来做的微服务，但是每个微服务里的代码绝壁不会是简单的几个文件就能搞定的，而且同一个微服务也不会只涉及到一张表。又加上我们是从PHP转go过来的，开发的脑子里想的还是离不开mvc。基于这方面的考虑，我最后根据go-micro的模板，增加了model目录和repository目录来组织我们的代码结构,以用户服务为例：
```
├── handler
│   └── user.go
├── model
│   └── user.go
├── proto
│   ├── user
│   │   └── user.proto
├── repository
│   └── user.go
├── subscriber
│   ├── subscriber.go
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── plugin.go
├── README.md
```

