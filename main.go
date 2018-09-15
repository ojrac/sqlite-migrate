package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ojrac/libmigrate"
)

var usage string = `Sqlite-migrate is a tool for managing a sqlite3 database's schema.

Usage:

	sqlite-migrate [version] [arguments]
	sqlite-migrate command [arguments]

The commands are:

	(none)		Migrate to the latest migration version on disk
	[number]	Migrate to the given migration version number
	create		Create a new up and down migration in the migration directory
	version		Print the current migration version
	pending		Print true if there are unapplied migrations, or else false

All commands require a database file path, which can be set by command-line
flags or environment variables:

	--db-path (env: DB_PATH): Path to sqlite3 db file
	--migrations-path (env: MIGRATIONS_PATH): Directory containing migration files
		(default: ./migrations)

`

var dbPath string
var migrationsPath string

func parseFlags() {
	migrationsPath = path.Join(".", "migrations")
	parseEnv(map[string]*string{
		"DB_PATH":         &dbPath,
		"MIGRATIONS_PATH": &migrationsPath,
	})

	flag.CommandLine.Usage = func() { fmt.Printf(usage) }
	flag.StringVar(&dbPath, "db-path", dbPath, "")
	flag.StringVar(&migrationsPath, "migrations-path", migrationsPath, "")
	flag.Parse()
}

func main() {
	parseFlags()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Error opening %s: %+v\n", dbPath, err)
		os.Exit(1)
	}
	defer db.Close()

	m := libmigrate.New(db, migrationsPath, libmigrate.ParamTypeQuestionMark)
	run(m, flag.Args())
}
