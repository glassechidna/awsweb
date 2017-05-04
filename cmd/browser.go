// Copyright © 2017 Aidan Steele <aidan.steele@glassechidna.com.au>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"time"
	"encoding/json"
	"net/url"
	"github.com/skratchdot/open-golang/open"
	"net/http"
	"github.com/glassechidna/awsweb/shared"
	"github.com/spf13/viper"
)

func init() {
	// browserCmd represents the browser command
	var browserCmd = &cobra.Command{
		Use:   "browser",
		Short: "Open browser window at AWS web console",
		Long: `Assumes the given role and logs you into the AWS web console
		in the role's default region.`,
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			mfaSecret := viper.GetString("mfa-secret")
			doBrowser(profile, mfaSecret)
		},
	}

	RootCmd.AddCommand(browserCmd)
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

func doBrowser(profile string, mfaSecret string) {
	creds, region := shared.GetCreds(profile, mfaSecret)

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

	destinationUrl := "https://" + region + ".console.aws.amazon.com/"
	loginUrl := "https://signin.aws.amazon.com/federation?Action=login&SigninToken=" + escapedSigninToken + "&Destination=" + destinationUrl
	open.Start(loginUrl)
}
