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

type ListConditions struct {
	After   string
	Before  string
	Filters string
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

func (client *Entries) List(entries *[]worklog.Entry, conds ListConditions) error {
	exp := sq.Select("*").From("entries").OrderBy("timestamp DESC")
	if len(conds.After) > 0 {
		exp = exp.Where(sq.Gt{"timestamp": conds.After})
	}
	if len(conds.Before) > 0 {
		exp = exp.Where(sq.Lt{"timestamp": conds.Before})
	}

	contains := []string{}
	for _, char := range conds.Filters {
		switch char {
		// special case for filtering by `important` entries,
		// irrelevant to proceeding category filtering
		case 'I':
			exp = exp.Where(sq.Eq{"important": true})

		case 'B':
			contains = append(contains, "bug")
		case 'F':
			contains = append(contains, "feature")
		case 'R':
			contains = append(contains, "fix")
		case 'M':
			contains = append(contains, "meeting")
		case 'N':
			contains = append(contains, "note")
		case 'C':
			contains = append(contains, "refactor")
		}
	}
	if len(contains) > 0 {
		exp = exp.Where(sq.Eq{"category": contains})
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
