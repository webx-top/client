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

package ueditor

import (
	"path"
	"regexp"
	"strings"

	uploadClient "github.com/webx-top/client/upload"
)

func init() {
	uploadClient.Register(`ueditor`, func() uploadClient.Client {
		return New()
	})
}

var FormField = `upfile`

func New() uploadClient.Client {
	client := &UEditor{}
	client.BaseClient = uploadClient.New(client, FormField)
	return client
}

type UEditor struct {
	*uploadClient.BaseClient
}

var callbackNameRegExp = regexp.MustCompile(`^[\w_]+$`)

func (a *UEditor) BuildResult() uploadClient.Client {
	var publicURL string
	if a.Form("immediate") == "1" {
		publicURL = "!" + a.Data.FileURL
	} else {
		publicURL = a.Data.FileURL
	}
	a.RespData = &Data{
		State:    `SUCCESS`,
		URL:      publicURL,
		Title:    a.Data.FileName,
		Original: a.Data.FileName,
		Type:     strings.ToLower(path.Ext(a.Data.FileName)),
		Size:     a.Data.FileSize,
	}
	callback := a.Query(`callback`)
	if len(callback) > 0 && callbackNameRegExp.MatchString(callback) {
		a.JSONPVarName = callback
	}
	return a
}
