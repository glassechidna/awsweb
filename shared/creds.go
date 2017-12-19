package shared

import (
	"time"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pquerna/otp/totp"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/glassechidna/awsweb/sneakyvendor/aws-shared-defaults"
	"github.com/go-ini/ini"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"github.com/pkg/errors"
)

type awsConfigFiles struct {
	cfg *ini.File
	cred *ini.File
}

type ProfileConfig struct {
	Name string
	Region string
	Credentials *credentials.Credentials
}

func profileConfig(cfg awsConfigFiles, profile string) (*ProfileConfig, error) {
	section, err := cfg.cfg.GetSection(fmt.Sprintf("profile %s", profile))
	if err != nil {
		section, err = cfg.cfg.GetSection(profile)
		if err != nil { return nil, err }
	}

	region := section.Key("region").String()

	sourceProfile, err := section.GetKey("source_profile")
	hasSourceProfile := err == nil

	if hasSourceProfile {
		sourceConfig, err := profileConfig(cfg, sourceProfile.String())
		if err != nil { return nil, err }

		sourceRegion := sourceConfig.Region
		if len(region) == 0 {
			region = sourceRegion
		}

		roleArn := section.Key("role_arn").String()
		if len(roleArn) == 0 { return nil, errors.New("empty role arn") }

		roleCreds, err := roleCredentials(sourceConfig.Credentials, roleArn)
		if err != nil { return nil, err }

		return &ProfileConfig{
			Name:        profile,
			Region:      region,
			Credentials: roleCreds,
		}, nil
	} else {
		credsSection, err := cfg.cred.GetSection(profile)
		if err != nil { return nil, err }

		id := credsSection.Key("aws_access_key_id").String()
		secret := credsSection.Key("aws_secret_access_key").String()
		token := credsSection.Key("aws_session_token").String()

		if len(id) == 0 { return nil, errors.New("empty access key id") }
		if len(secret) == 0 { return nil, errors.New("empty secret access key") }
		creds := credentials.NewStaticCredentials(id, secret, token)

		mfaSerial := section.Key("mfa_serial").String()
		if len(mfaSerial) > 0 {
			mfaSecret := credsSection.Key("mfa_secret").String()
			creds, err = mfaAuthenticatedCredentials(creds, mfaSerial, mfaProvider(mfaSecret))
			if err != nil { return nil, err }
		}

		return &ProfileConfig{
			Name:        profile,
			Region:      region,
			Credentials: creds,
		}, nil
	}
}

func roleCredentials(sourceCreds *credentials.Credentials, roleArn string) (*credentials.Credentials, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: sourceCreds,
	}))

	api := sts.New(sess)

	roleSessionName := fmt.Sprintf("awsweb-%d", time.Now().Second())
	resp, err := api.AssumeRole(&sts.AssumeRoleInput{
		RoleArn: aws.String(roleArn),
		RoleSessionName: &roleSessionName,
	})
	if err != nil { return nil, err }

	c := resp.Credentials
	return credentials.NewStaticCredentials(*c.AccessKeyId, *c.SecretAccessKey, *c.SessionToken), nil
}

func mfaProvider(mfaSecret string) func() (string, error) {
	if len(mfaSecret) == 0 {
		return stscreds.StdinTokenProvider
	} else {
		return func() (string, error) {
			return totp.GenerateCode(mfaSecret, time.Now())
		}
	}
}

func loadConfig() awsConfigFiles {
	cfgIni, _ := ini.Load(shareddefaults.SharedConfigFilename())
	credIni, _ := ini.Load(shareddefaults.SharedCredentialsFilename())
	return awsConfigFiles{cfg: cfgIni, cred: credIni}
}

func GetCreds(profile string) ProfileConfig {
	if len(profile) == 0 {
		// in the special case that no profile name has been passed in,
		// it's because the user wants to use the current (i.e. stored
		// in the the env vars) session

		region := os.Getenv("AWS_REGION")
		if len(region) == 0 { region = os.Getenv("AWS_DEFAULT_REGION") }
		if len(region) == 0 { region = "us-east-1" }

		return ProfileConfig{
			Name: "env",
			Region: region,
			Credentials: credentials.NewEnvCredentials(),
		}
	}

	cfg, err := profileConfig(loadConfig(), profile)
	if err != nil {
		log.Panicf(err.Error())
	}

	if len(cfg.Region) == 0 {
		cfg.Region = "us-east-1"
	}

	return *cfg
}
