package model

type SysRoleUser struct {
	RoleId int64 `gorm:"column:role_id"`
	UserId int64 `gorm:"column:user_id"`
}

// TableName 指定表名
func (SysRoleUser) TableName() string {
	return "sys_role_user"
}
