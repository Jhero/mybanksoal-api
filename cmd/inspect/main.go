package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jovan/mybanksoal-api/config"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT CONSTRAINT_NAME, TABLE_NAME FROM information_schema.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA = ? AND REFERENCED_TABLE_NAME IS NOT NULL", cfg.Database.Name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cName, tName string
		if err := rows.Scan(&cName, &tName); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Table: %s, Constraint: %s\n", tName, cName)
	}
}
