package relational

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

var config mysql.Config = mysql.Config{
	User:   os.Getenv("DB_USER"),
	Passwd: os.Getenv("DB_PASSWORD"),
	Net:    os.Getenv("DB_NET"),
	Addr:   os.Getenv("DB_ADDRESS"),
	DBName: os.Getenv("DB_NAME"),
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
