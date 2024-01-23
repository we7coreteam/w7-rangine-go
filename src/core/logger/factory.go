package logger

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go-support/src/logger"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger/driver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
	"time"
)

type Factory struct {
	driverResolverMap map[string]func(config logger.Config) (logger.Driver, error)
	loggerResolverMap map[string]func() (*zap.Logger, error)
	loggerMap         map[string]*zap.Logger
	lock              sync.Mutex
}

func NewLoggerFactory() *Factory {
	factory := &Factory{
		loggerMap:         make(map[string]*zap.Logger),
		loggerResolverMap: make(map[string]func() (*zap.Logger, error)),
		driverResolverMap: make(map[string]func(config logger.Config) (logger.Driver, error)),
	}

	factory.RegisterDriver("console", driver.NewConsoleDriver)
	factory.RegisterDriver("file", driver.NewFileDriver)
	factory.RegisterDriver("stack", driver.NewStackDriver(func(channel string) (zapcore.Core, error) {
		logger, err := factory.Channel(channel)
		if err != nil {
			return nil, err
		}

		return logger.Core(), nil
	}))

	return factory
}

func (factory *Factory) MakeDriver(config logger.Config) (logger.Driver, error) {
	driverResolver, exists := factory.driverResolverMap[config.Driver]
	if !exists {
		return nil, errors.New("logger driver " + config.Driver + " not exists")
	}

	return driverResolver(config)
}

func (factory *Factory) MakeLogger(drivers ...logger.Driver) *zap.Logger {
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
	//factory.lock.RLock()  //暂不用读写锁控制，这里可能会导致死锁, 比如新的channel 需要依赖其他channel, 写锁打开后，读锁获取不到，死锁
	logger, exists := factory.loggerMap[channel]
	//factory.lock.RUnlock()
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
		logger = logger.Named(channel)
		factory.loggerMap[channel] = logger
	}

	return logger, nil
}

func (factory *Factory) RegisterDriver(driver string, resolver func(config logger.Config) (logger.Driver, error)) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterLogger(channel string, loggerResolver func() (*zap.Logger, error)) {
	_, exists := factory.loggerMap[channel]
	if exists {
		delete(factory.loggerMap, channel)
	}

	factory.loggerResolverMap[channel] = loggerResolver
}

func (factory *Factory) Register(conf map[string]logger.Config) {
	for channelName, channel := range conf {
		func(channelName string, driver logger.Config) {
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
