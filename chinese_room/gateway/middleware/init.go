package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware 错误处理
func ErrorMiddleware()gin.HandlerFunc{
	return func(c *gin.Context){
		defer func(){
			if r:=recover();r!=nil{
				c.JSON(200,gin.H{
					"code":404,
					"msg":fmt.Sprintf("%s",r),
				})
				c.Abort()
			}
			c.Next()

		}()
	}
}

// InitMiddleware 接受服务实例，并存到gin.Key中
func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		// 将实例存在gin.Keys中
		context.Keys = make(map[string]interface{})
		context.Keys["userService"] = service[0]
		context.Keys["taskService"] = service[1]
		context.Next()
	}
}