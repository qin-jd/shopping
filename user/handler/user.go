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