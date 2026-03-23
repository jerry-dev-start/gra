package files

import "gra/internal/public"

// SysFile 文件元数据表
type SysFile struct {
	public.BaseModel
	FileName  string `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	FileMd5   string `gorm:"column:file_md5;type:char(32);uniqueIndex:idx_md5;not null" json:"file_md5"`
	FilePath  string `gorm:"column:file_path;type:varchar(512);not null" json:"file_path"`
	FileSize  uint64 `gorm:"column:file_size;type:bigint unsigned;not null;default:0" json:"file_size"`
	Extension string `gorm:"column:extension;type:varchar(20)" json:"extension"`
}

// TableName 指定表名
func (SysFile) TableName() string {
	return "sys_file"
}
