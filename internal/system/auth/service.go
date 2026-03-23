package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gra/global"
	"gra/internal/middleware"
	"gra/pkg/config"
)

// UserQuerier auth 模块对 user 的依赖接口（接口隔离）
type UserQuerier interface {
	GetByUsername(username string) (UserInfo, error)
}

// UserInfo auth 模块需要的用户信息（不依赖 user.User 具体类型）
type UserInfo struct {
	ID       int64
	Username string
	Password string
	Status   int8
}

type Service struct {
	userQ UserQuerier
}

func NewService(userQ UserQuerier) *Service {
	return &Service{userQ: userQ}
}

// refreshTokenKey 生成 Redis key
func refreshTokenKey(userID int64) string {
	return fmt.Sprintf("auth:refresh_token:%d", userID)
}

func (s *Service) Login(req *LoginReq) (*LoginResp, error) {
	u, err := s.userQ.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if u.Status != 1 {
		return nil, errors.New("账号已禁用")
	}

	// 签发 access token
	accessToken, accessExp, err := generateToken(u.ID, u.Username, "access", config.Cfg.JWT.Expire)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 签发 refresh token
	refreshExpire := config.Cfg.JWT.RefreshExpire
	refreshToken, refreshExp, err := generateToken(u.ID, u.Username, "refresh", refreshExpire)
	if err != nil {
		return nil, errors.New("生成RefreshToken失败")
	}

	// 存入 Redis，TTL 与 refresh token 有效期一致
	err = global.Rdb.Set(
		context.Background(),
		refreshTokenKey(u.ID),
		refreshToken,
		time.Duration(refreshExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, errors.New("存储RefreshToken失败")
	}

	return &LoginResp{
		Token:        accessToken,
		ExpireAt:     accessExp,
		RefreshToken: refreshToken,
		RefreshExpAt: refreshExp,
	}, nil
}

// RefreshToken 通过 refresh token 换取新的 access token
func (s *Service) RefreshToken(req *RefreshReq) (*RefreshResp, error) {
	// 解析 refresh token
	claims := &middleware.Claims{}
	t, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})
	if err != nil || !t.Valid {
		return nil, errors.New("RefreshToken无效或已过期")
	}

	// 必须是 refresh 类型
	if claims.TokenType != "refresh" {
		return nil, errors.New("Token类型错误")
	}

	// 校验 Redis 中是否存在（防止已吊销的 token 继续使用）
	key := refreshTokenKey(claims.UserID)
	stored, err := global.Rdb.Get(context.Background(), key).Result()
	if err != nil {
		return nil, errors.New("RefreshToken已失效，请重新登录")
	}
	if stored != req.RefreshToken {
		// token 不匹配，可能是旧 token 被重放，直接吊销
		global.Rdb.Del(context.Background(), key)
		return nil, errors.New("RefreshToken已被使用，请重新登录")
	}

	// 校验用户状态
	u, err := s.userQ.GetByUsername(claims.Username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if u.Status != 1 {
		global.Rdb.Del(context.Background(), key)
		return nil, errors.New("账号已禁用")
	}

	// 签发新的 access token
	accessToken, accessExp, err := generateToken(u.ID, u.Username, "access", config.Cfg.JWT.Expire)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// Token 轮换：签发新的 refresh token，旧的立即作废
	refreshExpire := config.Cfg.JWT.RefreshExpire
	newRefreshToken, newRefreshExp, err := generateToken(u.ID, u.Username, "refresh", refreshExpire)
	if err != nil {
		return nil, errors.New("生成RefreshToken失败")
	}
	err = global.Rdb.Set(
		context.Background(),
		key,
		newRefreshToken,
		time.Duration(refreshExpire)*time.Second,
	).Err()
	if err != nil {
		return nil, errors.New("存储RefreshToken失败")
	}

	return &RefreshResp{
		Token:        accessToken,
		ExpireAt:     accessExp,
		RefreshToken: newRefreshToken,
		RefreshExpAt: newRefreshExp,
	}, nil
}

// Logout 主动吊销 refresh token
func (s *Service) Logout(userID int64) error {
	return global.Rdb.Del(context.Background(), refreshTokenKey(userID)).Err()
}

// generateToken 生成指定类型的 JWT token
func generateToken(userID int64, username, tokenType string, expireSec int) (string, int64, error) {
	expire := time.Duration(expireSec) * time.Second
	expireAt := time.Now().Add(expire)

	claims := &middleware.Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		return "", 0, err
	}
	return token, expireAt.Unix(), nil
}
