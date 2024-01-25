package database

import (
	"errors"
	"fmt"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/database"
	loggerFactory "github.com/we7coreteam/w7-rangine-go/pkg/support/logger"
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
	buf := fmt.Appendf(nil, info, vs...)
	dbLogger.logger.Info(string(buf))
}

type Factory struct {
	driverResolverMap map[string]func(config database.Config) (gorm.Dialector, error)
	dbResolverMap     map[string]func() (*gorm.DB, error)
	dbMap             map[string]*gorm.DB
	loggerFactory     loggerFactory.Factory
	lock              sync.RWMutex
	debug             bool
}

func NewDatabaseFactory() *Factory {
	factory := &Factory{
		dbMap:             make(map[string]*gorm.DB),
		dbResolverMap:     make(map[string]func() (*gorm.DB, error)),
		driverResolverMap: make(map[string]func(config database.Config) (gorm.Dialector, error)),
	}

	factory.RegisterDriver("mysql", factory.MakeMysqlDriver)
	factory.RegisterDriver("sqlite", factory.MakeSqliteDriver)

	return factory
}

func (factory *Factory) SetDebug() {
	factory.debug = true
}

func (factory *Factory) SetLoggerFactory(loggerFactory loggerFactory.Factory) {
	factory.loggerFactory = loggerFactory
}

func (factory *Factory) MakeMysqlDriver(config database.Config) (gorm.Dialector, error) {
	if config.Port == 0 {
		config.Port = 3306
	}
	fields := helper.ValidateAndGetErrFields(config)
	if len(fields) > 0 {
		return nil, errors.New("database config error, reason: fields: " + strings.Join(fields, ","))
	}

	dns := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(int(config.Port)) + ")/" + config.DbName + "?charset=" + config.Charset + "&parseTime=True&loc=Local"
	return mysql.New(mysql.Config{
		DSN:                       dns,   // data source name
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), nil
}

func (factory *Factory) MakeSqliteDriver(config database.Config) (gorm.Dialector, error) {
	if config.Host == "" {
		config.Host = "-"
	}
	if config.Username == "" {
		config.Username = "-"
	}
	if config.Password == "" {
		config.Password = "-"
	}
	fields := helper.ValidateAndGetErrFields(config)
	if len(fields) > 0 {
		return nil, errors.New("database config error, reason: fields: " + strings.Join(fields, ","))
	}

	absPath, _ := filepath.Abs(config.DbName)
	dsn := fmt.Sprintf("file://%s?cache=shared&mode=rwc", absPath)
	return sqlite.Open(dsn), nil
}

func (factory *Factory) MakeDriver(config database.Config) (gorm.Dialector, error) {
	driverResolver, exists := factory.driverResolverMap[config.Driver]
	if !exists {
		return nil, errors.New("db driver " + config.Driver + " not exists")
	}

	return driverResolver(config)
}

func (factory *Factory) MakeDb(config database.Config, driver gorm.Dialector) (*gorm.DB, error) {
	//可根据配置开启日志
	var dbLogger logger.Interface = nil

	if factory.loggerFactory != nil {
		var loggerChannel = "default"
		val, exists := config.Options["logger"]
		if exists {
			loggerChannel = val.(string)
		}
		log, err := factory.loggerFactory.Channel(loggerChannel)
		if err == nil {
			if config.SlowThreshold <= 0 {
				config.SlowThreshold = uint64(200 * time.Millisecond)
			}
			dbLogger = logger.New(
				&DbLogger{
					logger: log,
				},
				logger.Config{
					SlowThreshold:             time.Duration(config.SlowThreshold), // Slow SQL threshold
					LogLevel:                  logger.Silent,                       // Log level
					IgnoreRecordNotFoundError: true,                                // Ignore ErrRecordNotFound error for logger
					Colorful:                  false,                               // Disable color
				},
			)
		}
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

	dbDriver.SetMaxIdleConns(int(config.MaxIdleConn))
	dbDriver.SetMaxOpenConns(int(config.MaxConn))

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
			return nil, errors.New("database resolve fail, channel:" + channel + ", error:" + err.Error())
		}
		factory.dbMap[channel] = db
	}

	return db, nil
}

func (factory *Factory) RegisterDriver(driver string, resolver func(config database.Config) (gorm.Dialector, error)) {
	factory.driverResolverMap[driver] = resolver
}

func (factory *Factory) RegisterDb(channel string, dbResolver func() (*gorm.DB, error)) {
	factory.dbResolverMap[channel] = dbResolver
}

func (factory *Factory) Register(maps map[string]database.Config) {
	for key, value := range maps {
		func(channel string, config database.Config) {
			factory.RegisterDb(channel, func() (*gorm.DB, error) {
				driver, err := factory.MakeDriver(config)
				if err != nil {
					return nil, err
				}
				return factory.MakeDb(config, driver)
			})
		}(key, value)

	}
}
