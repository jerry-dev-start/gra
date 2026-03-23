package files

import (
	"gra/global"
	"gra/pkg/response"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// Upload 上传小文件
func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}

	// --- 1. 生成基于日期的保存路径 ---
	// 格式：uploads/2023/10/27/
	dateDir := time.Now().Format("2006/01/02")
	uploadDir := filepath.Join("uploads", dateDir)
	uploadAllPath := filepath.Join(global.Conf.FileUploadConfig.Dir, "uploads", dateDir)
	// 确保文件夹存在 (相当于 mkdir -p)
	if err := os.MkdirAll(uploadAllPath, os.ModePerm); err != nil {
		response.FailMsg(c, "无法创建存储目录")
		return
	}
	// --- 2. 业内最佳重命名方案 ---
	// 使用 UUID/雪花ID + 原始后缀名
	extension := filepath.Ext(file.Filename) // 获取 .jpg, .png 等
	newFileName := uuid.New().String() + extension
	// 最终保存的完整相对路径
	dst := filepath.Join(uploadAllPath, newFileName)

	// --- 3. 执行保存 ---
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.FailMsg(c, "文件保存失败")
		return
	}
	response.OK(c, &FileUploadRes{
		Url:      filepath.Join(uploadDir, newFileName),
		FileName: newFileName,
	})
}

// RegisterRoutes 注册需认证的路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/upload", h.Upload)
	}
}
