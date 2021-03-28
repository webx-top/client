package ueditor

import (
	"errors"
	"fmt"

	"github.com/webx-top/echo"
)

func NewAction() *Action {
	a := &Action{}
	a.Config = func(echo.Context) (interface{}, error) {
		return DefaultConfig, nil
	}
	a.Upload.Image = EmptyAction
	a.Upload.Scrawl = EmptyAction
	a.Upload.Video = EmptyAction
	a.List.File = EmptyAction
	a.List.Image = EmptyAction
	a.Catch.Image = EmptyAction
	return a
}

type Action struct {
	Config func(echo.Context) (interface{}, error)
	Upload ActionUpload
	List   ActionList
	Catch  ActionCatch
}

var (
	ErrUnsupportedAction = errors.New("unsupported action")
	EmptyAction          = func(echo.Context) (interface{}, error) {
		return nil, ErrUnsupportedAction
	}
)

type ActionUpload struct {
	Image  func(echo.Context) (interface{}, error)
	Scrawl func(echo.Context) (interface{}, error)
	Video  func(echo.Context) (interface{}, error)
	File   func(echo.Context) (interface{}, error)
}

type ActionList struct {
	Image func(echo.Context) (interface{}, error)
	File  func(echo.Context) (interface{}, error)
}

type ActionCatch struct {
	Image func(echo.Context) (interface{}, error)
}

func (a *Action) Handle(c echo.Context) (err error) {
	action := c.Form(`action`)
	var result interface{}
	switch action {
	case `config`:
		result, err = a.Config(c)
	case `uploadimage`: //上传图片
		result, err = a.Upload.Image(c)
	case `uploadscrawl`: //上传涂鸦
		result, err = a.Upload.Scrawl(c)
	case `uploadvideo`: //上传视频
		result, err = a.Upload.Video(c)
	case `uploadfile`: //上传文件
		result, err = a.Upload.File(c)
	case `listimage`: //列出图片
		result, err = a.List.Image(c)
	case `listfile`: //列出文件
		result, err = a.List.File(c)
	case `catchimage`: //抓取远程文件
		result, err = a.Catch.Image(c)
	}

	if err != nil {
		return fmt.Errorf(`%s: %w`, action, err)
	}
	callback := c.Query(`callback`)
	if len(callback) > 0 && callbackNameRegExp.MatchString(callback) {
		return c.JSONP(callback, result)
	}
	return c.JSON(result)
}
