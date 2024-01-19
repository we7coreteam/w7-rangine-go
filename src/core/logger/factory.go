package logger

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger/config"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger/driver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

type Factory struct {
	driverResolverMap map[string]func(config config.Driver) (driver.Driver, error)
	loggerResolverMap map[string]func() (*zap.Logger, error)
	loggerMap         map[string]*zap.Logger
	lock              sync.RWMutex
}

func NewLoggerFactory() *Factory {
	factory := &Factory{
		loggerMap:         make(map[string]*zap.Logger),
		loggerResolverMap: make(map[string]func() (*zap.Logger, error)),
		driverResolverMap: make(map[string]func(config config.Driver) (driver.Driver, error)),
	}

	factory.RegisterDriverResolver("console", driver.NewConsoleDriver)
	factory.RegisterDriverResolver("file", driver.NewFileDriver)
	factory.RegisterDriverResolver("stack", driver.NewStackDriver(func(channel string) (zapcore.Core, error) {
		logger, err := factory.Channel(channel)
		if err != nil {
			return nil, err
		}

		return logger.Core(), nil
	}))

	return factory
}

func (factory *Factory) MakeDriver(config config.Driver) (driver.Driver, error) {
	driverResolver, exists := factory.driverResolverMap[config.Driver]
	if !exists {
		return nil, errors.New("logger driver " + config.Driver + " not exists")
	}

	return driverResolver(config)
}

func (factory *Factory) MakeLogger(drivers ...driver.Driver) *zap.Logger {
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

	core := NewDefaultLogger(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		drivers,                                  // 打印到控制台和文件
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

func (factory *Factory) RegisterDriverResolver(driver string, resolver func(config config.Driver) (driver.Driver, error)) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterLogger(channel string, loggerResolver func() (*zap.Logger, error)) {
	_, exists := factory.loggerMap[channel]
	if exists {
		delete(factory.loggerMap, channel)
	}

	factory.loggerResolverMap[channel] = loggerResolver
}

func (factory *Factory) Register(conf map[string]config.Driver) {
	for channelName, channel := range conf {
		func(channelName string, driver config.Driver) {
			factory.RegisterLogger(channelName, func() (*zap.Logger, error) {
				driverHandler, err := factory.MakeDriver(driver)
				if err != nil {
					return nil, err
				}

				return factory.MakeLogger(driverHandler), nil
			})
		}(channelName, channel)
	}
}
