package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	controller_room "wchatv1/controller/room"
)

var _ Router = &Tenant{}

type Tenant struct {
}

func (t *Tenant) Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println(ctx)

		//请求方法
		method := ctx.Request.Method

		// 允许任何源
		ctx.Header("Access-Control-Allow-Origin", "*")
		//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		//ctx.Header("Access-Control-Allow-Headers", "Token,Content-Type")
		ctx.Header("Access-Control-Allow-Headers", "*")
		// 跨域关键设置 让浏览器可以解析
		ctx.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		ctx.Next() //  处理请求
	}
}

func (t *Tenant) Route(engine *gin.Engine) {
	room := engine.Group("/room")
	{
		room.POST("/makePassToken", controller_room.MakePassToken)
	}
}
