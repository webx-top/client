package upload

import (
	"bytes"
	"io"
	"io/ioutil"
	"path"

	"github.com/webx-top/client/upload/watermark"
	"github.com/webx-top/echo"
	"github.com/webx-top/image"
)

type Options struct {
	ClientName       string
	Result           *Result
	Storer           Storer
	WatermarkOptions *image.WatermarkOptions
	Checker          func(*Result) error
	Callback         func(*Result, io.Reader, io.Reader) error
}

type OptionsSetter func(options *Options)

func OptClientName(clientName string) OptionsSetter {
	return func(options *Options) {
		options.ClientName = clientName
	}
}

func OptResult(result *Result) OptionsSetter {
	return func(options *Options) {
		options.Result = result
	}
}

func OptStorer(storer Storer) OptionsSetter {
	return func(options *Options) {
		options.Storer = storer
	}
}

func OptWatermarkOptions(wmOpt *image.WatermarkOptions) OptionsSetter {
	return func(options *Options) {
		options.WatermarkOptions = wmOpt
	}
}

func OptChecker(checker func(*Result) error) OptionsSetter {
	return func(options *Options) {
		options.Checker = checker
	}
}

func OptCallback(callback func(*Result, io.Reader, io.Reader) error) OptionsSetter {
	return func(options *Options) {
		options.Callback = callback
	}
}

type ReaderAndSizer interface {
	io.Reader
	Sizer
}

func Upload(ctx echo.Context, opts ...OptionsSetter) Client {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	client := Get(options.ClientName)
	client.Init(ctx, options.Result)
	body, err := client.Body()
	if err != nil {
		return client.SetError(err)
	}
	defer body.Close()
	if options.Checker != nil {
		err = options.Checker(options.Result)
		if err != nil {
			return client.SetError(err)
		}
	}
	dstFile, err := options.Result.GenFileName()
	if err != nil {
		return client.SetError(err)
	}

	var readerAndSizer ReaderAndSizer = body

	if options.Result.FileType.String() == `image` {
		if options.WatermarkOptions != nil && options.WatermarkOptions.IsEnabled() {
			var b []byte
			b, err = ioutil.ReadAll(body)
			if err != nil {
				return client.SetError(err)
			}
			b, err = watermark.Bytes(b, path.Ext(options.Result.FileName), options.WatermarkOptions)
			if err != nil {
				return client.SetError(err)
			}
			byteReader := bytes.NewReader(b)
			readerAndSizer = byteReader
		}
	}
	options.Result.SavePath, options.Result.FileURL, err = options.Storer.Put(dstFile, readerAndSizer, readerAndSizer.Size())
	if err != nil {
		return client.SetError(err)
	}
	if options.Callback != nil {
		err = options.Callback(options.Result, body, readerAndSizer)
		if err != nil {
			options.Storer.Delete(dstFile)
			return client.SetError(err)
		}
	}
	return client
}
