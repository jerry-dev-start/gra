package files

import "mime/multipart"

type CheckFileReq struct {
	Hash        string `json:"hash" form:"hash"`
	FileName    string `json:"fileName" form:"fileName"`
	TotalChunks int    `json:"totalChunks" form:"totalChunks"`
}

type ChunkUploadReq struct {
	Hash        string                `json:"hash" form:"hash"`
	ChunkIndex  int64                 `json:"chunkIndex" form:"chunkIndex"`
	TotalChunks int64                 `json:"totalChunks" form:"totalChunks"`
	File        *multipart.FileHeader `form:"file"`
}

type ChunkMergeReq struct {
	Hash        string `json:"hash" form:"hash"`
	FileName    string `json:"fileName" form:"fileName"`
	TotalChunks int64  `json:"totalChunks" form:"totalChunks"`
}
