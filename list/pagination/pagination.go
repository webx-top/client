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
package pagination

import (
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	X "github.com/webx-top/webx"
	formcommon "github.com/webx-top/webx/lib/forms/common"
)

type Paginationer interface {
	//设置导航条模板
	SetTmpl(tmpl string) Paginationer

	//设置从URL中获取的参数名：页码参数名，每页数据量参数名，数据总量参数名
	SetQueryName(args ...string) Paginationer

	//设置各项参数
	SetPage(args ...int64) Paginationer

	//准备渲染所必须的数据
	Ready(linkNum int64, args ...string) Paginationer

	//ajax导航时更新的容器ID
	SetContainerId(containerId string) Paginationer

	SetRowsOptions(rows ...int64) Paginationer

	//唯一ID。用于在相同页面内区分多个导航条
	UniqueId() string

	//设置唯一ID
	SetUniqueId(uniqueId string) Paginationer

	//渲染输出导航条HTML对象
	Render() template.HTML

	//渲染输出导航条字符串
	String() string

	//组合数据并渲染输出导航条HTML对象
	HTML(linkNum int64, args ...string) template.HTML

	//数据库查询所需的游标
	Limit() (limit, offset int64)

	//设置计算数据总量的操作函数
	SetCount(fn func() int64) Paginationer

	//登记查询出来的数据结果，并检查是否为空
	CheckData(data interface{}) Paginationer

	//获取登记的数据结果
	Data() interface{}
	Pagebar() *Pagebar

	//URL中的“页码”变量标识
	PageKey() string

	//URL中的“每页数据量”变量标识
	PageRowsKey() string

	//URL中的“数据总量”变量标识
	TotalRowsKey() string

	CleanURL(string, ...string) string
}

type Pages struct {
	Data interface{}
	Page Paginationer
}

func NewPagebar() *Pagebar {
	return &Pagebar{}
}

type Pagebar struct {
	Url         string  //网址
	Rows        int64   //每页显示行数
	RowsOptions []int64 //可切换的每页显示行数
	PageRowsUrl string  //页面行数切换网址
	Page        int64   //页码
	Total       int64   //页码总数
	NextPage    int64   //下一页页码
	PrevPage    int64   //上一页页码
	ContainerId string
	UniqueId    string //唯一ID
	Links       int64  //页码链接数量
}

func (p *Pagebar) Pages(args ...int64) []int64 {
	n := p.Links
	if len(args) > 0 {
		n = args[0]
	}

	half := int64(n / 2)
	start := p.Page - half

	half = n - half
	end := p.Page + half

	if start < 1 {
		end = end - start
		start = 1
	} else {
		end--
	}
	if end > p.Total {
		remain := end - p.Total
		end = p.Total
		if start > remain {
			start -= remain
			if start < 1 {
				start = 1
			}
		}
	}
	r := make([]int64, 0)
	for i := start; i <= end; i++ {
		r = append(r, i)
	}
	return r
}

type Pagination struct {
	*X.Context   `json:"-" xml:"-"`
	pageNum      int64
	pageRows     int64
	maxPageRows  int64
	TotalRows    int64
	countFunc    func() int64 //统计TotalRows的函数
	pageKey      string
	pageRowsKey  string
	totalRowsKey string
	Tmpl         string
	pagebar      *Pagebar
	emptyData    bool
	data         interface{}
	paramsLoaded bool
	rowsOptions  []int64 //可切换的每页显示行数
}

func New(ctx *X.Context) Paginationer {
	b := &Pagination{
		Context:      ctx,
		TotalRows:    -1,
		pageKey:      `page`,
		pageRowsKey:  `pagerows`,
		totalRowsKey: `totalrows`,
		maxPageRows:  1000,
		pageRows:     10,
	}
	//b.SetTmpl(`default`).SetPage(args...)
	return b
}

func (this *Pagination) SetTmpl(tmpl string) Paginationer {
	if tmpl == `` {
		return this
	}
	switch tmpl[0] {
	case '#':
		this.Tmpl = filepath.Join(this.Context.Server.RootDir(), `data`, `widgets`, tmpl[1:]+`.html`)
	case '@':
		tmpl = tmpl[1:]
		var tdir string
		if pos := strings.Index(tmpl, `:`); pos > 0 {
			module := this.Context.Server.Module(tmpl[0:pos])
			tmpl = tmpl[pos+1:]
			tdir = module.Renderer.TmplDir()
		} else {
			tdir = this.Context.Module.Renderer.TmplDir()
		}
		this.Tmpl = filepath.Join(tdir, tmpl+`.html`)
	default:
		this.Tmpl = filepath.Join(this.Context.Module.Renderer.TmplDir(), tmpl+`.html`)
	}
	return this
}

//SetQueryName 1-页码变量名 2-每页数据量变量名 3-总数据量变量名
func (this *Pagination) SetQueryName(args ...string) Paginationer {
	switch len(args) {
	case 3:
		if len(args[2]) > 0 {
			this.totalRowsKey = args[2]
		}
		fallthrough
	case 2:
		if len(args[1]) > 0 {
			this.pageRowsKey = args[1]
		}
		fallthrough
	case 1:
		if len(args[0]) > 0 {
			this.pageKey = args[0]
		}
	}
	return this
}

func (this *Pagination) QueryValue() Paginationer {
	if v := this.Context.Form(this.pageKey); len(v) > 0 {
		if num, err := strconv.Atoi(v); err == nil && num > 1 {
			this.pageNum = int64(num)
		}
	}
	if v := this.Context.Form(this.pageRowsKey); len(v) > 0 {
		if num, err := strconv.Atoi(v); err == nil && num > 0 && int64(num) <= this.maxPageRows {
			this.pageRows = int64(num)
		}
	}
	if v := this.Context.Form(this.totalRowsKey); len(v) > 0 {
		if num, err := strconv.Atoi(v); err == nil && num > 0 {
			this.TotalRows = int64(num)
		}
	}
	return this
}

//SetPage 1-总数据量 2-每页数据量 3-页码
func (this *Pagination) SetPage(args ...int64) Paginationer {
	switch len(args) {
	case 4:
		this.maxPageRows = args[3]
		fallthrough
	case 3:
		this.pageNum = args[2]
		fallthrough
	case 2:
		if args[1] > 0 && args[1] <= this.maxPageRows {
			this.pageRows = args[1]
		}
		fallthrough
	case 1:
		this.TotalRows = args[0]
	default:
		this.QueryValue()
	}
	if this.pageNum < 1 {
		this.pageNum = 1
	}
	this.paramsLoaded = true
	return this
}

func (this *Pagination) Pages(total, limit int64) int64 {
	if total <= 0 {
		return 1
	}
	x := total % limit
	if x > 0 {
		return total/limit + 1
	}
	return total / limit
}

func (this *Pagination) SetRowsOptions(rows ...int64) Paginationer {
	this.rowsOptions = rows
	return this
}

func (this *Pagination) Ready(linkNum int64, args ...string) Paginationer {
	if !this.paramsLoaded {
		this.SetPage()
	}
	pageRows := this.pageRows
	totalRows := this.TotalRows
	if totalRows < 0 {
		if this.countFunc != nil {
			totalRows = this.countFunc()
		} else {
			totalRows = 0
		}
		this.TotalRows = totalRows
	}

	var urlFormat string
	switch len(args) {
	case 2:
		this.SetTmpl(args[1])
		fallthrough
	case 1:
		if len(args[0]) == 0 {
			urlFormat = this.ReqURI()
		} else {
			urlFormat = args[0]
		}
	default:
		urlFormat = this.ReqURI()
	}
	page := this.pageNum

	if this.emptyData && totalRows > 0 && page > 1 {
		this.Context.Go(urlFormat)
		this.pagebar = nil
		return this
	}
	totalPages := this.Pages(totalRows, pageRows)
	if page < 1 {
		page = 1
	}
	if page > totalPages && page > 1 {
		this.Context.Go(urlFormat)
		this.pagebar = nil
		return this
	}
	separator := `&`
	if strings.Contains(urlFormat, `?`) == false {
		urlFormat += `?`
	} else {
		urlFormat = strings.Replace(urlFormat, `%`, `%%`, -1)
		urlFormat += separator
	}

	urlFormat += this.pageKey + `=%v`
	if totalRows > 0 && len(this.totalRowsKey) > 0 {
		urlFormat += separator + this.totalRowsKey + `=` + fmt.Sprintf(`%d`, totalRows)
	}
	pg := NewPagebar()
	pg.Rows = this.pageRows
	pg.PageRowsUrl = fmt.Sprintf(urlFormat, 1) + separator + this.pageRowsKey + `=`
	pg.Page = page
	if len(this.rowsOptions) == 0 {
		pg.RowsOptions = []int64{10, 20, 50, 100, 200, 300, 500, 1000}
	} else {
		pg.RowsOptions = this.rowsOptions
	}
	pg.PrevPage = page - 1
	pg.NextPage = page + 1
	pg.Total = totalPages
	if pageRows > 0 && len(this.pageRowsKey) > 0 {
		urlFormat += separator + this.pageRowsKey + `=` + fmt.Sprintf(`%d`, this.pageRows)
	}
	pg.Url = urlFormat
	pg.Links = linkNum
	this.pagebar = pg
	return this
}

func (this *Pagination) SetContainerId(containerId string) Paginationer {
	if this.pagebar != nil {
		this.pagebar.ContainerId = containerId
		this.pagebar.UniqueId = fmt.Sprintf(`%x`, containerId)
	}
	return this
}

func (this *Pagination) UniqueId() (r string) {
	if this.pagebar != nil {
		r = this.pagebar.UniqueId
	}
	return
}

func (this *Pagination) SetUniqueId(uniqueId string) Paginationer {
	if this.pagebar != nil {
		this.pagebar.UniqueId = uniqueId
	}
	return this
}

func (this *Pagination) Render() template.HTML {
	return template.HTML(this.String())
}

func (this *Pagination) String() string {
	if this.pagebar == nil {
		return ``
	}
	funcMap := this.Context.Server.FuncMap
	contexFuncMap := this.Context.Funcs()
	if contexFuncMap != nil {
		for name, fn := range contexFuncMap {
			funcMap[name] = fn
		}
	}
	return formcommon.ParseTmpl(this.pagebar, funcMap, nil, this.Tmpl)
}

func (this *Pagination) HTML(linkNum int64, args ...string) template.HTML {
	return this.Ready(linkNum, args...).Render()
}

func (this *Pagination) Limit() (limit, offset int64) {
	pn := this.pageNum
	if pn < 1 {
		pn = 1
	}
	limit = this.pageRows
	offset = (pn - 1) * limit
	return
}

func (this *Pagination) SetCount(fn func() int64) Paginationer {
	this.countFunc = fn
	return this
}

func (this *Pagination) CheckData(data interface{}) Paginationer {
	if data == nil || fmt.Sprint(data) == `[]` {
		this.emptyData = true
	}
	this.data = data
	return this
}

func (this *Pagination) Data() interface{} {
	return this.data
}

func (this *Pagination) Pagebar() *Pagebar {
	return this.pagebar
}

func (this *Pagination) PageNum() int64 {
	return this.pageNum
}

func (this *Pagination) PageRows() int64 {
	return this.pageRows
}

func (this *Pagination) PageKey() string {
	return this.pageKey
}

func (this *Pagination) PageRowsKey() string {
	return this.pageRowsKey
}

func (this *Pagination) TotalRowsKey() string {
	return this.totalRowsKey
}

func (this *Pagination) ReqURI() (r string) {
	r = this.Context.Request().URL().Path()
	q := this.Context.Request().URL().Query()
	q.Del(this.pageKey)
	q.Del(this.pageRowsKey)
	q.Del(this.totalRowsKey)
	if len(this.Context.Query(`_pjax`)) > 0 {
		q.Del(`_pjax`)
	}
	s := q.Encode()
	if len(s) > 0 {
		r += `?` + s
	}
	return
}

func (this *Pagination) CleanURL(q string, keys ...string) (r string) {
	cr, err := regexp.Compile(`(&|^)(` + strings.Join(keys, `|`) + `)=[^&]*`)
	if err != nil {
		return err.Error()
	}
	r = cr.ReplaceAllString(q, ``)
	r = strings.Trim(r, `&`)
	return
}
