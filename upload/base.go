package upload

import (
	"fmt"
	"net/http"

	"github.com/webx-top/echo"
)

func New(object Client, formFields ...string) *BaseClient {
	formField := DefaultFormField
	if len(formFields) > 0 {
		formField = formFields[0]
	}
	return &BaseClient{Object: object, FormField: formField, Results: Results{}}
}

var DefaultFormField = `filedata`

type BaseClient struct {
	Data *Result
	echo.Context
	Object       Client
	FormField    string // 表单文件字段名
	Code         int    // HTTP code
	ContentType  string
	JSONPVarName string
	err          error
	RespData     interface{}
	Results      Results
}

func (a *BaseClient) Init(ctx echo.Context, res *Result) {
	a.Context = ctx
	a.Data = res
}

func (a *BaseClient) Name() string {
	if len(a.FormField) == 0 {
		return DefaultFormField
	}
	return a.FormField
}

func (a *BaseClient) SetError(err error) Client {
	a.err = err
	return a
}

func (a *BaseClient) GetError() error {
	return a.err
}

func (a *BaseClient) ErrorString() string {
	if a.err != nil {
		return a.err.Error()
	}
	return ``
}

func (a *BaseClient) Body() (file ReadCloserWithSize, err error) {
	file, a.Data.FileName, err = Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	a.Data.FileSize = file.Size()
	a.Data.Md5, err = file.Md5()
	return
}

func (a *BaseClient) BuildResult() Client {
	data := a.Context.Data()
	data.SetData(echo.H{
		`Url`: a.Data.FileURL,
		`Id`:  a.Data.FileIdString(),
	}, 1)
	if a.err != nil {
		data.SetError(a.err)
	}
	a.RespData = data
	return a
}

func (a *BaseClient) GetRespData() interface{} {
	return a.RespData
}

func (a *BaseClient) GetBatchUploadResults() Results {
	return a.Results
}

func (a *BaseClient) SetRespData(data interface{}) Client {
	a.RespData = data
	return a
}

func (a *BaseClient) Response() error {
	if a.Object != nil {
		a.Object.BuildResult()
	} else {
		a.BuildResult()
	}
	if a.Code > 0 {
		return a.responseContentType(a.Code)
	}
	return a.responseContentType(http.StatusOK)
}

func (a *BaseClient) responseContentType(code int) error {
	switch a.ContentType {
	case `string`:
		return a.String(fmt.Sprint(a.RespData), code)
	case `xml`:
		return a.XML(a.RespData, code)
	case `redirect`:
		return a.Redirect(fmt.Sprint(a.RespData), code)
	default:
		if len(a.JSONPVarName) > 0 {
			return a.JSONP(a.JSONPVarName, a.RespData, code)
		}
		return a.JSON(a.RespData, code)
	}
}
