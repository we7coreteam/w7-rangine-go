package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
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

func (factory *Factory) MakeFileDriver(config Config) (zapcore.WriteSyncer, error) {
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
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // ???????????????
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC ????????????
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // ??????????????????
		EncodeName:     zapcore.FullNameEncoder,
	}

	// ??????????????????
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // ???????????????
		zapcore.NewMultiWriteSyncer(ws...),    // ???????????????????????????
		atomicLevel,                           // ????????????
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(4), zap.AddStacktrace(zap.FatalLevel+1))
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
			return nil, err
		}
		factory.loggerMap[channel] = logger
	}

	return logger, nil
}

func (factory *Factory) RegisterDriverResolver(driver string, resolver func(config Config) (zapcore.WriteSyncer, error)) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterLogger(channel string, loggerResolver func() (*zap.Logger, error)) {
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
