package business

import (
	"gorm.io/gorm"

	"gra/internal/system"
)

// Handlers 业务域所有 Handler 集合
type Handlers struct {
	// 后续扩展：
	// Order   *order.Handler
	// Product *product.Handler
}

// Init 业务域统一初始化入口
// sysSvc 用于跨域依赖注入（通过接口隔离）
func Init(db *gorm.DB, sysSvc *system.Services) *Handlers {
	// 示例：跨域注入
	// orderHandler := order.Init(db, sysSvc.User)

	_ = db
	_ = sysSvc

	return &Handlers{}
}
