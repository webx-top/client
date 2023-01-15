package upload_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/admpub/log"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/client/upload"
	"github.com/webx-top/client/upload/utesting"
	"github.com/webx-top/com"
	"github.com/webx-top/com/ratelimit"
	"github.com/webx-top/echo/testing/test"
)

var speedBytes int64 = 3 * 1024 * 1024 // 3Mb/s
var counters = upload.NewCounters()
var delayMerge bool

func init() {
	path := "../_testdata/"
	os.RemoveAll(path)
	log.SetLevel(`Debug`)
	fileTarget := log.NewFileTarget()
	fileTarget.FileName = path + `/test.log`
	log.SetTarget(log.NewConsoleTarget(), fileTarget)
	log.Sync()
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

func TestRealFile(t *testing.T) {
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
				if os.IsNotExist(err) {
					return
				}
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
	cc := counters.Map()
	com.Dump(cc)
	for _, v := range cc {
		assert.Equal(t, 1, v)
	}
}

func uploadTestFile(t *testing.T, subdir string, readSeeker io.ReadSeeker, totalSize int64, fileName string, chunks int, chunkSize int) {
	tempDir := `../_testdata` + subdir + `/chunk_temp`
	saveDir := `../_testdata` + subdir + `/chunk_merged`
	cu := &upload.ChunkUpload{
		TempDir:    tempDir,
		SaveDir:    saveDir,
		DelayMerge: delayMerge,
	}
	cu.OnBeforeMerge(func(ctx context.Context, info upload.ChunkInfor, filename string) error {
		counters.Add(filename, 1)
		return nil
	})
	utesting.UploadTestFile(t, cu, readSeeker, totalSize, fileName, chunks, chunkSize)
	utesting.VerifyUploadedTestFile(t, cu, fileName, totalSize)
}

func TestChunkUploadMergeAll(t *testing.T) {
	testChunkUpload(t)
	delayMerge = true
	testChunkUpload(t)
	delayMerge = false
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
	ci := &upload.ChunkInfo{}
	found := ci.ParseHeaderValue(`bytes 500-999/67589`)
	com.Dump(ci)
	test.True(t, found)
	test.Eq(t, uint64(500), ci.ChunkOffsetBytes)
	test.Eq(t, uint64(999), ci.ChunkEndBytes)
	test.Eq(t, uint64(67589), ci.FileTotalBytes)
	test.Eq(t, uint64(2), ci.ChunkIndex)
	test.Eq(t, uint64(136), ci.FileTotalChunks)
	test.Eq(t, uint64(500), ci.CurrentSize)
	test.Eq(t, uint64(500), ci.FileChunkBytes)

	found = ci.ParseHeaderValue(`bytes 1000-1499/67589`)
	com.Dump(ci)
	test.True(t, found)
	test.Eq(t, uint64(1000), ci.ChunkOffsetBytes)
	test.Eq(t, uint64(1499), ci.ChunkEndBytes)
	test.Eq(t, uint64(67589), ci.FileTotalBytes)
	test.Eq(t, uint64(3), ci.ChunkIndex)
	test.Eq(t, uint64(136), ci.FileTotalChunks)
	test.Eq(t, uint64(500), ci.CurrentSize)
	test.Eq(t, uint64(500), ci.FileChunkBytes)

	fileUUID := `abc123-32321s-122222`
	ci.FileUUID = fileUUID
	ci.FixedUUID()
	test.Eq(t, fileUUID, ci.FileUUID)
	ci.FileUUID = `..abc123-32321s-122222`
	ci.FixedUUID()
	test.Empty(t, ci.FileUUID)

	fileUUID = `abc123_32321s_122222`
	ci.FileUUID = fileUUID
	ci.FixedUUID()
	test.Eq(t, fileUUID, ci.FileUUID)
}
