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

		sourceCredsSection, err := cfg.cred.GetSection(sourceProfile.String())
		if err != nil { return nil, err }

		mfaSecret := sourceCredsSection.Key("mfa_secret")
		var mfaProvider func() (string, error)

		if len(mfaSecret.String()) == 0 {
			mfaProvider = stscreds.StdinTokenProvider
		} else {
			mfaProvider = func() (string, error) {
				return totp.GenerateCode(mfaSecret.String(), time.Now())
			}
		}

		mfaSerial, err := section.GetKey("mfa_serial")
		if err != nil { return nil, err }

		roleArn, err := section.GetKey("role_arn")
		if err != nil { return nil, err }

		creds, _ := mfaAuthenticatedCredentials(sourceConfig.Credentials, mfaSerial.String(), mfaProvider)
		sess := session.Must(session.NewSession(&aws.Config{
			Credentials: creds,
		}))

		api := sts.New(sess)

		roleSessionName := fmt.Sprintf("awsweb-%d", time.Now().Second())
		resp, err := api.AssumeRole(&sts.AssumeRoleInput{
			RoleArn: aws.String(roleArn.String()),
			RoleSessionName: &roleSessionName,
		})
		if err != nil { return nil, err }

		c := resp.Credentials
		roleCreds := credentials.NewStaticCredentials(*c.AccessKeyId, *c.SecretAccessKey, *c.SessionToken)

		return &ProfileConfig{
			Name: profile,
			Region: region,
			Credentials: roleCreds,
		}, nil
	} else {
		credsSection, err := cfg.cred.GetSection(profile)
		if err != nil { return nil, err }

		id, err := credsSection.GetKey("aws_access_key_id")
		if err != nil { return nil, err }

		secret, err := credsSection.GetKey("aws_secret_access_key")
		if err != nil { return nil, err }

		token := credsSection.Key("aws_session_token")

		return &ProfileConfig{
			Name: profile,
			Region: region,
			Credentials: credentials.NewStaticCredentials(
				id.String(),
				secret.String(),
				token.String(),
			),
		}, nil
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
