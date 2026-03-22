package dept

import (
	"errors"
	"fmt"
	"gra/internal/public"
	"sort"
)

type DeptUserQuerier interface {
	CheckDeptHasUsers(id int64) (bool, error)
}
type Service struct {
	repo  *Repository
	userQ DeptUserQuerier
}

func NewService(repo *Repository, userQ DeptUserQuerier) *Service {
	return &Service{repo: repo, userQ: userQ}
}
func (s *Service) Create(req DeptReq) error {
	isExist, err := s.repo.CheckDeptNameExist(req.Name, req.ID)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("部门名称 [%s] 已存在", req.Name)
	}
	deptDb := SysDept{
		BaseModel: public.BaseModel{},
		ParentID:  req.ParentID,
		DeptName:  req.Name,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		SortOrder: req.Sort,
		Status:    req.Status,
	}
	err = s.repo.Create(deptDb)
	return err
}

func (s *Service) DeptTree(req DeptQueryReq) ([]*DeptTree, error) {
	//查询部门的数据
	result, err := s.repo.SelectDeptList(req)
	if err != nil {
		return nil, err
	}
	//构建树结构
	tree := buildTree(result, 1)
	return tree, err
}

func (s *Service) GetDeptInfo(id int64) (*SysDept, error) {
	return s.repo.GetDeptInfoById(id)
}

func (s *Service) Update(req *DeptReq) error {
	isExist, err := s.repo.CheckDeptNameExist(req.Name, req.ID)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("部门名称 [%s] 已存在", req.Name)
	}
	deptDb := SysDept{
		ID:        req.ID,
		ParentID:  req.ParentID,
		DeptName:  req.Name,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		SortOrder: req.Sort,
		Status:    req.Status,
	}
	return s.repo.Update(&deptDb)
}

func (s *Service) DeleteDept(id int64) error {
	//删除前检查是否存在下级
	exist, err := s.repo.HasChildren(id)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("该部门下存在子部门不能删除！")
	}

	//判断该部门下是否存在用户
	exist, err = s.userQ.CheckDeptHasUsers(id)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("该部门下存在用户不能删除！")
	}
	return s.repo.DeleteDeptById(id)
}

func buildTree(list []SysDept, rootParentID int64) []*DeptTree {
	nodeMap := make(map[int64]*DeptTree, len(list))

	// 第一步：初始化所有节点
	for i := range list {
		nodeMap[list[i].ID] = &DeptTree{
			SysDept:  list[i],
			Children: []*DeptTree{},
		}
	}

	var roots []*DeptTree
	// 第二步：构建父子关系
	for i := range list {
		node := nodeMap[list[i].ID]
		if list[i].ParentID == rootParentID {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[list[i].ParentID]; ok {
			parent.Children = append(parent.Children, node)
		}
	}

	// 第三步：对根节点和所有子节点进行递归排序
	// 先排根节点
	sortDeptTrees(roots)

	// 对每个节点的 Children 进行排序
	for _, node := range nodeMap {
		if len(node.Children) > 0 {
			sortDeptTrees(node.Children)
		}
	}

	if roots == nil {
		roots = []*DeptTree{}
	}
	return roots
}

// 辅助排序函数：根据 sortOrder 升序排列
func sortDeptTrees(nodes []*DeptTree) {
	sort.Slice(nodes, func(i, j int) bool {
		// 如果 sortOrder 相同，可以增加第二个排序条件，比如 ID
		if nodes[i].SortOrder == nodes[j].SortOrder {
			return nodes[i].ID < nodes[j].ID
		}
		return nodes[i].SortOrder < nodes[j].SortOrder
	})
}
