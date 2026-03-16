package user

import "gorm.io/gorm"

// Init 用户模块依赖注入入口
func Init(db *gorm.DB) (*Handler, *Service, *Repository) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)
	return h, svc, repo
}
