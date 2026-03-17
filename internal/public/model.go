package public

import (
	"time"

	"gorm.io/gorm"

	"gra/pkg/id"
)

type PageReq struct {
	Page int `form:"page,default=1" binding:"min=1"`
	Size int `form:"size,default=10" binding:"min=1,max=100"`
}

func (p *PageReq) Offset() int {
	return (p.Page - 1) * p.Size
}

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
