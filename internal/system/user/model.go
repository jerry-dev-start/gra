package user

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 公共字段
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// User 用户模型
type User struct {
	BaseModel
	Username string `json:"username" gorm:"size:64;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:128;not null"`
	Nickname string `json:"nickname" gorm:"size:64"`
	Avatar   string `json:"avatar" gorm:"size:255"`
	Email    string `json:"email" gorm:"size:128"`
	Phone    string `json:"phone" gorm:"size:20"`
	Status   int8   `json:"status" gorm:"default:1;comment:1-启用 0-禁用"`
	RoleID   uint   `json:"role_id" gorm:"index"`
}

func (User) TableName() string { return "sys_user" }

// DTO

type CreateReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id"`
}

type UpdateReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Status   *int8  `json:"status"`
	RoleID   uint   `json:"role_id"`
}

type PageReq struct {
	Page int `form:"page,default=1" binding:"min=1"`
	Size int `form:"size,default=10" binding:"min=1,max=100"`
}

func (p *PageReq) Offset() int {
	return (p.Page - 1) * p.Size
}
