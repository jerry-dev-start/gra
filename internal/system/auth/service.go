package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

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

	expire := time.Duration(config.Cfg.JWT.Expire) * time.Second
	expireAt := time.Now().Add(expire)

	claims := &middleware.Claims{
		UserID:   u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		return nil, errors.New("生成Token失败")
	}
	return &LoginResp{Token: token, ExpireAt: expireAt.Unix()}, nil
}
