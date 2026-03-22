package dept

import "gorm.io/gorm"

// Init 菜单模块依赖注入入口
func Init(db *gorm.DB, querier DeptUserQuerier) *Handler {
	repo := NewRepository(db)
	svc := NewService(repo, querier)
	h := NewHandler(svc)
	return h
}
