" Notelog plugin
" currently integrated with OrgMode by sharing the same file extention: org

autocmd FileType org set grepprg=notelog\ -c\ search

" Notes autocomplete for Notelog
:fun NotesList(A,L,P)
:    let res = system("notelog -c list")
:    if v:shell_error != 0
:       echoerr l:red
:       return ''
:    endif
:    return l:res
:endfun

" Returns full path to notes in Notelog
:fun NotesFullPath(name)
:    let res = system('notelog -c path ' . a:name)
:    if v:shell_error != 0
:       echoerr l:res
:       return ''
:    endif
:    return l:res
:endfun

" Executes notleg search
:fun NotesDoSearch(terms)
:   execute 'grep ' . a:terms
:   copen
:endfun

" Inserts org-mode link to another note under the cursor
:func NotesDoInsertLink(name)
:   try
:       let path = NotesFullPath(a:name)
:       let link = '[[' . l:path . '][' . a:name . ']]'
:       execute "normal! i" . l:link . "\<Esc>"
:   catch
:       echoerr v:exception
:   endtry
:endfun

" Opens an existing note with Notelog
autocmd FileType org command! -nargs=1 -complete=custom,NotesList NotesOpen execute ':e ' NotesFullPath(<f-args>)

" Sorts todos with notelog
autocmd FileType org command! -range=% NotesSortTodos :<line1>,<line2>!notelog -c sort-todos

" Inserts link to another note under the cursor pos
autocmd FileType org command! -nargs=1 -complete=custom,NotesList NotesInsertLink :call NotesDoInsertLink(<f-args>)

" Performs search
autocmd FileType org command! -nargs=1 NotesSearch :call NotesDoSearch(<f-args>)

" Installs notelog binaries
autocmd FileType org command! NoteInstallBinaries :!go get github.com/nchern/notelog/...
" Updates notelog binaries
autocmd FileType org command! NoteUpdateBinaries :!go get -u github.com/nchern/notelog/...
