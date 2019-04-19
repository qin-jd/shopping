package repository

import (
"github.com/jinzhu/gorm"
"shopping/product/model"
)

type Repository interface {
	Find(id int32) (*model.Product, error)
	Create(*model.Product) error
	Update(*model.Product, int64) (*model.Product, error)
	FindByField(string, string, string) (*model.Product, error)
}

type Product struct {
	Db *gorm.DB
}

func (repo *Product) Find(id uint32) (*model.Product, error) {
	product :=  &model.Product{}
	product.ID = uint(id)
	if err := repo.Db.First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *Product) Create(product *model.Product) error {
	if err := repo.Db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (repo *Product) Update(product *model.Product) (*model.Product, error) {
	if err := repo.Db.Model(&product).Updates(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (repo *Product) FindByField(key string, value string, fields string) (*model.Product, error) {
	if len(fields) == 0 {
		fields = "*"
	}
	product :=  &model.Product{}
	if err := repo.Db.Select(fields).Where(key+" = ?", value).First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}
