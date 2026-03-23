package files

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) CheckFileInfoByHash(hash string) (*SysFile, error) {
	var file SysFile
	err := r.db.Where("file_md5 = ?", hash).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
