# worklog

An opinionated note-taking tool for the developer's day-to-day.

## Idea

After working at Hellosaurus for two years, I can really only remember the big features I built. But what about the tricky bugs solved, legacy code rewrites, or instrumental meetings attended? Those moments, arguably just as important as feature work, stuck somewhere in memory though seemingly inaccessible forever.

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
  $ worklog add -m "struggling with GraphQL"
  $ worklog add -c fix -im "solved the >1000 RPCs causing Twitter being slow"

OPTIONS
  -m, --message=<MSG>
    Use the given <MSG> as the entry body. No limit on length for now, though
    ideally should be kept short (think tweet size).

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

  -i, --important
    Mark/flag the entry as important. Useful if entry is a known "big-deal" and
    warrants special attention when reviewing entries later.
```

### `$ worklog list`

Show recorded entries. Current idea is for output to look similar to `git log`.

```console
$ worklog list --help

EXAMPLES
  $ worklog list
  $ worklog list --before="2022-01-01"
  $ worklog list --filter=FGR --after="2022-06-05"

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
