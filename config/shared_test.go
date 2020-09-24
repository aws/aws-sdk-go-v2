package config

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
)

const ecsResponse = `{
  "Code": "Success",
  "Type": "AWS-HMAC",
  "AccessKeyId": "ecs-access-key",
  "SecretAccessKey": "ecs-secret-key",
  "Token": "token",
  "Expiration": "2100-01-01T00:00:00Z",
  "LastUpdated": "2009-11-23T0:00:00Z"
}`

const ec2MetadataResponse = `{
  "Code": "Success",
  "Type": "AWS-HMAC",
  "AccessKeyId": "ec2-access-key",
  "SecretAccessKey": "ec2-secret-key",
  "Token": "token",
  "Expiration": "2100-01-01T00:00:00Z",
  "LastUpdated": "2009-11-23T0:00:00Z"
}`

const assumeRoleRespMsg = `
<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
    <AssumeRoleResult>
        <AssumedRoleUser>
            <Arn>arn:aws:sts::account_id:assumed-role/role/session_name</Arn>
            <AssumedRoleId>AKID:session_name</AssumedRoleId>
        </AssumedRoleUser>
        <Credentials>
            <AccessKeyId>AKID</AccessKeyId>
            <SecretAccessKey>SECRET</SecretAccessKey>
            <SessionToken>SESSION_TOKEN</SessionToken>
            <Expiration>%s</Expiration>
        </Credentials>
    </AssumeRoleResult>
    <ResponseMetadata>
        <RequestId>request-id</RequestId>
    </ResponseMetadata>
</AssumeRoleResponse>
`

type mockHTTPClient func(*http.Request) (*http.Response, error)

func (m mockHTTPClient) Do(r *http.Request) (*http.Response, error) {
	return m(r)
}

func initConfigTestEnv() (oldEnv []string) {
	oldEnv = awstesting.StashEnv()
	os.Setenv("AWS_CONFIG_FILE", "file_not_exists")
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "file_not_exists")

	return oldEnv
}
