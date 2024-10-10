package db

import (
	"fmt"
	"ginorm/config"
	"ginorm/model/db/driver"
	"gorm.io/gorm"
)

type Database interface {
	Connect(dsn string) (*gorm.DB, error)
}

func LoadDB() {
	dbConfig := config.Conf.Get("database")
	if dbConfig == nil {
		return
	}
	if databases, ok := dbConfig.([]interface{}); ok {
		for _, db := range databases {
			if conf, ok := db.(map[string]interface{}); ok {
				LoadDriver(conf)
			} else {
				fmt.Println("Error loading db driver")
			}
		}
	} else {
		fmt.Println("No database configured")
	}
}

func LoadDriver(conf map[string]interface{}) Database {
	switch conf["driver"] {
	case "mysql":
		return driver.MySQLDatabase{}
	}
}
