package initialize

import (
	"context"
	"fmt"
	"go-gin-rest-api/pkg/global"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"
)

// DailyWriter 按天分割的日志写入器
type DailyWriter struct {
	logDir      string
	appName     string
	currentFile *os.File
	currentDate string
	mutex       sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewDailyWriter 创建按天分割的日志写入器
func NewDailyWriter(logDir, appName string) *DailyWriter {
	dw := &DailyWriter{
		logDir:  logDir,
		appName: appName,
	}
	dw.ctx, dw.cancel = context.WithCancel(context.Background())

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	dw.rotateFile()
	dw.startRotationTicker()

	return dw
}

// Write 实现 io.Writer 接口
func (dw *DailyWriter) Write(p []byte) (n int, err error) {
	dw.mutex.RLock()
	file := dw.currentFile
	dw.mutex.RUnlock()

	return file.Write(p)
}

// rotateFile 切换日志文件
func (dw *DailyWriter) rotateFile() {
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s_%s.log", dw.appName, currentDate)
	fileName := filepath.Join(dw.logDir, logFileName)

	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	// 关闭旧文件
	if dw.currentFile != nil {
		dw.currentFile.Close()
	}

	// 打开新文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("打开日志文件失败: %v", err))
	}

	dw.currentFile = file
	dw.currentDate = currentDate
}

// startRotationTicker 启动定时轮转
func (dw *DailyWriter) startRotationTicker() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // 每天检查一次
		defer ticker.Stop()

		for {
			select {
			case <-dw.ctx.Done():
				return
			case <-ticker.C:
				// 检查日期是否改变
				currentDate := time.Now().Format("2006-01-02")
				if currentDate != dw.currentDate {
					dw.rotateFile()
				}
			}
		}
	}()
}

// Close 关闭写入器
func (dw *DailyWriter) Close() error {
	dw.cancel()

	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.currentFile != nil {
		return dw.currentFile.Close()
	}
	return nil
}

// 初始化日志
func Logger() {
	dailyWriter := NewDailyWriter(global.Conf.Logs.Path, global.Conf.System.AppName)

	logOpts := slog.HandlerOptions{
		AddSource: true,
		Level:     global.Conf.Logs.Level,
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, dailyWriter), &logOpts))
	slog.SetDefault(logger)
	global.Log = logger
	global.Logger = slog.NewLogLogger(slog.NewJSONHandler(io.MultiWriter(os.Stdout, dailyWriter), &logOpts), slog.LevelInfo)
	global.Log.Info("初始化日志完成")

	panicFile, err := os.Create(fmt.Sprintf("%s/panic_%s.log", global.Conf.Logs.Path, time.Now().Format("2006-01-02")))
	if err != nil {
		global.Log.Info(fmt.Sprint("初始化panic日志完成错误", err.Error()))
		panic(err)
	}
	debug.SetCrashOutput(panicFile, debug.CrashOptions{})
	global.Log.Info("初始化panic日志完成")
}
