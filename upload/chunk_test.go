package upload

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
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
	"github.com/webx-top/com/ratelimit"
	"github.com/webx-top/echo/testing/test"
)

var speedBytes int64 = 3 * 1024 * 1024 // 3Mb/s

func init() {
	log.SetLevel(`Warn`)
	log.Sync()
	path := "../_testdata/"
	os.RemoveAll(path)
}

func testChunkUpload(t *testing.T, index ...int) {
	subdir := `/mergeAll`
	path := "../_testdata" + subdir + "/" //要上传文件所在路径
	os.MkdirAll(path, os.ModePerm)
	if len(index) > 0 {
		path += fmt.Sprintf(`test_%d.txt`, index[0])
	} else {
		path += `test.txt`
	}
	var file *os.File
	var err error
	if _, err = os.Stat(path); err != nil && os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			t.Error(err)
		}

		for i := 1; i <= 1500; i++ {
			file.WriteString(fmt.Sprintf("%05d github.com/webx-top/client\n", i))
		}
	} else {
		file, err = os.Open(path)
		if err != nil {
			t.Error(err)
		}
	}
	chunks := 15
	limitReader := ratelimit.New(speedBytes).NewReadSeeker(file)
	fi, err := file.Stat()
	assert.NoError(t, err)
	uploadTestFile(t, subdir, limitReader, fi.Size(), file.Name(), chunks, 0)
	file.Close()
	//os.RemoveAll("../_testdata")
}

func _TestRealFile(t *testing.T) {
	subdir := `/realfile`
	path := "../_testdata" + subdir + "/" //要上传文件所在路径
	os.MkdirAll(path, os.ModePerm)
	chunkSize := 1048576 * 2 // 2M
	fileList := []string{
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_windows_386.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_windows_amd64.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_linux_386.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_linux_amd64.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_linux_arm64.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_linux_arm-5.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_linux_arm-6.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_linux_arm-7.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_darwin_amd64.tar.gz`,
		`/Users/hank/go/src/github.com/admpub/nging/dist/nging_darwin_arm64.tar.gz`,
	}
	wg := sync.WaitGroup{}
	wg.Add(len(fileList))
	for _, filePath := range fileList {
		go func(filePath string) {
			defer wg.Done()
			file, err := os.Open(filePath)
			if err != nil {
				t.Error(err)
			}
			limitReader := ratelimit.New(speedBytes).NewReadSeeker(file)
			fi, err := file.Stat()
			assert.NoError(t, err)
			uploadTestFile(t, subdir, limitReader, fi.Size(), file.Name(), 0, chunkSize)
			file.Close()
		}(filePath)
	}
	wg.Wait()
}

func uploadTestFile(t *testing.T, subdir string, readSeeker io.ReadSeeker, totalSize int64, fileName string, chunks int, chunkSize int) {
	if chunks > 0 {
		chunkSize = int(totalSize) / chunks
	} else {
		chunks = int(TotalChunks(uint64(totalSize), uint64(chunkSize)))
	}
	tempDir := `../_testdata` + subdir + `/chunk_temp`
	saveDir := `../_testdata` + subdir + `/chunk_merged`
	wg := &sync.WaitGroup{}
	wg.Add(chunks)
	upload := func(r io.Reader, chunkIndex int, chunkSize int) {
		cu := &ChunkUpload{
			TempDir: tempDir,
			SaveDir: saveDir,
		}
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
		req.Form.Add(`chunkIndex`, fmt.Sprintf(`%d`, chunkIndex))
		req.Form.Add(`fileTotalChunks`, fmt.Sprintf(`%d`, chunks))
		req.Form.Add(`fileChunkBytes`, fmt.Sprintf(`%d`, chunkSize))
		req.Form.Add(`fileTotalBytes`, fmt.Sprintf(`%d`, totalSize))
		n, err := cu.Upload(req)
		test.Eq(t, nil, err)
		test.NotEq(t, 0, n)
		wg.Done()
		log.Warn(`Post: ` + fileName + ` chunk(` + fmt.Sprintf(`%d`, chunkIndex) + `) elapsed: ` + time.Since(chunkStartTime).String())
	}
	startTime := time.Now()
	readSeeker.Seek(0, 0)
	for i := 0; i < chunks; i++ {
		offset := i * chunkSize
		if i == chunks-1 {
			chunkSize = int(totalSize) - chunkSize*(chunks-1)
		}
		data := make([]byte, chunkSize)
		fmt.Printf("%v => chunkIndex: %d offset: %d (%d) chunkSize: %d\n", fileName, i, offset, offset+chunkSize, chunkSize)
		n, err := readSeeker.Read(data)
		if err == io.EOF {
			wg.Done()
			continue
		}
		buf := bytes.NewBuffer(data[:n])
		go upload(buf, i, chunkSize)
	}
	wg.Wait()
	log.Warn(fileName + ` elapsed: ` + time.Since(startTime).String())
	savePath, err := genSavePath(saveDir, filepath.Base(fileName), DefaultNameGenerator)
	assert.NoError(t, err)
	fi, err := os.Stat(savePath)
	test.Eq(t, nil, err)
	test.Eq(t, totalSize, fi.Size())
}

func TestChunkUploadMergeAll(t *testing.T) {
	testChunkUpload(t)
}

func TestChunkUploadSyncMergeAllBatch(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			testChunkUpload(t, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestChunkUploadParseHeader(t *testing.T) {
	ci := &ChunkInfo{}
	found := ci.parseHeader(`bytes 500-999/67589`)
	com.Dump(ci)
	test.True(t, found)
	test.Eq(t, uint64(500), ci.ChunkOffsetBytes)
	test.Eq(t, uint64(999), ci.ChunkEndBytes)
	test.Eq(t, uint64(67589), ci.FileTotalBytes)
	test.Eq(t, uint64(2), ci.ChunkIndex)
	test.Eq(t, uint64(136), ci.FileTotalChunks)
	test.Eq(t, uint64(500), ci.CurrentSize)
	test.Eq(t, uint64(500), ci.FileChunkBytes)

	found = ci.parseHeader(`bytes 1000-1499/67589`)
	com.Dump(ci)
	test.True(t, found)
	test.Eq(t, uint64(1000), ci.ChunkOffsetBytes)
	test.Eq(t, uint64(1499), ci.ChunkEndBytes)
	test.Eq(t, uint64(67589), ci.FileTotalBytes)
	test.Eq(t, uint64(3), ci.ChunkIndex)
	test.Eq(t, uint64(136), ci.FileTotalChunks)
	test.Eq(t, uint64(500), ci.CurrentSize)
	test.Eq(t, uint64(500), ci.FileChunkBytes)
}
