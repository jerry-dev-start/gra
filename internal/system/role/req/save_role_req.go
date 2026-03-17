package req

type SaveRoleReq struct {

	// RoleName 角色名称
	// binding:"required" 是 Gin 等框架常用的校验标签，确保字段不为空
	RoleName string `json:"roleName"`

	// RoleCode 角色编码
	RoleCode string `json:"roleCode"`

	// Description 角色描述
	Description string `json:"description"`

	// SortOrder 显示排序
	SortOrder int `json:"sortOrder" default:"0"`

	// Status 角色状态 (1:正常, 0:停用)
	Status int8 `json:"status"`
}
