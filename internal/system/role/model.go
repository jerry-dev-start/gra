package role

import "gra/internal/public"

// SysRole 角色信息表
type SysRole struct {
	public.BaseModel
	RoleName    string `gorm:"column:role_name;type:varchar(50);not null;comment:角色名称" json:"roleName"`
	RoleCode    string `gorm:"column:role_code;type:varchar(50);unique;not null;comment:角色编码" json:"roleCode"`
	Description string `gorm:"column:description;type:varchar(255);comment:描述" json:"description"`
	SortOrder   int    `gorm:"column:sort_order;type:int;default:0;comment:显示排序" json:"sortOrder"`
	Status      int8   `gorm:"column:status;type:tinyint(1);default:1;comment:状态（1正常 0停用）" json:"status"`
	IsReadonly  int8   `gorm:"column:is_readonly;type:tinyint(1);default:0;comment:是否系统内置（1是 0否）" json:"isReadonly"`
}

// TableName 指定表名
func (SysRole) TableName() string {
	return "sys_role"
}
