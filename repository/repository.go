package repository

import (
	"ginorm/db"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

// 构造函数，获取数据库连接
func NewRepository(name ...string) *Repository {
	conn := "default"
	if len(name) > 0 {
		conn = name[0]
	}
	return &Repository{
		DB: db.Conn.GetDB(conn),
	}
}
