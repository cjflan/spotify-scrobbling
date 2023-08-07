package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	username string
	password string
	address  string
	port     string
	db_name  string
}

func getDSN(db DB) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.username, db.password, db.address, db.port, db.db_name)

}
func Connect() *sql.DB {

	db_info := DB{
		username: os.Getenv("MYSQL_USER"),
		password: os.Getenv("MYSQL_PASSWORD"),
		address:  "127.0.0.1",
		port:     "3306",
		db_name:  os.Getenv("MYSQL_DATABASE"),
	}

	db, err := sql.Open("mysql", getDSN(db_info))

	if err != nil {
		log.Fatal(err.Error())
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	create_db := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", db_info.db_name)
	_, err = db.Exec(create_db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL and database initialized!")

	return db
}
