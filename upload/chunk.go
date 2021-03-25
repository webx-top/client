package upload

import "time"

type Chunked struct {
	TempDir      string
	TempLifetime time.Duration
}
