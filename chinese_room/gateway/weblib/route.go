package weblib

import (
	"gateway/weblib/handler"
	"gateway/weblib/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{})*gin.Engine{
	engine:=gin.Default()
	engine.Use(middleware.Cors(),middleware.InitMiddleware(service),middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte("something-very-secret"))
	engine.Use(sessions.Sessions("mysessions",store))
	v1:=engine.Group("api/v1")
	{
		//测试
		v1.GET("ping",func(c *gin.Context){
			c.JSON(200,"success")
		})
		//用户服务
		v1.POST("/user/register",handler.UserRegister)
		v1.POST("/user/login",handler.UserLogin)
		// 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware.Jwt())
		{


		}
	}
	return engine

}
