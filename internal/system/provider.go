package system

import (
	"gorm.io/gorm"

	"gra/internal/system/user"
)

// Handlers 系统域所有 Handler 集合
type Handlers struct {
	User *user.Handler
	// 后续扩展：
	// Role *role.Handler
	// Menu *menu.Handler
	// Dept *dept.Handler
}

// Services 系统域对外暴露的 Service 集合
// 供业务域跨域调用（通过接口隔离）
type Services struct {
	User *user.Service
	// 后续扩展：
	// Role *role.Service
}

// Init 系统域统一初始化入口
func Init(db *gorm.DB) (*Handlers, *Services) {
	userHandler, userSvc := user.Init(db)

	handlers := &Handlers{
		User: userHandler,
	}
	services := &Services{
		User: userSvc,
	}
	return handlers, services
}
