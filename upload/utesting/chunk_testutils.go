package utesting

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/admpub/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/client/upload"
	"github.com/webx-top/echo/testing/test"
)

func UploadTestFile(t *testing.T, parentCU *upload.ChunkUpload, readSeeker io.ReadSeeker, totalSize int64, fileName string, chunks int, chunkSize int) {
	if chunks > 0 {
		chunkSize = int(totalSize) / chunks
	} else {
		chunks = int(upload.TotalChunks(uint64(totalSize), uint64(chunkSize)))
	}
	fileUUID := uuid.New().String()
	wg := &sync.WaitGroup{}
	wg.Add(chunks)
	uploadChunk := func(r io.Reader, chunkIndex int, chunkSize int) error {
		cu := parentCU.Clone()
		chunkStartTime := time.Now()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		filename := fileName
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			writer.Close()
			t.Error(err)
		}
		io.Copy(part, r)
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Form = make(url.Values)
		req.Form.Add(`fileUUID`, fileUUID)
		req.Form.Add(`chunkIndex`, fmt.Sprintf(`%d`, chunkIndex))
		req.Form.Add(`fileTotalChunks`, fmt.Sprintf(`%d`, chunks))
		req.Form.Add(`fileChunkBytes`, fmt.Sprintf(`%d`, chunkSize))
		req.Form.Add(`fileTotalBytes`, fmt.Sprintf(`%d`, totalSize))
		_, err = cu.Upload(req)
		log.Warn(`Post: ` + fileName + ` chunk(` + fmt.Sprintf(`%d`, chunkIndex) + `) elapsed: ` + time.Since(chunkStartTime).String())
		return err
	}
	startTime := time.Now()
	readSeeker.Seek(0, 0)
	doChunk := func(i int, wg *sync.WaitGroup) {
		offset := i * chunkSize
		if i == chunks-1 {
			chunkSize = int(totalSize) - chunkSize*(chunks-1)
		}
		data := make([]byte, chunkSize)
		fmt.Printf("%v => chunkIndex: %d offset: %d (%d) chunkSize: %d\n", fileName, i, offset, offset+chunkSize, chunkSize)
		n, err := readSeeker.Read(data)
		if err == io.EOF {
			wg.Done()
			return
		}
		up := func(n int, chunkIndex int, chunkSize int) {
			buf := bytes.NewBuffer(data[:n])
			err := uploadChunk(buf, chunkIndex, chunkSize)
			test.Eq(t, nil, err)
			if chunkIndex == 0 {
				buf = bytes.NewBuffer(data[:n])
				err = uploadChunk(buf, chunkIndex, chunkSize)
				assert.ErrorIs(t, err, upload.ErrChunkUploadCompleted)
			}
			wg.Done()
		}
		go up(n, i, chunkSize)
	}
	if chunks > 0 {
		firstWg := &sync.WaitGroup{}
		firstWg.Add(1)
		doChunk(0, firstWg)
		firstWg.Wait()
		wg.Done()
	}
	for i := 1; i < chunks; i++ {
		doChunk(i, wg)
	}
	wg.Wait()
	log.Warn(fileName + ` elapsed: ` + time.Since(startTime).String())
}

func VerifyUploadedTestFile(t *testing.T, parentCU *upload.ChunkUpload, fileName string, totalSize int64) {
	savePath, err := upload.GenSavePath(parentCU.SaveDir, filepath.Base(fileName), parentCU.FileNameGenerator())
	assert.NoError(t, err)
	fi, err := os.Stat(savePath)
	test.Eq(t, nil, err)
	if err == nil {
		test.Eq(t, totalSize, fi.Size())
	}
}
