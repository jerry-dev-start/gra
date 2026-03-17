package menus

import (
	"gra/pkg/validate"
	"strconv"

	"github.com/gin-gonic/gin"

	"gra/pkg/response"
)

// 登录校验规则
var createMenuRules = validate.Rules{
	"Name": {validate.Required("菜单名称不能为空")},
	"Path": {validate.Required("菜单路径不能为空")},
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailMsg(c, "参数错误: "+err.Error())
		return
	}
	if err := validate.Check(req, createMenuRules); err != nil {
		response.FailMsg(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Create(&req); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OKMsg(c, "创建成功")
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效ID")
		return
	}
	m, err := h.svc.GetByID(id)
	if err != nil {
		response.Fail(c, 404, "菜单不存在")
		return
	}
	response.OK(c, m)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, 400, "无效ID")
		return
	}
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Update(id, &req); err != nil {
		response.Fail(c, 500, err.Error())
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

func (h *Handler) ListTree(c *gin.Context) {
	tree, err := h.svc.ListTree()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.OK(c, tree)
}
func (h *Handler) UserMenuTree(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, 401, "用户未登录")
		return
	}
	userID, ok := val.(int64)
	if !ok {
		response.Fail(c, 401, "用户信息异常")
		return
	}

	result, err := h.svc.UserMenuTree(userID)
	if err != nil {
		response.FailMsg(c, "获取用户菜单失败")
		return
	}
	response.OK(c, result)
}

// RegisterRoutes 注册需认证路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	m := r.Group("/menus")
	{
		m.POST("", h.Create)
		m.GET("/tree", h.ListTree)
		m.GET("/user/tree", h.UserMenuTree)
		m.GET("/:id", h.GetByID)
		m.PUT("/:id", h.Update)
		m.DELETE("/:id", h.Delete)
	}
}
