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

	stmts := []string{
		"ALTER TABLE user_level_scores DROP FOREIGN KEY user_level_scores_ibfk_2",
		"ALTER TABLE user_level_scores DROP INDEX user_id",
		// "ALTER TABLE user_level_scores DROP INDEX crossword_id", // Might fail if implicitly dropped or used by other things, let's try dropping column directly, usually safer to drop FK first.
		"ALTER TABLE user_level_scores DROP COLUMN crossword_id",
		"ALTER TABLE user_level_scores ADD UNIQUE(user_id, level_id)",
		"DROP TABLE IF EXISTS crosswords",
	}

	for _, stmt := range stmts {
		fmt.Printf("Executing: %s\n", stmt)
		_, err := db.Exec(stmt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			// Don't exit fatal, try next (e.g. index might not exist)
		}
	}
}
