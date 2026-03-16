package user

import "gorm.io/gorm"

// Init 用户模块依赖注入入口
// 内部接线完全自包含，外部只需传入 db
func Init(db *gorm.DB) (*Handler, *Service) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)
	return h, svc
}
