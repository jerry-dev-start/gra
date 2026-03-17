package role

import (
	"gra/internal/public"
	"gra/internal/system/role/req"
	"gra/pkg/response"
	"gra/pkg/validate"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Create 保存角色信息
func (h *Handler) RoleList(c *gin.Context) {
	var req public.PageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailMsg(c, "参数错误: "+err.Error())
		return
	}
	roles, total, err := h.svc.RoleList(&req)

	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OKPage(c, roles, total, req.Page, req.Size)
}
func (h *Handler) Create(c *gin.Context) {
	var roleReq req.SaveRoleReq

	if err := c.ShouldBind(&roleReq); err != nil {
		response.FailMsg(c, "参数错误:"+err.Error())
		return
	}
	//校验
	if err := validate.Check(roleReq, validate.Rules{
		"RoleName": {validate.Required("角色名称不能为空")},
		"RoleCode": {validate.Required("角色编码不能为空")},
	}); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	if err := h.svc.Create(&roleReq); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OKMsg(c, "创建成功")
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	roles := r.Group("/roles")
	{
		roles.POST("", h.Create)
		roles.GET("", h.RoleList)
	}
}
