package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var Db *sql.DB

// This function will make a connection to the database only once.
func init() {
	var err error

	connStr := viper.GetString("DATABASE_URL")
	Db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = Db.Ping(); err != nil {
		panic(err)
	}
	// this will be printed in the terminal, confirming the connection to the database
	fmt.Println("The database is connected")
}
