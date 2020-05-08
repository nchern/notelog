package env

import (
	"fmt"
	"os"
	"strings"
)

var (
	settings = settingsBag{}
)

// Get returns value of env var with given name. If it's empty, returns defaultVal
func Get(name string, defaultVal string) string {
	res := defaultVal

	defer func() { settings[name] = res }()

	val := os.Getenv(name)
	if val == "" {
		return res
	}
	res = val
	return res
}

// Vars returns string dump of all env variables along with their values
func Vars() string {
	return settings.String()
}

type settingsBag map[string]string

func (b settingsBag) String() string {
	res := []string{}
	for k, v := range b {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(res, "\n")
}
