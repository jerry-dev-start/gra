package public

// SysRoleMenu 角色信息表
type SysRoleMenu struct {
	RoleId int64 `gorm:"column:role_id"`
	MenuId int64 `gorm:"column:menu_id"`
}

// TableName 指定表名
func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
