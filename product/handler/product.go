package handler

import (
	"context"
	"fmt"
	"shopping/product/model"
	"shopping/product/repository"

	"github.com/micro/go-log"

	product "shopping/product/proto/product"
)

type Product struct {
	Repo *repository.Product
}

func (e *Product) Search (ctx context.Context , req *product.SearchRequest , rsp *product.SearchResponse) error{
	log.Log("Received Product.Search request")
	var products []*product.Product
	if err := e.Repo.Db.Where("name like ?" , "%"+req.Name+"%").Limit(10).Find(&products).Error; err != nil{
		return err
	}

	rsp.Code = "200"
	rsp.Msg = fmt.Sprintf("搜索结果共%v条" , len(products))
	rsp.Products = products
	return nil
}


func (e *Product) Detail (ctx context.Context , req *product.DetailRequest , rsp *product.DetailResponse) error{
	log.Log("Received Product.Detail request")
	var product = &product.Product{}
	if err := e.Repo.Db.Where("id = ?" , req.Id).First(&product).Error; err != nil{
		return err
	}

	rsp.Code = "200"
	rsp.Msg = "商品详情如下："
	rsp.Product = product

	return nil
}


func (e *Product) ReduceNumber (ctx context.Context , req *product.ReduceNumberRequest , rsp *product.Response) error{
	log.Log("Received Product.Detail request")

	var product = model.Product{}
	product.ID = uint(req.Id)

	if err := e.Repo.Db.First(&product).Error; err != nil{
		return err
	}

	product.Number -= 1
	log.Log("库存数量为:" , product.Number)
	if err := e.Repo.Db.Model(&model.Product{}).Where("id = ?", product.ID).Update("number" , product.Number).Error; err != nil{
		return err
	}

	rsp.Code = "200"
	rsp.Msg = fmt.Sprintf("库存更新成功,更新后的数量为%v" , product.Number)
	return nil
}