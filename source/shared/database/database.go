package database

import (
	"database/sql"
	"log"

	// Import mysql driver
	_ "github.com/go-sql-driver/mysql"

	"fmt"
)

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// DB is a sql database handle for use with this package
var DB *sql.DB

// Connect opens a database with specified information
func Connect(dbinfo *MySQLInfo) {
	var err error
	// Open the database
	DB, err = sql.Open("mysql", DSN(dbinfo))
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
}

// Terminate closes database and releases resources
func Terminate() error {
	return DB.Close()
}

// DSN returns the Data Source Name
func DSN(ci *MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}
