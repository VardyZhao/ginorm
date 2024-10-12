package migration

import (
	"ginorm/model"
	"gorm.io/gorm"
)

var migrateMap = map[string][]interface{}{
	"default": {
		&model.User{},
	},
}

func Run(name string, db *gorm.DB) {
	if modelList, exists := migrateMap[name]; exists {
		_ = db.AutoMigrate(modelList...)
	}
}
