package log

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger     *logrus.Logger
	once       sync.Once
	currentDay string
)

func SetupLogrusWithDailyRotation() {
	once.Do(func() {
		logger = logrus.New()
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		updateLoggerOutput()
		go monitorDayChange()
	})
}

func updateLoggerOutput() {
	currentDay = time.Now().Format("2006-01-02")
	logDir := filepath.Join("logs", currentDay)
	logFile := filepath.Join(logDir, "app.log")

	// 创建日志目录
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logger.Errorf("无法创建日志目录: %v", err)
		return
	}

	// 使用 lumberjack 实现日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // 每个日志文件最大10MB
		MaxBackups: 5,    // 最多保留5个备份日志文件
		MaxAge:     7,    // 最长保留7天的日志
		Compress:   true, // 启用压缩
	}

	logger.SetOutput(lumberjackLogger)
}

func monitorDayChange() {
	for {
		time.Sleep(time.Hour * 1)
		newDay := time.Now().Format("2006-01-02")
		if newDay != currentDay {
			updateLoggerOutput()
		}
	}
}

func Debug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warning(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func init() {
	SetupLogrusWithDailyRotation()
}
