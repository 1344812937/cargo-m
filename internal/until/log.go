package until

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/petermattis/goid"
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
	logPath := "./logs/app"
	logDir := filepath.Dir(logPath)
	os.MkdirAll(logDir, os.ModePerm) // 自动创建目录

	var options []rotatelogs.Option
	options = append(options, rotatelogs.WithRotationTime(24*time.Hour)) // 每日切割
	options = append(options, rotatelogs.WithMaxAge(7*24*time.Hour))     // 保留7天

	// 在非Windows环境或者有权限的环境下使用符号链接
	if runtime.GOOS != "windows" {
		options = append(options, rotatelogs.WithLinkName(logPath))
	}

	rotateLogger, err := rotatelogs.New(
		logPath+".%Y%m%d", // 每日切割文件格式
		options...,
	)
	if err != nil {
		fmt.Printf("failed to create rotate logs: %v\n", err)
		// 如果创建rotateLogger失败，我们使用标准输出作为日志输出，避免程序无法记录日志
		logger.Out = os.Stdout
	} else {
		// 添加hook实现写入切割文件
		logger.Hooks.Add(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.PanicLevel: rotateLogger,
				logrus.FatalLevel: rotateLogger,
				logrus.ErrorLevel: rotateLogger,
				logrus.WarnLevel:  rotateLogger,
				logrus.InfoLevel:  rotateLogger,
				logrus.DebugLevel: rotateLogger,
			},
			&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05.000",
			},
		))
	}

	// 添加goroutine ID的hook
	logger.AddHook(&GoroutineIDHook{})

	// 添加堆栈捕获hook
	logger.AddHook(&StackTraceHook{})

	return logger
}

// GoroutineIDHook 添加goroutine ID的hook
type GoroutineIDHook struct{}

func (h *GoroutineIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *GoroutineIDHook) Fire(entry *logrus.Entry) error {
	entry.Data["goroutine"] = goid.Get()
	return nil
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
	buf := make([]byte, 1<<16) // 64KB buffer
	stackSize := runtime.Stack(buf, false)
	stackTrace := string(buf[0:stackSize])

	// 简化堆栈路径（可选）
	stackTrace = simplifyPaths(stackTrace)

	entry.Data["stack"] = stackTrace
	return nil
}

// 简化文件路径（可选）
func simplifyPaths(stack string) string {
	lines := strings.Split(stack, "\n")
	for i := 1; i < len(lines); i += 2 { // 跳过goroutine行
		if i >= len(lines) {
			break
		}
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
		// 替换当前目录路径?
	}
	return strings.Join(lines, "\n")
}
