package upload

import (
	"sync"
	"time"
)

type ChunkUpload struct {
	TempDir      string
	SaveDir      string
	TempLifetime time.Duration
	lock         sync.RWMutex
}
