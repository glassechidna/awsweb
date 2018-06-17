package awsweb

import (
	"net/http"
	"time"
	"encoding/json"
	"net/url"
)

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

func GetLoginUrl(provider CredRegionProvider) string {
	creds, _ := provider.Retrieve()

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

	destinationUrl := "https://" + provider.Region() + ".console.aws.amazon.com/"
	loginUrl := "https://signin.aws.amazon.com/federation?Action=login&SigninToken=" + escapedSigninToken + "&Destination=" + destinationUrl

	return loginUrl
}
