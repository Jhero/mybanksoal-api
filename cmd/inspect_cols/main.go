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

	tables := []string{"crossword_questions", "user_level_scores"}
	for _, table := range tables {
		rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", table))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		fmt.Printf("Columns for %s:\n", table)
		for rows.Next() {
			var field, typ, null, key, def, extra sql.NullString
			rows.Scan(&field, &typ, &null, &key, &def, &extra)
			fmt.Printf("- %s (%s)\n", field.String, typ.String)
		}
	}
}
