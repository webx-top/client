package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/admpub/log"
)

// 合并切片文件
func (c *ChunkUpload) merge(info ChunkInfor, fileName, savePath string) error {
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
	uid := c.GetUIDString()
	fileChunkBytes := int64(info.GetFileChunkBytes())
	if fileChunkBytes <= 0 {
		if info.GetChunkIndex() == 0 {
			return err
		}
		// 分片大小获取
		fi, err := os.Stat(filepath.Join(c.TempDir, uid, fileName+"_0"))
		if err != nil {
			return err
		}
		fileChunkBytes = fi.Size()
	}
	// 设置文件写入偏移量
	file.Seek(fileChunkBytes*int64(info.GetChunkIndex()), 0)

	chunkFilePath := filepath.Join(c.TempDir, uid, fmt.Sprintf(`%s_%d`, fileName, info.GetChunkIndex()))
	log.Debug("分片路径: ", chunkFilePath)

	chunkFileObj, err := os.Open(chunkFilePath)
	if err != nil {
		return fmt.Errorf("打开分片文件失败: %w", err)
	}

	_, err = WriteTo(chunkFileObj, file)

	chunkFileObj.Close()

	if err != nil {
		return fmt.Errorf("文件上传失败: %w", err)
	}
	// 删除文件 需要先关闭该文件
	err = os.Remove(chunkFilePath)
	if err != nil {
		return fmt.Errorf("临时记录文件删除失败: %w", err)
	}
	log.Debug("文件复制完毕")
	return err
}

// 合并切片文件
func (c *ChunkUpload) Merge(info ChunkInfor, saveFileName string) (savePath string, err error) {
	var saveName string
	saveName, err = c.FileNameGenerator()(saveFileName)
	if err != nil {
		return
	}
	savePath = filepath.Join(c.SaveDir, saveName)
	c.savePath = savePath
	err = c.merge(info, saveFileName, savePath)
	return
}

// 合并某个文件的所有切片
func (c *ChunkUpload) MergeAll(chunkFileNames []string, saveFileName string) (savePath string, err error) {
	if err = os.MkdirAll(c.SaveDir, os.ModePerm); err != nil {
		return
	}
	var saveName string
	saveName, err = c.FileNameGenerator()(saveFileName)
	if err != nil {
		return
	}
	savePath = filepath.Join(c.SaveDir, saveName)
	c.savePath = savePath
	// 打开之前上传文件
	var file *os.File
	file, err = os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("打开之前上传文件不存在: %w", err)
		return
	}
	defer file.Close()
	uid := c.GetUIDString()
	for _, chunkFileName := range chunkFileNames {
		chunkFilePath := filepath.Join(c.TempDir, uid, chunkFileName)
		cfile, cerr := os.Open(chunkFilePath)
		if cerr != nil {
			err = fmt.Errorf("分片文件“%s”打开失败: %w", chunkFilePath, cerr)
			return
		}
		_, err = WriteTo(cfile, file)

		cfile.Close()

		if err != nil {
			err = fmt.Errorf("分片文件合并失败: %w", err)
			return
		}
		// 删除文件 需要先关闭该文件
		err = os.Remove(chunkFilePath)
		if err != nil {
			err = fmt.Errorf("分片文件删除失败: %w", err)
			return
		}
	}
	return
}
