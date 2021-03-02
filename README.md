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
 - [ ] in-note macros:
   - [ ] when adding lines, format them according to a given template
 - [ ] multiple temporary drafts - when open a draft, this should not be the same file every time
 - [ ] have notes on `.md` format and not only in `.org`
 - [ ] create dir structure in one go during init phase. Consider fixing existing incomplete structure.
 - [ ] archive: a note can:
   - [X] be put into archive, so it will not stay in the main note list
         the curren search behaviour depends on a search command. E.g. `sift` does not search in .directories
         question still is: should it be available or not for the search?
   - [ ] be restored from the archive (eventually)
 - [ ] sub-notes: notes that exist only in a context of a main note (?)
   - example notelog -c subnote <notename> <sub-notename>
 - [ ] attachments to notes (?)
   - [ ] notelog -c attach <notename> <filepath> - puts <filepath> into note directory
   - [ ] notelog -c attach-open <notename> <attach-name> - opens attach
   - [ ] integrate with search?
 - [ ] note templates (?)
 - [ ] cross-linking: you can fetch all the references from other notes to a given note
   - [ ] embed cross links in notes?
 - [ ] migrate to cobra subcommand parser (?)
 - [ ] vim plugin
   - [x] MVP
   - [ ] with ftdetect to '.note' (???)
 - [x] smart bash auto-completion for subcommands
 - [x] search capabilities
   - [x] search browsing: quick jump to search results from command line
   - [x] search aggregations:  - group by files(—titles-only)  (? - do we need to return all lines or just docs?)
   - [x] search over note names
 - [x] search: add simple query language, e.g. "foo bar -buzz" == "foo OR bar EXCLUDE buzz"
 - [x] add rm feature
 - [x] add rename feature
 - [x] cloud capabilities:
   - [x] via simple rsync
   - [x] via git
