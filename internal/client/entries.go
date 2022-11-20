package client

import (
	"database/sql"

	"github.com/cedricium/worklog"
	_ "github.com/mattn/go-sqlite3"
)

type Entries struct {
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

func (client *Entries) Initialize() error {
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

func NewClient() (*Entries, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initializeStmt)
	if err != nil {
		return nil, err
	}

	return &Entries{Database: db}, nil
}

func (client *Entries) Add(entry worklog.Entry) error {
	if _, err := client.Database.Exec(insertStmt, entry.ID, entry.Timestamp.Format(worklog.ISO8601),
		entry.Important, entry.Category, entry.Message); err != nil {
		return err
	}

	return nil
}

func (client *Entries) List(entries *[]worklog.Entry) error {
	rows, err := client.Database.Query(listStmt)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		entry := worklog.Entry{}
		if err := rows.Scan(&entry.ID, &entry.Timestamp, &entry.Important,
			&entry.Category, &entry.Message); err != nil {
			return err
		}

		*entries = append(*entries, entry)
	}

	return nil
}

func (client *Entries) Clear() error {
	if _, err := client.Database.Exec(clearStmt); err != nil {
		return err
	}

	return nil
}
