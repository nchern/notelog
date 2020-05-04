if exists("b:note_syntax")
    finish
endif


" highlights deadlines(the usual usecase is todo item): <by Mon, 3d>
syntax match notelogItemDeadline /\v\<by\s.*\>/ containedin=ALL
highlight notelogItemDeadline ctermfg=DarkRed term=bold cterm=bold gui=bold


" highlights warnings: !Warning!
syntax match notelogWarning /\v!.*!/
highlight notelogWarning    ctermfg=Red term=bold,underline cterm=bold,underline


" highlights a block with backgound color: %this block draws attention%
syntax match notelogMarker  /\v\%.*\%/
highlight notelogMarker     ctermbg=LightYellow ctermfg=DarkBlue


let b:note_syntax = "notelog"
