package driver

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type BaseDatabase struct {
	Name            string
	Host            string
	Port            string
	Username        string
	Password        string
	Migrate         bool
	Charset         string
	MaxConnections  int
	IdleConnections int
	ConnectTimeout  time.Duration
}

func (b *BaseDatabase) Init() {

}

func (b *BaseDatabase) Connect() (*gorm.DB, error) {
	return nil, nil
}

func (b *BaseDatabase) ConfigureConnectionPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
