package searcher

import (
	"regexp"
)

// Taken from here: https://github.com/acarl005/stripansi/blob/master/stripansi.go
const ansiTermColors = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var (
	termEscapeSequence = regexp.MustCompile(ansiTermColors)
)
