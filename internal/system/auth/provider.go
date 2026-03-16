package auth

// Init auth 模块依赖注入入口
// userQ: 通过接口隔离依赖 user 模块
func Init(userQ UserQuerier) *Handler {
	svc := NewService(userQ)
	return NewHandler(svc)
}
