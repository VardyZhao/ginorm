package driver

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDatabase struct {
	BaseDatabase
}

func (m *MySQLDatabase) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	ConfigureConnectionPool(db)
	return db, nil
}
