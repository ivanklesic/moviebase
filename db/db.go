package db

import (
	"database/sql"
	"moviebase/moviebase/config"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

func GetConnection() *sql.DB {
	dbConfig := config.GetConfig().Db
	db, err := sql.Open(dbConfig["driver"], dbConfig["user"] + ":" + dbConfig["pass"] + "@/" + dbConfig["name"])
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}