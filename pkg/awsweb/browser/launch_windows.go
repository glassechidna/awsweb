package browser

import (
	"os/exec"
	"strings"
)

func winCmd(name, url string, extraArgs []string) {
	args := []string{"cmd", "/c", "start", name, strings.Replace(url, "&", "^&", -1)}
	args = append(args, extraArgs...)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Start()
}

func (b *GoogleChrome) Launch(url, profile string) {
	winCmd("chrome", url, b.openArgs(profile))
}

func (b *MozillaFirefox) Launch(url, profile string) {
	winCmd("firefox", url, b.openArgs(profile))
}
