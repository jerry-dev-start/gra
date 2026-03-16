package menus

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(m *Menus) error {
	return r.db.Create(m).Error
}

func (r *Repository) GetByID(id int64) (*Menus, error) {
	var m Menus
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *Repository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&Menus{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&Menus{}).Error
}

// DeleteByParentID 删除某父级下所有子菜单
func (r *Repository) DeleteByParentID(parentID int64) error {
	return r.db.Where("parent_id = ?", parentID).Delete(&Menus{}).Error
}

// ListAll 查询全部菜单，按 sort 升序
func (r *Repository) ListAll() ([]Menus, error) {
	var list []Menus
	err := r.db.Order("sort ASC, id ASC").Find(&list).Error
	return list, err
}

// HasChildren 检查是否有子菜单
func (r *Repository) HasChildren(parentID int64) (bool, error) {
	var count int64
	err := r.db.Model(&Menus{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count > 0, err
}
