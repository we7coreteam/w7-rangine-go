package database

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DbLogger struct {
	logger.Writer
	logger *zap.Logger
}

func (dbLogger *DbLogger) Printf(info string, vs ...interface{}) {
	var _vs []zap.Field
	for k, v := range vs {
		_vs = append(_vs, zap.Reflect(string(rune(k)), v))
	}
	dbLogger.logger.Info(info, _vs...)
}

type Factory struct {
	driverResolverMap map[string]func(config Config) gorm.Dialector
	dbResolverMap     map[string]func() *gorm.DB
	dbMap             map[string]*gorm.DB
	logger            *zap.Logger
	once              sync.Once
	debug             bool
}

func NewDatabaseFactory() *Factory {
	factory := &Factory{
		dbMap:             make(map[string]*gorm.DB),
		dbResolverMap:     make(map[string]func() *gorm.DB),
		driverResolverMap: make(map[string]func(config Config) gorm.Dialector),
	}

	factory.RegisterDriverResolver("mysql", factory.MakeMysqlDriver)

	return factory
}

func (factory *Factory) SetDebug() {
	factory.debug = true
}

func (factory *Factory) SetLogger(logger *zap.Logger) {
	factory.logger = logger
}

func (factory *Factory) MakeMysqlDriver(databaseConfig Config) gorm.Dialector {
	dns := databaseConfig.User + ":" + databaseConfig.Password + "@tcp(" + databaseConfig.Host + ":" + strconv.Itoa(databaseConfig.Port) + ")/" + databaseConfig.DbName + "?charset=" + databaseConfig.Charset + "&parseTime=True&loc=Local"

	return mysql.New(mysql.Config{
		DSN:                       dns,   // data source name
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})
}

func (factory *Factory) MakeDriver(databaseConfig Config) gorm.Dialector {
	driverResolver, exists := factory.driverResolverMap[databaseConfig.Driver]
	if !exists {
		panic("db driver " + databaseConfig.Driver + " not exists")
	}

	return driverResolver(databaseConfig)
}

func (factory *Factory) MakeDb(databaseConfig Config, driver gorm.Dialector) *gorm.DB {
	//可根据配置开启日志
	var dbLogger logger.Interface = nil
	if factory.logger != nil {
		if databaseConfig.SlowThreshold <= 0 {
			databaseConfig.SlowThreshold = int64(200 * time.Millisecond)
		}
		dbLogger = logger.New(
			&DbLogger{
				logger: factory.logger,
			},
			logger.Config{
				SlowThreshold:             time.Duration(databaseConfig.SlowThreshold), // Slow SQL threshold
				LogLevel:                  logger.LogLevel(factory.logger.Level()),     // Log level
				IgnoreRecordNotFoundError: true,                                        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,                                       // Disable color
			},
		)
	}

	db, err := gorm.Open(driver, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   databaseConfig.Prefix,
			SingularTable: true,
		},
		Logger: dbLogger,
	})
	if err != nil {
		panic(err)
	}

	dbDriver, err := db.DB()
	if err != nil {
		panic(err)
	}

	dbDriver.SetMaxIdleConns(databaseConfig.MaxIdleConn)
	dbDriver.SetMaxOpenConns(databaseConfig.MaxConn)

	if factory.debug {
		db = db.Debug()
	}

	return db
}

func (factory *Factory) Channel(channel string) (*gorm.DB, error) {
	db, exists := factory.dbMap[channel]
	if exists {
		return db, nil
	}

	var err error = nil
	factory.once.Do(func() {
		dbResolver, exists := factory.dbResolverMap[channel]
		if !exists {
			err = errors.New("db channel " + channel + " not exists")
			return
		}

		factory.dbMap[channel] = dbResolver()
	})
	if err != nil {
		return nil, err
	}

	return factory.dbMap[channel], nil
}

func (factory *Factory) RegisterDriverResolver(driver string, resolver func(config Config) gorm.Dialector) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterDb(channel string, dbResolver func() *gorm.DB) {
	factory.dbResolverMap[channel] = dbResolver
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			factory.RegisterDb(channel, func() *gorm.DB {
				return factory.MakeDb(config, factory.MakeDriver(config))
			})
		}(key, value)

	}
}
