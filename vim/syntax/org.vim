if exists("b:note_syntax")
    finish
endif


" highlights deadlines(the usual usecase is todo item): '<by Mon, 3d>'
syntax match notelogItemDeadline /\v\<by\s.*\>/ containedin=ALL
highlight notelogItemDeadline ctermfg=DarkRed term=bold cterm=bold gui=bold


let b:note_syntax = "notelog"
