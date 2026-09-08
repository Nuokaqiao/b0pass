package app

import (
	"b0go/apps/pass/lib/audit"
	"b0go/apps/pass/lib/chat"
	"b0go/apps/pass/lib/files"
	"b0go/apps/pass/lib/keys"
	"b0go/apps/pass/lib/stream"
	"b0go/core/engine"
	"b0go/core/tools/cmd"
	"b0go/core/tools/nets"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Ping 主电脑连通性测试
func Ping(c *gin.Context) {
	engine.OK("OK", true, c)
}

// AuthStatus 是否启用登录鉴权（公开）
func AuthStatus(c *gin.Context) {
	engine.OK("OK", gin.H{"enabled": engine.AuthEnabled()}, c)
}

// Login 共享口令登录，返回 JWT（公开）
func Login(c *gin.Context) {
	if !engine.AuthEnabled() {
		engine.OK("鉴权未启用", gin.H{"token": "", "enabled": false}, c)
		return
	}
	password := strings.TrimSpace(c.PostForm("password"))
	if password == "" {
		password = strings.TrimSpace(c.Query("password"))
	}
	if password == "" {
		var body struct {
			Password string `json:"password"`
		}
		_ = c.ShouldBindJSON(&body)
		password = strings.TrimSpace(body.Password)
	}
	if !engine.CheckPassword(password) {
		audit.Log(c, "login-fail", "")
		engine.JSON(401, "口令错误", nil, c)
		return
	}
	token, err := engine.CreateToken("share")
	if err != nil {
		engine.ERR("签发失败", c)
		return
	}
	audit.Log(c, "login-ok", "")
	engine.OK("OK", gin.H{"token": token, "enabled": true, "expireHours": 24}, c)
}

// TextHistory 传内容最近历史（最多 10 条且 72 小时内）
func TextHistory(c *gin.Context) {
	if textHub == nil || textHub.History() == nil {
		engine.OK("OK", gin.H{"items": []chat.HistoryItem{}, "backend": "none"}, c)
		return
	}
	store := textHub.History()
	list, err := store.List()
	if err != nil {
		engine.ERR(err.Error(), c)
		return
	}
	if list == nil {
		list = []chat.HistoryItem{}
	}
	engine.OK("OK", gin.H{"items": list, "backend": store.Backend()}, c)
}

// ReadConfig 读取配置
func ReadConfig(c *gin.Context) {
	engine.OK("OK", config, c)
}

// ReadIP 读取配置
func ReadIP(c *gin.Context) {
	ip := nets.GetOutBoundIP()
	engine.OK("OK", ip, c)
}

// CmdOpen 使用命令行打开
func CmdOpen(c *gin.Context) {
	RootPath := strings.ReplaceAll(config.Path, "\\", "/")
	f := c.Query("f")
	f = strings.ReplaceAll(f, "\\", "/")
	ext := strings.ToUpper(path.Ext(f))
	engine.Println(ext)
	if ext == ".BAT" || ext == ".CMD" || ext == ".EXE" {
		engine.ERR("该文件暂不支持打开", c)
	} else {
		audit.Log(c, "cmd-open", "path="+f)
		cmd.Open(RootPath + f)
	}
}

// CmdKey 主电脑键盘
func CmdKey(c *gin.Context) {
	k := c.Query("k")
	audit.Log(c, "cmd-key", "k="+k)
	keys.SendKey(k)
	engine.OK("OK", nil, c)
}

// NodeTree 目录树结构
func NodeTree(c *gin.Context) {
	f := c.Query("f")
	listMap := files.NodeTree(config.Path, f)
	engine.OK("OK", listMap, c)
}

// NodeRename 重命名目录
func NodeRename(c *gin.Context) {
	RootPath := config.Path
	f := c.Query("f")
	n := c.Query("n")
	if f == "" || n == "" {
		engine.ERR("路径不能为空", c)
		return
	}
	err := files.NodeRename(RootPath+f, RootPath+n)
	if err != nil {
		engine.ERR(err.Error(), c)
		return
	}
	audit.Log(c, "node-rename", "from="+f+" to="+n)
	engine.OK("OK", nil, c)
}

// NodeRemove 删除目录
func NodeRemove(c *gin.Context) {
	RootPath := config.Path
	f := c.Query("f")
	if f == "" {
		engine.ERR("路径不能为空", c)
		return
	}
	// 统一成相对路径，避免出现 "files"+"/x" 之外的意外拼接
	f = strings.ReplaceAll(f, "\\", "/")
	for strings.HasPrefix(f, "/") {
		f = strings.TrimPrefix(f, "/")
	}
	if f == "" || strings.Contains(f, "..") {
		engine.ERR("非法路径", c)
		return
	}
	filePath := filepath.Join(RootPath, f)
	log.Println("::NodeRemove::", filePath)
	err := files.NodeRemove(filePath)
	if err != nil {
		engine.ERR(err.Error(), c)
		return
	}
	files.GetExpireStore(config.Path).DeleteMeta("/" + f)
	audit.Log(c, "node-delete", "path=/"+f)
	engine.OK("OK", nil, c)
}

// NodeAdd 添加目录
func NodeAdd(c *gin.Context) {
	RootPath := config.Path
	f := c.Query("f")
	if f == "" {
		engine.ERR("路径不能为空:(f = 结尾带“/”为创建目录,否则为创建文件)", c)
		return
	}
	filePath := RootPath + f
	filePath = strings.ReplaceAll(filePath, "//", "/")
	log.Println("::NodeAdd::", filePath)
	err := files.NodeAdd(filePath)
	if err != nil {
		engine.ERR(err.Error(), c)
	}
	audit.Log(c, "node-add", "path="+f)
	engine.OK("OK", nil, c)
}

func FileCount(c *gin.Context) {
	RootPath := config.Path
	listMap := files.GetCounts(RootPath)
	engine.OK("OK", listMap, c)
}

// FileList 文件列表
func FileList(c *gin.Context) {
	RootPath := config.Path
	f := c.Query("f")
	t := c.DefaultQuery("t", "")
	// 列表前顺手清理一次过期文件
	files.GetExpireStore(RootPath).PurgeExpired()
	listMap := files.GetDirTree(RootPath, RootPath+f, "", t)
	store := files.GetExpireStore(RootPath)
	now := time.Now().Unix()
	for _, m := range listMap {
		p, _ := m["path"].(string)
		exp := store.Get(p)
		m["expire"] = exp
		if exp > 0 {
			m["expireAt"] = time.Unix(exp, 0).Format("01-02 15:04")
			left := exp - now
			if left < 0 {
				left = 0
			}
			m["expireLeft"] = left
		} else {
			m["expireAt"] = ""
			m["expireLeft"] = int64(0)
		}
	}
	engine.OK("OK", listMap, c)
}

// FileContent 文件内容
func FileContent(c *gin.Context) {
	RootPath := config.Path
	f := c.Query("f") //!strings.HasPrefix(f, RootPath)
	if f == "" {
		engine.ERR("路径不能为空", c)
		return
	}
	RootPath = path.Clean(RootPath + f)
	data, err := files.GetData(RootPath)
	if err != nil {
		engine.ERR(err.Error(), c)
		return
	}
	engine.OK("OK", string(data), c)
}

// FileDownload 文件下载
func FileDownload(c *gin.Context) {
	RootPath := config.Path
	f := c.Query("f")
	if f == "" {
		engine.ERR("路径不能为空", c)
		return
	}
	filePath := RootPath + f
	//获取文件的名称
	fileName := path.Base(filePath)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "inline;filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	audit.Log(c, "file-download", "path="+f)
	c.File(filePath)
}

// FileUpload 上传文件
func FileUpload(c *gin.Context) {
	lens := 0
	if vals := c.Request.Header["Content-Length"]; len(vals) > 0 {
		lens, _ = strconv.Atoi(vals[0])
	}
	log.Println("FileUpload::::", lens)
	audit.Log(c, "file-upload-start", fmt.Sprintf("bytes=%d dir=%s", lens, c.DefaultQuery("f", "/")))
	if lens > 4096 {
		FileUploadBig(c)
	} else {
		FileUploadTiny(c)
	}
}

// FileUploadTiny 小文件上传
func FileUploadTiny(c *gin.Context) {
	RootPath := config.Path
	f := c.DefaultQuery("f", "/")
	RootPath = RootPath + f
	if err := files.NodeAdd(RootPath); err != nil {
		engine.ERR("创建目录失败: "+err.Error(), c)
		return
	}
	log.Println("FileUploadTiny:::", RootPath)
	file, err := c.FormFile("file")
	if err != nil {
		engine.ERR("请先选择文件: "+err.Error(), c)
		return
	}
	if err := c.SaveUploadedFile(file, RootPath+file.Filename); err != nil {
		engine.ERR("保存文件失败: "+err.Error(), c)
		return
	}
	applyUploadExpire(c, f, file.Filename)
	audit.Log(c, "file-upload-ok", "file="+file.Filename+" dir="+f)
	engine.OK("上传成功", "", c)
}

// FileUploadBig 大文件上传
func FileUploadBig(c *gin.Context) {

	RootPath := config.Path
	dirQuery := c.DefaultQuery("f", "/")
	RootPath = RootPath + dirQuery
	if err := files.NodeAdd(RootPath); err != nil {
		engine.ERR("创建目录失败: "+err.Error(), c)
		return
	}
	log.Println("FileUploadBig:::", RootPath)

	content_type_, has_key := c.Request.Header["Content-Type"]
	if !has_key {
		engine.ERR("请先选择文件", c)
		return
	}
	content_type := content_type_[0]
	const BOUNDARY string = "; boundary="
	loc := strings.Index(content_type, BOUNDARY)
	if len(content_type_) != 1 || loc == -1 {
		engine.ERR("请先选择文件", c)
		return
	}
	boundary := []byte(content_type[(loc + len(BOUNDARY)):])
	log.Printf("file boundary: [%s]\n", boundary)
	read_data := make([]byte, 1024*12)
	var read_total int = 0
	for {
		//解析文件头信息
		file_header, file_data, err := stream.ParseFromHead(
			read_data, read_total,
			append(boundary, []byte("\r\n")...), c.Request.Body,
		)
		if err != nil {
			engine.ERR("parse from fail: "+err.Error(), c)
			return
		}
		//创建保存文件
		log.Printf("save file: [%s]\n", RootPath+file_header.FileName)
		out, err := os.Create(RootPath + file_header.FileName)
		if err != nil {
			engine.ERR("create file fail: "+err.Error(), c)
			return
		}
		out.Write(file_data)
		file_data = nil

		//搜索boundary
		temp_data, reach_end, err := stream.ReadToBoundary(boundary, c.Request.Body, out)
		out.Close()
		if err != nil {
			engine.ERR("search boundary fail: "+err.Error(), c)
			return
		}
		applyUploadExpire(c, dirQuery, file_header.FileName)
		if reach_end {
			break
		} else {
			copy(read_data[0:], temp_data)
			read_total = len(temp_data)
			continue
		}
	}
	//上传成功
	audit.Log(c, "file-upload-ok", "dir="+dirQuery+" mode=big")
	engine.OK("上传成功", "", c)
}

// FileExpire 设置/清除文件过期时间；expire=0 表示不过期
func FileExpire(c *gin.Context) {
	f := c.Query("f")
	if f == "" {
		engine.ERR("路径不能为空", c)
		return
	}
	rel := files.NormalizeRelPath(f)
	if rel == "/" || strings.Contains(rel, "..") {
		engine.ERR("非法路径", c)
		return
	}
	abs := filepath.Join(config.Path, filepath.FromSlash(strings.TrimPrefix(rel, "/")))
	if _, err := os.Stat(abs); err != nil {
		engine.ERR("文件不存在", c)
		return
	}
	expire, _ := strconv.ParseInt(c.DefaultQuery("expire", "0"), 10, 64)
	if expire < 0 {
		expire = 0
	}
	if expire > 0 && expire < time.Now().Unix() {
		engine.ERR("过期时间不能早于现在", c)
		return
	}
	if err := files.GetExpireStore(config.Path).Set(rel, expire); err != nil {
		engine.ERR("保存过期设置失败: "+err.Error(), c)
		return
	}
	audit.Log(c, "file-expire", fmt.Sprintf("path=%s expire=%d", rel, expire))
	engine.OK("OK", gin.H{"path": rel, "expire": expire}, c)
}

func applyUploadExpire(c *gin.Context, dirQuery, filename string) {
	expire, _ := strconv.ParseInt(c.DefaultQuery("expire", "0"), 10, 64)
	if expire <= 0 {
		return
	}
	rel := files.NormalizeRelPath(path.Join(dirQuery, filename))
	if err := files.GetExpireStore(config.Path).Set(rel, expire); err != nil {
		log.Println("apply upload expire fail:", rel, err)
	}
}
