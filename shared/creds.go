package shared

import (
	"time"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pquerna/otp/totp"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"log"
)

func GetCreds(profile string, mfaSecret string) (credentials.Value, string) {
	stscreds.DefaultDuration = 3600 * time.Second // TODO: issue #11 make this configurable

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 profile,
		AssumeRoleTokenProvider: func() (string, error) {
			return totp.GenerateCode(mfaSecret, time.Now())
		},
	}))

	creds, err := sess.Config.Credentials.Get()

	if err != nil {
		log.Panicf(err.Error())
	}

	region := sess.Config.Region
	if len(*region) == 0 {
		defaultRegion := "us-east-1"
		region = &defaultRegion
	}

	return creds, *region
}
