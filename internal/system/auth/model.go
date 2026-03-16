package auth

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" label:"用户名"`
	Password string `json:"password" label:"密码"`
}

// LoginResp 登录响应
type LoginResp struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
}
