package engine

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

/****** Cors MiddleWare *******/

// 处理跨域请求,支持options访问
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

/****** Token MiddleWare *******/

var (
	TokenExpire = time.Hour * 24
	TokenSecret = []byte("01xda2d8f6x9n4x8")

	authMu       sync.RWMutex
	authPassword string
)

type MyClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// SetAuthPassword 由 gateway 在加载配置后注入共享口令；空字符串表示关闭鉴权
func SetAuthPassword(password string) {
	authMu.Lock()
	defer authMu.Unlock()
	authPassword = strings.TrimSpace(password)
}

// AuthEnabled 是否启用登录鉴权
func AuthEnabled() bool {
	authMu.RLock()
	defer authMu.RUnlock()
	return authPassword != ""
}

// CheckPassword 校验共享口令
func CheckPassword(password string) bool {
	authMu.RLock()
	defer authMu.RUnlock()
	return authPassword != "" && password == authPassword
}

// CreateToken 签发 JWT
func CreateToken(user string) (string, error) {
	if user == "" {
		user = "guest"
	}
	claims := MyClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(TokenSecret)
}

// ExtractToken 从 Header / Query / Cookie 取 token
func ExtractToken(c *gin.Context) string {
	if t := strings.TrimSpace(c.Request.Header.Get("token")); t != "" {
		return t
	}
	if t := strings.TrimSpace(c.Query("token")); t != "" {
		return t
	}
	if t, err := c.Cookie("token"); err == nil {
		if t = strings.TrimSpace(t); t != "" {
			return t
		}
	}
	return ""
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return TokenSecret, nil
		})
	if err != nil {
		return nil, err
	}
	// 校验token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// EnsureAuth 若启用鉴权则校验 token；失败时已 Abort。返回是否放行。
func EnsureAuth(c *gin.Context) bool {
	if !AuthEnabled() {
		return true
	}
	authHeader := ExtractToken(c)
	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请求缺少token信息"})
		c.Abort()
		return false
	}
	mc, err := ParseToken(authHeader)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请求的token信息无效"})
		c.Abort()
		return false
	}
	c.Set("user", mc.User)
	return true
}

// JWTMiddleware 基于JWT的认证中间件（Password 为空时自动放行）
func JWTMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if !EnsureAuth(c) {
			return
		}
		c.Next()
	}
}

// PathAuthMiddleware 对指定路径前缀强制鉴权（如 /files、/ws）
func PathAuthMiddleware(prefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, p := range prefixes {
			if path == p || strings.HasPrefix(path, p+"/") || strings.HasPrefix(path, p) && (len(path) == len(p) || path[len(p)] == '/') {
				if !EnsureAuth(c) {
					return
				}
				break
			}
		}
		c.Next()
	}
}
