package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jovan/mybanksoal-api/config"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "up", "Command to run: up, down, step")
	var step int
	flag.IntVar(&step, "step", 1, "Number of steps for step command")
	var create string
	flag.StringVar(&create, "create", "", "Create a new migration with the given name")
	var dir string
	flag.StringVar(&dir, "dir", "db/migrations", "Migrations directory")
	flag.Parse()

	if create != "" {
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
		return
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Construct DSN for migrate
	// Format: mysql://user:password@tcp(host:port)/dbname?query
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	m, err := migrate.New(
		"file://db/migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	if len(os.Args) > 1 {
		// If arguments are provided without flags, assume the first one is the command
		// This allows running `go run cmd/migrate/main.go up`
		if os.Args[1] == "up" || os.Args[1] == "down" || os.Args[1] == "force" || os.Args[1] == "version" {
			command = os.Args[1]
		}
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while running up migration: %v", err)
		}
		log.Println("Migration up completed successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while running down migration: %v", err)
		}
		log.Println("Migration down completed successfully")
	case "step":
		if err := m.Steps(step); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("An error occurred while running step migration: %v", err)
		}
		log.Println("Migration step completed successfully")
	case "force":
		// Expects version as second arg, simplified here
		log.Println("Force command requires version implementation")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("An error occurred while getting version: %v", err)
		}
		log.Printf("Version: %v, Dirty: %v\n", v, dirty)
	default:
		log.Println("Usage: go run cmd/migrate/main.go [up|down|step|version]")
	}
}
