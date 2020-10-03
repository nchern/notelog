" Notelog plugin
" currently integrated with OrgMode by sharing the same file extention: org

autocmd FileType org set grepprg=notelog\ -c\ search

" Notes autocomplete for Notelog
fun NotesList(A,L,P)
    let res = system("notelog -c list")
    if v:shell_error != 0
       echoerr l:red
       return ''
    endif
    return l:res
endfun

" Returns full path to notes in Notelog
fun NotesFullPath(name)
    let res = system('notelog -c path ' . a:name)
    if v:shell_error != 0
       echoerr l:res
       return ''
    endif
    return l:res
endfun

" Executes notleg search
fun NotesDoSearch(terms)
   execute 'grep ' . a:terms
   copen
endfun

" Inserts org-mode link to another note under the cursor
func NotesDoInsertLink(name)
   try
       let path = NotesFullPath(a:name)
       let link = '[[' . l:path . '][' . a:name . ']]'
       execute "normal! i" . l:link . "\<Esc>"
   catch
       echoerr v:exception
   endtry
endfun

" NotesBrowseGroupDirectory calls an external command if a word under cursor
" is a reference to a person. This word is passed to the command. This command
" can call an external program to browse info for this person
fun NotesBrowseGroupDirectory()
    let person_class = 'notelogPerson'

    let is_person = 0

    for id in synstack(line('.'), col('.'))
        if (synIDattr(id, 'name') == l:person_class)
            let l:is_person = 1
            break
        endif
    endfor

    if !l:is_person
        return
    endif
    let name = substitute(trim(expand('<cWORD>'), '@'), '\.', ' ', 'g')

    if !exists('g:nl_gd_browse_command')
        let g:nl_gd_browse_command = 'nl_gd_browse'
    endif

    execute ':silent !' . g:nl_gd_browse_command . ' "' . l:name . '"'
endfun


" Calls an external command to search info on a person
autocmd FileType org nnoremap <Localleader>gd :call NotesBrowseGroupDirectory()<CR>

" Opens an existing note with Notelog
command! -nargs=1 -complete=custom,NotesList NLOpen execute ':e ' NotesFullPath(<f-args>)

" Sorts todos with notelog
autocmd FileType org command! -range=% NLSortTodos :<line1>,<line2>!notelog -c sort-todos

" Inserts link to another note under the cursor pos
autocmd FileType org command! -nargs=1 -complete=custom,NotesList NLInsertLink :call NotesDoInsertLink(<f-args>)

" Performs search
autocmd FileType org command! -nargs=1 NLSearch :call NotesDoSearch(<f-args>)

" Installs notelog binaries
autocmd FileType org command! NLInstallBinaries :!go get github.com/nchern/notelog/...
" Updates notelog binaries
autocmd FileType org command! NLUpdateBinaries :!go get -u github.com/nchern/notelog/...
