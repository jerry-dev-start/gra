package role

import "gorm.io/gorm"

// Init 菜单模块依赖注入入口
func Init(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)
	return h
}
