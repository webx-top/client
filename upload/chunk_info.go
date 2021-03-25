package upload

import (
	"github.com/webx-top/echo/param"
)

type ChunkInfor interface {
	// 当前分片
	GetChunkIndex() uint64       // 当前分片索引编码
	GetChunkOffsetBytes() uint64 // 分片内的偏移字节

	// 总计
	GetTotalChunks() uint64 // 文件分片数量
	GetChunkBytes() uint64  // 文件分片尺寸
	GetTotalBytes() uint64  // 文件总尺寸
	GetUUID() string        // UUID
}

var _ ChunkInfor = &ChunkInfo{}

type ChunkInfo struct {
	UUID             string
	ChunkIndex       uint64
	TotalBytes       uint64
	ChunkBytes       uint64
	TotalChunks      uint64
	ChunkOffsetBytes uint64
}

func (c *ChunkInfo) BatchSet(m param.Store) {
	c.UUID = m.String(`UUID`)
	c.ChunkIndex = m.Uint64(`ChunkIndex`)
	c.TotalBytes = m.Uint64(`TotalBytes`)
	c.ChunkBytes = m.Uint64(`ChunkBytes`)
	c.TotalChunks = m.Uint64(`TotalChunks`)
	c.ChunkOffsetBytes = m.Uint64(`ChunkOffsetBytes`)
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
func (c *ChunkInfo) GetTotalChunks() uint64 {
	return c.TotalChunks
}

// 文件分片尺寸
func (c *ChunkInfo) GetChunkBytes() uint64 {
	return c.ChunkBytes
}

// 文件总尺寸
func (c *ChunkInfo) GetTotalBytes() uint64 {
	return c.TotalBytes
}

// UUID
func (c *ChunkInfo) GetUUID() string {
	return c.UUID
}
