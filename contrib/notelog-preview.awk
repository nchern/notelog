#!/usr/bin/awk -f

BEGIN {
    LN=1;
    CTX_LINES=10;
    name_var = "NOTE_NAME"
    if (name_var in ENVIRON) {
        name = ENVIRON[name_var]
        if (split(name, toks, ":") > 1) {
            LN = int(toks[2])
        }
        name = toks[1]
    }
    print "\033[38;5;244m" name ":\033[0m\n"
}

NR >= LN-CTX_LINES && NR <= LN+CTX_LINES && NR != LN { print "\033[38;5;244m" NR "\033[0m " $0 }

# print selected line in bold; \033[1mBOLD TEXT\033[0m
NR == LN { print NR " \033[1m" $0 "\033[0m"}
