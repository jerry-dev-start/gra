package role_menu

type Service struct {
	repo *Repository
}

func (s *Service) GetRoleBindMenuId(id string) ([]string, error) {
	return s.repo.GetRoleBindMenuId(id)
}

func (s *Service) SaveRoleMenu(id string, req SaveRoleMenuReq) error {
	return s.repo.SaveRoleMenu(id, req.MenuIds)
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
