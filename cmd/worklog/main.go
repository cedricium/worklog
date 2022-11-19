package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/teris-io/shortid"
	"github.com/urfave/cli/v2"
)

const ISO8601 = "2006-01-02 15:04:05"

type Entry struct {
	ID        string
	Timestamp time.Time
	Important bool
	Category  string
	Message   string
}

func main() {
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
					fmt.Println("TODO")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
