package logger

import (
	"fmt"
	"ginorm/config"
	"ginorm/constant"
	"ginorm/util"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

type loggerConfig struct {
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"mxx_size"`
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

// Load 初始化 Logger，将不同日志级别输出到不同文件
func Load() {
	if err := config.Conf.UnmarshalKey("log", &lc); err != nil {
		log.Fatalf("Error unmarshaling databases config: %v", err)
	}

	infoCore := newCore(zapcore.InfoLevel)   // INFO 级别日志
	errorCore := newCore(zapcore.ErrorLevel) // ERROR 级别日志

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
	var ctx *gin.Context

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]any:
			dataMap = v
		case *gin.Context:
			ctx = v
		}
	}

	if dataMap != nil {
		fields = append(fields, zap.Any("context", dataMap))
	}

	if ctx != nil {
		// 假设GetString方法存在
		fields = append(fields, zap.String(constant.TraceId, ctx.GetString(constant.TraceId)))
	}

	Writer.Log(level, msg, fields...)
}

func newCore(level zapcore.Level) zapcore.Core {
	filename, ok := levelFileMap[level]
	if !ok {
		filename = levelFileMap[zapcore.InfoLevel]
	}
	filename = util.GetAbsPath("/" + lc.Path + "/" + filename)
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    lc.MaxSize,    // 每个日志文件最大 10MB
		MaxBackups: lc.MaxBackups, // 最多保留 5 个旧文件
		MaxAge:     lc.MaxAge,     // 日志文件最长保留 30 天
		Compress:   lc.Compress,   // 是否压缩旧日志文件
	})

	// 配置日志编码
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"                        // 添加时间字段
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 使用 ISO8601 格式化时间
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	return zapcore.NewCore(encoder, writer, level)
}
