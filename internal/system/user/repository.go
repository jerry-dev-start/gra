package user

import (
	"gra/internal/public"
	"gra/internal/system/model"
	"os/user"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(u *User, roleIds []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		if len(roleIds) == 0 {
			return nil
		}
		roleUser := make([]model.SysRoleUser, 0, len(roleIds))
		for _, roleId := range roleIds {
			roleUser = append(roleUser, model.SysRoleUser{
				RoleId: roleId,
				UserId: u.ID,
			})
		}
		return tx.Create(&roleUser).Error
	})
}

func (r *Repository) GetByID(id int64) (*UserDetail, error) {
	var u User
	err := r.db.Preload("Roles").First(&u, id).Error
	if err != nil {
		return nil, err
	}
	roleIds := make([]public.StringInt64, 0, len(u.Roles))
	for _, role := range u.Roles {
		roleIds = append(roleIds, public.StringInt64(role.ID))
	}
	userRes := UserDetail{
		BaseModel: public.BaseModel{
			ID: u.ID,
		},
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Email:    u.Email,
		Phone:    u.Phone,
		Status:   u.Status,
		RoleIds:  roleIds,
		DeptId:   u.DeptId,
	}
	return &userRes, nil
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	var u User
	err := r.db.Where("username = ?", username).First(&u).Error
	return &u, err
}

func (r *Repository) Update(id int64, updates map[string]interface{}, roleIds []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", id).Delete(&model.SysRoleUser{}).Error; err != nil {
			return err
		}

		if len(roleIds) == 0 {
			return nil
		}

		roleUser := make([]model.SysRoleUser, 0, len(roleIds))
		for _, roleId := range roleIds {
			roleUser = append(roleUser, model.SysRoleUser{
				RoleId: roleId,
				UserId: id,
			})
		}
		return tx.Create(&roleUser).Error
	})
}

func (r *Repository) Delete(id int64) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *Repository) List(offset, limit int, page *DeptQueryReq) ([]User, int64, error) {
	var users []User
	var total int64

	db := r.db.Model(&User{})

	if page.DeptId != 0 {
		db = db.Where("dept_id = ? OR dept_id IN (SELECT id FROM sys_dept WHERE FIND_IN_SET(?, ancestors))",
			page.DeptId, page.DeptId)
	}

	if page.UserName != "" {
		db = db.Where("username LIKE ?", "%"+page.UserName+"%")
	}

	if page.Phone != "" {
		db = db.Where("phone = ?", page.Phone)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *Repository) CheckDeptHasUsers(id int64) (bool, error) {
	var exists bool
	err := r.db.Model(&user.User{}).
		Select("count(*) > 0").
		Where("dept_id = ?", id).
		Find(&exists).
		Error
	if err != nil {
		return false, err
	}
	return exists, err
}
