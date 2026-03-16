package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gra/internal/middleware"
	"gra/pkg/config"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(req *LoginReq) (*LoginResp, error) {
	u, err := s.repo.GetByUsername(req.Username)
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

func (s *Service) Create(req *CreateReq) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}
	u := &User{
		Username: req.Username,
		Password: string(hash),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
		Status:   1,
	}
	return s.repo.Create(u)
}

func (s *Service) GetByID(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id uint, req *UpdateReq) error {
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.RoleID != 0 {
		updates["role_id"] = req.RoleID
	}
	if len(updates) == 0 {
		return nil
	}
	return s.repo.Update(id, updates)
}

func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *Service) List(page *PageReq) ([]User, int64, error) {
	return s.repo.List(page.Offset(), page.Size)
}
