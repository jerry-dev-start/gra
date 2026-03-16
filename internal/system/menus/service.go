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

// buildTree 递归构建菜单树
func buildTree(list []Menus, parentID int64) []*MenuTree {
	var tree []*MenuTree
	for i := range list {
		if list[i].ParentID == parentID {
			node := &MenuTree{
				Menus:    list[i],
				Children: buildTree(list, list[i].ID),
			}
			if node.Children == nil {
				node.Children = []*MenuTree{}
			}
			tree = append(tree, node)
		}
	}
	return tree
}
