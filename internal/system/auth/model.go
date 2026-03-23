package auth

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" label:"用户名"`
	Password string `json:"password" label:"密码"`
}

// LoginResp 登录响应
type LoginResp struct {
	Token        string `json:"token"`
	ExpireAt     int64  `json:"expire_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshExpAt int64  `json:"refresh_exp_at"`
}

// RefreshReq 刷新Token请求
type RefreshReq struct {
	RefreshToken string `json:"refresh_token" label:"刷新Token"`
}

// RefreshResp 刷新Token响应
type RefreshResp struct {
	Token        string `json:"token"`
	ExpireAt     int64  `json:"expire_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshExpAt int64  `json:"refresh_exp_at"`
}
