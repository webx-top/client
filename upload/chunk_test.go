package upload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestChunkUpload(t *testing.T) {
	path := "../_testdata/test.txt" //要上传文件所在路径
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	var file *os.File
	var err error
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			t.Error(err)
		}

		for i := 1; i <= 150; i++ {
			file.WriteString(fmt.Sprintf("%03d github.com/webx-top/client\n", i))
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
		TempDir: `../_testdata/tmp`,
		SaveDir: `../_testdata/upload/`,
	}
	wg := &sync.WaitGroup{}
	wg.Add(chunks)
	upload := func(r io.Reader, chunkIndex int) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(path))
		if err != nil {
			writer.Close()
			t.Error(err)
		}
		io.Copy(part, r)
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()
		req.Form = make(url.Values)
		req.Form.Add(`chunkIndex`, fmt.Sprintf(`%d`, chunkIndex))
		req.Form.Add(`totalChunks`, fmt.Sprintf(`%d`, chunks))
		req.Form.Add(`chunkBytes`, fmt.Sprintf(`%d`, chunkSize))
		req.Form.Add(`totalBytes`, fmt.Sprintf(`%d`, len(b)))
		n, err := cu.Upload(req, nil)
		test.Eq(t, nil, err)
		test.NotEq(t, 0, n)
		if res.Code != http.StatusOK {
			t.Error("not 200")
		}

		t.Log(res.Body.String())
		wg.Done()
	}
	file.Seek(0, 0)
	for i := 0; i < chunks; i++ {
		if i == chunks-1 {
			chunkSize = len(b) - chunkSize*(chunks-1)
		}
		data := make([]byte, chunkSize)
		fmt.Printf("offset: %d (%d)\n", i*chunkSize, i*chunkSize+chunkSize)
		n, err := file.Read(data)
		if err == io.EOF {
			//wg.Done()
			//continue
		}
		buf := bytes.NewBuffer(data[:n])
		go upload(buf, i)
	}
	wg.Wait()
	uploaded, err := ioutil.ReadFile(cu.SaveDir + `test.txt`)
	test.Eq(t, nil, err)
	test.Eq(t, string(b), string(uploaded))
	os.RemoveAll("../_testdata")
}
