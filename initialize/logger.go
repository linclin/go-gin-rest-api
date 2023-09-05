package initialize

import (
	"fmt"
	"go-gin-rest-api/pkg/global"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
)

// 初始化日志
func Logger() {
	now := time.Now()
	fileName := fmt.Sprintf("%s/%04d-%02d-%02d.log", global.Conf.Logs.Path, now.Year(), now.Month(), now.Day())
	logFile := &lumberjack.Logger{
		Filename:   fileName,                    // 日志文件路径
		MaxSize:    global.Conf.Logs.MaxSize,    // 最大尺寸, M
		MaxBackups: global.Conf.Logs.MaxBackups, // 备份数
		MaxAge:     global.Conf.Logs.MaxAge,     // 存放天数
		Compress:   global.Conf.Logs.Compress,   // 是否压缩
	}
	logOpts := slog.HandlerOptions{
		AddSource: true,
		Level:     global.Conf.Logs.Level,
	}
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), &logOpts))
	slog.SetDefault(logger)
	global.Log = logger
	global.Logger = slog.NewLogLogger(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), &logOpts), slog.LevelInfo)
	global.Log.Info("初始化日志完成")
}
