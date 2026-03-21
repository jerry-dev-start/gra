package dept

import (
	"gra/internal/public"
)

// SysDept 部门信息表
type SysDept struct {
	public.BaseModel
	ID        int64  `gorm:"column:id;primaryKey;comment:部门ID" json:"id,string"`
	ParentID  int64  `gorm:"column:parent_id;default:0;index:idx_parent_id;comment:上级部门ID (0表示顶级部门)" json:"parentId,string"`
	DeptName  string `gorm:"column:dept_name;type:varchar(50);not null;comment:部门名称" json:"deptName"`
	Leader    string `gorm:"column:leader;type:varchar(20);comment:负责人" json:"leader"`
	Phone     string `gorm:"column:phone;type:varchar(11);comment:联系电话" json:"phone"`
	Email     string `gorm:"column:email;type:varchar(50);comment:邮箱" json:"email"`
	SortOrder int    `gorm:"column:sort_order;default:0;comment:显示排序" json:"sortOrder"`
	Status    string `gorm:"column:status;type:char(1);default:0;comment:部门状态" json:"status"`
	Ancestors string `gorm:"column:ancestors;type:varchar(255);comment:祖级部门Id" json:"ancestors"`
}

// TableName 指定表名
func (SysDept) TableName() string {
	return "sys_dept"
}

// 1. 让你的模型实现接口
func (m SysDept) GetID() int64       { return m.ID }
func (m SysDept) GetParentID() int64 { return m.ParentID }

// DeptParams 对应前端的请求参数
type DeptReq struct {
	ID       int64  `json:"id,string"`        // 部门ID (可选，用于修改)
	ParentID int64  `json:"parentId,string" ` // 上级部门ID (必填)
	Name     string `json:"name"`             // 部门名称 (必填)
	Leader   string `json:"leader"`           // 负责人 (可选)
	Phone    string `json:"phone"`            // 联系电话 (可选)
	Email    string `json:"email"`            // 邮箱 (可选)
	Sort     int    `json:"sort"`             // 显示排序
	Status   string `json:"status"`           // 部门状态 (0正常 1停用)
}

type DeptQueryReq struct {
}

// MenuTree 树形结构返回
type DeptTree struct {
	SysDept
	Children []*DeptTree `json:"children"`
}
