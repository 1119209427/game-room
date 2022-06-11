package core

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"user/model"
	"user/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}
func(us *UserService)UserRegister(ctx context.Context ,req *pb.UserRequest)(*pb.UserDetailResponse,error){

	var resp pb.UserDetailResponse
	var count int64
	resp.Code=200
	if req.PassWord!=req.PassWordConfirm{
		resp.Code=400
		return &resp,errors.New("两次密码不同，请重新输入")
	}

	if err:=model.DB.Model(&model.User{}).Where("user_name",req.UserName).Count(&count).Error;err!=nil{
		/*if err.Error()=="record not found"{//说明不重复，则注册

		}*/
		resp.Code=400
		return &resp,err
	}
	if count>0{
		resp.Code=500
		err:=errors.New("用户名已存在")
		return &resp, err
	}
	user:=model.User{UserName: req.UserName}
	user.SetPassword(req.PassWord)
	if err:=model.DB.Create(&user).Error;err!=nil{
		resp.Code=400
		return &resp,err
	}
	model:=BuildUser(user)
	resp.UserDetail=model
	return &resp,nil
	return &resp,nil


}
func(us *UserService)UserLogin(ctx context.Context ,req *pb.UserRequest)(*pb.UserDetailResponse,error){
	var user model.User
	var resp pb.UserDetailResponse
	resp.Code=200
	if err:=model.DB.Where("user_name",req.UserName).First(&user).Error;err!=nil{
		if gorm.ErrRecordNotFound==err{
			resp.Code=400
			err=errors.New("请注册后再登录")
			return &resp,nil
		}
		resp.Code=500
		return &resp,err
	}
	flag:=user.CheckPassword(req.PassWord)
	if flag==false{
		resp.Code=400
		err:=errors.New("密码错误，请重新输入")
		return &resp,err
	}
	model:=BuildUser(user)
	resp.UserDetail=model
	return &resp,nil

}

