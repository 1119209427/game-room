package handler

import (
	"context"
	"gateway/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(c *gin.Context){
	var userReq pb.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	//从gin.key中取出服务
	userService:=c.Keys["userService"].(pb.UserServiceServer)
	userResp,err:=userService.UserRegister(context.Background(),&userReq)
	PanicIfUserError(err)
	c.JSON(http.StatusOK,gin.H{"msg":userResp})

}
func UserLogin(c *gin.Context){
	var userReq pb.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	//从gin.key中取出服务
	userService:=c.Keys["userService"].(pb.UserServiceServer)
	userResp,err:=userService.UserLogin(context.Background(),&userReq)
	PanicIfUserError(err)
	c.JSON(http.StatusOK,gin.H{"msg":userResp})

}
