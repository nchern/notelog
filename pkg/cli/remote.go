package cli

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/nchern/notelog/pkg/editor"
	"github.com/nchern/notelog/pkg/note"
	"github.com/nchern/notelog/pkg/remote"
)

func handleNoRemoteConfig(err error) error {
	if os.IsNotExist(err) {
		configPath := note.NewList().MetadataFilename(remote.ConfigName)
		if err := os.MkdirAll(path.Dir(configPath), editor.DefaultDirPerms); err != nil {
			return err
		}
		if err := ioutil.WriteFile(configPath, []byte(remote.DefaultConfig), editor.DefaultFilePerms); err != nil {
			return err
		}

		return editor.Shellout(configPath).Run()
	}
	return err
}
