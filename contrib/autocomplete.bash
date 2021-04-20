# Bash autocompletion for notelog. Completes notes
complete -C 'notelog -c autocomplete' notelog
complete -W "\`notes-select -ls\`" notes-select
