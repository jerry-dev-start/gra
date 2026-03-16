package router

import (
	"github.com/gin-gonic/gin"

	"gra/internal/business"
)

// registerBusiness 注册业务域路由
func registerBusiness(auth *gin.RouterGroup, h *business.Handlers) {
	// 新增模块 = 加一行：
	// h.Order.RegisterRoutes(auth)
	// h.Product.RegisterRoutes(auth)
	_ = auth
	_ = h
}
