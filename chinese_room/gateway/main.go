package main

import (
	"gateway/chat_room/core"
	"gateway/handler"
	"gateway/middleware"
	"gateway/pb"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

)
const (
	address = "localhost:8090"
)

func main(){
	//建立连接
	conn,err:=grpc.Dial(address,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err!=nil{
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c:=pb.NewUserServiceClient(conn)
	go core.H.Run()
	engine:=gin.Default()
	engine.Use(middleware.Cors(),middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte("something-very-secret"))
	engine.Use(sessions.Sessions("mysessions",store))
	engine.GET("ws",core.ServerWs)
	engine.GET("chess",core.ChessWs)
	v1:=engine.Group("api/v1")
	{
		//测试
		v1.GET("ping",func(c *gin.Context){
			c.JSON(200,"success")
		})
		//用户服务
		v1.POST("/user/register",func(ctx *gin.Context){
			var userReq pb.UserRequest
			handler.PanicIfUserError(ctx.Bind(&userReq))
			userResp,err:=handler.UserRegister(c,&userReq)
			handler.PanicIfUserError(err)
			ctx.JSON(http.StatusOK,gin.H{"msg":userResp})
		})
		v1.POST("/user/login",func(ctx *gin.Context){
			var userReq pb.UserRequest
			handler.PanicIfUserError(ctx.Bind(&userReq))
			userResp,err:=handler.UserLogin(c,&userReq)
			handler.PanicIfUserError(err)
			ctx.JSON(http.StatusOK,gin.H{"msg":userResp})
		})
		// 需要登录保护
		authed := v1.Group("/")
		/*authed.GET("ws",core.ServerWs)//消息读取
		authed.GET("chess",core.ChessWs)*/
		authed.Use(middleware.Jwt())
		{


	}

	}
	if err := engine.Run(":8052"); err != nil { 
		log.Fatalf("could not run server: %v", err)
	}
}
