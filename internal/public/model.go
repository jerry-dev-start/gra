package public

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	ID        int64          `json:"id,string" gorm:"column:id;primaryKey;comment:主键ID"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate GORM 钩子：创建前自动生成雪花ID
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == 0 {
		b.ID = id.Generate()
	}
	return nil
}

// StringInt64 JSON 序列化为 string，反序列化兼容 string 和 number
// 用于雪花ID等超出 JS Number.MAX_SAFE_INTEGER 的场景
type StringInt64 int64

func (s StringInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(s), 10))
}

func (s *StringInt64) UnmarshalJSON(data []byte) error {
	// 兼容 "123" 和 123 两种格式
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid StringInt64 value: %s", str)
		}
		*s = StringInt64(v)
		return nil
	}
	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return fmt.Errorf("invalid StringInt64 value: %s", string(data))
	}
	*s = StringInt64(num)
	return nil
}

// Int64 返回底层 int64 值
func (s StringInt64) Int64() int64 {
	return int64(s)
}

// ToStringInt64Slice 将 []StringInt64 转为 []int64
func ToStringInt64Slice(ids []StringInt64) []int64 {
	result := make([]int64, len(ids))
	for i, v := range ids {
		result[i] = v.Int64()
	}
	return result
}
