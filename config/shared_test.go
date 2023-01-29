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
  "LastUpdated": "2009-11-23T00:00:00Z"
}`

const ec2MetadataResponse = `{
  "Code": "Success",
  "Type": "AWS-HMAC",
  "AccessKeyId": "ec2-access-key",
  "SecretAccessKey": "ec2-secret-key",
  "Token": "token",
  "Expiration": "2100-01-01T00:00:00Z",
  "LastUpdated": "2009-11-23T00:00:00Z"
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

var assumeRoleWithWebIdentityResponse = `<AssumeRoleWithWebIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleWithWebIdentityResult>
    <SubjectFromWebIdentityToken>amzn1.account.AF6RHO7KZU5XRVQJGXK6HB56KR2A</SubjectFromWebIdentityToken>
    <Audience>client.5498841531868486423.1548@apps.example.com</Audience>
    <AssumedRoleUser>
      <Arn>arn:aws:sts::123456789012:assumed-role/FederatedWebIdentityRole/app1</Arn>
      <AssumedRoleId>AROACLKWSDQRAOEXAMPLE:app1</AssumedRoleId>
    </AssumedRoleUser>
    <Credentials>
      <AccessKeyId>WEB_IDENTITY_AKID</AccessKeyId>
      <SecretAccessKey>WEB_IDENTITY_SECRET</SecretAccessKey>
      <SessionToken>WEB_IDENTITY_SESSION_TOKEN</SessionToken>
      <Expiration>%s</Expiration>
    </Credentials>
    <Provider>www.amazon.com</Provider>
  </AssumeRoleWithWebIdentityResult>
  <ResponseMetadata>
    <RequestId>request-id</RequestId>
  </ResponseMetadata>
</AssumeRoleWithWebIdentityResponse>
`

const getRoleCredentialsResponse = `{
  "roleCredentials": {
    "accessKeyId": "SSO_AKID",
    "secretAccessKey": "SSO_SECRET_KEY",
    "sessionToken": "SSO_SESSION_TOKEN",
    "expiration": %d
  }
}`

const ssoTokenCacheFile = `{
  "accessToken": "ssoAccessToken",
  "expiresAt": "%s"
}`

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
