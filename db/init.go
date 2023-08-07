package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Username string
	Password string
	Address  string
	Port     string
	DB_name  string
}

type rolandDB struct {
	db *sql.DB
}

func getDSN(db DB) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.Username, db.Password, db.Address, db.Port, db.DB_name)

}
func (db_info DB) Connect() rolandDB {

	db, err := sql.Open("mysql", getDSN(db_info))

	roland := rolandDB{db: db}

	if err != nil {
		log.Fatal(err.Error())
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	create_db := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", db_info.DB_name)
	_, err = db.Exec(create_db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL and database initialized!")

	return roland
}

func (r rolandDB) Close() error {
	return r.db.Close()
}

func (r rolandDB) Ping() error {
	return r.db.Ping()
}
