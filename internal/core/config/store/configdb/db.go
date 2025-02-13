// Code generated by gowebx, DO AVOID EDIT.
package configdb

import (
	"github.com/gowvp/gb28181/internal/core/config"
	"gorm.io/gorm"
)

var _ config.Storer = DB{}

// DB Related business namespaces
type DB struct {
	db *gorm.DB
}

// NewDB instance object
func NewDB(db *gorm.DB) DB {
	return DB{db: db}
}

// Config Get business instance
func (d DB) Config() config.ConfigStorer {
	return Config(d)
}

// AutoMigrate sync database
func (d DB) AutoMigrate(ok bool) DB {
	if !ok {
		return d
	}
	if err := d.db.AutoMigrate(
		new(config.Config),
	); err != nil {
		panic(err)
	}
	return d
}
