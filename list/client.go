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
package list

import (
	"github.com/admpub/nging/application/library/common"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/reflectx"
	X "github.com/webx-top/webx"
)

func New() *BaseClient {
	return &BaseClient{SearchPK: true}
}

type BaseClient struct {
	*X.Context
	List       *common.List
	SearchPK   bool
	Middleware func(r db.Result) db.Result
}

func (a *BaseClient) Init(ctx *X.Context, ls *common.List) Client {
	a.Context = ctx
	a.List = ls
	return a
}

//Apply 数据
func (a *BaseClient) Apply(args ...string) error {
	if a.Middleware == nil {
		a.Middleware = func(r db.Result) db.Result {
			return r.Where(a.Build()).OrderBy(a.Sorts()...)
		}
	}
	a.List.AddMiddleware(a.Middleware)
	r, err := a.List.DataTable(a.Context)
	if len(args) > 0 {
		a.Context.Assign(args[0], r)
	} else {
		a.Context.AssignX(&r)
	}
	return err
}

//Sorts 排序
func (a *BaseClient) Sorts(...func(*reflectx.FieldInfo, string) string) []interface{} {
	return []interface{}{}
}

//UsePK 是否自动搜索主键字段
func (a *BaseClient) PrimaryKey(on bool) Client {
	a.SearchPK = on
	return a
}

//Build 生成搜索条件
func (a *BaseClient) Build(defaultFields ...string) *db.Compounds {
	return db.NewCompounds()
}

//Client 客户端接口
type Client interface {
	//初始化数据
	Init(*X.Context, *common.List) Client

	//结果数据
	Apply(...string) error

	//排序方式
	Sorts(...func(*reflectx.FieldInfo, string) string) []interface{}

	//是否自动搜索主键字段
	PrimaryKey(bool) Client

	//生成搜索条件
	Build(defaultFields ...string) *db.Compounds
}
