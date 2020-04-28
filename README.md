# notelog

Notelog is a simple organiser to maintain personal notes.

It keeps notebase in your home directory as files under a single folder.
Integrates with your favourite text editor(tested with VIM).

### Installation
```bash
go get github.com/nchern/notelog/...
```

### Examples

```bash
# Opens "my-note" in editor
$ notelog my-note

# adds a line "foo bar" to "my-note" directly from command line
$ notelog my-note foo bar

# Prints help
$ notelog -h
```

## Roadmap

 - [ ] note templates
 - [ ] cross-linking: you can fetch all the references from other notes to a given note
   - [ ] embed cross links in notes?
 - [ ] search browsing: quick jump to search results from command line
 - [ ] vim plugin
   - [x] mvp
   - [ ] with ftdetect to '.note'?
