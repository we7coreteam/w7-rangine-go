package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
)

type Factory struct {
	driverResolverMap map[string]func(config Config) zapcore.WriteSyncer
	loggerResolverMap map[string]func() *zap.Logger
	loggerMap         map[string]*zap.Logger
	once              sync.Once
}

func NewLoggerFactory() *Factory {
	factory := &Factory{
		loggerMap:         make(map[string]*zap.Logger),
		loggerResolverMap: make(map[string]func() *zap.Logger),
		driverResolverMap: make(map[string]func(config Config) zapcore.WriteSyncer),
	}

	factory.RegisterDriverResolver("file", factory.MakeFileDriver)

	return factory
}

func (factory *Factory) ConvertLevel(level string) zapcore.Level {
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

func (factory *Factory) MakeFileDriver(config Config) zapcore.WriteSyncer {
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

func (factory *Factory) MakeDriver(config Config) zapcore.WriteSyncer {
	driverResolver, exists := factory.driverResolverMap[config.Driver]
	if !exists {
		panic("logger driver " + config.Driver + " not exists")
	}

	return driverResolver(config)
}

func (factory *Factory) MakeLogger(level zapcore.Level, ws ...zapcore.WriteSyncer) *zap.Logger {
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

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2), zap.AddStacktrace(zap.FatalLevel+1))
}

func (factory *Factory) Channel(channel string) (*zap.Logger, error) {
	logger, exists := factory.loggerMap[channel]
	if exists {
		return logger, nil
	}

	var err error = nil
	factory.once.Do(func() {
		loggerResolver, exists := factory.loggerResolverMap[channel]
		if !exists {
			err = errors.New("logger channel " + channel + " not exists")
			return
		}

		factory.loggerMap[channel] = loggerResolver()
	})
	if err != nil {
		return nil, err
	}

	return factory.loggerMap[channel], nil
}

func (factory *Factory) RegisterDriverResolver(driver string, resolver func(config Config) zapcore.WriteSyncer) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterLogger(channel string, loggerResolver func() *zap.Logger) {
	factory.loggerResolverMap[channel] = loggerResolver
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			factory.RegisterLogger(channel, func() *zap.Logger {
				return factory.MakeLogger(factory.ConvertLevel(config.Level), factory.MakeDriver(config))
			})
		}(key, value)
	}
}