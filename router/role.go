package router

import (
	"github.com/bearllflee/go_shop/api"

	"github.com/gin-gonic/gin"
)

type RoleGroup struct{}

func (r *RoleGroup) InitRoleRouters(engine *gin.Engine) {
	roleRouters := engine.Group("role")
	roleRouters.POST("list", api.RoleList)
	roleRouters.POST("", api.RoleCreate)
}
