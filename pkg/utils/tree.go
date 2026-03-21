package utils

type TreeItem interface {
	GetID() int64
	GetParentID() int64
}

type TreeNode[T any] struct {
	Data     T              `json:"data"`     // 原始数据
	Children []*TreeNode[T] `json:"children"` // 子节点
}

// BuildTree 通用构建树结构工具方法
// T 必须实现 TreeItem 接口
func BuildTree[T TreeItem](list []T, rootParentID int64) []*TreeNode[T] {
	// 预分配内存，提高性能
	nodeMap := make(map[int64]*TreeNode[T], len(list))
	roots := make([]*TreeNode[T], 0)

	// 第一遍遍历：初始化所有节点并存入 map
	for _, item := range list {
		nodeMap[item.GetID()] = &TreeNode[T]{
			Data:     item,
			Children: make([]*TreeNode[T], 0),
		}
	}

	// 第二遍遍历：构建父子关系
	for _, item := range list {
		id := item.GetID()
		pid := item.GetParentID()
		node := nodeMap[id]

		if pid == rootParentID {
			// 如果是根节点
			roots = append(roots, node)
		} else {
			// 如果有父节点，则加入父节点的 Children 中
			if parent, ok := nodeMap[pid]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}
