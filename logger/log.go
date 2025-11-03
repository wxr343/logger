package logger

import (
	"github.com/wxr343/logger/config"
	"github.com/wxr343/logger/global"
	"github.com/wxr343/logger/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	level   zapcore.Level // zap 日志等级
	options []zap.Option  // zap 配置项
)

func init() {
	InitializeConfig()
}

func InitializeLog(logConfig ...config.Configuration) *zap.Logger {
	var SetConfig config.Log
	if len(logConfig) == 0 {
		SetConfig = global.App.Config.Log
	} else {
		SetConfig = logConfig[0].Log
	}
	// 设置日志等级
	setLogLevel(SetConfig.Level)
	// 设置日志目录
	createRootDir(SetConfig.RootDir)
	// 设置日志是否显示调用行
	if SetConfig.ShowLine {
		options = append(options, zap.AddCaller())
	}

	return zap.New(getZapCore(SetConfig), options...)
}

func createRootDir(LogPath ...string) {
	var logPath string
	if len(LogPath) == 0 {
		logPath = global.App.Config.Log.RootDir
	} else {
		logPath = LogPath[0]
	}

	if ok, _ := utils.PathExists(logPath); !ok {
		_ = os.Mkdir(logPath, os.ModePerm)
	}
}

func setLogLevel(SetLevel ...string) {
	var logLevel string
	if len(SetLevel) == 0 {
		logLevel = global.App.Config.Log.Level
	} else {
		logLevel = SetLevel[0]
	}
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
		options = append(options, zap.AddStacktrace(level))
	case "panic":
		level = zap.PanicLevel
		options = append(options, zap.AddStacktrace(level))
	case "fatal":
		level = zap.FatalLevel
		options = append(options, zap.AddStacktrace(level))
	default:
		level = zap.InfoLevel
	}
}

// 扩展 Zap
func getZapCore(logConfig config.Log) zapcore.Core {
	var encoder zapcore.Encoder
	var logCore zapcore.Core
	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("logger" + "." + l.String())
	}

	// 设置编码器
	if logConfig.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	logCore = zapcore.NewCore(encoder, getLogWriter(logConfig), level)
	return logCore
}

// 使用 lumberjack 作为日志写入器
func getLogWriter(logConfig config.Log) zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   logConfig.RootDir + "/" + logConfig.Filename[0],
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		Compress:   logConfig.Compress,
	}

	return zapcore.AddSync(file)
}
