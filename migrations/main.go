package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var create string
	flag.StringVar(&create, "create", "", "Create a new migration with the given name")
	var dir string
	flag.StringVar(&dir, "dir", "db/migrations", "Migrations directory")
	flag.Parse()

	if create == "" {
		log.Println("Usage: go run migrations/main.go -create <name> [-dir db/migrations]")
		return
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to ensure migrations dir: %v", err)
	}
	ts := time.Now().Format("20060102150405")
	name := strings.TrimSpace(strings.ToLower(create))
	name = strings.ReplaceAll(name, " ", "_")
	base := fmt.Sprintf("%s_%s", ts, name)
	up := filepath.Join(dir, base+".up.sql")
	down := filepath.Join(dir, base+".down.sql")
	upSQL := fmt.Sprintf("CREATE TABLE %s (\n    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n", name)
	downSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", name)
	if err := os.WriteFile(up, []byte(upSQL), 0644); err != nil {
		log.Fatalf("Failed to write up migration: %v", err)
	}
	if err := os.WriteFile(down, []byte(downSQL), 0644); err != nil {
		log.Fatalf("Failed to write down migration: %v", err)
	}
	log.Printf("Created migrations:\n%s\n%s\n", up, down)
}
