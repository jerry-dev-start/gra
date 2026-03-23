package files

import (
	"gra/global"
	"gra/pkg/response"
	"gra/pkg/validate"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

// ChunkMerge 合并分片的请求
func (h *Handler) ChunkMerge(c *gin.Context) {
	var mergeReq ChunkMergeReq
	if err := c.ShouldBindJSON(&mergeReq); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	if err := validate.Check(mergeReq, validate.Rules{
		"Hash":        {validate.Required("Hash不能为空")},
		"TotalChunks": {validate.Required("TotalChunks不能为空")},
	}); err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	//拼接Chunk全路径
	chunkPath := filepath.Join(global.Conf.FileUploadConfig.Dir, "chunk", mergeReq.Hash)
	exists := isDirExists(chunkPath)
	if !exists {
		response.FailMsg(c, "分片目录不存在")
		return
	}
	// 1. 读取目录下所有文件
	entries, err := os.ReadDir(chunkPath)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}

	var chunks []string
	for _, entry := range entries {
		// 排除文件夹和正在写入的临时文件
		if !entry.IsDir() && !strings.HasSuffix(entry.Name(), ".part") {
			chunks = append(chunks, entry.Name())
		}
	}

	// 2. 关键步骤：按数字顺序排序 (避免 1, 10, 2 的情况)
	sort.Slice(chunks, func(i, j int) bool {
		a, _ := strconv.Atoi(chunks[i])
		b, _ := strconv.Atoi(chunks[j])
		return a < b
	})

	// 3. 创建目标文件
	extension := filepath.Ext(mergeReq.FileName) // 获取 .jpg, .png 等
	newFileName := uuid.New().String() + extension
	dateDir := time.Now().Format("2006/01/02")
	uploadDir := filepath.Join("uploads", dateDir)
	uploadAllPath := filepath.Join(global.Conf.FileUploadConfig.Dir, "uploads", dateDir, newFileName)
	targetFile, err := os.Create(uploadAllPath)
	if err != nil {
		response.FailMsg(c, err.Error())
		return
	}
	defer targetFile.Close()

	// 4. 逐个读取分片并写入目标文件
	for _, chunkName := range chunks {
		chunkPath := filepath.Join(chunkPath, chunkName)

		// 开启分片文件
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			response.FailMsg(c, err.Error())
			return
		}

		// 使用 io.Copy 实现流式拷贝，节省内存
		if _, err := io.Copy(targetFile, chunkFile); err != nil {
			chunkFile.Close()
			return
		}

		chunkFile.Close() // 读完立即关闭
	}

	// 5. 合并成功后，清理分片目录
	os.RemoveAll(chunkPath)

	response.OK(c, &FileUploadRes{
		FileUrl:        filepath.Join(uploadDir, newFileName),
		FileName:       mergeReq.FileName,
		UploadedChunks: nil,
		Uploaded:       false,
	})
}

// IsDirExists 判断路径是否存在且是一个文件夹
func isDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// 如果错误是“不存在”，直接返回 false
		if os.IsNotExist(err) {
			return false
		}
		// 其他错误（如权限不足）也视为不可用，返回 false
		return false
	}
	// 核心步骤：判断查找到的路径是否为目录
	return info.IsDir()
}

// RegisterRoutes 注册需认证的路由
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/upload", h.Upload)
		fileGroup.GET("/chunk/check", h.ChunkCheck)
		fileGroup.POST("/chunk/upload", h.ChunkUpload)
		fileGroup.POST("/chunk/merge", h.ChunkMerge)
	}
}
