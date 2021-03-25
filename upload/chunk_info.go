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
	GetTotalChunks() uint64 // 文件分片数量
	GetChunkBytes() uint64  // 文件分片尺寸
	GetTotalBytes() uint64  // 文件总尺寸
	GetUUID() string        // UUID
}

var _ ChunkInfor = &ChunkInfo{}

type ChunkInfo struct {
	ChunkOffsetBytes uint64
	ChunkIndex       uint64

	TotalBytes  uint64
	TotalChunks uint64
	ChunkBytes  uint64
	UUID        string
	Mapping     map[string]string
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
	c.UUID = m.String(c.getFormField(`uuid`))
	c.ChunkIndex = m.Uint64(c.getFormField(`chunkIndex`))
	c.TotalBytes = m.Uint64(c.getFormField(`totalBytes`))
	c.ChunkBytes = m.Uint64(c.getFormField(`chunkBytes`))
	c.TotalChunks = m.Uint64(c.getFormField(`totalChunks`))
	c.ChunkOffsetBytes = m.Uint64(c.getFormField(`chunkOffsetBytes`))
}

func (c *ChunkInfo) BatchSetByURLValues(m url.Values) {
	c.UUID = m.Get(c.getFormField(`uuid`))
	c.ChunkIndex = param.AsUint64(m.Get(c.getFormField(`chunkIndex`)))
	c.TotalBytes = param.AsUint64(m.Get(c.getFormField(`totalBytes`)))
	c.ChunkBytes = param.AsUint64(m.Get(c.getFormField(`chunkBytes`)))
	c.TotalChunks = param.AsUint64(m.Get(c.getFormField(`totalChunks`)))
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
func (c *ChunkInfo) GetTotalChunks() uint64 {
	return c.TotalChunks
}

// 文件分片尺寸
func (c *ChunkInfo) GetChunkBytes() uint64 {
	if c.ChunkBytes > 0 {
		return c.ChunkBytes
	}
	if c.GetTotalChunks() <= 0 {
		return 0
	}
	return c.GetTotalBytes() / c.GetTotalChunks()
}

// 文件总尺寸
func (c *ChunkInfo) GetTotalBytes() uint64 {
	return c.TotalBytes
}

// UUID
func (c *ChunkInfo) GetUUID() string {
	return c.UUID
}
