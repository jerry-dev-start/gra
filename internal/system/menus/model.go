package menus

import (
	"time"

	"gorm.io/gorm"

	"gra/pkg/id"
)

// 定义菜单类型和状态的别名
type MenuType string
type MenuStatus int8

const (
	MenuStatusDisabled MenuStatus = 0
	MenuStatusEnabled  MenuStatus = 1
)

const (
	MenuTypeDir    MenuType = "directory" // 目录
	MenuTypeMenu   MenuType = "menu"      // 菜单
	MenuTypeButton MenuType = "button"    // 按钮
)

type Menus struct {
	// 使用 ,string 标签将 int64 在 JSON 序列化时转为 string，解决前端精度丢失问题
	ID         int64      `json:"id,string" gorm:"primaryKey;comment:主键ID"`
	ParentID   int64      `json:"parentId,string" gorm:"index;default:0;comment:父级ID"`
	Name       string     `json:"name" gorm:"size:64;not null;comment:菜单名称"`
	Type       MenuType   `json:"type" gorm:"size:20;comment:菜单类型"`
	Path       string     `json:"path" gorm:"size:255;comment:路由路径"`
	Component  string     `json:"component" gorm:"size:255;comment:组件路径"`
	Icon       string     `json:"icon" gorm:"size:128;comment:图标"`
	Permission string     `json:"permission" gorm:"size:128;comment:权限标识"`
	Sort       int        `json:"sort" gorm:"default:0;comment:排序"`
	Visible    bool       `json:"visible" gorm:"default:true;comment:是否可见"`
	Status     MenuStatus `json:"status" gorm:"default:1;comment:状态"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func (Menus) TableName() string { return "menus" }

// BeforeCreate GORM 钩子：创建前自动生成雪花ID
func (m *Menus) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = id.Generate()
	}
	return nil
}

// DTO

type CreateReq struct {
	ParentID   int64    `json:"parentId,string"`
	Name       string   `json:"name" binding:"required,max=64"`
	Type       MenuType `json:"type" binding:"required"`
	Path       string   `json:"path"`
	Component  string   `json:"component"`
	Icon       string   `json:"icon"`
	Permission string   `json:"permission"`
	Sort       int      `json:"sort"`
	Visible    *bool    `json:"visible"`
	Status     *int8    `json:"status"`
}

type UpdateReq struct {
	ParentID   *int64   `json:"parentId,string"`
	Name       string   `json:"name" binding:"omitempty,max=64"`
	Type       MenuType `json:"type"`
	Path       string   `json:"path"`
	Component  string   `json:"component"`
	Icon       string   `json:"icon"`
	Permission string   `json:"permission"`
	Sort       *int     `json:"sort"`
	Visible    *bool    `json:"visible"`
	Status     *int8    `json:"status"`
}

// MenuTree 树形结构返回
type MenuTree struct {
	Menus
	Children []*MenuTree `json:"children"`
}
