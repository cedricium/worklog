# worklog

An opinionated note-taking tool for the developer's day-to-day.

## Idea

After working at Hellosaurus for two years, I can only remember the big features I built. But what about the tricky bugs solved, legacy code rewrites, or instrumental meetings attended? Those moments, arguably just as important as feature work, stuck somewhere in memory though seemingly inaccessible forever.

Wouldn't it be nice if there was a way to simply and quickly document one's day-to-day while at work? Something akin to the captain's log in Star Trek, with a way for easy perusal at a later time.

Inspired by `git`, I envision a CLI tool and package for documenting and parsing: tasks, notes, meetings, bug fixes, etc., automatically grouped by day in order to maintain a record of the day-to-day on the job.

## Usage

This is a work-in-progress and currently being used to drive development. Very likely to change in the future.

### `$ worklog`

```console
$ worklog --help
usage: worklog [--help] <command> [<options>]

COMMANDS
  add   add entries to the log
  list  show recorded entries
```

### `$ worklog add`

Add entries to the log.

```console
$ worklog add --help

EXAMPLES
$ worklog add --c fix -it "solved the >1000 RPCs causing Twitter being slow"
$ worklog add -t "struggling with GraphQL" -d 2022-06-05

OPTIONS
  -c, --category={bug|fix|feature|meeting|note|refactor}
    Choose a category for the entry. Variants are as follows:

    bug
      For documenting currently-unresolved issues within the codebase.

    feature
      New development work.

    fix
      Work that solves previous issues.

    meeting
      Helpful for remembering key takeaways/outcomes.

    note (default)
      General-purpose entry, useful for note-taking.

    refactor
      Legacy code rewrites or code improvements/enhancements.

  -d, --date=<date>
    Override the default date (default being NOW) for the entry being recorded. Needs to be an ISO 8601-like date string (e.g., "2022-11-04").

  -i, --important
    Mark/flag the entry as important. Useful if entry is a known "big-deal" and warrants special attention when reviewing entries later.

  -t, --text=<txt>
    Use the given <txt> as the entry body. No limit on length for now, though ideally should be kept short (think tweet size).
```

### `$ worklog add`

Show recorded entries. Current idea is for output to look similar to `git log`.

```console
$ worklog list --help

EXAMPLES
  $ worklog list # show all entries from newest to oldest
  $ worklog list --before="2022-01-01" # show entries recorded before January 1, 2022.
  $ worklog list --filter=FGR # show only entries categorized as Fixes (F), Features (G), and Refactors (R)

OPTIONS
  -f, --filter=[(B|F|G|M|N|R|I)…[*]]
    Select only entries that are categorized:
      bugs (B)
      fixes (F)
      features (G)
      meetings (M)
      note (N)
      refactors (R)
      important (I)

  -a, --after=<date>
  -s, --since=<date>
    Show entries newer than the given date. Requires ISO 8601-like date string.

  -b, --before=<date>
  -u, --until=<date>
    Show entries older than the given date. Requires ISO 8601-like date string.
```
