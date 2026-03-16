package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *Repository) GetByID(id int64) (*User, error) {
	var u User
	err := r.db.First(&u, id).Error
	return &u, err
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	var u User
	err := r.db.Where("username = ?", username).First(&u).Error
	return &u, err
}

func (r *Repository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) Delete(id int64) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *Repository) List(offset, limit int) ([]User, int64, error) {
	var users []User
	var total int64

	db := r.db.Model(&User{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(offset).Limit(limit).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
