package user

import (
	"gra/internal/public"
	"gra/internal/system/role"
	"time"
)

// User 用户模型
type User struct {
	public.BaseModel
	Username string `json:"username" gorm:"size:64;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:128;not null"`
	Nickname string `json:"nickname" gorm:"size:64"`
	Avatar   string `json:"avatar" gorm:"size:255"`
	Email    string `json:"email" gorm:"size:128"`
	Phone    string `json:"phoneNumber" gorm:"size:20"`
	Status   int8   `json:"status" gorm:"default:1;comment:1-启用 0-禁用"`
	RoleID   int64  `json:"role_id,string" gorm:"index"`
	DeptId   int64  `json:"dept_id,string" gorm:"index"`
	//注意：这里的 foreignKey 指的是 User 里的字段名，References 指的是 Role 里的字段名
	Roles         []role.SysRole `json:"roles" gorm:"many2many:sys_role_user;foreignKey:ID;joinForeignKey:UserId;References:ID;joinReferences:RoleId"`
	LastLoginDate *time.Time     `json:"lastLoginDate" gorm:"comment:最后一次登录日期"`
}

func (User) TableName() string { return "users" }

// DTO

type CreateReq struct {
	Username string               `json:"username" binding:"required,min=3,max=64"`
	Password string               `json:"password" binding:"required,min=6"`
	Nickname string               `json:"nickname"`
	Email    string               `json:"email" binding:"omitempty,email"`
	Phone    string               `json:"phone"`
	DeptId   int64                `json:"deptId,string" `
	RoleIds  []public.StringInt64 `json:"roleIds"`
}

type UpdateReq struct {
	Nickname string               `json:"nickname"`
	Email    string               `json:"email"`
	Phone    string               `json:"phoneNumber"`
	Status   *int8                `json:"status"`
	DeptId   int64                `json:"deptId,string" `
	RoleIds  []public.StringInt64 `json:"roleIds"`
}

type PageReq struct {
	Page int `form:"page,default=1" binding:"min=1"`
	Size int `form:"size,default=10" binding:"min=1,max=100"`
}

func (p *PageReq) Offset() int {
	return (p.Page - 1) * p.Size
}

type DeptQueryReq struct {
	DeptId   int64  `form:"deptId"`
	UserName string `form:"username"`
	Phone    string `form:"phoneNumber"`
	// 嵌套分页结构体
	PageReq
}

// 获取用户信息返回的结构体
type UserInfoRes struct {
	UserInfo UserInfoResponse `json:"userInfo"`
}

type UserInfoResponse struct {
	ID       int64  `json:"id,string"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type UserDetail struct {
	public.BaseModel
	Username string               `json:"username"`
	Password string               `json:"-"`
	Nickname string               `json:"nickname"`
	Avatar   string               `json:"avatar"`
	Email    string               `json:"email"`
	Phone    string               `json:"phoneNumber"`
	Status   int8                 `json:"status"`
	DeptId   int64                `json:"deptId,string"`
	RoleIds  []public.StringInt64 `json:"roleIds"`
}
