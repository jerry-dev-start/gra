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

// 刷新Token校验规则
var refreshRules = validate.Rules{
	"RefreshToken": {validate.Required("RefreshToken不能为空")},
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

// RefreshToken 通过 refresh token 换取新的 access token
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := validate.Check(req, refreshRules); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}
	resp, err := h.svc.RefreshToken(&req)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}
	response.OK(c, resp)
}

// Logout 退出登录，主动吊销 refresh token
func (h *Handler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, 401, "未登录")
		return
	}
	if err := h.svc.Logout(userID.(int64)); err != nil {
		response.Fail(c, 500, "退出失败")
		return
	}
	response.OK(c, nil)
}

// RegisterPublicRoutes 注册公开路由（无需认证）
func (h *Handler) RegisterPublicRoutes(r *gin.RouterGroup) {
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.RefreshToken)
}

// RegisterRoutes 注册需认证的路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/auth/logout", h.Logout)
}
