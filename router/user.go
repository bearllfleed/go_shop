package router

import (
	"github.com/bearllflee/go_shop/api"

	"github.com/gin-gonic/gin"
)

type UserGroup struct{}

func (u *UserGroup) InitUserRouters(engine *gin.Engine) {
	userRouters := engine.Group("user")
	userRouters.POST("login", api.Login)
	userRouters.POST("register", api.Register)
	userRouters.POST("list", api.UserList)
	userRouters.GET("ws", api.OnlineTool)
}
