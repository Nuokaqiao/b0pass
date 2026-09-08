package audit

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// ClientInfo 请求方设备摘要
type ClientInfo struct {
	IP     string
	UA     string
	Device string
}

// FromGin 从 Gin 上下文提取 IP / UA / 设备类型
func FromGin(c *gin.Context) ClientInfo {
	if c == nil || c.Request == nil {
		return ClientInfo{IP: "-", UA: "-", Device: "unknown"}
	}
	ua := strings.TrimSpace(c.Request.UserAgent())
	return ClientInfo{
		IP:     clientIP(c),
		UA:     ua,
		Device: guessDevice(ua),
	}
}

func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	if xri := strings.TrimSpace(c.GetHeader("X-Real-IP")); xri != "" {
		return xri
	}
	ip := c.ClientIP()
	if ip == "" {
		return "-"
	}
	return ip
}

func guessDevice(ua string) string {
	u := strings.ToLower(ua)
	switch {
	case strings.Contains(u, "iphone"), strings.Contains(u, "ipad"), strings.Contains(u, "ipod"):
		return "iOS"
	case strings.Contains(u, "android"):
		return "Android"
	case strings.Contains(u, "windows"):
		return "Windows"
	case strings.Contains(u, "macintosh"), strings.Contains(u, "mac os"):
		return "macOS"
	case strings.Contains(u, "linux"):
		return "Linux"
	case ua == "":
		return "unknown"
	default:
		return "other"
	}
}

func shortUA(ua string) string {
	ua = strings.TrimSpace(ua)
	if ua == "" {
		return "-"
	}
	if len(ua) > 120 {
		return ua[:117] + "..."
	}
	return ua
}

// Log 记录操作与设备信息
func Log(c *gin.Context, op string, extra string) {
	info := FromGin(c)
	if extra != "" {
		log.Printf("[audit] op=%s ip=%s device=%s %s ua=%q", op, info.IP, info.Device, extra, shortUA(info.UA))
		return
	}
	log.Printf("[audit] op=%s ip=%s device=%s ua=%q", op, info.IP, info.Device, shortUA(info.UA))
}

// LogFields 无 gin 上下文时直接打日志（如 WebSocket）
func LogFields(op, ip, device, ua, extra string) {
	if ip == "" {
		ip = "-"
	}
	if device == "" {
		device = "unknown"
	}
	if extra != "" {
		log.Printf("[audit] op=%s ip=%s device=%s %s ua=%q", op, ip, device, extra, shortUA(ua))
		return
	}
	log.Printf("[audit] op=%s ip=%s device=%s ua=%q", op, ip, device, shortUA(ua))
}

// GuessDevice 导出给 WS 等使用
func GuessDevice(ua string) string {
	return guessDevice(ua)
}
