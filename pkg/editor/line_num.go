package editor

import "strconv"

// LineNumber encapsulates a line number in editing file that is given from the UI in string form
type LineNumber string

// ToInt converts this line number to integer
func (ln LineNumber) ToInt() (n int64, err error) {
	if ln == "" {
		return
	}
	n, err = strconv.ParseInt(string(ln), 10, 64)
	return
}
