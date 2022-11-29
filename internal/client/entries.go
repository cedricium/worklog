package client

import (
	"database/sql"
	"os"
	"path"

	sq "github.com/Masterminds/squirrel"
	"github.com/cedricium/worklog"
	_ "github.com/mattn/go-sqlite3"
)

type Entries struct {
	Database *sql.DB
}

type ListFilters struct {
	After  string
	Before string
	/*
		TODO: implement filtering based on array of chars mapping to categories, e.g.
		{"B": "bugs"}
		{"F": "features"}
		{"R": "fix/(repair)"}
		{"M": "meeting"}
		{"N": "note"}
		{"C": "refactor/(cleanup)"}
		{"I": "important"}
		// categories string
	*/
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
)

func (client *Entries) Initialize() error {
	worklogDataDir := path.Join(os.Getenv("HOME"), ".local", "share", "worklog")
	if _, err := os.Stat(worklogDataDir); os.IsNotExist(err) {
		os.Mkdir(worklogDataDir, 0755)
	}

	db, err := sql.Open("sqlite3", path.Join(worklogDataDir, dbFile))
	if err != nil {
		return err
	}

	if _, err = db.Exec(initializeStmt); err != nil {
		return err
	}

	client.Database = db
	return nil
}

func (client *Entries) Add(entry worklog.Entry) error {
	if _, err := sq.Insert("entries").Values(entry.ID, entry.Timestamp.Format(worklog.ISO8601),
		entry.Important, entry.Category, entry.Message).RunWith(client.Database).Exec(); err != nil {
		return err
	}

	return nil
}

func (client *Entries) List(entries *[]worklog.Entry, filters ListFilters) error {
	exp := sq.Select("*").From("entries").OrderBy("timestamp DESC")
	if len(filters.After) > 0 {
		exp = exp.Where(sq.Gt{"timestamp": filters.After})
	}
	if len(filters.Before) > 0 {
		exp = exp.Where(sq.Lt{"timestamp": filters.Before})
	}

	rows, err := exp.RunWith(client.Database).Query()
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
	if _, err := sq.Delete("entries").RunWith(client.Database).Exec(); err != nil {
		return err
	}

	return nil
}
