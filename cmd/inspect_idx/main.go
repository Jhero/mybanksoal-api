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

	rows, err := db.Query("SHOW INDEX FROM user_level_scores")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Indexes on user_level_scores:")
	for rows.Next() {
		var table, nonUnique, keyName, seq, colName, collation, cardinality, subPart, packed, null, indexType, comment, indexComment, visible, expression sql.NullString
		rows.Scan(&table, &nonUnique, &keyName, &seq, &colName, &collation, &cardinality, &subPart, &packed, &null, &indexType, &comment, &indexComment, &visible, &expression)
		fmt.Printf("- %s: %s\n", keyName.String, colName.String)
	}
}
