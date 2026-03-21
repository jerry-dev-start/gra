package dept

import (
	"gra/pkg/response"
	"gra/pkg/validate"

	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// DeptTree 获取到部门列表树
func (h *Handler) DeptTree(r *gin.Context) {
	var deptReq DeptQueryReq
	if err := r.ShouldBindQuery(&deptReq); err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	tree, err := h.svc.DeptTree(deptReq)
	if err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	response.OK(r, tree)
}

func (h *Handler) Create(r *gin.Context) {
	var deptReq DeptReq
	if err := r.ShouldBindJSON(&deptReq); err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	err := validate.Check(deptReq, validate.Rules{
		"ParentID": {validate.Required("父级节点不能为空")},
		"Name":     {validate.Required("部门节点不能为空")},
	})
	if err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	if err := h.svc.Create(deptReq); err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	response.OKMsg(r, "保存成功")
}

func (h *Handler) GetDeptInfo(c *gin.Context) {
	// 1. 获取字符串类型的 ID
	idStr := c.Param("id")
	// 2. 如果你的数据库 ID 是 bigint/int64，需要进行类型转换
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	result, err := h.svc.GetDeptInfo(id)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *Handler) Update(r *gin.Context) {
	var deptReq DeptReq
	if err := r.ShouldBindJSON(&deptReq); err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	err := validate.Check(deptReq, validate.Rules{
		"ParentID": {validate.Required("父级节点不能为空")},
		"Name":     {validate.Required("部门节点不能为空")},
		"ID":       {validate.Required("ID不能为空")},
	})
	if err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	err = h.svc.Update(&deptReq)
	if err != nil {
		response.FailMsg(r, err.Error())
		return
	}
	response.OKMsg(r, "保存成功")
}

func (h *Handler) DeleteDept(c *gin.Context) {
	idStr := c.Param("id")
	// 2. 如果你的数据库 ID 是 bigint/int64，需要进行类型转换
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	err = h.svc.DeleteDept(id)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OKMsg(c, "删除成功")
}

// RegisterRoutes 注册需认证路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	m := r.Group("/depts")
	{
		m.POST("", h.Create)
		m.PUT("", h.Update)
		m.GET("/tree", h.DeptTree)
		m.GET("/:id", h.GetDeptInfo)
		m.DELETE("/:id", h.DeleteDept)
	}
}
