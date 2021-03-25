package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/admpub/errors"
	"github.com/admpub/log"
	"github.com/webx-top/com"
)

// 合并切片文件
func (c *ChunkUpload) Merge(info ChunkInfor, fileName, savePath string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	saveDir := filepath.Dir(savePath)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return err
	}
	// 打开之前上传文件
	file, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("打开之前上传文件不存在: %w", err)
	}
	defer file.Close()
	chunkSize := int64(info.GetChunkBytes())
	if chunkSize <= 0 {
		if info.GetChunkIndex() == 0 {
			return err
		}
		// 分片大小获取
		fi, err := os.Stat(filepath.Join(c.TempDir, fileName+"_0"))
		if err != nil {
			return err
		}
		chunkSize = fi.Size()
	}
	// 设置文件写入偏移量
	file.Seek(chunkSize*int64(info.GetChunkIndex()), 0)
	chunkFilePath := filepath.Join(c.TempDir, fmt.Sprintf(`%s_%d`, fileName, info.GetChunkIndex()))
	log.Debug("分片路径: ", chunkFilePath)
	chunkFileObj, err := os.Open(chunkFilePath)
	if err != nil {
		return fmt.Errorf("打开分片文件失败: %w", err)
	}

	// 上传总数
	_, err = WriteTo(chunkFileObj, file)

	chunkFileObj.Close()

	if err != nil {
		return fmt.Errorf("文件上传失败: %w", err)
	}
	// 删除文件 需要先关闭该文件
	err = os.Remove(chunkFilePath)
	if err != nil {
		log.Debug("临时记录文件删除失败: ", err)
	}
	log.Debug("文件复制完毕")
	return err
}

// 分片上传
func (c *ChunkUpload) Upload(r *http.Request, mapping map[string]string) (int64, error) {
	// 获取上传文件
	upFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		return 0, errors.New("上传文件错误")
	}

	info := &ChunkInfo{Mapping: mapping}
	info.BatchSetByURLValues(r.Form)
	com.Dump(info)

	fileName := fileHeader.Filename
	if err := os.MkdirAll(c.TempDir, os.ModePerm); err != nil {
		return 0, err
	}
	// 新文件创建
	filePath := filepath.Join(c.TempDir, fmt.Sprintf("%s_%d", fileName, info.GetChunkIndex()))
	// 获取现在文件大小
	fi, err := os.Stat(filePath)
	var size int64
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
	} else {
		size = fi.Size()
	}
	// 判断文件是否传输完成
	if size == fileHeader.Size {
		return 0, errors.New("文件已存在, 不继续上传: " + filepath.Base(filePath))
	}
	start := size
	saveStart := size
	offset := int64(info.GetChunkOffsetBytes())
	if offset > 0 {
		start = 0
		saveStart = offset
	}

	// 进行断点上传
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, errors.New("打开之前上传文件不存在")
	}

	defer file.Close()

	// 将数据写入文件
	total, err := uploadFile(upFile, start, file, saveStart)
	if err == nil && total == fileHeader.Size {
		err = c.Merge(info, fileName, filepath.Join(c.SaveDir, fileName))
		if err != nil {
			log.Error(err)
		}
	}
	return total, err
}

// 上传文件
func uploadFile(upfile multipart.File, upSeek int64, file *os.File, fSeek int64) (int64, error) {
	// 设置上传偏移量
	upfile.Seek(upSeek, 0)
	// 设置文件偏移量
	file.Seek(fSeek, 0)
	return WriteTo(upfile, file)
}

func WriteTo(r io.Reader, w io.Writer) (int64, error) {
	var fileSzie int64
	data := make([]byte, 1024)
	for {
		n, err := r.Read(data)
		if err != nil {
			if err == io.EOF {
				return fileSzie, nil
			}
			return 0, err
		}
		len, err := w.Write(data[:n])
		if err != nil {
			return 0, err
		}
		// 记录上传长度
		fileSzie += int64(len)
	}
}
