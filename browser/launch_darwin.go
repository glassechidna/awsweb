package browser

import (
	"os/exec"
)

func macCmd(url string, bundleId string, extraArgs []string) {
	args := []string{"open", "-n", "-b", bundleId, url, "--args"}
	args = append(args, extraArgs...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Start()
}

func (b *GoogleChrome) Launch(url, profile string) {
	macCmd(url, "com.google.Chrome", b.openArgs(profile))
}

func (b *MozillaFirefox) Launch(url, profile string) {
	macCmd(url, "org.mozilla.firefox", b.openArgs(profile))
}
