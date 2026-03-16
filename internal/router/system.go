package router

import (
	"github.com/gin-gonic/gin"

	"gra/internal/system"
)

// registerSystem 注册系统管理域路由
func registerSystem(auth *gin.RouterGroup, h *system.Handlers) {
	h.User.RegisterRoutes(auth)
	// 新增模块 = 加一行：
	// h.Role.RegisterRoutes(auth)
	// h.Menu.RegisterRoutes(auth)
	// h.Dept.RegisterRoutes(auth)
}
