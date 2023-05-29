# Bash autocompletion for notelog. Completes notes
complete -C 'notelog autocomplete' notelog
complete -C 'notelog autocomplete --note-names' notelog-cat
complete -C 'notelog autocomplete --note-names' nlg-batch

complete -W "\`notes-select -ls\`" notes-select
