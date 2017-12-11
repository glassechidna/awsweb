package browser

const MozillaFirefoxName = "firefox"
const GoogleChromeName = "chrome"
//const AppleSafariName = "AppleSafari"
//const MicrosoftIEName = "MicrosoftIE"
//const MicrosoftEdgeName = "MicrosoftEdge"

type Browser interface {
	Name() string
	Launch(url, profile string)
}

