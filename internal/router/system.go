package router

import (
	"github.com/gin-gonic/gin"

	"gra/internal/system"
)

// registerSystem 注册系统管理域路由
func registerSystem(auth *gin.RouterGroup, h *system.Handlers) {
	h.Auth.RegisterRoutes(auth)
	h.User.RegisterRoutes(auth)
	h.Menu.RegisterRoutes(auth)
	h.Role.RegisterRoutes(auth)
	h.RoleMenu.RegisterRoutes(auth)
	h.Dept.RegisterRoutes(auth)
}
