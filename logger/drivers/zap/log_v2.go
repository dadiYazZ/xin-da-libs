package zap

import (
	"context"
	"fmt"
	"github.com/dadiYazZ/xin-da-libs/helper"
	"github.com/dadiYazZ/xin-da-libs/logger/contract"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LoggerV2 struct {
	Logger *zap.Logger
	sugar  *zap.SugaredLogger
	ctx    context.Context
}

func NewLoggerV2(config *LoggerConfig) (logger contract.LoggerInterface, err error) {

	zapLogger, err := newZapLoggerV2(config)

	if err != nil {
		return nil, err
	}

	defer zapLogger.Sync() // flushes buffer, if any

	logger = &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}

	return logger, err
}

func (log *LoggerV2) WithContext(ctx context.Context) contract.LoggerInterface {

	if log.ctx != nil {
		return log
	}
	log.ctx = ctx

	traceID := helper.TraceIDFromContext(ctx)
	if len(traceID) > 0 {
		log.sugar = log.sugar.With(traceKey, traceID)
	}

	spanID := helper.SpanIDFromContext(log.ctx)
	if len(spanID) > 0 {
		log.sugar = log.sugar.With(spanKey, spanID)
	}

	return log
}

func (log *LoggerV2) Debug(msg string, v ...interface{}) {
	log.sugar.Debugw(msg, v...)
}
func (log *LoggerV2) Info(msg string, v ...interface{}) {
	log.sugar.Infow(msg, v...)
}
func (log *LoggerV2) Warn(msg string, v ...interface{}) {
	log.sugar.Warnw(msg, v...)
}
func (log *LoggerV2) Error(msg string, v ...interface{}) {
	log.sugar.Errorw(msg, v...)
}
func (log *LoggerV2) Panic(msg string, v ...interface{}) {
	log.sugar.Panicw(msg, v...)
}
func (log *LoggerV2) Fatal(msg string, v ...interface{}) {
	log.sugar.Fatalw(msg, v...)
}

func (log *LoggerV2) DebugF(format string, args ...interface{}) {
	log.sugar.Debugf(format, args...)
}
func (log *LoggerV2) InfoF(format string, args ...interface{}) {
	log.sugar.Infof(format, args...)
}
func (log *LoggerV2) WarnF(format string, args ...interface{}) {
	log.sugar.Warnf(format, args...)
}
func (log *LoggerV2) ErrorF(format string, args ...interface{}) {
	log.sugar.Errorf(format, args...)
}
func (log *LoggerV2) PanicF(format string, args ...interface{}) {
	log.sugar.Panicf(format, args...)
}
func (log *LoggerV2) FatalF(format string, args ...interface{}) {
	log.sugar.Fatalf(format, args...)
}

// LoggerConfig logger配置
type LoggerConfig struct {
	Env     string // "development" | "production"
	Level   string // "debug" | "info" | "warn" | "error"
	Stdout  bool   // true=控制台输出  false=文件输出
	LogDir  string // 文件日志目录，如 "./logs"
	FileExt string // 文件扩展名，如 ".log"
}

// ------------------------------------------------------------------
// newZapLogger 外部只暴露这个入口
func newZapLoggerV2(cfg *LoggerConfig) (*zap.Logger, error) {
	// ---------- 基础 zap.Config ----------
	zapCfg := baseZapConfig(cfg.Env)

	// ---------- Writer 工厂 ----------
	newWriteSyncer := func(name string) zapcore.WriteSyncer {
		if cfg.Stdout {
			return zapcore.AddSync(os.Stdout)
		}

		// 创建目录
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			panic("无法创建日志目录: " + err.Error())
		}

		logPath := filepath.Join(cfg.LogDir, name+cfg.FileExt)
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			panic(fmt.Sprintf("无法打开日志文件 %s: %v", logPath, err))
		}
		return zapcore.AddSync(file)
	}

	// ---------- 按级别分别创建 Core ----------
	cores := []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapCfg.EncoderConfig),
			newWriteSyncer("debug"),
			levelEnabler(zapcore.DebugLevel, cfg.Level),
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapCfg.EncoderConfig),
			newWriteSyncer("info"),
			levelEnabler(zapcore.InfoLevel, cfg.Level),
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapCfg.EncoderConfig),
			newWriteSyncer("warn"),
			levelEnabler(zapcore.WarnLevel, cfg.Level),
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zapCfg.EncoderConfig),
			newWriteSyncer("error"),
			levelEnabler(zapcore.ErrorLevel, cfg.Level),
		),
	}

	// ---------- 拼接为 Tee ----------
	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),      // 打印调用方
		zap.AddCallerSkip(1), // 让日志指向真正的业务层
	)

	return logger, nil
}

// ------------------------------------------------------------------
// 公共函数：生成基础 zap.Config
func baseZapConfig(env string) zap.Config {
	var c zap.Config
	if env == "production" {
		c = zap.NewProductionConfig()
	} else {
		c = zap.NewDevelopmentConfig()
	}

	enc := &c.EncoderConfig
	enc.TimeKey = "timestamp"
	enc.LevelKey = "level"
	enc.MessageKey = "message"
	enc.CallerKey = "caller"

	enc.LineEnding = zapcore.DefaultLineEnding
	enc.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	enc.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	enc.EncodeCaller = zapcore.FullCallerEncoder
	enc.EncodeDuration = zapcore.SecondsDurationEncoder

	return c
}

// ------------------------------------------------------------------
// 根据最小等级(minLevel)判断当前 core 是否输出
func levelEnabler(this zapcore.Level, minLevelStr string) zap.LevelEnablerFunc {
	minLevel := zapcore.InfoLevel
	switch strings.ToLower(minLevelStr) {
	case "debug":
		minLevel = zapcore.DebugLevel
	case "info":
		minLevel = zapcore.InfoLevel
	case "warn":
		minLevel = zapcore.WarnLevel
	case "error":
		minLevel = zapcore.ErrorLevel
	}
	return func(lvl zapcore.Level) bool { return lvl == this && lvl >= minLevel }
}
