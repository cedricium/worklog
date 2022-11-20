package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

const (
	ISO8601         string = "2006-01-02 15:04:05"
	EMPTY_ARG_USAGE string = " "
)

var categories []string = []string{"bug", "feature", "fix", "meeting", "note", "refactor"}

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

func configureCommands(client *Client) []*cli.Command {
	return []*cli.Command{
		{
			Name:      "add",
			Usage:     "Add entries to the log",
			ArgsUsage: EMPTY_ARG_USAGE,
			Action: func(ctx *cli.Context) error {
				return client.Add(ctx)
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
					Usage: fmt.Sprintf("Choose a category for the entry. `TAG` must be one of: %v",
						"["+strings.Join(categories, "|")+"]"),
					Value: "note",
					Action: func(ctx *cli.Context, input string) error {
						for _, valid := range categories {
							if input == valid {
								return nil
							}
						}
						return fmt.Errorf("flag category value '%v' is not valid. options are: %v",
							input, "["+strings.Join(categories, "|")+"]")
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
			Name:      "list",
			Usage:     "Show recorded entries",
			ArgsUsage: EMPTY_ARG_USAGE,
			Action: func(ctx *cli.Context) error {
				return client.List(ctx)
			},
		},
		{
			Name:      "clear",
			Usage:     "Delete all recorded entries",
			ArgsUsage: EMPTY_ARG_USAGE,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "force",
					Aliases: []string{"f"},
					Usage:   "Skip confirmation and forcefully delete all entries.",
					Value:   false,
				},
			},
			Action: func(ctx *cli.Context) error {
				return client.Clear(ctx)
			},
		}}
}

func main() {
	client := Client{}
	if err := client.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer client.Database.Close()

	app := &cli.App{
		Name:                   "worklog",
		Usage:                  "An opinionated note-taking tool for the developer's day-to-day.",
		ArgsUsage:              EMPTY_ARG_USAGE,
		Version:                "0.0.1",
		UseShortOptionHandling: true,
		Commands:               configureCommands(&client),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
