package files

import (
	"gra/global"
	"gra/pkg/response"
	"gra/pkg/validate"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
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
		FileUrl:  filepath.Join(uploadDir, newFileName),
		FileName: newFileName,
	})
}

// ChunkCheck 分片上传前置检查
// 如果检查到已经上传文件就直接秒传返回
// 如果没有上传文件 创建后续分片需要文件夹
// 如果文件已经传分片了，就返回所传的分片序号
func (h *Handler) ChunkCheck(c *gin.Context) {
	var req CheckFileReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	err := validate.Check(req, validate.Rules{
		"Hash": {validate.Required("文件MD5必必须传入")},
	})
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	result, err := h.svc.ChunkCheck(req.Hash)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OK(c, result)
}

// ChunkUpload 上传分片
func (h *Handler) ChunkUpload(c *gin.Context) {
	var req ChunkUploadReq
	if err := c.ShouldBind(&req); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	if err := validate.Check(req, validate.Rules{
		"Hash":        {validate.Required("Hash不存在")},
		"ChunkIndex":  {validate.Required("分片号不存在")},
		"TotalChunks": {validate.Required("总分片不存在")},
		"File":        {validate.Required("文件不存在")},
	}); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	// 拼接保存分片的目录
	chunkDir := filepath.Join(global.Conf.FileUploadConfig.Dir, "chunk", req.Hash)
	if err := os.MkdirAll(chunkDir, os.ModePerm); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	//写入文件的全地址
	partFilePath := filepath.Join(chunkDir, strconv.FormatInt(req.ChunkIndex, 10)+".part")
	finalFilePath := filepath.Join(chunkDir, strconv.FormatInt(req.ChunkIndex, 10))
	if err := c.SaveUploadedFile(req.File, partFilePath); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	if err := os.Rename(partFilePath, finalFilePath); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	response.OKMsg(c, "保存分片完成")
}

// RegisterRoutes 注册需认证的路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/upload", h.Upload)
		fileGroup.GET("/chunk/check", h.ChunkCheck)
		fileGroup.POST("/chunk/upload", h.ChunkUpload)
	}
}
