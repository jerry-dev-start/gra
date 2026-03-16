package public

import (
	"time"

	"gorm.io/gorm"

	"gra/pkg/id"
)

// BaseModel 公共字段
type BaseModel struct {
	ID        int64          `json:"id,string" gorm:"primaryKey;comment:主键ID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate GORM 钩子：创建前自动生成雪花ID
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		b.ID = id.Generate()
	}
	return nil
}
