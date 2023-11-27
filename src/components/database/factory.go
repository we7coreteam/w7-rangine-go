package database

import (
	"errors"
	"fmt"
	"github.com/we7coreteam/w7-rangine-go/src/core/helper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DbLogger struct {
	logger.Writer
	logger *zap.Logger
}

func (dbLogger *DbLogger) Printf(info string, vs ...any) {
	var _vs []zap.Field
	for k, v := range vs {
		_vs = append(_vs, zap.Reflect(string(rune(k)), v))
	}
	dbLogger.logger.Info(info, _vs...)
}

type Factory struct {
	driverResolverMap map[string]func(config Config) (gorm.Dialector, error)
	dbResolverMap     map[string]func() (*gorm.DB, error)
	dbMap             map[string]*gorm.DB
	logger            *zap.Logger
	lock              sync.RWMutex
	debug             bool
}

func NewDatabaseFactory() *Factory {
	factory := &Factory{
		dbMap:             make(map[string]*gorm.DB),
		dbResolverMap:     make(map[string]func() (*gorm.DB, error)),
		driverResolverMap: make(map[string]func(config Config) (gorm.Dialector, error)),
	}

	factory.RegisterDriverResolver("mysql", factory.MakeMysqlDriver)
	factory.RegisterDriverResolver("sqlite", factory.MakeSqliteDriver)

	return factory
}

func (factory *Factory) SetDebug() {
	factory.debug = true
}

func (factory *Factory) SetLogger(logger *zap.Logger) {
	factory.logger = logger
}

func (factory *Factory) MakeMysqlDriver(config Config) (gorm.Dialector, error) {
	var dns = ""
	if config.DSN != "" {
		dns = config.DSN
	} else {
		dns = config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DbName + "?charset=" + config.Charset + "&parseTime=True&loc=Local"
	}

	return mysql.New(mysql.Config{
		DSN:                       dns,   // data source name
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), nil
}

func (factory *Factory) MakeSqliteDriver(config Config) (gorm.Dialector, error) {
	var dsn = ""
	if config.DSN != "" {
		dsn = config.DSN
	} else {
		absPath, _ := filepath.Abs(config.DbName)
		dsn = fmt.Sprintf("file://%s?cache=shared&mode=rwc", absPath)
	}
	return sqlite.Open(dsn), nil
}

func (factory *Factory) MakeDriver(config Config) (gorm.Dialector, error) {
	fmt.Printf("%v \n", config)

	driverResolver, exists := factory.driverResolverMap[config.Driver]
	if !exists {
		return nil, errors.New("db driver " + config.Driver + " not exists")
	}

	return driverResolver(config)
}

func (factory *Factory) MakeDb(config Config, driver gorm.Dialector) (*gorm.DB, error) {
	//可根据配置开启日志
	var dbLogger logger.Interface = nil
	if factory.logger != nil {
		if config.SlowThreshold <= 0 {
			config.SlowThreshold = int64(200 * time.Millisecond)
		}
		dbLogger = logger.New(
			&DbLogger{
				logger: factory.logger,
			},
			logger.Config{
				SlowThreshold:             time.Duration(config.SlowThreshold),     // Slow SQL threshold
				LogLevel:                  logger.LogLevel(factory.logger.Level()), // Log level
				IgnoreRecordNotFoundError: true,                                    // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,                                   // Disable color
			},
		)
	}

	db, err := gorm.Open(driver, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Prefix,
			SingularTable: true,
		},
		Logger: dbLogger,
	})
	if err != nil {
		return nil, err
	}

	dbDriver, err := db.DB()
	if err != nil {
		return nil, err
	}

	dbDriver.SetMaxIdleConns(config.MaxIdleConn)
	dbDriver.SetMaxOpenConns(config.MaxConn)

	if factory.debug {
		db = db.Debug()
	}

	return db, nil
}

func (factory *Factory) Channel(channel string) (*gorm.DB, error) {
	//map非线程安全  https://github.com/golang/go/blob/master/src/runtime/map.go#L579
	// double check rlock  https://launchdarkly.com/blog/golang-pearl-thread-safe-writes-and-double-checked-locking-in-go/
	factory.lock.RLock()
	db, exists := factory.dbMap[channel]
	factory.lock.RUnlock()
	if exists {
		return db, nil
	}

	factory.lock.Lock()
	defer factory.lock.Unlock()

	db, exists = factory.dbMap[channel]
	if !exists {
		dbResolver, exists := factory.dbResolverMap[channel]
		if !exists {
			return nil, errors.New("db channel " + channel + " not exists")
		}

		var err error = nil
		db, err = dbResolver()
		if err != nil {
			return nil, err
		}
		factory.dbMap[channel] = db
	}

	return db, nil
}

func (factory *Factory) RegisterDriverResolver(driver string, resolver func(config Config) (gorm.Dialector, error)) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterDb(channel string, dbResolver func() (*gorm.DB, error)) {
	factory.dbResolverMap[channel] = dbResolver
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			factory.RegisterDb(channel, func() (*gorm.DB, error) {
				fields := helper.ValidateAndGetErrFields(config)
				if len(fields) > 0 {
					panic("database config error, channel: " + channel + ", fields: " + strings.Join(fields, ","))
				}

				driver, err := factory.MakeDriver(config)
				if err != nil {
					return nil, err
				}
				return factory.MakeDb(config, driver)
			})
		}(key, value)

	}
}
