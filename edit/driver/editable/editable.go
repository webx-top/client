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

package editable

import (
	editClient "github.com/webx-top/client/edit"
	"github.com/webx-top/echo"
)

func init() {
	editClient.Register(`editable`, func() editClient.Client {
		return New()
	})
}

func New() editClient.Client {
	return &Editable{
		BaseClient: editClient.New(),
	}
}

type Editable struct {
	*editClient.BaseClient
}

func (a *Editable) Init(ctx echo.Context, m interface{}) editClient.Client {
	a.Context = ctx
	a.Model = m
	a.StructField = ctx.Form(`name`)
	a.Value = ctx.Form(`value`)
	a.Primary = ctx.Form(`pk`)
	return a
}
