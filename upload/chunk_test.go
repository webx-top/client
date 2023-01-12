package upload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/testing/test"
)

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
	uploadTestFile(t, subdir, file, chunks, 0)
	file.Close()
	//os.RemoveAll("../_testdata")
}

func _TestRealFile(t *testing.T) {
	subdir := `/realfile`
	path := "../_testdata" + subdir + "/" //要上传文件所在路径
	os.MkdirAll(path, os.ModePerm)
	file, err := os.Open(`/Users/hank/go/src/github.com/admpub/nging/dist/nging_windows_amd64.tar.gz`)
	if err != nil {
		t.Error(err)
	}
	chunkSize := 1048576 * 2 // 2M
	uploadTestFile(t, subdir, file, 0, chunkSize)
	file.Close()
}

func uploadTestFile(t *testing.T, subdir string, file *os.File, chunks int, chunkSize int) {
	file.Seek(0, 0)
	b, err := ioutil.ReadAll(file)
	test.Eq(t, nil, err)
	cu := &ChunkUpload{
		TempDir: `../_testdata` + subdir + `/chunk_temp`,
		SaveDir: `../_testdata` + subdir + `/chunk_merged`,
	}
	if chunks > 0 {
		chunkSize = len(b) / chunks
	} else {
		chunks = int(TotalChunks(uint64(len(b)), uint64(chunkSize)))
	}
	wg := &sync.WaitGroup{}
	wg.Add(chunks)
	upload := func(r io.Reader, chunkIndex int, chunkSize int) {
		//chunkStartTime := time.Now()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		filename := file.Name()
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
		req.Form.Add(`fileTotalBytes`, fmt.Sprintf(`%d`, len(b)))
		n, err := cu.Upload(req)
		test.Eq(t, nil, err)
		test.NotEq(t, 0, n)
		chunkTempFile := cu.ChunkFilename(chunkIndex)
		if _, err := os.Stat(chunkTempFile); err != nil {
			t.Log(err)
		}
		wg.Done()
		//log.Warn(subdir + ` chunk(` + fmt.Sprintf(`%d`, chunkIndex) + `) elapsed: ` + time.Since(chunkStartTime).String())
	}
	startTime := time.Now()
	file.Seek(0, 0)
	for i := 0; i < chunks; i++ {
		offset := i * chunkSize
		if i == chunks-1 {
			chunkSize = len(b) - chunkSize*(chunks-1)
		}
		data := make([]byte, chunkSize)
		fmt.Printf("chunkIndex: %d offset: %d (%d) chunkSize: %d\n", i, offset, offset+chunkSize, chunkSize)
		n, err := file.Read(data)
		if err == io.EOF {
			wg.Done()
			continue
		}
		buf := bytes.NewBuffer(data[:n])
		go upload(buf, i, chunkSize)
	}
	wg.Wait()
	log.Warn(subdir + ` elapsed: ` + time.Since(startTime).String())
	uploaded, err := ioutil.ReadFile(cu.GetSavePath())
	test.Eq(t, nil, err)
	test.Eq(t, len(b), len(uploaded))
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
