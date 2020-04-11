Vim plugin to organise your notes.

The best way to integrate `notelog` with vim.

As of now, this plugin works nicely together with [vim-orgmode plugin](https://github.com/jceb/vim-orgmode) . It will work even without `vim-orgmode` but this plugin is highly recommended to maintain features like todos, lists, handling links between notes etc.
Notes extension is anyway set to `.org`

## Installation

First, install [vim-orgmode](https://github.com/jceb/vim-orgmode). This is highly recommended.

### vim-plug
Add the following line to your .vimrc:

```vimrc
Plug 'nchern/notelog', { 'rtp': 'vim' }
```

And then update your packages by running `:PlugInstall`.

`notelog` binary should be already [installed](https://github.com/nchern/notelog#installation) or you can install it by running `:NoteInstallBinaries` command in the open `org` buffer.
