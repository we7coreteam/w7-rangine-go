package database

import (
	"gorm.io/gorm"
)

type Factory interface {
	Channel(channel string) (*gorm.DB, error)
	RegisterDriver(driver string, resolver func(config Config) (gorm.Dialector, error))
	RegisterDb(channel string, dbResolver func() (*gorm.DB, error))
}
