package app

import (
	"b0go/apps/pass/lib/chat"
	"b0go/apps/pass/lib/files"
	"b0go/core/engine"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
)

// APP:AppConfig
type AppConfig struct {
	Live          bool
	Path          string //文件根目录路径
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// APP:VAR
var (
	app    *engine.AppConfig
	appId  = "pass"
	config = new(AppConfig)

	uiPath string
	//go:embed ui/dist
	uiFS embed.FS

	textHub *chat.Hub
)

// APP:INIT
func init() {
	uiDist, _ := fs.Sub(uiFS, "ui/dist")
	app = &engine.AppConfig{
		Name:   appId,
		Type:   engine.APP_APP,
		Config: config,
		UIFS:   uiDist,
		Run:    run,
	}
	engine.AppInstall(app)
	uiPath = filepath.Join(app.Dir, "ui", "dist")
}

// APP:RUN
func run() {
	engine.Print(aurora.Green("App pass loaded"), aurora.BrightCyan(config))
	engine.Gin.Use(engine.CorsMiddleware())
	engine.Gin.Use(engine.PathAuthMiddleware("/files", "/ws"))
	routeStatic(config.Live)
	routeApi()
	routeWs()
	putDll()
	// 默认不过期；有设置过期的文件由后台定期清理
	files.StartExpireCleaner(config.Path, time.Minute)
}

// 注册静态路由
func routeStatic(live bool) {
	if live {
		viteProxy := createViteProxy()
		if viteProxy != nil {
			engine.Gin.Any("/app/pass", viteProxy)
			engine.Gin.Any("/app/pass/*path", viteProxy)
			engine.Gin.Any("/src/*path", viteProxy)
			engine.Gin.Any("/node_modules/*path", viteProxy)
			engine.Gin.Any("/@vite/*path", viteProxy)
			engine.Gin.Any("/@fs/*path", viteProxy)
			engine.Gin.Any("/@id/*path", viteProxy)
			engine.Print(aurora.Yellow("Vue 开发模式：代理到 Vite (http://localhost:5173)"))
			engine.Print(aurora.Yellow("请先运行: cd apps/pass/ui/vue/stonePass && npm run dev"))
		} else {
			engine.Gin.Static("/app/pass", uiPath)
			engine.Print(aurora.Yellow("Vite 未启动，回退到 ui/dist 静态文件"))
		}
	} else {
		uiDist, _ := fs.Sub(uiFS, "ui/dist")
		engine.Gin.StaticFS("/app/pass", http.FS(uiDist))
	}

	engine.Gin.GET("/", func(c *gin.Context) {
		c.Writer.Write([]byte("<script>location.href='/app/pass/'</script>"))
	})
	engine.Gin.StaticFS("/files", http.Dir(config.Path))
}

func createViteProxy() gin.HandlerFunc {
	viteURL, err := url.Parse("http://127.0.0.1:5173")
	if err != nil {
		return nil
	}
	proxy := httputil.NewSingleHostReverseProxy(viteURL)
	proxy.Director = func(req *http.Request) {
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.URL.Scheme = viteURL.Scheme
		req.URL.Host = viteURL.Host
		path := req.URL.Path
		if strings.HasPrefix(path, "/app/pass") {
			req.URL.Path = path[len("/app/pass"):]
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
	}

	client := &http.Client{Timeout: http.DefaultClient.Timeout}
	resp, err := client.Get("http://127.0.0.1:5173")
	if err != nil {
		return nil
	}
	resp.Body.Close()
	if resp.StatusCode >= 500 {
		return nil
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// 注册应用路由
func routeApi() {
	// 公开：登录与鉴权状态
	engine.GET(appId, "/auth-status", "{}", "鉴权是否启用", AuthStatus)
	engine.POST(appId, "/login", "{password}", "共享口令登录", Login)

	GETX("/ping", "{}", "连通性测试", Ping)
	GETX("/read-config", "{}", "读取配置", ReadConfig)
	GETX("/read-ip", "{}", "读取IP", ReadIP)
	GETX("/cmd-open", "{}", "命令行打开", CmdOpen)
	GETX("/cmd-key", "{}", "主电脑键盘", CmdKey)

	GETX("/node-tree", "{}", "目录树结构", NodeTree)
	GETX("/node-add", "{f=相对路径(结尾带“/”为创建目录,否则为创建文件)}", "添加目录", NodeAdd)
	GETX("/node-rename", "{f=原路径,n=新路径}", "重命名节点", NodeRename)
	GETX("/node-delete", "{f=相对路径}", "删除节点", NodeRemove)

	GETX("/file-count", "{f=相对路径}", "文件数量", FileCount)
	GETX("/file-list", "{f=相对路径,[t=需要的类型]}", "文件列表", FileList)
	GETX("/file-content", "{f=相对路径,文件名称}", "文件内容", FileContent)
	GETX("/file-download", "{f=相对路径}", "文件下载", FileDownload)
	GETX("/file-expire", "{f=相对路径,expire=unix秒(0不过期)}", "设置文件过期时间", FileExpire)
	GETX("/text-history", "{}", "传内容最近历史", TextHistory)

	POSTX("/file-upload", "{post file,[expire=unix秒]}", "大文件上传", FileUpload)
}

/***** HttpRequest *****/

// GET
func GET(url, param, title string, handle ...gin.HandlerFunc) {
	engine.Router(appId, "GET", url, param, title, handle...)
}

// POST
func POST(url, param, title string, handle ...gin.HandlerFunc) {
	engine.Router(appId, "POST", url, param, title, handle...)
}

// GETX
func GETX(url, param, title string, handle gin.HandlerFunc) {
	engine.Router(appId, "GET", url, param, "(Auth)"+title, engine.JWTMiddleware(), handle)
}

// POSTX
func POSTX(url, param, title string, handle gin.HandlerFunc) {
	engine.Router(appId, "POST", url, param, "(Auth)"+title, engine.JWTMiddleware(), handle)
}

// 注册ws路由
func routeWs() {
	store := chat.NewHistoryStore(chat.RedisOptions{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	textHub = chat.NewHub(store)
	go textHub.Run()
	engine.Gin.GET("/ws", func(c *gin.Context) {
		chat.ServeWs(textHub, c)
	})
}

// 释放dll文件
func putDll() {
	if runtime.GOOS == "windows" && !config.Live {
		_, errNow := os.ReadFile("zlib1.dll")
		if errNow != nil {
			dll, err := uiFS.ReadFile("ui/dist/dll/zlib1.dll")
			if err != nil {
				log.Println("zlib1.dll err:", err)
			} else {
				os.WriteFile("zlib1.dll", dll, 0777)
			}
		}
	}
}
