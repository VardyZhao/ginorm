package migration

import (
	"fmt"
	"ginorm/logger"
	"ginorm/util"
	"gorm.io/gorm"
	"os"
	"strings"
)

func Run(name string, db *gorm.DB) {
	filePath := util.GetAbsPath("/db/migration/" + name + ".sql")
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		logger.Writer.Error(fmt.Sprintf("failed to read SQL file: %v", err))
	}

	sqlStatements := strings.Split(string(sqlContent), ";")
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if len(stmt) > 0 {
			tx := db.Exec(stmt)
			if tx.Error != nil {
				logger.Writer.Error(fmt.Sprintf("failed to execute SQL: %v", tx.Error))
			}
		}
	}
}
