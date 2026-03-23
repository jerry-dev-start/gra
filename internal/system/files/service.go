package files

import (
	"gra/global"
	"os"
	"path/filepath"
)

type DeptUserQuerier interface {
	CheckDeptHasUsers(id int64) (bool, error)
}
type Service struct {
	repo *Repository
}

func (s *Service) ChunkCheck(hash string) (*FileUploadRes, error) {
	fileInfo, err := s.repo.CheckFileInfoByHash(hash)
	if err != nil {
		return nil, err
	}
	// --- 场景 A：文件已存在且已完成上传 (秒传) ---
	if fileInfo != nil {
		return &FileUploadRes{
			FileUrl:        fileInfo.FilePath,
			FileName:       fileInfo.FileName,
			UploadedChunks: nil,
			Uploaded:       true,
		}, nil
	}
	// --- 场景 B：文件不存在或上传未完成 ---
	//1. 拼接文件保存路径
	filePath := filepath.Join(global.Conf.FileUploadConfig.Dir, "chunk", hash)
	if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
		return nil, err
	}
	//2. 读取目录下的文件
	files, err := readDirFiles(filePath)
	if err != nil {
		return nil, err
	}
	return &FileUploadRes{
		FileUrl:        "",
		FileName:       "",
		UploadedChunks: files,
		Uploaded:       false,
	}, nil
}

// readDirFiles 读取指定目录下的所有文件，并排除子目录
// 参数 dirname: 目标目录的路径
// 返回值: 文件名切片 ([]string) 和 可能发生的错误
func readDirFiles(dirname string) ([]string, error) {
	dir, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	fileList := make([]string, 0, 0)
	for _, f := range dir {
		if !f.IsDir() {
			fileList = append(fileList, f.Name())
		}
	}
	return fileList, nil
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
