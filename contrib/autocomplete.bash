# Bash autocompletion for notelog. Completes notes
complete -C 'notelog autocomplete' notelog
complete -C 'notelog autocomplete' notelog-cat
complete -C 'notelog autocomplete' nlg-batch

complete -W "\`notes-select -ls\`" notes-select
