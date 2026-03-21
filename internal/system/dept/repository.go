package dept

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// CheckDeptNameExist 检查部门名称是否存在
func (r *Repository) CheckDeptNameExist(name string, excludeID int64) (bool, error) {
	var count int64

	db := r.db.Model(&SysDept{}).Where("dept_name = ?", name)

	// 如果 excludeID > 0，说明是编辑操作，需要排除掉当前正在修改的这行数据
	if excludeID > 0 {
		db = db.Where("id <> ?", excludeID)
	}

	err := db.Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) Create(dept SysDept) error {
	if dept.ParentID != 1 {
		var parent SysDept
		err := r.db.First(&parent, dept.ParentID).Error
		if err != nil {
			return fmt.Errorf("父部门不存在：%w", err)
		}
		dept.Ancestors = parent.Ancestors + "," + strconv.FormatInt(parent.ID, 10)
	} else {
		dept.Ancestors = "1"
	}
	return r.db.Create(&dept).Error
}

func (r *Repository) SelectDeptList(req DeptQueryReq) ([]SysDept, error) {
	var deptList []SysDept
	err := r.db.Find(&deptList).Order("sort_order asc").Error
	return deptList, err
}

func (r *Repository) GetDeptInfoById(id int64) (*SysDept, error) {
	var dept SysDept
	err := r.db.First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dept, nil
}

func (r *Repository) Update(s *SysDept) error {
	return r.db.Where("id = ?", s.ID).Updates(&s).Error
}

func (r *Repository) HasChildren(id int64) (bool, error) {
	//查询出该部门的祖级
	var currentDept SysDept
	if err := r.db.Where("id = ?", id).First(&currentDept).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("没有找到当前部门的信息")
		}
		return false, err
	}
	//拼级祖级Id
	ancestors := currentDept.Ancestors + "," + strconv.FormatInt(currentDept.ID, 10)
	//查询数据库中左模糊有没有类似的数据
	var count int64
	err := r.db.Model(&SysDept{}).Where("ancestors like ?", "%"+ancestors).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) DeleteDeptById(id int64) error {
	return r.db.Where("id = ?", id).Delete(&SysDept{}).Error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
