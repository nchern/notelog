package cli

import "strings"

type commandList []string

func (l commandList) String() string {
	return strings.Join(l, ", ")
}

func c(s string) string {
	commands = append(commands, s)
	return s
}
