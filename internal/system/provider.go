package system

import (
	"gra/internal/system/role"

	"gorm.io/gorm"

	"gra/internal/system/auth"
	"gra/internal/system/menus"
	"gra/internal/system/user"
)

// Handlers 系统域所有 Handler 集合
type Handlers struct {
	Auth *auth.Handler
	User *user.Handler
	Menu *menus.Handler
	Role *role.Handler
}

// Services 系统域对外暴露的 Service 集合
// 供业务域跨域调用（通过接口隔离）
type Services struct {
	User *user.Service
}

// userAdapter 适配 auth.UserQuerier 接口
// 放在 provider 层接线，不污染 user 模块
type userAdapter struct {
	repo *user.Repository
}

func (a *userAdapter) GetByUsername(username string) (auth.UserInfo, error) {
	u, err := a.repo.GetByUsername(username)
	if err != nil {
		return auth.UserInfo{}, err
	}
	return auth.UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Status:   u.Status,
	}, nil
}

// Init 系统域统一初始化入口
func Init(db *gorm.DB) (*Handlers, *Services) {
	// user 模块
	userHandler, userSvc, userRepo := user.Init(db)

	// auth 模块（通过适配器依赖 user）
	authHandler := auth.Init(&userAdapter{repo: userRepo})

	// menus 模块
	menuHandler := menus.Init(db)
	roleHandler := role.Init(db)

	handlers := &Handlers{
		Auth: authHandler,
		User: userHandler,
		Menu: menuHandler,
		Role: roleHandler,
	}
	services := &Services{
		User: userSvc,
	}
	return handlers, services
}
