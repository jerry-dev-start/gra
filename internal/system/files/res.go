package files

type FileUploadRes struct {
	FileUrl        string   `json:"fileUrl"`
	FileName       string   `json:"fileName"`
	UploadedChunks []string `json:"uploadedChunks"`
	Uploaded       bool     `json:"uploaded"`
}
