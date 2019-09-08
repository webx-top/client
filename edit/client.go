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
package edit

import (
	"errors"

	"github.com/webx-top/db/lib/sqlbuilder"
	X "github.com/webx-top/webx"
)

func New() *BaseClient {
	return &BaseClient{}
}

type BaseClient struct {
	*X.Context
	Model       interface{}
	StructField string
	TableField  string
	Value       string
	Primary     string
}

func (a *BaseClient) Init(ctx *X.Context, m interface{}) Client {
	a.Context = ctx
	a.Model = m
	a.StructField = ctx.Form(`field`)
	a.Value = ctx.Form(`value`)
	a.Primary = ctx.Form(`id`)
	return a
}

func (a *BaseClient) Do(fn func(string, string, string) error, validField ...bool) error {
	if len(a.StructField) < 1 {
		return errors.New(`Invalid field name: missing paramter field`)
	}
	if len(a.Primary) == 0 {
		return errors.New(`Primary key value is invalid: missing paramter id`)
	}
	if len(validField) < 1 || validField[0] {
		_, ok := sqlbuilder.Mapper().StructMap(a.Model).Find(a.StructField, true)
		if !ok {
			return errors.New(`Invalid field name: ` + a.StructField)
		}
	}
	a.TableField = a.StructField
	a.Context.MapData(a.Model, map[string][]string{
		a.StructField: []string{a.Value},
	})
	return fn(a.Primary, a.TableField, a.Value)
}

type Client interface {
	//初始化
	Init(*X.Context, interface{}) Client
	Do(func(string, string, string) error, ...bool) error
}
