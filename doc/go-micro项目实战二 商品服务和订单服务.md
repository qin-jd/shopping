### 商品服务
提供商品列表、详情、库存更新等服务

### 订单服务
提供订单提交、订单查询、状态变更等服务

### 创建服务
创建服务的过程和创建用户服务的过程一样。这里就不赘述了。如需查看源码可以参考https://github.com/qin-jd/shopping

### 服务间交互
例如在订单服务中，提交订单时需要先去商品服务那查询库存数量，然后生成订单，再对商品进行减库存操作。
我们需要在订单服务的handler里拿到商品服务的客户端然后进行RPC调用。

第一个坑：因为我们是把所有的服务都合并到了一个文件夹管理，同时又没有把该服务传到github,完全在我本地开发。当我在本地的订单服务想要调取商品服务里的接口时，我不知道怎么用go mod来处理代码的依赖。直接在import 里写`product "shopping/product/proto/product"`不生效
解决方案：在go.mod里对shopping/product进行了replace。指到我本地的代码里。即如下所示：
```
require (
	...
	shopping/product v0.0.0
)

replace shopping/product => ../product
```

第二个坑：我看到网上所有的例子都是利用protoc-gen-micro生成的代码里都有对应服务的客户端，比如ProductServiceClient，为啥我的product.micro.go里面没有xxxClint这个结构体呢，难道我装了假的代码生成器？
解决方案：我查了protoc-gen-micro代码提交版本发现了一个Move Client to Service的PR。This PR moves anything with Client to Service.原来新版的代码把所有的client操作都移动到了service里。所以新版本生成的代码没办法使用ProductServiceClient来进行RPC调用，我们依然用ProductService。即如下所示
在main.go里添加商品服务的客户端
```
//product-srv client
productCli := product.NewProductService("go.micro.srv.product" , service.Client())
```
然后再把商品服务的客户端在注册订单服务的handler时传进去
```
// Register Handler
order.RegisterOrderServiceHandler(service.Server(), &handler.Order{repo , productCli})
```
接下来就能在handler里使用商品服务的客户端了
```
func (h *Order) Submit (ctx context.Context , req *order.SubmitRequest, rsp *order.Response) error{
	log.Log("Received Product.Search request")

	//查询商品的库存数量
	productDetail,err := h.ProductCli.Detail(context.TODO() , &product.DetailRequest{Id:req.ProductId})
	if productDetail.Product.Number < 1 {
		return errors.BadRequest("go.micro.srv.order" , "库存不足")
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}

	// Generate a snowflake ID.
	orderId := node.Generate().String()

	order := &model.Order{
		Status:1,
		OrderId:orderId,
		ProductId:req.ProductId,
	}

	if err = h.Order.Create(order); err != nil{
		return err
	}

	//减库存
	reduce,err := h.ProductCli.ReduceNumber(context.TODO() , &product.ReduceNumberRequest{Id:req.ProductId})
	if reduce == nil || reduce.Code != "200" {
		return errors.BadRequest("go.micro.srv.order" , err.Error())
	}

	rsp.Code = "200"
	rsp.Msg = "订单提交成功"
	return nil
}
```

