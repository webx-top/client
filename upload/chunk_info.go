package upload

import (
	"net/url"

	"github.com/webx-top/echo/param"
)

type ChunkInfor interface {
	// 当前分片
	GetChunkIndex() uint64       // 当前分片索引编码
	GetChunkOffsetBytes() uint64 // 分片内的偏移字节

	// 总计
	GetFileTotalChunks() uint64 // 文件分片数量
	GetFileChunkBytes() uint64  // 文件分片尺寸
	GetFileTotalBytes() uint64  // 文件总尺寸
	GetFileUUID() string        // UUID
	GetFileName() string        // 文件名
	GetCurrentSize() uint64     // 本次上传尺寸
}

var _ ChunkInfor = &ChunkInfo{}

type ChunkInfo struct {
	ChunkOffsetBytes uint64 // chunk offset bytes
	ChunkIndex       uint64 // index of chunk
	CurrentSize      uint64 // 当前上传切片总尺寸  // 从上传中自动获取

	FileTotalBytes  uint64 // 文件总尺寸(字节)
	FileTotalChunks uint64 // 文件分割分片数量
	FileChunkBytes  uint64 // 文件每个分片尺寸(字节)
	FileUUID        string // 文件唯一标识
	FileName        string // 文件路径名   // 从上传中自动获取
	Mapping         map[string]string
}

func (c *ChunkInfo) getFormField(field string) string {
	if c.Mapping == nil {
		return field
	}
	if v, y := c.Mapping[field]; y {
		return v
	}
	return field
}

func (c *ChunkInfo) BatchSet(m param.Store) {
	c.FileUUID = m.String(c.getFormField(`fileUUID`))
	c.ChunkIndex = m.Uint64(c.getFormField(`chunkIndex`))
	c.FileTotalBytes = m.Uint64(c.getFormField(`fileTotalBytes`))
	c.FileChunkBytes = m.Uint64(c.getFormField(`fileChunkBytes`))
	c.FileTotalChunks = m.Uint64(c.getFormField(`fileTotalChunks`))
	c.ChunkOffsetBytes = m.Uint64(c.getFormField(`chunkOffsetBytes`))
}

func (c *ChunkInfo) BatchSetByURLValues(m url.Values) {
	c.FileUUID = m.Get(c.getFormField(`fileUUID`))
	c.ChunkIndex = param.AsUint64(m.Get(c.getFormField(`chunkIndex`)))
	c.FileTotalBytes = param.AsUint64(m.Get(c.getFormField(`fileTotalBytes`)))
	c.FileChunkBytes = param.AsUint64(m.Get(c.getFormField(`fileChunkBytes`)))
	c.FileTotalChunks = param.AsUint64(m.Get(c.getFormField(`fileTotalChunks`)))
	c.ChunkOffsetBytes = param.AsUint64(m.Get(c.getFormField(`chunkOffsetBytes`)))
}

// - 当前分片 -

// 当前分片索引编码
func (c *ChunkInfo) GetChunkIndex() uint64 {
	return c.ChunkIndex
}

// 分片内的偏移字节
func (c *ChunkInfo) GetChunkOffsetBytes() uint64 {
	return c.ChunkOffsetBytes
}

// - 总计 -

// 文件分片数量
func (c *ChunkInfo) GetFileTotalChunks() uint64 {
	return c.FileTotalChunks
}

// 文件分片尺寸
func (c *ChunkInfo) GetFileChunkBytes() uint64 {
	if c.FileChunkBytes > 0 {
		return c.FileChunkBytes
	}
	if c.GetFileTotalChunks() <= 0 {
		return 0
	}
	return c.GetFileTotalBytes() / c.GetFileTotalChunks()
}

// 文件总尺寸
func (c *ChunkInfo) GetFileTotalBytes() uint64 {
	return c.FileTotalBytes
}

// UUID
func (c *ChunkInfo) GetFileUUID() string {
	return c.FileUUID
}

// GetFileName 文明名称
func (c *ChunkInfo) GetFileName() string {
	return c.FileName
}

// GetCurrentSize 当前上传切片总尺寸
func (c *ChunkInfo) GetCurrentSize() uint64 {
	return c.CurrentSize
}
