package menus

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *CreateReq) error {
	m := &Menus{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Type:       req.Type,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Permission: req.Permission,
		Sort:       req.Sort,
		Visible:    true,
		Status:     MenuStatusEnabled,
	}
	if req.Visible != nil {
		m.Visible = *req.Visible
	}
	if req.Status != nil {
		m.Status = MenuStatus(*req.Status)
	}
	return s.repo.Create(m)
}

func (s *Service) GetByID(id int64) (*Menus, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Update(id int64, req *UpdateReq) error {
	updates := make(map[string]interface{})
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Path != "" {
		updates["path"] = req.Path
	}
	if req.Component != "" {
		updates["component"] = req.Component
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Permission != "" {
		updates["permission"] = req.Permission
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Visible != nil {
		updates["visible"] = *req.Visible
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		return nil
	}
	return s.repo.Update(id, updates)
}

func (s *Service) Delete(id int64) error {
	has, err := s.repo.HasChildren(id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("存在子菜单，无法删除")
	}
	return s.repo.Delete(id)
}

// ListTree 获取全部菜单并构建树形结构
func (s *Service) ListTree() ([]*MenuTree, error) {
	list, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return buildTree(list, 0), nil
}

func (s *Service) UserMenuTree(userID int64) ([]*MenuTree, error) {
	menuList, err := s.repo.GetMenusByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(menuList) == 0 {
		return []*MenuTree{}, nil
	}
	return buildTree(menuList, 0), nil
}

// buildTree 使用 map 索引构建菜单树，时间复杂度 O(n)
func buildTree(list []Menus, rootParentID int64) []*MenuTree {
	nodeMap := make(map[int64]*MenuTree, len(list))
	for i := range list {
		nodeMap[list[i].ID] = &MenuTree{
			Menus:    list[i],
			Children: []*MenuTree{},
		}
	}

	var roots []*MenuTree
	for i := range list {
		node := nodeMap[list[i].ID]
		if list[i].ParentID == rootParentID {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[list[i].ParentID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}
	if roots == nil {
		roots = []*MenuTree{}
	}
	return roots
}
