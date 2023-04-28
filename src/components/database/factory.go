package database

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strconv"
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

func (factory *Factory) MakeMysqlDriver(databaseConfig Config) (gorm.Dialector, error) {
	var dns = ""
	if databaseConfig.DSN != "" {
		dns = databaseConfig.DSN
	} else {
		dns = databaseConfig.User + ":" + databaseConfig.Password + "@tcp(" + databaseConfig.Host + ":" + strconv.Itoa(databaseConfig.Port) + ")/" + databaseConfig.DbName + "?charset=" + databaseConfig.Charset + "&parseTime=True&loc=Local"
	}

	return mysql.New(mysql.Config{
		DSN:                       dns,   // data source name
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), nil
}

func (factory *Factory) MakeSqliteDriver(databaseConfig Config) (gorm.Dialector, error) {
	return sqlite.Open(databaseConfig.DSN), nil
}

func (factory *Factory) MakeDriver(databaseConfig Config) (gorm.Dialector, error) {
	driverResolver, exists := factory.driverResolverMap[databaseConfig.Driver]
	if !exists {
		return nil, errors.New("db driver " + databaseConfig.Driver + " not exists")
	}

	return driverResolver(databaseConfig)
}

func (factory *Factory) MakeDb(databaseConfig Config, driver gorm.Dialector) (*gorm.DB, error) {
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
		return nil, err
	}

	dbDriver, err := db.DB()
	if err != nil {
		return nil, err
	}

	dbDriver.SetMaxIdleConns(databaseConfig.MaxIdleConn)
	dbDriver.SetMaxOpenConns(databaseConfig.MaxConn)

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
				driver, err := factory.MakeDriver(config)
				if err != nil {
					return nil, err
				}
				return factory.MakeDb(config, driver)
			})
		}(key, value)

	}
}
