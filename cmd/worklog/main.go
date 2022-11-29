package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cedricium/worklog"
	"github.com/cedricium/worklog/internal/client"
	"github.com/teris-io/shortid"
	"github.com/urfave/cli/v2"
)

const (
	EMPTY_ARG_USAGE string = " "
	CLEAR_WARNING   string = `CAUTION! This is a destructive action and connect be
undone. To proceed, type 'continue' or 'q' to quit:

> `
)

var CATEGORIES []string = []string{"bug", "feature", "fix", "meeting", "note", "refactor"}

func configureCommands(entriesClient *client.Entries) []*cli.Command {
	return []*cli.Command{
		{
			Name:      "add",
			Usage:     "Add entries to the log",
			ArgsUsage: EMPTY_ARG_USAGE,
			Action: func(ctx *cli.Context) error {
				entry := worklog.Entry{
					ID:        shortid.MustGenerate(),
					Timestamp: time.Now(),
					Category:  ctx.String("category"),
					Important: ctx.Bool("important"),
					Message:   ctx.String("message"),
				}

				if err := entriesClient.Add(entry); err != nil {
					return err
				}

				fmt.Println(entry)
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
					Usage: `Choose a category for the entry. 'value' must be one of:
		[bug|feature|fix|meeting|note|refactor]`,
					Value: "note",
					Action: func(ctx *cli.Context, input string) error {
						for _, valid := range CATEGORIES {
							if input == valid {
								return nil
							}
						}

						return fmt.Errorf("flag category value '%v' is not valid. options are: %v",
							input, "["+strings.Join(CATEGORIES, "|")+"]")
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
				after := ctx.String("after")
				before := ctx.String("before")
				filters := ctx.String("filter")

				entries := []worklog.Entry{}
				conditions := client.ListConditions{After: after, Before: before, Filters: filters}

				if err := entriesClient.List(&entries, conditions); err != nil {
					return nil
				}

				for _, entry := range entries {
					fmt.Println(entry)
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "after",
					Aliases: []string{"a", "since"},
					Usage:   "Show entries newer than given `DATE`.",
				},
				&cli.StringFlag{
					Name:    "before",
					Aliases: []string{"b", "until"},
					Usage:   "Show entries older than given `DATE`.",
				},
				&cli.StringFlag{
					Name:    "filter",
					Aliases: []string{"f"},
					Usage: `Select only entries that are categorized:
		I marked important
		B bugs
		F features
		R repairs/fixes
		M meetings
		N notes
		C cleanups/refactors`,
				},
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
				force := ctx.Bool("force")
				if !force {
					fmt.Print(CLEAR_WARNING)

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

				if err := entriesClient.Clear(); err != nil {
					return err
				}

				return nil
			},
		}}
}

func main() {
	entriesClient := client.Entries{}
	if err := entriesClient.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer entriesClient.Database.Close()

	app := &cli.App{
		Name:                   "worklog",
		Usage:                  "An opinionated note-taking tool for the developer's day-to-day.",
		ArgsUsage:              EMPTY_ARG_USAGE,
		Version:                "0.0.1",
		UseShortOptionHandling: true,
		Commands:               configureCommands(&entriesClient),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
