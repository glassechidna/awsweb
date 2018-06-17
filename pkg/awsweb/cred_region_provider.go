package awsweb

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/glassechidna/awscredcache"
	"github.com/pquerna/otp/totp"
	"os"
	"time"
)

func GetProvider(profile string) CredRegionProvider {
	if len(profile) == 0 {
		return &NoopProvider{}
	}
	p := awscredcache.NewAwsCacheCredProvider(profile)
	p.MfaCodeProvider = func(mfaSecret string) (string, error) {
		if len(mfaSecret) == 0 {
			return stscreds.StdinTokenProvider()
		} else {
			return totp.GenerateCode(mfaSecret, time.Now())
		}
	}
	return p
}

type CredRegionProvider interface {
	Retrieve() (credentials.Value, error)
	IsExpired() bool
	Region() string
}

type NoopProvider struct{}

func (p *NoopProvider) Retrieve() (credentials.Value, error) {
	env := credentials.EnvProvider{}
	return env.Retrieve()
}

func (p *NoopProvider) IsExpired() bool {
	env := credentials.EnvProvider{}
	return env.IsExpired()
}

func (p *NoopProvider) Region() string {
	region := os.Getenv("AWS_REGION")
	if len(region) == 0 {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if len(region) == 0 {
		region = "us-east-1"
	}
	return region
}
