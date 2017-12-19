package shared

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/glassechidna/awsweb/sneakyvendor/aws-shared-defaults"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"path/filepath"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"context"
	"github.com/aws/aws-sdk-go/aws/request"
)

type cachedSessionTokenResponse struct {
	MfaSerialNumber string
	Credentials struct {
		AccessKeyId string
		SecretAccessKey string
		SessionToken string
		Expiration time.Time
	}
	ResponseMetadata struct {
		RetryAttempts int
		HTTPStatusCode int
		RequestId string
		HTTPHeaders map[string]string
	}
}

func mfaAuthenticatedCredentials(sourceCreds *credentials.Credentials, mfaSerial string, mfaProvider func() (string, error)) (*credentials.Credentials, error) {
	sess, _ := session.NewSession(&aws.Config{
		Credentials: sourceCreds,
	})
	api := sts.New(sess)

	cached := cachedMfaAuthenticatedCredentials(mfaSerial)

	if cached == nil {
		code, _ := mfaProvider()

		input := &sts.GetSessionTokenInput{
			SerialNumber: &mfaSerial,
			TokenCode: &code,
			DurationSeconds: aws.Int64(3600),
		}

		statusCode := 0
		requestId := ""
		headers := map[string]string{}

		resp, _ := api.GetSessionTokenWithContext(context.Background(), input, func(r *request.Request) {
			r.Handlers.Complete.PushBack(func(req *request.Request) {
				statusCode = req.HTTPResponse.StatusCode
				requestId = req.RequestID

				for key, val := range req.HTTPResponse.Header {
					headers[key] = val[0]
				}
			})
		})

		c := resp.Credentials

		cached = &cachedSessionTokenResponse{
			MfaSerialNumber: mfaSerial,
			Credentials: struct {
				AccessKeyId     string
				SecretAccessKey string
				SessionToken    string
				Expiration      time.Time
			}{
				AccessKeyId:     *c.AccessKeyId,
				SecretAccessKey: *c.SecretAccessKey,
				SessionToken:    *c.SessionToken,
				Expiration:      time.Now().Add(time.Hour),
			},
			ResponseMetadata: struct {
				RetryAttempts  int
				HTTPStatusCode int
				RequestId      string
				HTTPHeaders    map[string]string
			}{
				RetryAttempts: 0,
				HTTPStatusCode: statusCode,
				RequestId: requestId,
				HTTPHeaders: headers,
			},
		}

		cachedBytes, _ := json.MarshalIndent(cached, "", "  ")
		path := cachePathForMfaSerial(mfaSerial)
		ioutil.WriteFile(path, cachedBytes, 0600)
	}

	c := cached.Credentials
	return credentials.NewStaticCredentials(c.AccessKeyId, c.SecretAccessKey, c.SessionToken), nil
}

func cachedMfaAuthenticatedCredentials(mfaSerial string) *cachedSessionTokenResponse {
	path := cachePathForMfaSerial(mfaSerial)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	resp := cachedSessionTokenResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil
	}

	if resp.Credentials.Expiration.Before(time.Now()) {
		return nil
	}

	return &resp
}

func cachePathForMfaSerial(mfaSerial string) string {
	dir := filepath.Join(shareddefaults.UserHomeDir(), ".aws", "awswebcache")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	// make name filesystem-friendly
	mfaSerial = strings.Replace(mfaSerial, ":", "-", -1)
	mfaSerial = strings.Replace(mfaSerial, "/", "-", -1)

	return filepath.Join(dir, mfaSerial) + ".json"
}

