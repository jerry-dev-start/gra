package role

import (
	"gra/internal/public"
	"gra/internal/system/role/req"
)

type Service struct {
	repo *Repository
}

func (s *Service) Create(r *req.SaveRoleReq) error {
	//判断该角色名称或编码是否存在
	if err := s.repo.CheckRoleNameOrRoleCodeExist(r.RoleName, r.RoleCode); err != nil {
		return err
	}
	role := &SysRole{
		RoleName:    r.RoleName,
		RoleCode:    r.RoleCode,
		Description: r.Description,
		SortOrder:   r.SortOrder,
		Status:      r.Status,
		IsReadonly:  0,
	}
	return s.repo.Create(role)
}

func (s *Service) RoleList(page *public.PageReq) ([]SysRole, int64, error) {
	return s.repo.List(page.Offset(), page.Size)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
