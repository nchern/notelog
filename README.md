[![Go Report Card](https://goreportcard.com/badge/github.com/nchern/notelog)](https://goreportcard.com/report/github.com/nchern/notelog)
# notelog

Notelog is a simple organiser to maintain personal notes.

It keeps notebase in your home directory as files under a single folder.
Integrates with your favourite text editor(tested with VIM).

## Installation
```bash
go get github.com/nchern/notelog/...
```

## Examples

```bash
# Opens "my-note" in editor
$ notelog my-note

# adds a line "foo bar" to "my-note" directly from command line
$ notelog my-note foo bar

# Prints help
$ notelog -h
```

## Roadmap
 - [x] search: add simple query language, e.g "foo bar -buzz" == "foo OR bar EXCLUDE buzz"
 - [x] add rm feature
 - [x] add rename feature
 - [ ] archive: note can:
   - [ ] be put into archive, so it will not stay in the main note list
   - [ ] restored from the archive (eventually)
 - [ ] cross-linking: you can fetch all the references from other notes to a given note
   - [ ] embed cross links in notes?
 - [x] search browsing: quick jump to search results from command line
 - [ ] cloud capabilities:
   - [x] via simple rsync
   - [ ] via git
 - [ ] vim plugin
   - [x] mvp
   - [ ] with ftdetect to '.note'?
 - [ ] note templates
