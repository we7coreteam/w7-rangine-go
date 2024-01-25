package test

import (
	databaseSpt "github.com/we7coreteam/w7-rangine-go/pkg/support/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)
import "testing"

func TestRegisterDbMap(t *testing.T) {
	factory := database.NewDatabaseFactory()

	factory.Register(map[string]databaseSpt.Config{
		"test1": {
			Driver: "sqlite_test",
			DbName: "./test1.db",
		},
		"test2": {
			Driver: "sqlite_test",
			DbName: "./test2.db",
		},
	})
	_, err := factory.Channel("test1")
	if err == nil {
		t.Error("database channel test1 driver error")
	}
	_, err = factory.Channel("test2")
	if err == nil {
		t.Error("database channel test2 driver error")
	}

	factory.RegisterDriver("sqlite_test", func(config databaseSpt.Config) (gorm.Dialector, error) {
		return sqlite.Open(config.DbName), nil
	})

	db1, err := factory.Channel("test1")
	if err != nil {
		t.Error(err)
	}
	db2, err := factory.Channel("test2")
	if err != nil {
		t.Error(err)
	}
	if _, ok := db1.Config.Dialector.(*sqlite.Dialector); !ok {
		t.Error("database channel test1 driver error")
	}
	if _, ok := db2.Config.Dialector.(*sqlite.Dialector); !ok {
		t.Error("database channel test2 driver error")
	}

	os.Remove("./test1.db")
	os.Remove("./test2.db")
}

func TestRegisterDb(t *testing.T) {
	factory := database.NewDatabaseFactory()
	factory.RegisterDriver("sqlite", func(config databaseSpt.Config) (gorm.Dialector, error) {
		return sqlite.Open(config.DbName), nil
	})

	factory.RegisterDb("sqlite1", func() (*gorm.DB, error) {
		driver, err := factory.MakeDriver(databaseSpt.Config{
			Driver: "sqlite",
			DbName: "./test1.db",
		})
		if err != nil {
			return nil, err
		}
		return factory.MakeDb(databaseSpt.Config{
			Prefix: "test1_",
		}, driver)
	})
	factory.RegisterDb("sqlite2", func() (*gorm.DB, error) {
		driver, err := factory.MakeDriver(databaseSpt.Config{
			Driver: "sqlite",
			DbName: "./test2.db",
		})
		if err != nil {
			return nil, err
		}
		return factory.MakeDb(databaseSpt.Config{
			Prefix: "test2_",
		}, driver)
	})

	test1Db, err := factory.Channel("sqlite1")
	if err != nil {
		t.Error(err)
	}
	test2Db, err := factory.Channel("sqlite2")
	if err != nil {
		t.Error(err)
	}
	if _, ok := test1Db.Config.Dialector.(*sqlite.Dialector); !ok {
		t.Error("database channel sqlite1 driver error")
	}
	if _, ok := test2Db.Config.Dialector.(*sqlite.Dialector); !ok {
		t.Error("database channel sqlite2 driver error")
	}

	os.Remove("./test1.db")
	os.Remove("./test2.db")
}
