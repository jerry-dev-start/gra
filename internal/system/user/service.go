package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
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

func (s *Service) GetByID(id int64) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id int64, req *UpdateReq) error {
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

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) List(page *PageReq) ([]User, int64, error) {
	return s.repo.List(page.Offset(), page.Size)
}
