package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/teris-io/shortid"
	"github.com/urfave/cli/v2"
)

type Client struct {
	Database *sql.DB
}

const (
	dbFile string = "worklog.db"

	initializeStmt string = `CREATE TABLE IF NOT EXISTS entries (
	id TEXT NOT NULL PRIMARY KEY,
	timestamp DATETIME NOT NULL,
	important INTEGER NOT NULL DEFAULT 0,
	category TEXT NOT NULL DEFAULT 'note',
	message TEXT NOT NULL
);`
	insertStmt string = `INSERT INTO entries(id, timestamp, important, category, message) values(?, ?, ?, ?, ?);`
	listStmt   string = `SELECT * FROM entries ORDER BY timestamp DESC;`
	clearStmt  string = `DELETE FROM entries;`
)

const (
	clearWarning string = `CAUTION! This is a destructive action and connect be
undone. To proceed, type 'continue' or 'q' to quit:

> `
)

func (client *Client) Initialize() error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(initializeStmt)
	if err != nil {
		return err
	}

	client.Database = db
	return nil
}

func NewClient() (*Client, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initializeStmt)
	if err != nil {
		return nil, err
	}

	return &Client{Database: db}, nil
}

func (client *Client) Add(context *cli.Context) error {
	entry := Entry{
		ID:        shortid.MustGenerate(),
		Timestamp: time.Now(),
		Category:  context.String("category"),
		Important: context.Bool("important"),
		Message:   context.String("message"),
	}

	if _, err := client.Database.Exec(insertStmt, entry.ID, entry.Timestamp.Format(ISO8601),
		entry.Important, entry.Category, entry.Message); err != nil {
		return err
	}

	fmt.Println(entry)
	return nil
}

func (client *Client) List(context *cli.Context) error {
	rows, err := client.Database.Query(listStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		entry := Entry{}
		if err := rows.Scan(&entry.ID, &entry.Timestamp, &entry.Important,
			&entry.Category, &entry.Message); err != nil {
			return err
		}

		fmt.Println(entry)
	}

	return nil
}

func (client *Client) Clear(context *cli.Context) error {
	force := context.Bool("force")
	if !force {
		fmt.Print(clearWarning)

		var input string
		fmt.Scanln(&input)

		switch input {
		case "q", "quit":
			return nil
		case "continue":
			break
		default:
			return fmt.Errorf("input value '%v' does not match 'continue'", input)
		}
	}

	if _, err := client.Database.Exec(clearStmt); err != nil {
		return err
	}

	return nil
}
