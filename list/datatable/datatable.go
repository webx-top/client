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

package datatable

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/admpub/nging/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/mysql"
	"github.com/webx-top/db/lib/reflectx"
	"github.com/webx-top/db/lib/sqlbuilder"
	X "github.com/webx-top/webx"
	listClient "github.com/webx-top/webx/lib/client/list"
	"github.com/webx-top/webx/lib/database"
)

func init() {
	_ = fmt.Sprint
	listClient.Reg(`dataTable`, func() listClient.Client {
		return New()
	})
}

func New() listClient.Client {
	return &DataTable{
		BaseClient: listClient.New(),
	}
}

type Sort struct {
	FieldInfo *reflectx.FieldInfo
	Sort      string
}

type Item struct {
	TableField string
	Keywords   string
}

type Sorts []*Sort

func (a Sorts) Each(f func(*reflectx.FieldInfo, string)) {
	for _, v := range a {
		if v != nil {
			f(v.FieldInfo, v.Sort)
		}
	}
}

func (a *Sorts) Insert(index int, fieldInfo *reflectx.FieldInfo, sort string) {
	length := len(*a)
	if length > index {
		(*a)[index] = &Sort{FieldInfo: fieldInfo, Sort: sort}
	} else if index <= 10 {
		for i := length; i <= index; i++ {
			if i == index {
				*a = append(*a, &Sort{FieldInfo: fieldInfo, Sort: sort})
			} else {
				*a = append(*a, nil)
			}
		}
	}
}

func (a Sorts) Sorts(args ...func(*reflectx.FieldInfo, string) string) (r []interface{}) {
	var fn func(*reflectx.FieldInfo, string) string
	if len(args) > 0 {
		fn = args[0]
	} else {
		fn = func(fi *reflectx.FieldInfo, sort string) string {
			if len(fi.Name) == 0 {
				return ``
			}
			if sort == `asc` {
				return fi.Name
			}
			return `-` + fi.Name
		}
	}
	a.Each(func(fi *reflectx.FieldInfo, sort string) {
		v := fn(fi, sort)
		if len(v) == 0 {
			return
		}
		r = append(r, v)
	})
	return
}

type DataTable struct {
	*listClient.BaseClient
	tableFields []string                        //字段搜索框查询的字段
	orders      Sorts                           //字段和排序方式
	search      string                          //搜索关键字
	searches    []*Item                         //搜索某个字段
	pkName      string                          //primary key name
	fieldsInfo  map[string]*reflectx.FindResult //键为表字段名
}

type fieldInfo struct {
	Index string
	Sort  string
}

func (a *DataTable) Init(c *X.Context, ls *common.List) listClient.Client {
	a.BaseClient.Init(c, ls)
	a.tableFields = make([]string, 0)
	a.searches = []*Item{}
	a.orders = Sorts{}
	a.fieldsInfo = make(map[string]*reflectx.FindResult)
	a.search = c.Form(`search[value]`)

	fm := []string{`columns[`, `][data]`} //strings.Split(`columns[0][data]`, `0`)
	typeMap := sqlbuilder.Mapper().StructMap(a.List.Model())
	/*
		for i, f := range typeMap.Tree.Children {
			fmt.Println(i, `struct============>`, f.Field.Name)
		}
		for name := range typeMap.Names {
			fmt.Println(`name============>`, name)
		}
		// */
	// ==========================
	// 获取客户端提交的字段名
	// ==========================
	fieldPaths := map[string]map[string]interface{}{}
	//com.Dump(c.Request().Form().All())

	var parseParent = func(field string, info *fieldInfo) {
		parse := strings.Split(field, `.`)
		var parent string
		if len(parse) >= 2 {
			parent = parse[0]
			field = parse[1]
		} else {
			parent = field
			field = ``
		}
		if _, exists := fieldPaths[parent]; !exists {
			fieldPaths[parent] = map[string]interface{}{}
		}
		old, exists := fieldPaths[parent][field]
		if !exists {
			fieldPaths[parent][field] = info
			return
		}
		history := old.(*fieldInfo)
		if len(info.Sort) == 0 && len(history.Sort) > 0 {
			info.Sort = history.Sort
		}
		fieldPaths[parent][field] = info
	}

	for k := range c.Request().Form().All() {
		if !strings.HasPrefix(k, fm[0]) || !strings.HasSuffix(k, fm[1]) {
			continue
		}
		idx := strings.TrimSuffix(k, fm[1])
		idx = strings.TrimPrefix(idx, fm[0])

		//要查询的所有字段
		field := c.Form(k)
		parseParent(field, &fieldInfo{Index: idx, Sort: ``})
		//要排序的字段
		fidx := c.Form(`order[` + idx + `][column]`)
		if len(fidx) == 0 {
			continue
		}
		field = c.Form(fm[0] + fidx + fm[1])
		if len(field) == 0 {
			continue
		}
		fieldInfo := &fieldInfo{
			Index: fidx,
			Sort:  c.Form(`order[` + idx + `][dir]`),
		}
		if fieldInfo.Sort != `asc` {
			fieldInfo.Sort = `desc`
		}
		parseParent(field, fieldInfo)
	}
	var pk []string
	a.fieldsInfo, pk = typeMap.FindTableFieldByMap(fieldPaths, true)
	if len(pk) > 0 {
		a.pkName = pk[0]
	}
	for fieldPath, findResult := range a.fieldsInfo {
		info := findResult.RawData.(*fieldInfo)
		a.tableFields = append(a.tableFields, fieldPath)
		//搜索本字段
		kw := c.Form(`columns[` + info.Index + `][search][value]`)
		if len(kw) > 0 {
			a.searches = append(a.searches, &Item{TableField: fieldPath, Keywords: kw})
		}
		if len(info.Sort) > 0 {
			a.orders.Insert(com.Int(info.Index), findResult.FieldInfo, info.Sort)
		}
	}
	//com.Dump(a.searches)
	//a.Form(`search[regex]`)=="false"
	//columns[0][search][regex]=false / columns[0][search][value]
	a.Middleware = func(r db.Result) db.Result {
		return r.Where(a.Build()).OrderBy(a.Sorts()...)
	}
	return a
}

// Sorts 获取排序字段
func (a *DataTable) Sorts(args ...func(*reflectx.FieldInfo, string) string) []interface{} {
	return a.orders.Sorts(args...)
}

// Build 生成搜索条件
func (a *DataTable) Build(defaultFields ...string) *db.Compounds {
	condition := db.NewCompounds()
	table := func(fr *reflectx.FindResult) string {
		fi := fr.Parent(0)
		if fi == nil {
			return ``
		}
		structName := fi.Field.Name
		return factory.DBIGet().TableName(structName)
	}
	build := func(field string, keywords string, idFields ...string) *db.Compounds {
		cond := db.NewCompounds()
		fr, ok := a.fieldsInfo[field]
		if !ok {
			return cond
		}
		tableName := table(fr)
		fi := fr.FieldInfo
		//fmt.Printf("%s => %v\n", fi.Field.Name, fi.Field.Type.Kind())
		fieldInfo, exists := factory.FieldFind(tableName, fi.Name)
		if !exists {
			return nil
		}
		switch fi.Field.Type.Kind() {
		case reflect.String:
			switch fieldInfo.DataType {
			case `enum`, `set`, `char`:
				return cond.Add(mysql.EqField(field, keywords))
			default:
				return mysql.SearchField(field, keywords)
			}
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			if strings.Contains(keywords, ` - `) {
				return mysql.GenDateRange(field, keywords)
			}
			if database.IsCompareField(keywords) {
				return cond.Add(mysql.CompareField(field, keywords))
			}
			return mysql.RangeField(field, keywords)
		}
		return cond
	}
	for _, item := range a.searches {
		condition.From(build(item.TableField, item.Keywords))
	}
	if condition.Size() > 0 {
		return condition
	}
	if len(a.search) > 0 && len(defaultFields) > 0 {
		if a.SearchPK {
			condition.From(mysql.SearchFields(defaultFields, a.search, a.pkName))
		} else {
			condition.From(mysql.SearchFields(defaultFields, a.search))
		}
		return condition
	}
	return condition
}
