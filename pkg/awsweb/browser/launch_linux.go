package browser

import (
	"os/exec"
)

func macCmd(url string, name string, extraArgs []string) {
	args := []string{name, url}
	args = append(args, extraArgs...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Start()
}

func (b *GoogleChrome) Launch(url, profile string) {
	macCmd(url, "google-chrome", b.openArgs(profile))
}

func (b *MozillaFirefox) Launch(url, profile string) {
	macCmd(url, "firefox", b.openArgs(profile))
}
