" highlights deadlines(the usual usecase is todo item): '(by Mon, 3d)'
syntax match notelogItemDeadline /\v\(by\s.*\)/ containedin=ALL

highlight notelogItemDeadline ctermfg=Red term=bold cterm=bold gui=bold     
