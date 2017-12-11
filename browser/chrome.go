package browser

import (
	"os"
	"path"
)

type GoogleChrome struct {}

func (b *GoogleChrome) Name() string {
	return GoogleChromeName
}

func (b *GoogleChrome) openArgs(profile string) []string {
	userDataDir := path.Join(os.TempDir(), "awsweb-chrome-" + profile)
	userDataDirFlag := "--user-data-dir=" + userDataDir
	return []string{userDataDirFlag, "--no-first-run"}
}
