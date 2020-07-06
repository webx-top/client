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

package markdown

import (
	"net/http"
	"net/url"
	"time"

	uploadClient "github.com/webx-top/client/upload"
)

func init() {
	uploadClient.Register(`markdown`, func() uploadClient.Client {
		return New()
	})
}

func New() uploadClient.Client {
	client := &Markdown{}
	client.BaseClient = uploadClient.New(client)
	return client
}

type Markdown struct {
	*uploadClient.BaseClient
}

func (a *Markdown) Name() string {
	return "editormd-image-file"
}

func (a *Markdown) Result() (r string) {
	succed := "0" // 0 表示上传失败，1 表示上传成功
	if a.GetError() != nil {
		succed = "1"
	}
	callback := a.Form(`callback`)
	dialogID := a.Form(`dialog_id`)
	if len(callback) > 0 && len(dialogID) > 0 {
		//跨域上传返回操作
		nextURL := callback + "?dialog_id=" + dialogID + "&temp=" + time.Now().Format(`20060102150405`) + "&success=" + succed + "&message=" + url.QueryEscape(a.Error()) + "&url=" + a.Data.FileURL
		a.Context.Response().Redirect(nextURL, http.StatusFound)
	} else {
		r = `{"success":` + succed + `,"message":"` + a.Error() + `","url":"` + a.Data.FileURL + `","id":"` + a.Data.FileIdString() + `"}`
	}
	return
}
