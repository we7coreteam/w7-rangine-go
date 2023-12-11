package logger

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go/src/core/helper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"sync"
	"time"
)

type Factory struct {
	driverResolverMap map[string]func(config Config) (zapcore.WriteSyncer, error)
	loggerResolverMap map[string]func() (*zap.Logger, error)
	loggerMap         map[string]*zap.Logger
	lock              sync.RWMutex
}

func NewLoggerFactory() *Factory {
	factory := &Factory{
		loggerMap:         make(map[string]*zap.Logger),
		loggerResolverMap: make(map[string]func() (*zap.Logger, error)),
		driverResolverMap: make(map[string]func(config Config) (zapcore.WriteSyncer, error)),
	}

	factory.RegisterDriverResolver("stream", factory.MakeFileStreamDriver)

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

func (factory *Factory) MakeFileStreamDriver(config Config) (zapcore.WriteSyncer, error) {
	fields := helper.ValidateAndGetErrFields(config)
	if len(fields) > 0 {
		return nil, errors.New("log config error, reason: fields: " + strings.Join(fields, ","))
	}

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
		Filename:   "./runtime/logs/" + config.Path,
		MaxSize:    int(config.MaxSize),
		MaxBackups: int(config.MaxBackups),
		MaxAge:     int(config.MaxDays),
		Compress:   false,
	}

	return zapcore.AddSync(&hook), nil
}

func (factory *Factory) MakeDriver(config Config) (zapcore.WriteSyncer, error) {
	driverResolver, exists := factory.driverResolverMap[config.Driver]
	if !exists {
		return nil, errors.New("logger driver " + config.Driver + " not exists")
	}

	return driverResolver(config)
}

func (factory *Factory) MakeLogger(level zapcore.Level, ws ...zapcore.WriteSyncer) *zap.Logger {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
	}
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,             // 小写编码器
		EncodeTime:     customTimeEncoder,              // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   customCallerEncoder,            // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		zapcore.NewMultiWriteSyncer(ws...),       // 打印到控制台和文件
		atomicLevel,                              // 日志级别
	)

	return zap.New(core)
}

func (factory *Factory) Channel(channel string) (*zap.Logger, error) {
	factory.lock.RLock()
	logger, exists := factory.loggerMap[channel]
	factory.lock.RUnlock()
	if exists {
		return logger, nil
	}

	factory.lock.Lock()
	defer factory.lock.Unlock()

	logger, exists = factory.loggerMap[channel]
	if !exists {
		loggerResolver, exists := factory.loggerResolverMap[channel]
		if !exists {
			return nil, errors.New("logger channel " + channel + " not exists")
		}

		var err error = nil
		logger, err = loggerResolver()
		if err != nil {
			return nil, errors.New("log resolve fail, channel:" + channel + ", error:" + err.Error())
		}
		factory.loggerMap[channel] = logger
	}

	return logger, nil
}

func (factory *Factory) RegisterDriverResolver(driver string, resolver func(config Config) (zapcore.WriteSyncer, error)) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterLogger(channel string, loggerResolver func() (*zap.Logger, error)) {
	_, exists := factory.loggerMap[channel]
	if exists {
		delete(factory.loggerMap, channel)
	}

	factory.loggerResolverMap[channel] = loggerResolver
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			factory.RegisterLogger(channel, func() (*zap.Logger, error) {
				driver, err := factory.MakeDriver(config)
				if err != nil {
					return nil, err
				}
				return factory.MakeLogger(factory.ConvertLevel(config.Level), driver), nil
			})
		}(key, value)
	}
}
