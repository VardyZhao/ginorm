package logger

import (
	"context"
	"fmt"
	"ginorm/config"
	"ginorm/constant"
	"ginorm/util"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

type loggerConfig struct {
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

var (
	Writer       *zap.Logger
	lc           loggerConfig
	levelFileMap = map[zapcore.Level]string{
		zapcore.InfoLevel:   "info.log",
		zapcore.WarnLevel:   "warn.log",
		zapcore.ErrorLevel:  "error.log",
		zapcore.DPanicLevel: "dpanic.log",
		zapcore.PanicLevel:  "panic.log",
		zapcore.FatalLevel:  "fatal.log",
	}
	LevelInfo   = zapcore.InfoLevel
	LevelWarn   = zapcore.WarnLevel
	LevelError  = zapcore.ErrorLevel
	LevelDPanic = zapcore.DPanicLevel
	LevelPanic  = zapcore.PanicLevel
	LevelFatal  = zapcore.FatalLevel
)

func Load() {
	if err := config.Conf.UnmarshalKey("log", &lc); err != nil {
		log.Fatalf("Error unmarshaling databases config: %v", err)
	}

	// INFO 级别日志核心，只记录 INFO 和 WARN 等较低级别日志
	infoCore := newCore(zapcore.InfoLevel, zapcore.WarnLevel)

	// ERROR 级别日志核心，只记录 ERROR 及以上级别的日志
	errorCore := newCore(zapcore.ErrorLevel, zapcore.FatalLevel)

	core := zapcore.NewTee(infoCore, errorCore)

	// 加入自定义字段
	fields := []zap.Field{
		zap.String("app_name", config.Conf.GetString("app.name")),
		zap.String("app_version", config.Conf.GetString("app.version")),
	}

	options := []zap.Option{
		zap.Fields(fields...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	Writer = zap.New(core, options...)
}

func NewError(err error) zap.Field {
	return zap.Error(err)
}

func Flush() {
	err := Writer.Sync()
	if err != nil {
		fmt.Println("Sync log to file err:", err)
		return
	}
}

func Log(level zapcore.Level, msg string, args ...interface{}) {
	var fields []zap.Field
	var dataMap map[string]any
	var ctx context.Context

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]any:
			dataMap = v
		case context.Context:
			ctx = v
		}
	}

	if dataMap != nil {
		fields = append(fields, zap.Any("context", dataMap))
	}

	if ctx != nil {
		fields = append(fields, zap.Any(constant.TraceId, ctx.Value(constant.TraceId)))
	}

	Writer.Log(level, msg, fields...)
}

func newCore(minLevel zapcore.Level, maxLevel zapcore.Level) zapcore.Core {
	filename, ok := levelFileMap[minLevel]
	if !ok {
		filename = levelFileMap[zapcore.InfoLevel]
	}
	filename = util.GetAbsPath("/" + lc.Path + "/" + filename)

	// 文件日志输出
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    lc.MaxSize,
		MaxBackups: lc.MaxBackups,
		MaxAge:     lc.MaxAge,
		Compress:   lc.Compress,
	})

	// 控制台日志输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 配置日志编码
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"                        // 添加时间字段
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 使用 ISO8601 格式化时间
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 设置 LevelEnabler 过滤级别
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= minLevel && lvl <= maxLevel
	})

	// 使用 zapcore.NewTee，将文件和控制台输出合并并限制日志级别
	return zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, levelEnabler),
		zapcore.NewCore(encoder, consoleWriter, levelEnabler),
	)
}
