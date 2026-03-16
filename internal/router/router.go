package router

import (
	"github.com/gin-gonic/gin"

	"gra/internal/business"
	"gra/internal/middleware"
	"gra/internal/system"
)

func Setup(r *gin.Engine, sysH *system.Handlers, bizH *business.Handlers) {
	r.Use(middleware.Cors())

	api := r.Group("/api")

	// 公开接口
	pub := api.Group("")
	registerPublic(pub, sysH)

	// 需认证接口
	auth := api.Group("", middleware.JWTAuth())
	registerSystem(auth, sysH)
	registerBusiness(auth, bizH)
}

// registerPublic 注册公开路由
func registerPublic(pub *gin.RouterGroup, h *system.Handlers) {
	h.Auth.RegisterPublicRoutes(pub)
}
