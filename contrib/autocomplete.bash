# Bash autocompletion for notelog. Completes notes
complete -C 'notelog do autocomplete' notelog
complete -W "\`notes-select -ls\`" notes-select
