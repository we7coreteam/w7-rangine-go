package database

import (
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

type DatabaseFactory struct {
	driverResolverMap map[string]func(config Config) gorm.Dialector
	dbResolverMap     map[string]func() *gorm.DB
	dbMap             map[string]*gorm.DB
	logger            *zap.Logger
	once              sync.Once
}

func NewDatabaseFactory() *DatabaseFactory {
	databaseFactory := &DatabaseFactory{
		dbMap:             make(map[string]*gorm.DB),
		dbResolverMap:     make(map[string]func() *gorm.DB),
		driverResolverMap: make(map[string]func(config Config) gorm.Dialector),
	}

	databaseFactory.RegisterDriverResolver("mysql", databaseFactory.MakeMysqlDriver)

	return databaseFactory
}

func (databaseFactory *DatabaseFactory) SetLogger(logger *zap.Logger) {
	databaseFactory.logger = logger
}

func (databaseFactory *DatabaseFactory) MakeMysqlDriver(databaseConfig Config) gorm.Dialector {
	dns := databaseConfig.User + ":" + databaseConfig.Password + "@tcp(" + databaseConfig.Host + ":" + strconv.Itoa(databaseConfig.Port) + ")/" + databaseConfig.DbName + "?charset=" + databaseConfig.Charset + "&parseTime=True&loc=Local"

	return mysql.New(mysql.Config{
		DSN:                       dns,   // data source name
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})
}

func (databaseFactory *DatabaseFactory) MakeDriver(databaseConfig Config) gorm.Dialector {
	driverResolver, exists := databaseFactory.driverResolverMap[databaseConfig.Driver]
	if !exists {
		panic("db driver " + databaseConfig.Driver + " not exists")
	}

	return driverResolver(databaseConfig)
}

func (databaseFactory *DatabaseFactory) MakeDb(databaseConfig Config, driver gorm.Dialector) *gorm.DB {
	//可根据配置开启日志
	var dbLogger logger.Interface = nil
	if databaseFactory.logger != nil {
		dbLogger = logger.New(
			&DbLogger{
				logger: databaseFactory.logger,
			},
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.Info,            // Log level
				IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,                  // Disable color
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

	return db
}

func (databaseFactory *DatabaseFactory) Channel(channel string) *gorm.DB {
	db, exists := databaseFactory.dbMap[channel]
	if exists {
		return db
	}

	databaseFactory.once.Do(func() {
		dbResolver, exists := databaseFactory.dbResolverMap[channel]
		if !exists {
			panic("db channel " + channel + " not exists")
		}

		databaseFactory.dbMap[channel] = dbResolver()
	})

	return databaseFactory.dbMap[channel]
}

func (databaseFactory *DatabaseFactory) RegisterDriverResolver(driver string, resolver func(config Config) gorm.Dialector) {
	databaseFactory.driverResolverMap[driver] = resolver
}

func (databaseFactory *DatabaseFactory) RegisterDb(channel string, dbResolver func() *gorm.DB) {
	databaseFactory.dbResolverMap[channel] = dbResolver
}

func (databaseFactory *DatabaseFactory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			databaseFactory.RegisterDb(channel, func() *gorm.DB {
				return databaseFactory.MakeDb(config, databaseFactory.MakeDriver(config))
			})
		}(key, value)

	}
}
