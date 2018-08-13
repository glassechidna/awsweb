package browser

import (
	"github.com/DHowett/go-plist"
	"github.com/mitchellh/go-homedir"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func cliString(inp ...string) string {
	cmd := exec.Command(inp[0], inp[1:]...)
	output, _ := cmd.CombinedOutput()
	return string(output)
}

func plistPath() string {
	verstr := cliString("uname", "-r")
	parts := strings.Split(verstr, ".")
	ver, _ := strconv.ParseInt(parts[0], 10, 32)
	path := ""

	if ver >= 14 {
		path = "~/Library/Preferences/com.apple.LaunchServices/com.apple.launchservices.secure.plist"
	} else {
		path = "~/Library/Preferences/com.apple.LaunchServices.plist"
	}

	path, _ = homedir.Expand(path)
	return path
}

type plistStruct struct {
	LSHandlers []struct {
		LSHandlerRoleAll   string
		LSHandlerURLScheme string
	}
}

func httpBundleId() string {
	path := plistPath()
	file, _ := os.Open(path)
	decoder := plist.NewDecoder(file)

	contents := plistStruct{}
	_ = decoder.Decode(&contents)

	for _, handler := range contents.LSHandlers {
		if handler.LSHandlerURLScheme == "http" {
			return handler.LSHandlerRoleAll
		}
	}

	return ""
}

func browserForBundleId(bundleId string) (Browser, error) {
	switch strings.ToLower(bundleId) {
	case "org.mozilla.firefox":
		return &MozillaFirefox{}, nil
	//case "com.apple.safari":
	//	return &AppleSafari{}, nil
	case "com.google.chrome":
		return &GoogleChrome{}, nil
	default:
		return nil, nil
	}
}

func DefaultBrowser() (Browser, error) {
	bundleId := httpBundleId()
	appPath := cliString("mdfind", "kMDItemCFBundleIdentifier", "=", bundleId)
	appPath = strings.TrimSpace(appPath)

	return browserForBundleId(bundleId)
}
