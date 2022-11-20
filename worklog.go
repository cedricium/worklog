package worklog

import (
	"fmt"
	"time"
)

const (
	ISO8601 string = "2006-01-02 15:04:05"
)

type Entry struct {
	ID        string
	Timestamp time.Time
	Important bool
	Category  string
	Message   string
}

func (entry Entry) String() string {
	importantIndicator := " "
	if entry.Important {
		importantIndicator = "*"
	}

	return fmt.Sprintf("%v\t%v\t%v  [%v]\t'%v'", entry.Timestamp.Format(ISO8601),
		entry.ID, importantIndicator, entry.Category, entry.Message)
}
