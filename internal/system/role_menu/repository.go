package role_menu

import (
	"fmt"
	"gra/internal/public"
	"strconv"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetRoleBindMenuId(id string) ([]string, error) {
	var roleMenu []public.SysRoleMenu
	if err := r.db.Model(public.SysRoleMenu{}).Where("role_id = ?", id).Find(&roleMenu).Error; err != nil {
		return nil, err
	}
	list := make([]string, 0, 10)
	for _, role := range roleMenu {
		list = append(list, strconv.FormatInt(role.MenuId, 10))
	}
	return list, nil
}

func (r *Repository) SaveRoleMenu(roleId string, menuId []string) error {

	roleIdInt, err := strconv.ParseInt(roleId, 10, 64)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleIdInt).Delete(&public.SysRoleMenu{}).Error; err != nil {
			return err
		}
		if len(menuId) == 0 {
			return nil
		}
		// 保存所有的值
		roleMenus := make([]public.SysRoleMenu, 0, len(menuId))
		for _, menuId := range menuId {
			menuIdInt, err := strconv.ParseInt(menuId, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid menuId %s: %w", menuId, err)
			}
			roleMenus = append(roleMenus, public.SysRoleMenu{
				RoleId: roleIdInt,
				MenuId: menuIdInt,
			})
		}
		err := tx.Create(&roleMenus).Error
		return err
	})
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
