package auth

import (
	"github.com/gin-gonic/gin"

	"gra/pkg/response"
	"gra/pkg/validate"
)

// 登录校验规则
var loginRules = validate.Rules{
	"Username": {validate.Required("用户名不能为空")},
	"Password": {validate.Required("密码不能为空")},
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := validate.Check(req, loginRules); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	resp, err := h.svc.Login(&req)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	response.OK(c, resp)
}

// RegisterPublicRoutes 注册公开路由（无需认证）
func (h *Handler) RegisterPublicRoutes(r *gin.RouterGroup) {
	r.POST("/auth/login", h.Login)
}
