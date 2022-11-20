package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/muesli/coral"
	"github.com/nchern/notelog/pkg/note"
)

const (
	scratchpadName = ".scratchpad"
)

var (
	notes = note.NewList()

	rootCmd = &coral.Command{
		Use:   "notelog",
		Short: "Efficient CLI personal note manager",
		Args:  coral.MinimumNArgs(0),

		SilenceUsage:  true,
		SilenceErrors: true,
	}

	defaultHelp = rootCmd.HelpFunc()

	mainConfPath = filepath.Join(os.Getenv("HOME"), note.DotNotelogDir+"-config.toml")
	conf         = Config{
		NoteFormat: defaultFormat,
		SkipLines:  defaultSkipLines,
	}
)

// Config represents a configuration of this app
type Config struct {
	NoteFormat          string `toml:"note_format"`
	SkipLines           uint   `toml:"skip_lines"`
	SkipLinesAfterMatch string `toml:"skip_lines_after_match"`
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *coral.Command, s []string) {
		defaultHelp(cmd, s)
	})
}

// Execute is an entry point of CLI
func Execute() error {
	const defaultDirPerms = 0700

	if err := loadConfig(); err != nil && !os.IsNotExist(err) {
		// TODO: possibly to main app log?
		log.Printf("WARN loadConfig failed: %s", err)
	}

	return rootCmd.Execute()
}

func loadConfig() error {
	_, err := toml.DecodeFile(mainConfPath, &conf)
	if err != nil {
		return err
	}

	ft, err := note.ParseFormat(conf.NoteFormat)
	if err != nil || ft == note.Unknown {
		log.Printf("WARN loadConfig: %s", err)
		conf.NoteFormat = defaultFormat
	}
	return nil
}

func parseNoteName(name string) (string, error) {
	// FIXME: this is a hack, need more elegant solution than double if
	if name == scratchpadName {
		return name, nil
	}
	name = strings.TrimSpace(name)
	if err := validateNoteName(name); err != nil {
		return "", err
	}
	return name, nil
}

func noteNameFromArgs(args []string) string {
	if len(args) < 1 {
		return scratchpadName
	}
	return args[0]
}
