" Notelog plugin
" currently integrated with OrgMode by sharing the same file extention: org

autocmd FileType org set grepprg=notelog\ do\ search

" Notes autocomplete for Notelog
fun NotesList(A,L,P)
    let res = system("notelog list")
    if v:shell_error != 0
       echoerr l:red
       return ''
    endif
    return l:res
endfun

" Returns full path to notes in Notelog
fun NotesFullPath(name)
    let res = system('notelog path ' . a:name)
    if v:shell_error != 0
       echoerr l:res
       return ''
    endif
    return l:res
endfun

" Executes notleg search
fun NotesDoSearch(terms)
   let cmd = 'sh -c "notelog search ' . escape(a:terms, '"') . ' | notermcolor"'
   cexpr system(l:cmd)
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

" NotesGoGet runs go get command on notelog util to install or update it
fun NotesGoGet(force_update)
    let up = a:force_update ? '-u' : ''
    execute ':!go get ' . l:up . ' github.com/nchern/notelog/...'
endfun

fun NotesArchive()
    w
    silent !notelog archive -f '%'
    bd
endfun

fun s:note_open(...)
    let l:name = a:0 > 0 ? a:1 : ''
    if l:name == ''
        :NLBrowse
        return
    endif
    execute ':e ' . NotesFullPath(l:name)
endfun

" Creates a new note with Notelog
command! -nargs=1 -complete=custom,NotesList NLNew execute ':silent !notelog touch <f-args>' | execute ':e ' NotesFullPath(<f-args>)

" Opens an existing note with Notelog
command! -nargs=? -complete=custom,NotesList NLOpen call s:note_open(<f-args>)
" Opens Notelog's scratchpad
command! -nargs=0 NLOpenScratch execute ':e ' NotesFullPath("")

" Adds quick record to existing note with Notelog
command! -nargs=+ -complete=custom,NotesList NLQuickLog :!notelog edit <args>

" Alias for NLQuickLog to see what works better
command! -nargs=+ -complete=custom,NotesList NLLog :NLQuickLog <args>

" Replaces selection with contents of a given note
command! -range -nargs=1 -complete=custom,NotesList NLPasteNote :<line1>,<line2>!notelog print <q-args>

" Runs fzf to search through note names and then opens selected note
" Requires FZF plugin
command! NLBrowse :call fzf#run({'source': 'notelog ls', 'sink': 'NLOpen'})

" Sorts checkboxed items with notelog
autocmd FileType org,markdown command! -range=% NLSortCheckList :<line1>,<line2>!notelog org-checklist-sort

" Inserts link to another note under the cursor pos
autocmd FileType org,markdown command! -nargs=1 -complete=custom,NotesList NLLinkNote :call NotesDoInsertLink(<f-args>)

" Performs search
autocmd FileType org,markdown command! -nargs=1 NLSearch :call NotesDoSearch(<f-args>)

" Syncs notes
autocmd FileType org,markdown command! -nargs=? NLSync :!notelog sync <q-args>

" Installs notelog binaries
autocmd FileType org,markdown command! NLInstallBinaries :call NotesGoGet(0)
" Updates notelog binaries
autocmd FileType org,markdown command! NLUpdateBinaries :call NotesGoGet(1)

" Calls an external command to search info on a person
autocmd FileType org,markdown nnoremap <Localleader>gd :call NotesBrowseGroupDirectory()<CR>

" Archives current note
autocmd FileType org,markdown command! NLArchive : call NotesArchive()

" Deletes the current note
autocmd FileType org,markdown command! NLDelete :w | execute "!notelog rm -f '%'" | bd
