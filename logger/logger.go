package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"saas-account/config"
	"strings"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	// DEBUG 调试级别
	DEBUG LogLevel = iota
	// INFO 信息级别
	INFO
	// WARN 警告级别
	WARN
	// ERROR 错误级别
	ERROR
	// FATAL 致命错误级别
	FATAL
)

// Logg 日志记录器
type Logg struct {
	level  LogLevel
	logger *log.Logger
	output io.Writer
}

var (
	// 全局日志记录器
	Logger *Logg
)

// InitLogger 初始化日志记录器
func InitLogger(level string, output string) {
	var logLevel LogLevel
	switch strings.ToLower(level) {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	case "fatal":
		logLevel = FATAL
	default:
		logLevel = INFO
	}

	var logOutput io.Writer
	switch strings.ToLower(output) {
	case "file":
		// 创建日志目录
		logDir := "logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("创建日志目录失败: %v", err)
		}

		// 创建日志文件
		logFile := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("20060102")))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("打开日志文件失败: %v", err)
		}
		logOutput = file
	case "console":
		logOutput = os.Stdout
	default:
		logOutput = os.Stdout
	}

	Logger = &Logg{
		level:  logLevel,
		logger: log.New(logOutput, "", log.LstdFlags),
		output: logOutput,
	}
}

// GetLogger 获取全局日志记录器
func GetLogger() *Logg {
	if Logger == nil {
		// 如果全局日志记录器未初始化，则使用默认配置初始化
		config := config.GetConfig()
		InitLogger(config.LogLevel, config.LogOutput)
	}
	return Logger
}

// Debug 记录调试级别日志
func (l *Logg) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", format, v...)
	}
}

// DebugWithContext 记录带上下文的调试级别日志
func (l *Logg) DebugWithContext(ctx context.Context, format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logWithContext(ctx, "DEBUG", format, v...)
	}
}

// Info 记录信息级别日志
func (l *Logg) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", format, v...)
	}
}

// InfoWithContext 记录带上下文的信息级别日志
func (l *Logg) InfoWithContext(ctx context.Context, format string, v ...interface{}) {
	if l.level <= INFO {
		l.logWithContext(ctx, "INFO", format, v...)
	}
}

// Warn 记录警告级别日志
func (l *Logg) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", format, v...)
	}
}

// WarnWithContext 记录带上下文的警告级别日志
func (l *Logg) WarnWithContext(ctx context.Context, format string, v ...interface{}) {
	if l.level <= WARN {
		l.logWithContext(ctx, "WARN", format, v...)
	}
}

// Error 记录错误级别日志
func (l *Logg) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.log("ERROR", format, v...)
	}
}

// ErrorWithContext 记录带上下文的错误级别日志
func (l *Logg) ErrorWithContext(ctx context.Context, format string, v ...interface{}) {
	if l.level <= ERROR {
		l.logWithContext(ctx, "ERROR", format, v...)
	}
}

// Fatal 记录致命错误级别日志，并终止程序
func (l *Logg) Fatal(format string, v ...interface{}) {
	if l.level <= FATAL {
		l.log("FATAL", format, v...)
	}
	os.Exit(1)
}

// FatalWithContext 记录带上下文的致命错误级别日志，并终止程序
func (l *Logg) FatalWithContext(ctx context.Context, format string, v ...interface{}) {
	if l.level <= FATAL {
		l.logWithContext(ctx, "FATAL", format, v...)
	}
	os.Exit(1)
}

// log 记录日志
func (l *Logg) log(level, format string, v ...interface{}) {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	// 提取文件名
	file = filepath.Base(file)

	// 格式化日志消息
	message := fmt.Sprintf(format, v...)
	logMessage := fmt.Sprintf("[%s] %s:%d %s", level, file, line, message)

	// 记录日志
	l.logger.Println(logMessage)
}

// logWithContext 记录带上下文的日志
func (l *Logg) logWithContext(ctx context.Context, level, format string, v ...interface{}) {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	// 提取文件名
	file = filepath.Base(file)

	// 获取请求ID、跟踪 ID 和跨度 ID
	requestID := ""
	traceID := ""
	spanID := ""

	// 从上下文中获取请求ID
	if ctx != nil {
		if id, ok := ctx.Value("X-Request-ID").(string); ok && id != "" {
			requestID = id
		}
		if id, ok := ctx.Value("X-Trace-ID").(string); ok && id != "" {
			traceID = id
		}
		if id, ok := ctx.Value("X-Span-ID").(string); ok && id != "" {
			spanID = id
		}
	}

	// 格式化日志消息
	message := fmt.Sprintf(format, v...)
	var logMessage string
	if requestID != "" {
		logMessage = fmt.Sprintf("[%s] [%s] [%s] [%s] %s:%d %s", level, requestID, traceID, spanID, file, line, message)
	} else {
		logMessage = fmt.Sprintf("[%s] %s:%d %s", level, file, line, message)
	}

	// 记录日志
	l.logger.Println(logMessage)
}

// Close 关闭日志记录器
func (l *Logg) Close() {
	if closer, ok := l.output.(io.Closer); ok {
		closer.Close()
	}
}
