# worklog 📝

An opinionated note-taking tool for the developer's day-to-day.

## Idea

After working at Hellosaurus for two years, I can really only remember the big features I built. But what about the tricky bugs solved, legacy code rewrites, or instrumental meetings attended? Those moments, arguably just as important as feature work, stuck somewhere in memory though seemingly inaccessible forever.

Wouldn't it be nice if there was a way to simply and quickly document one's day-to-day while at work? Something akin to the captain's log in Star Trek, with a way for easy perusal at a later time.

Inspired by `git`, I envision a CLI tool and package for documenting and parsing: tasks, notes, meetings, bug fixes, etc., automatically grouped by day in order to maintain a record of the day-to-day on the job.

## Getting Started

### Build

We can use `make` to run common/frequent commands, such as building an executable and cleaning the project.

```console
$ make
```

This will compile an up-to-date executable called `worklog` in the root directory based on the latest code changes from across the project (`worklog.go`, `internal/client`, `cmd/worklog`). To use, see next section.

## Usage

This is a work-in-progress and currently being used to drive development. Very likely to change in the future.

### `$ worklog`

```console
$ worklog --help
NAME:
   worklog - An opinionated note-taking tool for the developer's day-to-day.

USAGE:
   worklog [global options] command [command options]

VERSION:
   0.0.1

COMMANDS:
   add      Add entries to the log
   list     Show recorded entries
   clear    Delete all recorded entries
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -h, --help     show help (default: false)
   -v, --version  print the version (default: false)
```

### `$ worklog add`

Add entries to the log.

```console
$ worklog add -im "got the project packaged" -c feature
2022-11-20 11:39:26	z5DWk2OVR	*  [feature]	'got the project packaged'

$ worklog add --help
NAME:
   worklog add - Add entries to the log

USAGE:
   worklog add [command options]

OPTIONS:
   -m MSG, --message MSG   Use the given MSG as the entry body
   -c TAG, --category TAG  Choose a category for the entry. TAG must be one of: [bug|feature|fix|meeting|note|refactor] (default: "note")
   -i, --important         Mark/flag the entry as important (default: false)
   -h, --help              show help (default: false)
```

### `$ worklog list`

Show recorded entries. Current idea is for output to look similar to `git log`.

```console
$ worklog list
2022-11-20 11:39:26	z5DWk2OVR	*  [feature]	'got the project packaged'
2022-11-19 19:01:50	ul-RHcdVR	   [refactor]	'starting to see worklog-package structure'
2022-11-19 19:01:12	Vu2eScO4g	*  [fix]	'got clear working finally'
2022-11-19 19:00:44	uDa3ScO4g	*  [feature]	'testing none defaults'
2022-11-19 19:00:35	zGd3S5OVR	   [note]	'adding entries from cli'

$ worklog list --help
NAME:
   worklog list - Show recorded entries

USAGE:
   worklog list [command options]

OPTIONS:
   -a DATE, --after DATE, --since DATE   Show entries newer than given DATE.
   -b DATE, --before DATE, --until DATE  Show entries older than given DATE.
   -h, --help                            show help (default: false)
```

## TODOs

- [x] centralize project once installed (e.g. `$HOME/.worklog`) to prevent `worklog.db` being created anywhere command is ran
- [x] `worklog list` filtering:
  - [x] filter by categor(y|ies)
  - [x] filter by date (`--after=DATE`, `--before=DATE`)
- [ ] validate `--after`, `--before` dates are ISO8601 instead of simply checking string length
