package until

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = initLogger()
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	logger.SetLevel(logrus.DebugLevel)

	// 配置每日日志切割
	logPath := "./logs/app.until"
	logDir := filepath.Dir(logPath)
	os.MkdirAll(logDir, os.ModePerm) // 自动创建目录

	rotateLogger, _ := rotatelogs.New(
		logPath+".%Y%m%d",                         // 每日切割文件格式
		rotatelogs.WithLinkName(logPath),          // 软链接指向最新日志
		rotatelogs.WithRotationTime(24*time.Hour), // 每日切割
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 保留7天
	)

	// 添加hook实现写入切割文件和堆栈捕获
	logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.PanicLevel: rotateLogger,
			logrus.FatalLevel: rotateLogger,
			logrus.ErrorLevel: rotateLogger,
			logrus.WarnLevel:  rotateLogger,
			logrus.InfoLevel:  rotateLogger,
			logrus.DebugLevel: rotateLogger,
		},
		&logrus.JSONFormatter{},
	))

	// 添加堆栈捕获hook
	logger.AddHook(&StackTraceHook{})

	return logger
}

// 堆栈追踪Hook实现
type StackTraceHook struct{}

func (h *StackTraceHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}
}

func (h *StackTraceHook) Fire(entry *logrus.Entry) error {
	if entry.Level <= logrus.ErrorLevel {
		buf := make([]byte, 1<<16) // 64KB buffer
		stackSize := runtime.Stack(buf, false)
		stackTrace := string(buf[0:stackSize])

		// 简化堆栈路径（可选）
		stackTrace = simplifyPaths(stackTrace)

		entry.Data["stack"] = stackTrace
	}
	return nil
}

// 简化文件路径（可选）
func simplifyPaths(stack string) string {
	lines := strings.Split(stack, "\n")
	for i := 1; i < len(lines); i += 2 { // 跳过goroutine行
		if lines[i] == "" {
			continue
		}
		// 去除GOROOT路径
		if strings.HasPrefix(lines[i], runtime.GOROOT()) {
			lines[i] = strings.Replace(lines[i], runtime.GOROOT(), "$GOROOT", 1)
		}
		// 缩短GOPATH路径
		if gopath := os.Getenv("GOPATH"); gopath != "" {
			lines[i] = strings.Replace(lines[i], gopath+"/src/", "", -1)
		}
	}
	return strings.Join(lines, "\n")
}
