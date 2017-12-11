package browser

import (
	"golang.org/x/sys/windows/registry"
)

func browserForRegistry(name string) Browser {
	prefix := name[0:5] // firefox has some crap on the end

	switch prefix {
	case "Chrom":
		return &GoogleChrome{}
	case "Firef":
		return &MozillaFirefox{}
	//case "IE.HT":
	//	return MicrosoftIE
	//case "AppXq":
	//	return MicrosoftEdge
	default:
		return nil
	}
}

func DefaultBrowser() (Browser, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\Shell\Associations\URLAssociations\http\UserChoice`, registry.QUERY_VALUE)
	if err != nil { return nil, err }
	defer k.Close()

	s, _, err := k.GetStringValue("ProgId")
	if err != nil { return nil, err }

	return browserForRegistry(s), nil
}
