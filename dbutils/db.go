package dbutils

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

//wsl: sudo service postgresql start
//sudo service postgresql status
//sudo service postgresql stop
//psql -U g2 -d cg4002 (g2)
// \l: show all db

//windows
//psql -U postgres
//\c cg4002
// username: postgres, password: postgres

//can consider batch update, maybe every session/ 5 minutes

var (
	once sync.Once
	db   *gorm.DB
)

func GetDB() *gorm.DB {
	var err error

	once.Do(func() {
		dsn := "host=localhost user=g2 password=g2 dbname=cg4002 port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})

	if err != nil {
		log.Fatalln("gorm db initialisation error: ", err.Error())
	}

	return db
}
