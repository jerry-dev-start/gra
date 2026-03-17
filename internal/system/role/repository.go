package role

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CheckRoleNameOrRoleCodeExist(name string, code string) error {
	var count int64
	if err := r.db.Model(&SysRole{}).Where("role_name = ?", name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色名称已存在")
	}

	if err := r.db.Model(&SysRole{}).Where("role_code = ?", code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("角色编码已存在")
	}
	return nil
}

func (r *Repository) Create(role *SysRole) error {
	return r.db.Create(role).Error
}

func (r *Repository) List(offset, limit int) ([]SysRole, int64, error) {
	var roles []SysRole
	var total int64

	db := r.db.Model(&SysRole{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}
