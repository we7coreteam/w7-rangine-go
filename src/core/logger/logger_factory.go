package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
)

type LoggerFactory struct {
	driverResolverMap map[string]func(config Config) zapcore.WriteSyncer
	loggerResolverMap map[string]func() *zap.Logger
	loggerMap         map[string]*zap.Logger
	once              sync.Once
}

func NewLoggerFactory() *LoggerFactory {
	loggerFactory := &LoggerFactory{
		loggerMap:         make(map[string]*zap.Logger),
		loggerResolverMap: make(map[string]func() *zap.Logger),
		driverResolverMap: make(map[string]func(config Config) zapcore.WriteSyncer),
	}

	loggerFactory.RegisterDriverResolver("file", loggerFactory.MakeFileDriver)

	return loggerFactory
}

func (loggerFactory *LoggerFactory) ConvertLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	case "Panic":
		return zap.PanicLevel
	default:
		return zap.DebugLevel
	}
}

func (loggerFactory *LoggerFactory) MakeFileDriver(config Config) zapcore.WriteSyncer {
	if config.MaxSize <= 0 {
		config.MaxSize = 2
	}
	if config.MaxDays <= 0 {
		config.MaxDays = 7
	}
	if config.MaxBackups <= 0 {
		config.MaxBackups = 1
	}
	hook := lumberjack.Logger{
		Filename:   "./runtimes/logs/" + config.Path,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxDays,
		Compress:   false,
	}

	return zapcore.AddSync(&hook)
}

func (loggerFactory *LoggerFactory) MakeDriver(config Config) zapcore.WriteSyncer {
	driverResolver, exists := loggerFactory.driverResolverMap[config.Driver]
	if !exists {
		panic("logger driver " + config.Driver + " not exists")
	}

	return driverResolver(config)
}

func (loggerFactory *LoggerFactory) MakeLogger(level zapcore.Level, ws ...zapcore.WriteSyncer) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		zapcore.NewMultiWriteSyncer(ws...),    // 打印到控制台和文件
		atomicLevel,                           // 日志级别
	)

	return zap.New(core)
}

func (loggerFactory *LoggerFactory) Channel(channel string) (*zap.Logger, error) {
	logger, exists := loggerFactory.loggerMap[channel]
	if exists {
		return logger, nil
	}

	var err error = nil
	loggerFactory.once.Do(func() {
		loggerResolver, exists := loggerFactory.loggerResolverMap[channel]
		if !exists {
			err = errors.New("logger channel " + channel + " not exists")
			return
		}

		loggerFactory.loggerMap[channel] = loggerResolver()
	})
	if err != nil {
		return nil, err
	}

	return loggerFactory.loggerMap[channel], nil
}

func (loggerFactory *LoggerFactory) RegisterDriverResolver(driver string, resolver func(config Config) zapcore.WriteSyncer) {
	loggerFactory.driverResolverMap[driver] = resolver
}

func (loggerFactory *LoggerFactory) RegisterLogger(channel string, loggerResolver func() *zap.Logger) {
	loggerFactory.loggerResolverMap[channel] = loggerResolver
}

func (loggerFactory *LoggerFactory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			loggerFactory.RegisterLogger(channel, func() *zap.Logger {
				return loggerFactory.MakeLogger(loggerFactory.ConvertLevel(config.Level), loggerFactory.MakeDriver(config))
			})
		}(key, value)
	}
}
