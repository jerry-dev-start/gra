package files

import "gorm.io/gorm"

// Init 菜单模块依赖注入入口
func Init(db *gorm.DB) *Handler {
	h := NewHandler()
	return h
}
