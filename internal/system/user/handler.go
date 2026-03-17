package user

import (
	"gra/pkg/validate"
	"strconv"

	"github.com/gin-gonic/gin"

	"gra/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	err := validate.Check(req, validate.Rules{
		"Username": {validate.Required("用户名不能为空")},
		"Password": {validate.Required("密码不能为空")},
		"Nickname": {validate.Required("昵称不能为空")},
	})
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	if err := h.svc.Create(&req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OKMsg(c, "创建成功")
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailMsg(c, "无效ID")
		return
	}
	u, err := h.svc.GetByID(id)
	if err != nil {
		response.FailMsg(c, "用户不存在")
		return
	}
	response.OK(c, u)
}

// GetInfo 通过Token获取用户信息
func (h *Handler) GetInfo(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.FailMsg(c, "解析Token未获取到UserId")
		return
	}
	var req UserInfoRes
	userInfo, err := h.svc.GetByID(userId.(int64))
	if err != nil {
		response.FailMsg(c, "获取用户信息失败")
		return
	}
	req.UserInfo = UserInfoResponse{
		ID:       userInfo.ID,
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.Avatar,
	}
	response.OK(c, req)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailMsg(c, "无效ID")
		return
	}
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailMsg(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Update(id, &req); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OKMsg(c, "更新成功")
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效ID")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OKMsg(c, "删除成功")
}

func (h *Handler) List(c *gin.Context) {
	var req PageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	users, total, err := h.svc.List(&req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OKPage(c, users, total, req.Page, req.Size)
}

// RegisterRoutes 注册需认证路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("", h.Create)
		users.GET("", h.List)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)

		users.GET("/profile", h.GetInfo)
	}
}
