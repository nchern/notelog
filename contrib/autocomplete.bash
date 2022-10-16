# Bash autocompletion for notelog. Completes notes
complete -C 'notelog do autocomplete' notelog
complete -C 'notelog do autocomplete' notelog-cat
complete -C 'notelog do autocomplete' nlg-batch

complete -W "\`notes-select -ls\`" notes-select
