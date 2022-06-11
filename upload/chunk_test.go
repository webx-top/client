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
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/testing/test"
)

func testChunkUpload(t *testing.T, graduallyMerge bool, asyncMergeAll ...bool) {
	log.SetLevel(`Debug`)
	log.Sync()
	var subdir string
	if graduallyMerge {
		subdir = `/graduallyMerge`
	} else if len(asyncMergeAll) > 0 && asyncMergeAll[0] {
		subdir = `/asyncMergeAll`
	} else {
		subdir = `/syncMergeAll`
	}
	os.RemoveAll("../_testdata" + subdir)
	path := "../_testdata" + subdir + "/test.txt" //要上传文件所在路径
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	var file *os.File
	var err error
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
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
	defer file.Close()
	file.Seek(0, 0)
	b, err := ioutil.ReadAll(file)
	test.Eq(t, nil, err)
	chunks := 15
	chunkSize := len(b) / chunks
	cu := &ChunkUpload{
		TempDir:        `../_testdata` + subdir + `/chunk_temp`,
		SaveDir:        `../_testdata` + subdir + `/chunk_merged`,
		GraduallyMerge: graduallyMerge,
	}
	if len(asyncMergeAll) > 0 {
		cu.SetAsyncMerge(asyncMergeAll[0])
	}
	wg := &sync.WaitGroup{}
	wg.Add(chunks)
	upload := func(r io.Reader, chunkIndex int) {
		//chunkStartTime := time.Now()
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		filename := filepath.Base(path)
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
		if !graduallyMerge {
			chunkTempFile := cu.ChunkFilename(chunkIndex)
			if _, err := os.Stat(chunkTempFile); err != nil {
				t.Log(err)
			}
		}
		wg.Done()
		//log.Warn(subdir + ` chunk(` + fmt.Sprintf(`%d`, chunkIndex) + `) elapsed: ` + time.Since(chunkStartTime).String())
	}
	startTime := time.Now()
	file.Seek(0, 0)
	for i := 0; i < chunks; i++ {
		if i == chunks-1 {
			chunkSize = len(b) - chunkSize*(chunks-1)
		}
		data := make([]byte, chunkSize)
		//fmt.Printf("offset: %d (%d)\n", i*chunkSize, i*chunkSize+chunkSize)
		n, err := file.Read(data)
		if err == io.EOF {
			wg.Done()
			continue
		}
		buf := bytes.NewBuffer(data[:n])
		go upload(buf, i)
	}
	wg.Wait()
	log.Warn(subdir + ` elapsed: ` + time.Since(startTime).String())
	uploaded, err := ioutil.ReadFile(cu.GetSavePath())
	/*
		if err != nil || string(uploaded) != string(b) {
			_, err = cu.MergeAll(uint64(chunks), uint64(chunkSize), filepath.Base(path), false)
			test.Eq(t, nil, err)
			uploaded, err = ioutil.ReadFile(cu.GetSavePath())
		}
		// */
	test.Eq(t, nil, err)
	test.Eq(t, string(b), string(uploaded))
	//os.RemoveAll("../_testdata")
}

func TestChunkUploadAsyncMergeAll(t *testing.T) {
	testChunkUpload(t, false, true)
}

func TestChunkUploadSyncMergeAll(t *testing.T) {
	testChunkUpload(t, false, false)
}

func TestChunkUploadGraduallyMerge(t *testing.T) {
	testChunkUpload(t, true)
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
	test.Eq(t, uint64(135), ci.FileTotalChunks)
	test.Eq(t, uint64(500), ci.CurrentSize)
	test.Eq(t, uint64(500), ci.FileChunkBytes)

	found = ci.parseHeader(`bytes 1000-1499/67589`)
	com.Dump(ci)
	test.True(t, found)
	test.Eq(t, uint64(1000), ci.ChunkOffsetBytes)
	test.Eq(t, uint64(1499), ci.ChunkEndBytes)
	test.Eq(t, uint64(67589), ci.FileTotalBytes)
	test.Eq(t, uint64(3), ci.ChunkIndex)
	test.Eq(t, uint64(135), ci.FileTotalChunks)
	test.Eq(t, uint64(500), ci.CurrentSize)
	test.Eq(t, uint64(500), ci.FileChunkBytes)
}
