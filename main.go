package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pquerna/otp/totp"
	"gopkg.in/urfave/cli.v1"
	"net/http"
	"net/url"
	"os"
	"time"
	"gopkg.in/urfave/cli.v1/altsrc"
	"github.com/skratchdot/open-golang/open"
	"github.com/mitchellh/go-homedir"
)

const appName = "awsweb"
const appVersion = "1.0.0"
const unexpandedYmlPath = "~/.awsweb.yml"

func yamlProvider() (altsrc.InputSourceContext, error) {
	path, _ := homedir.Expand(unexpandedYmlPath)
	return altsrc.NewYamlSourceFromFile(path)
}

func main() {

	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Usage = "AWS web console shortcut tool"

	flags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{Name: "mfa-secret"}),
	}

	app.Action = func(c *cli.Context) {
		mfa_secret := c.String("mfa-secret")
		profile := c.Args().Get(0)
		login(profile, mfa_secret)
	}

	ymlPath, _ := homedir.Expand(unexpandedYmlPath)
	if _, err := os.Stat(ymlPath); err == nil {
		app.Before = altsrc.InitInputSource(flags, yamlProvider)
	}

	app.Flags = flags

	app.Run(os.Args)

}

type SigninTokenResponse struct {
	SigninToken string
}

func getJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}

	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func login(profile string, mfa_secret string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 profile,
		AssumeRoleTokenProvider: func() (string, error) {
			return totp.GenerateCode(mfa_secret, time.Now())
		},
	}))

	creds, _ := sess.Config.Credentials.Get()

	sessionJsonMap := map[string]string{
		"sessionId":    creds.AccessKeyID,
		"sessionKey":   creds.SecretAccessKey,
		"sessionToken": creds.SessionToken,
	}

	sessionJson, _ := json.Marshal(sessionJsonMap)
	sessionJsonEscaped := url.QueryEscape(string(sessionJson))

	getSigninTokenUrl := "https://signin.aws.amazon.com/federation?Action=getSigninToken&SessionType=json&Session=" + sessionJsonEscaped
	signinTokenResponse := new(SigninTokenResponse)
	getJson(getSigninTokenUrl, signinTokenResponse)
	escapedSigninToken := url.QueryEscape(signinTokenResponse.SigninToken)

	region := sess.Config.Region
	if len(*region) == 0 {
		defaultRegion := "us-east-1"
		region = &defaultRegion
	}

	destinationUrl := "https://" + *region + ".console.aws.amazon.com/"
	loginUrl := "https://signin.aws.amazon.com/federation?Action=login&SigninToken=" + escapedSigninToken + "&Destination=" + destinationUrl
	open.Start(loginUrl)
}
