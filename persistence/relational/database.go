package relational

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

var config mysql.Config = mysql.Config{
	User:   "root",
	Passwd: "root",
	Net:    "tcp",
	Addr:   "localhost:3306",
	DBName: "blog",
}

func getDatabase() *sql.DB {
	var dbOpenError error
	var newConnection *sql.DB

	if db == nil {
		newConnection, dbOpenError = sql.Open("mysql", config.FormatDSN())

		if dbOpenError != nil {
			log.Fatal("Cannot connect to database", dbOpenError.Error())
		}

		if dbPingError := newConnection.Ping(); dbPingError != nil {
			log.Fatal(dbPingError.Error())
		}

		db = newConnection

		log.Println("Conected to database")
	}

	return db
}
