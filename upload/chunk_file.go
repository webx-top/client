package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/admpub/errors"
	"github.com/admpub/log"
)

// 合并切片文件
func mergeFile(index int64, fileName, filePath, tmpFileDir string) error {
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("打开之前上传文件不存在: %w", err)
	}
	defer file.Close()
	// 分片大小获取
	fi, _ := os.Stat(tmpFileDir + fileName + "_0")
	chunkSize := fi.Size()
	// 设置文件写入偏移量
	file.Seek(chunkSize*index, 0)
	chunkFilePath := tmpFileDir + fileName + "_" + fmt.Sprintf(`%d`, index)
	log.Debug("分片路径: ", chunkFilePath)
	chunkFileObj, err := os.Open(chunkFilePath)
	if err != nil {
		return fmt.Errorf("打开分片文件失败: %w", err)
	}
	defer chunkFileObj.Close()

	// 上传总数
	var totalLen int
	// 写入数据
	data := make([]byte, 1024)
	for {
		n, err := chunkFileObj.Read(data)
		if err == io.EOF {
			// 删除文件 需要先关闭改文件
			chunkFileObj.Close()
			err := os.Remove(chunkFilePath)
			if err != nil {
				log.Debug("临时记录文件删除失败: ", err)
			}
			log.Debug("文件复制完毕")
			break
		}
		len, err := file.Write(data[:n])
		if err != nil {
			return fmt.Errorf("文件上传失败: %w", err)
		}
		totalLen += len
	}
	return nil
}

// 分片上传
func chunkUpload(r *http.Request, tmpFileDir string) (int, error) {
	// 分片序号
	chunkIndex := r.FormValue("chunkindex")
	// 获取上传文件
	upFile, fileHeader, err := r.FormFile("file")

	if err != nil {
		return 0, errors.New("上传文件错误")
	}

	// 新文件创建
	filePath := tmpFileDir + fileHeader.Filename + "_" + chunkIndex
	// 获取现在文件大小
	fi, _ := os.Stat(filePath)
	// 判断文件是否传输完成
	if fi.Size() == fileHeader.Size {
		return 0, errors.New("文件已存在, 不继续上传")
	}
	start := fi.Size()

	// 进行断点上传
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, errors.New("打开之前上传文件不存在")
	}
	defer file.Close()

	// 将数据写入文件
	total, err := uploadFile(upFile, start, file, start)
	return total, err
}

// 上传文件
func uploadFile(upfile multipart.File, upSeek int64, file *os.File, fSeek int64) (int, error) {
	// 上传文件大小记录
	fileSzie := 0
	// 设置上传偏移量
	upfile.Seek(upSeek, 0)
	// 设置文件偏移量
	file.Seek(fSeek, 0)
	data := make([]byte, 1024)
	for {
		total, err := upfile.Read(data)
		if err == io.EOF {
			//fmt.Println("文件复制完毕")
			break
		}
		len, err := file.Write(data[:total])
		if err != nil {
			return 0, errors.New("文件上传失败")
		}
		// 记录上传长度
		fileSzie += len
	}
	return fileSzie, nil
}
