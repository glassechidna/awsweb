package browser

import (
	"path"
	"os"
)

type MozillaFirefox struct {}

func (b *MozillaFirefox) Name() string {
	return MozillaFirefoxName
}

func (b *MozillaFirefox) openArgs(profile string) []string {
	profileDir := path.Join(os.TempDir(), "awsweb-firefox-" + profile)
	return []string{"-no-remote", "-profile", profileDir}
}

