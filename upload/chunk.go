package upload

import (
	"time"

	"github.com/webx-top/echo/param"
)

type ChunkUpload struct {
	TempDir           string
	SaveDir           string
	TempLifetime      time.Duration
	UID               interface{} // number or string
	fileNameGenerator FileNameGenerator
	savePath          string
}

func (c *ChunkUpload) GetUIDString() string {
	uid := param.AsString(c.UID)
	if len(uid) == 0 {
		uid = `0`
	}
	return uid
}

func (c *ChunkUpload) SetFileNameGenerator(generator FileNameGenerator) *ChunkUpload {
	c.fileNameGenerator = generator
	return c
}

func (c *ChunkUpload) FileNameGenerator() FileNameGenerator {
	if c.fileNameGenerator == nil {
		return DefaultNameGenerator
	}
	return c.fileNameGenerator
}

func (c *ChunkUpload) GetSavePath() string {
	return c.savePath
}
