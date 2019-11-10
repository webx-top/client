package ueditor

// DefaultConfig ueditor配置
var DefaultConfig = &Config{
	// 上传图片配置项
	ImageActionName:     "uploadimage",                                     /* 执行上传图片的action名称 */
	ImageFieldName:      "upfile",                                          /* 提交的图片表单名称 */
	ImageMaxSize:        2048000,                                           /* 上传大小限制，单位B */
	ImageAllowFiles:     []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 上传图片格式显示 */
	ImageCompressEnable: true,                                              /* 是否压缩图片,默认是true */
	ImageCompressBorder: 1600,                                              /* 图片压缩最长边限制 */
	ImageInsertAlign:    "none",                                            /* 插入的图片浮动方式 */
	ImageURLPrefix:      "",                                                /* 图片访问路径前缀 */

	/* 涂鸦图片上传配置项 */
	ScrawlActionName:  "uploadscrawl", /* 执行上传涂鸦的action名称 */
	ScrawlFieldName:   "upfile",       /* 提交的图片表单名称 */
	ScrawlMaxSize:     2048000,        /* 上传大小限制，单位B */
	ScrawlURLPrefix:   "",             /* 图片访问路径前缀 */
	ScrawlInsertAlign: "none",

	/* 截图工具上传 */
	SnapscreenActionName:  "uploadimage", /* 执行上传截图的action名称 */
	SnapscreenURLPrefix:   "",            /* 图片访问路径前缀 */
	SnapscreenInsertAlign: "none",        /* 插入的图片浮动方式 */

	/* 抓取远程图片配置 */
	CatcherLocalDomain: []string{"127.0.0.1", "localhost", "img.baidu.com"},
	CatcherActionName:  "catchimage",                                      /* 执行抓取远程图片的action名称 */
	CatcherFieldName:   "source",                                          /* 提交的图片列表表单名称 */
	CatcherURLPrefix:   "",                                                /* 图片访问路径前缀 */
	CatcherMaxSize:     2048000,                                           /* 上传大小限制，单位B */
	CatcherAllowFiles:  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 抓取图片格式显示 */

	/* 上传视频配置 */
	VideoActionName: "uploadvideo", /* 执行上传视频的action名称 */
	VideoFieldName:  "upfile",      /* 提交的视频表单名称 */
	VideoURLPrefix:  "",            /* 视频访问路径前缀 */
	VideoMaxSize:    102400000,     /* 上传大小限制，单位B，默认100MB */
	VideoAllowFiles: []string{
		".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
		".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid"},

	/* 上传文件配置 */
	FileActionName: "uploadfile", /* controller里,执行上传视频的action名称 */
	FileFieldName:  "upfile",     /* 提交的文件表单名称 */
	FileURLPrefix:  "",           /* 文件访问路径前缀 */
	FileMaxSize:    51200000,     /* 上传大小限制，单位B，默认50MB */
	FileAllowFiles: []string{
		".png", ".jpg", ".jpeg", ".gif", ".bmp",
		".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
		".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
		".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
		".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
	},

	/* 列出指定目录下的图片 */
	ImageManagerActionName:  "listimage",                                       /* 执行图片管理的action名称 */
	ImageManagerListSize:    20,                                                /* 每次列出文件数量 */
	ImageManagerURLPrefix:   "",                                                /* 图片访问路径前缀 */
	ImageManagerInsertAlign: "none",                                            /* 插入的图片浮动方式 */
	ImageManagerAllowFiles:  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 列出的文件类型 */

	/* 列出指定目录下的文件 */
	FileManagerActionName: "listfile", /* 执行文件管理的action名称 */
	FileManagerURLPrefix:  "",         /* 文件访问路径前缀 */
	FileManagerListSize:   20,         /* 每次列出文件数量 */
	FileManagerAllowFiles: []string{".png", ".jpg", ".jpeg", ".gif", ".bmp",
		".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg",
		".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid",
		".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso",
		".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml",
	},
}

// Config ueditor配置
type Config struct {
	// 上传图片配置项
	ImageActionName     string   `json:"imageActionName"`     // 执行上传图片的action名称
	ImageFieldName      string   `json:"imageFieldName"`      // 提交的图片表单名称
	ImageMaxSize        int64    `json:"imageMaxSize"`        // 上传大小限制，单位B
	ImageAllowFiles     []string `json:"imageAllowFiles"`     // 上传图片格式显示
	ImageCompressEnable bool     `json:"imageCompressEnable"` // 是否压缩图片,默认是true
	ImageCompressBorder int64    `json:"imageCompressBorder"` // 图片压缩最长边限制
	ImageInsertAlign    string   `json:"imageInsertAlign"`    // 插入的图片浮动方式
	ImageURLPrefix      string   `json:"imageUrlPrefix"`      // 图片访问路径前缀

	/* 涂鸦图片上传配置项 */
	ScrawlActionName  string `json:"scrawlActionName"` // 执行上传涂鸦的action名称 */
	ScrawlFieldName   string `json:"scrawlFieldName"`  // 提交的图片表单名称 */
	ScrawlMaxSize     int64  `json:"scrawlMaxSize"`    // 上传大小限制，单位B
	ScrawlURLPrefix   string `json:"scrawlUrlPrefix"`  // 图片访问路径前缀
	ScrawlInsertAlign string `json:"scrawlInsertAlign"`

	/* 截图工具上传 */
	SnapscreenActionName  string `json:"snapscreenActionName"`  // 执行上传截图的action名称
	SnapscreenURLPrefix   string `json:"snapscreenUrlPrefix"`   // 图片访问路径前缀
	SnapscreenInsertAlign string `json:"snapscreenInsertAlign"` // 插入的图片浮动方式

	/* 抓取远程图片配置 */
	CatcherLocalDomain []string `json:"catcherLocalDomain"` //: ["127.0.0.1", "localhost", "img.baidu.com"],
	CatcherActionName  string   `json:"catcherActionName"`  // 执行抓取远程图片的action名称
	CatcherFieldName   string   `json:"catcherFieldName"`   // 提交的图片列表表单名称
	CatcherURLPrefix   string   `json:"catcherUrlPrefix"`   // 图片访问路径前缀
	CatcherMaxSize     int64    `json:"catcherMaxSize"`     // 上传大小限制，单位B
	CatcherAllowFiles  []string `json:"catcherAllowFiles"`  //: [".png", ".jpg", ".jpeg", ".gif", ".bmp"], // 抓取图片格式显示

	/* 上传视频配置 */
	VideoActionName string   `json:"videoActionName"` // 执行上传视频的action名称
	VideoFieldName  string   `json:"videoFieldName"`  // 提交的视频表单名称
	VideoURLPrefix  string   `json:"videoUrlPrefix"`  // 视频访问路径前缀
	VideoMaxSize    int64    `json:"videoMaxSize"`    // 上传大小限制，单位B，默认100MB
	VideoAllowFiles []string `json:"videoAllowFiles"` // 上传视频格式显示

	/* 上传文件配置 */
	FileActionName string   `json:"fileActionName"` // controller里,执行上传视频的action名称
	FileFieldName  string   `json:"fileFieldName"`  // 提交的文件表单名称
	FileURLPrefix  string   `json:"fileUrlPrefix"`  // 文件访问路径前缀
	FileMaxSize    int64    `json:"fileMaxSize"`    // 上传大小限制，单位B，默认50MB
	FileAllowFiles []string `json:"fileAllowFiles"` // 上传文件格式显示

	/* 列出指定目录下的图片 */
	ImageManagerActionName  string   `json:"imageManagerActionName"`  // 执行图片管理的action名称
	ImageManagerListSize    int64    `json:"imageManagerListSize"`    // 每次列出文件数量
	ImageManagerURLPrefix   string   `json:"imageManagerUrlPrefix"`   // 图片访问路径前缀
	ImageManagerInsertAlign string   `json:"imageManagerInsertAlign"` // 插入的图片浮动方式
	ImageManagerAllowFiles  []string `json:"imageManagerAllowFiles"`  // 列出的文件类型

	/* 列出指定目录下的文件 */
	FileManagerActionName string   `json:"fileManagerActionName"` // 执行文件管理的action名称
	FileManagerURLPrefix  string   `json:"fileManagerUrlPrefix"`  // 文件访问路径前缀
	FileManagerListSize   int64    `json:"fileManagerListSize"`   // 每次列出文件数量
	FileManagerAllowFiles []string `json:"fileManagerAllowFiles"`
}
