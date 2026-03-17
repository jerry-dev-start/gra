package role_menu

import (
	"gra/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RoleCheckMenuId 获取到角色已邦定的菜单Id
func (h *Handler) RoleCheckMenuId(c *gin.Context) {
	roleId := c.Param("roleId")
	result, err := h.svc.GetRoleBindMenuId(roleId)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *Handler) SaveRoleMenu(c *gin.Context) {
	roleId := c.Param("roleId")
	var saveRoleMenuReq SaveRoleMenuReq
	if err := c.ShouldBindJSON(&saveRoleMenuReq); err != nil {
		response.FailMsg(c, err.Error())
		return
	}

	if err := h.svc.SaveRoleMenu(roleId, saveRoleMenuReq); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OKMsg(c, "设置成功")
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	roles := r.Group("/roles")
	{
		roles.GET("/:roleId/menus", h.RoleCheckMenuId)
		roles.PUT("/:roleId/menus", h.SaveRoleMenu)
	}
}
