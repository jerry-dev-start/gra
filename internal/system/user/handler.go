package user

import (
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
	if err := h.svc.Create(&req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OKMsg(c, "创建成功")
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效ID")
		return
	}
	u, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, 404, "用户不存在")
		return
	}
	response.OK(c, u)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效ID")
		return
	}
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Update(uint(id), &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OKMsg(c, "更新成功")
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效ID")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
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
	}
}
