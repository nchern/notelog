if exists("b:note_syntax")
    finish
endif


" highlights deadlines(the usual usecase is todo item): <by Mon, 3d>
syntax      match               notelogItemDeadline /\v\<(by|till|until)\s.*\>/ containedin=ALL
highlight   notelogItemDeadline ctermfg=DarkRed term=bold cterm=bold gui=bold


" highlights warnings: !Warning!
syntax      match               notelogWarning /\v!.*!/ containedin=ALL
highlight   notelogWarning      ctermfg=Red term=bold,underline cterm=bold,underline


" highlights a block with backgound color: % this block is highlighted %
" mind the spaces BEFORE and AFTER %
syntax      match           notelogMarker  /\v\% [0-9a-zA-Z _:;\%]* \%/ containedin=ALL
highlight   notelogMarker   ctermbg=Yellow ctermfg=Black term=bold cterm=bold gui=bold


" highlight references to people: @John.Doe
syntax      match           notelogPerson   /\v\@[a-zA-Z.]+/ containedin=ALL
highlight   notelogPerson   ctermfg=Blue    term=bold cterm=bold gui=bold


let b:note_syntax = "notelog"
