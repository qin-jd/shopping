### 用户服务
用户服务，提供登录、注册、修改密码等功能。

### 新建服务
`micro new shopping/user`

### 整理结构
增加model目录和repository目录，删掉proto里默认的example文件夹，创建user文件夹。
开发步骤：1.定义接口 -> 2.生成接口代码 -> 3.编写model层代码 -> 4.编写repository数据操作代码 -> 5.实现接口 -> 6.修改main.go 
因为这是第一个微服务，所以会尽可能地详细描述创建微服务的过程

### 定义用户服务接口
```
syntax = "proto3";

package go.micro.srv.user;

service UserService {
	rpc Register (RegisterRequest) returns (Response){}
	rpc Login (LoginRequest) returns (Response){}
	rpc UpdatePassword (UpdatePasswordRequest) returns (Response){}
}

message User {
	uint32 id = 1;
	string name = 2;
	string phone = 3;
	string password = 4;
}

message RegisterRequest{
	User user = 1;
}

message LoginRequest{
	string phone = 1;
	string password = 2;
}

message UpdatePasswordRequest{
	uint32 uid = 1;
	string oldPassword = 2;
	string newPassword = 3;
	string confirmPassword = 4;
}

message Response {
	string code = 1;
	string msg = 2;
}


```
### 生成接口代码
执行命令`protoc --proto_path=. --micro_out=. --go_out=. proto/user/user.proto`

### 编写model.user
```
package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name string
	Phone string `gorm:"type:char(11);`
	Password string
}
```

### 编写repository.user
```
package repository

import (
	"github.com/jinzhu/gorm"
	"shopping/user/model"
)

type Repository interface {
	Find(id int32) (*model.User, error)
	Create(*model.User) error
	Update(*model.User, int64) (*model.User, error)
	FindByField(string, string, string) (*model.User, error)
}

type User struct {
	Db *gorm.DB
}

func (repo *User) Find(id uint32) (*model.User, error) {
	user :=  &model.User{}
	user.ID = uint(id)
	if err := repo.Db.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) Create(user *model.User) error {
	if err := repo.Db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *User) Update(user *model.User) (*model.User, error) {
	if err := repo.Db.Model(user).Updates(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *User) FindByField(key string, value string, fields string) (*model.User, error) {
	if len(fields) == 0 {
		fields = "*"
	}
	user :=  &model.User{}
	if err := repo.Db.Select(fields).Where(key+" = ?", value).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

```

### 实现接口handler.user
```
package handler

import (
	"context"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/errors"
	"golang.org/x/crypto/bcrypt"
	"shopping/user/model"
	"shopping/user/repository"

	user "shopping/user/proto/user"
)

type User struct{
	Repo *repository.User
}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) Register(ctx context.Context, req *user.RegisterRequest, rsp *user.Response) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:req.User.Name,
		Phone:req.User.Phone,
		Password:string(hashedPwd),
	}

	if err := e.Repo.Create(user);err != nil{
		log.Log("create error")
		return err
	}

	rsp.Code = "200"
	rsp.Msg = "注册成功"

	return nil
}

func (e *User)Login(ctx context.Context, req *user.LoginRequest, rsp *user.Response) error {
	user , err := e.Repo.FindByField("phone" , req.Phone , "id , password")
	if err != err{
		return err
	}

	if user == nil{
		return errors.Unauthorized("go.micro.srv.user.login", "该手机号不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.login", "密码错误")
	}

	rsp.Code = "200"
	rsp.Msg = "登录成功"

	return nil
}

func (e *User)UpdatePassword(ctx context.Context, req *user.UpdatePasswordRequest, rsp *user.Response) error {
	user,err := e.Repo.Find(req.Uid)

	if user == nil{
		return errors.Unauthorized("go.micro.srv.user.login", "该用户不存在")
	}

	if err != nil {
		return err
	}
	//验证老密码是否正常
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.login", "旧密码认证失败")
	}

	//验证通过后，对新密码hash存下来
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPwd)
	e.Repo.Update(user)

	rsp.Code = "200"
	rsp.Msg =user.Name+"，您的密码更新成功"

	return nil
}
```

### 修改main.go
因为需要在main.go里建立数据库连接，在根目录下创建database.go
```
package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CreateConnection() (*gorm.DB, error) {
	host := "192.168.0.111"
	user := "mytestroot"
	dbName := "shopping"
	password := "mytestroot"
	return gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, dbName,
	),
	)
}

```
然后修改main.go
```
package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-grpc"
	"shopping/user/handler"
	"shopping/user/model"
	"shopping/user/repository"

	user "shopping/user/proto/user"
)

func main() {

	db,err := CreateConnection()
	defer db.Close()

	db.AutoMigrate(&model.User{})

	if err != nil {
		log.Fatalf("connection error : %v \n" , err)
	}

	repo := &repository.User{db}

	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	user.RegisterUserServiceHandler(service.Server(), &handler.User{repo})

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.user", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	//micro.RegisterSubscriber("go.micro.srv.user", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

```
### 启动服务，验证接口
`go run main.go database.go`和`micro api --namespace=go.micro.srv`
注册接口
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190424141213508.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)
登录接口
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190424141230667.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)
修改密码接口
![在这里插入图片描述](https://img-blog.csdnimg.cn/20190424141245331.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3UwMTM3MDUwNjY=,size_16,color_FFFFFF,t_70)