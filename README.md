[![Go Report Card](https://goreportcard.com/badge/github.com/nchern/notelog)](https://goreportcard.com/report/github.com/nchern/notelog)
# notelog

Notelog is a simple organiser to maintain personal notes. The notes are kept on
the file system under a certain folder `NOTELOG_HOME`. You can control this location
via corresponding environment variable.
Default location for your notes is `$HOME/notes/personal`.

Notelog seamlessly integrates with your favorite text editor(tested with VIM).

Vim plugin is [available](vim/README.md)

## Installation
```bash
make install
```

## Examples and features

```bash
# Opens "my-note" in editor. Creates a note if the note does not exist
$ notelog my-note

# Adds a line "foo bar" to "my-note" directly from command line
$ notelog my-note foo bar

# Archive a note: the note becomes unavailable in this collection
# for direct edits, also is not visible for search by default
$ notelog do archive my-note

# Lists all notes in collection sorted by last modified date
$ notelog do list --by-date

# Prints help
$ notelog help
```

## Advanced features

### Search

Notelog supports search over notes collection.

```bash
# Search all lines that contain "foo" over all notes
$ notelog do search --interactive foo
1. noteA:1:foo bar
2. noteB:10:hello foo

# Open the second search result note from previous search in editor
$ notelog do search-browse 2
```

### Integrate with git

Since your notes are just a bunch of text files in a subtree on a file system,
you can add your note collection to a git repo and have them version controlled.
Notelog has a couple of commands to simplify this task though currently this functionality is pretty basic.

```bash
# Initialize a repo in the current notes collection
$ notelog do init-repo

# Synchronize the repo with origin: add all the changes, commit, update from the origin and push
$ notelog do sync
```

**Please note that the "origin" remote is not automatically added**.
As of now you have to add it manually using standard `git remote add ...`

## Roadmap
 - [ ] add more description for `contrib` scripts
 - [ ] (?) in-note macros:
   - [ ] when adding lines, format them according to a given template
 - [ ] (?) multiple temporary drafts - when open a draft, this should not be the same file every time
 - [ ] have notes on `.md` format and not only in `.org`
 - [ ] archive: a note can:
   - [X] be put into archive, so it will not stay in the main note list
         Current behavior: no search in the archive. Only through actual notes
   - [ ] enable search in archive?
   - [ ] be restored from the archive (eventually)
 - [ ] sub-notes: notes that exist only in a context of a main note (?)
   - example notelog do subnote <notename> <sub-notename>
 - [ ] attachments to notes (?)
   - [ ] notelog do attach <notename> <filepath> - puts <filepath> into note directory
   - [ ] notelog do attach-open <notename> <attach-name> - opens attach
   - [ ] integrate with search?
 - [ ] note templates (?)
 - [ ] cross-linking: you can fetch all the references from other notes to a given note
   - [ ] embed cross links in notes?
 - [X] add more examples, hints, use cases and script recipes
 - [X] integration with fzf: search results
 - [X] vim plugin
   - [X] MVP
 - [X] create dir structure in one go during init phase. Consider fixing existing incomplete structure.
   - [X] create .notelog at least
 - [X] migrate to cobra subcommand parser
 - [X] create and populate .gitignore if NOTELOG_HOME is considered as a git repo
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
