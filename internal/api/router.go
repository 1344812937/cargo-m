package api

import (
	"cargo-m/internal/until"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GinLogWriter struct{}

func (w *GinLogWriter) Write(p []byte) (n int, err error) {
	log := until.Log
	// 去除尾部换行符
	msg := strings.TrimSpace(string(p))

	// 根据日志级别转发到 logrus
	switch {
	case strings.Contains(msg, "[WARNING]"):
		log.Warn(msg)
	case strings.Contains(msg, "[DEBUG]"):
		log.Debug(msg)
	case strings.Contains(msg, "[ERROR]"):
		log.Error(msg)
	default:
		log.Info(msg)
	}
	return len(p), nil
}

func initRouter() *gin.Engine {
	log := until.Log
	webEngine := gin.Default()
	// 重定向Gin日志
	gin.DefaultWriter = &GinLogWriter{}
	gin.DefaultErrorWriter = &GinLogWriter{}
	// 更换默认的日志输出方式
	webEngine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.WithFields(logrus.Fields{
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"user_agent": param.Request.UserAgent(),
		}).Info("HTTP Request")

		return ""
	}))
	return webEngine
}

// NewRouter 路由初始化
func NewRouter(mavenRepoHandler *MavenRepoHandler) *gin.Engine {
	webEngine := initRouter()
	// 注册路由
	mavenRepoHandler.Register(webEngine.Group("maven-repo"))
	return webEngine
}
