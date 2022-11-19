package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/teris-io/shortid"
	"github.com/urfave/cli/v2"
)

const ISO8601 string = "2006-01-02 15:04:05"

type Entry struct {
	ID        string
	Timestamp time.Time
	Important bool
	Category  string
	Message   string
}

const (
	db_filename string = "worklog.db"

	initialize_db string = `
CREATE TABLE IF NOT EXISTS entries (
	id 				TEXT NOT NULL PRIMARY KEY,
	timestamp DATETIME NOT NULL,
	important INTEGER NOT NULL DEFAULT 0,
	category 	TEXT NOT NULL DEFAULT 'note',
	message 	TEXT NOT NULL
);`
	insert_entry string = `
INSERT INTO entries(id, timestamp, important, category, message) values(?, ?, ?, ?, ?);
`
	list_sorted string = `
SELECT * FROM entries ORDER BY timestamp DESC;
`
)

func main() {
	db, err := sql.Open("sqlite3", db_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(initialize_db); err != nil {
		log.Fatal(err)
	}

	categories := []string{"bug", "feature", "fix", "meeting", "note", "refactor"}

	app := &cli.App{
		Usage:                  "An opinionated note-taking tool for the developer's day-to-day.",
		Version:                "0.0.1",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add entries to the log",
				Action: func(c *cli.Context) error {
					entry := Entry{
						ID:        shortid.MustGenerate(),
						Timestamp: time.Now(),
						Category:  c.String("category"),
						Important: c.Bool("important"),
						Message:   c.String("message"),
					}

					importantIndicator := " "
					if entry.Important {
						importantIndicator = "*"
					}

					if _, err := db.Exec(insert_entry, entry.ID, entry.Timestamp.Format(ISO8601), entry.Important, entry.Category, entry.Message); err != nil {
						log.Fatal(err)
						return err
					}
					fmt.Printf("%v\t%v\t%v  [%v]\t\t'%v'\n", entry.Timestamp.Format(ISO8601), entry.ID, importantIndicator, entry.Category, entry.Message)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "message",
						Aliases:  []string{"m"},
						Usage:    "Use the given `MSG` as the entry body",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Usage:   fmt.Sprintf("Choose a category for the entry. `TAG` must be one of: %v", "["+strings.Join(categories, "|")+"]"),
						Value:   "note",
						Action: func(ctx *cli.Context, input string) error {
							for _, valid := range categories {
								if input == valid {
									return nil
								}
							}
							return fmt.Errorf("flag category value '%v' is not valid. options are: %v", input, "["+strings.Join(categories, "|")+"]")
						},
					},
					&cli.BoolFlag{
						Name:    "important",
						Aliases: []string{"i"},
						Usage:   "Mark/flag the entry as important",
						Value:   false,
					},
				},
			},
			{
				Name:  "list",
				Usage: "Show recorded entries",
				Action: func(ctx *cli.Context) error {
					rows, err := db.Query(list_sorted)
					if err != nil {
						return err
					}
					defer rows.Close()

					for rows.Next() {
						entry := Entry{}
						err = rows.Scan(&entry.ID, &entry.Timestamp, &entry.Important, &entry.Category, &entry.Message)
						if err != nil {
							return err
						}

						importantIndicator := " "
						if entry.Important {
							importantIndicator = "*"
						}
						fmt.Printf("%v\t%v\t%v  [%v]\t\t'%v'\n", entry.Timestamp.Format(ISO8601), entry.ID, importantIndicator, entry.Category, entry.Message)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
