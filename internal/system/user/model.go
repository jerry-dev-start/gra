package user

import (
	"gra/internal/public"
	"time"
)

// User 用户模型
type User struct {
	public.BaseModel
	Username      string    `json:"username" gorm:"size:64;uniqueIndex;not null"`
	Password      string    `json:"-" gorm:"size:128;not null"`
	Nickname      string    `json:"nickname" gorm:"size:64"`
	Avatar        string    `json:"avatar" gorm:"size:255"`
	Email         string    `json:"email" gorm:"size:128"`
	Phone         string    `json:"phoneNumber" gorm:"size:20"`
	Status        int8      `json:"status" gorm:"default:1;comment:1-启用 0-禁用"`
	RoleID        int64     `json:"role_id,string" gorm:"index"`
	LastLoginDate time.Time `json:"lastLoginDate" gorm:"comment:最后一次登录日期"`
}

func (User) TableName() string { return "users" }

// DTO

type CreateReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	RoleID   int64  `json:"role_id,string"`
}

type UpdateReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Status   *int8  `json:"status"`
	RoleID   int64  `json:"role_id,string"`
}

type PageReq struct {
	Page int `form:"page,default=1" binding:"min=1"`
	Size int `form:"size,default=10" binding:"min=1,max=100"`
}

func (p *PageReq) Offset() int {
	return (p.Page - 1) * p.Size
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
