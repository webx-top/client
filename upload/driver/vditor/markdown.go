/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package vditor

import (
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

func init() {
	uploadClient.Register(`markdown`, func() uploadClient.Client {
		return New()
	})
}

var FormField = `file[]`

func New() uploadClient.Client {
	client := &Vditor{}
	client.BaseClient = uploadClient.New(client, FormField)
	return client
}

type Vditor struct {
	*uploadClient.BaseClient
}

type Response struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ResponseUpload struct {
	ErrFiles []string          `json:"errFiles"`
	SuccMap  map[string]string `json:"succMap"`

	// for save remote image
	OriginalURL string `json:"originalURL"`
	URL         string `json:"url"`
}

/*
{
 "msg": "",
 "code": 0,
 "data": {
 "errFiles": ['filename', 'filename2'],
 "succMap": {
   "filename3": "filepath3",
   "filename3": "filepath3"
   }
 }
}
*/

func (a *Vditor) BuildResult() {
	resp := &Response{}
	uploadResult := &ResponseUpload{}
	if a.GetError() != nil {
		uploadResult.ErrFiles = append(uploadResult.ErrFiles, a.Data.FileName)
		resp.Msg = a.GetError().Error()
		resp.Code = 1
	} else {
		uploadResult.SuccMap = map[string]string{
			a.Data.FileName: a.Data.FileURL,
		}
		if m, y := a.Data.Addon.(echo.H); y {
			uploadResult.OriginalURL = m.String(`originalURL`)
			uploadResult.URL = a.Data.FileURL
		}
	}
	resp.Data = uploadResult
	a.RespData = resp
}
