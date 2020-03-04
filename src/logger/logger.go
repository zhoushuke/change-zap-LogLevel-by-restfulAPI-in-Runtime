package logger

import (
        "path/filepath"
	"strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
        "net/http"
	"os"
)

const (
	port = ":9090"
)

var logger *zap.SugaredLogger
var atomicLevel = zap.NewAtomicLevel()

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func init() {
        http.HandleFunc("/handle/level", atomicLevel.ServeHTTP)
        go func() {
            if err := http.ListenAndServe(port, nil); err != nil {
                panic(err)
            }
        }()

	filePath := getFilePath()
	level := getLoggerLevel("info")
	log := NewLogger(filePath, level, 256, 10, 7, true, "main")
        // defer log.Sync()
	logger = log.Sugar()
	logger.Sync()
	// SugaredLogger transfer back to Logger object
	//plain := logger.Desugar()
}


func NewLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool, serviceName string) *zap.Logger {
	core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.Development(), zap.Fields(zap.String("serviceName", serviceName)))
}

func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,   // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,      // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,  //
		EncodeCaller:   zapcore.ShortCallerEncoder,      // 短路径编码器
		// EncodeCaller:   zapcore.FullCallerEncoder,    // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}


type Level int8

const (
	DebugLevel Level = iota - 1

	InfoLevel

	WarnLevel

	ErrorLevel

	DPanicLevel

	PanicLevel

	FatalLevel
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Info(err)
	}
	return dir
}

func getFilePath() string {
	logfile := getCurrentDirectory() + "/" + getAppname() + ".log"
	return logfile
}

func getAppname() string {
	full := os.Args[0]
	splits := strings.Split(full, "/")
	if len(splits) >= 1 {
		name := splits[len(splits)-1]
		return name
	}
	return ""
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}
 
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(args ... interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}
