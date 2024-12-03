package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

const (
	rdsAuthTokenID    = "rds-db"
	rdsClusterTokenID = "dsql"
	emptyPayloadHash  = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	userAction        = "DbConnect"
	adminUserAction   = "DbConnectAdmin"
)

// BuildAuthTokenOptions is the optional set of configuration properties for BuildAuthToken
type BuildAuthTokenOptions struct {
	ExpiresIn time.Duration
}

// BuildAuthToken will return an authorization token used as the password for a DB
// connection.
//
// * endpoint - Endpoint consists of the hostname and port needed to connect to the DB. <host>:<port>
// * region - Region is the location of where the DB is
// * dbUser - User account within the database to sign in with
// * creds - Credentials to be signed with
//
// The following example shows how to use BuildAuthToken to create an authentication
// token for connecting to a MySQL database in RDS.
//
//	authToken, err := BuildAuthToken(dbEndpoint, awsRegion, dbUser, awsCreds)
//
//	// Create the MySQL DNS string for the DB connection
//	// user:password@protocol(endpoint)/dbname?<params>
//	connectStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true&tls=rds",
//	   dbUser, authToken, dbEndpoint, dbName,
//	)
//
//	// Use db to perform SQL operations on database
//	db, err := sql.Open("mysql", connectStr)
//
// See http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
// for more information on using IAM database authentication with RDS.
func BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, creds aws.CredentialsProvider, optFns ...func(options *BuildAuthTokenOptions)) (string, error) {
	_, port := validateURL(endpoint)
	if port == "" {
		return "", fmt.Errorf("the provided endpoint is missing a port, or the provided port is invalid")
	}

	values := url.Values{
		"Action": []string{"connect"},
		"DBUser": []string{dbUser},
	}

	return generateAuthToken(ctx, endpoint, region, values, rdsAuthTokenID, creds, optFns...)
}

// GenerateDbConnectAuthToken will return an authorization token as the password for a
// DB connection.
//
// This is the regular user variant, see [GenerateDBConnectSuperUserAuthToken] for the superuser variant
//
// * endpoint - Endpoint is the hostname and optional port to connect to the DB
// * region - Region is the location of where the DB is
// * creds - Credentials to be signed with
func GenerateDbConnectAuthToken(ctx context.Context, endpoint, region string, creds aws.CredentialsProvider, optFns ...func(options *BuildAuthTokenOptions)) (string, error) {
	values := url.Values{
		"Action": []string{userAction},
	}
	return generateAuthToken(ctx, endpoint, region, values, rdsClusterTokenID, creds, optFns...)
}

// GenerateDBConnectSuperUserAuthToken will return an authorization token as the password for a
// DB connection.
//
// This is the superuser user variant, see [GenerateDBConnectSuperUserAuthToken] for the regular user variant
//
// * endpoint - Endpoint is the hostname and optional port to connect to the DB
// * region - Region is the location of where the DB is
// * creds - Credentials to be signed with
func GenerateDBConnectSuperUserAuthToken(ctx context.Context, endpoint, region string, creds aws.CredentialsProvider, optFns ...func(options *BuildAuthTokenOptions)) (string, error) {
	values := url.Values{
		"Action": []string{adminUserAction},
	}
	return generateAuthToken(ctx, endpoint, region, values, rdsClusterTokenID, creds, optFns...)
}

// All generate token functions are presigned URLs behind the scenes with the scheme stripped.
// This function abstracts generating this for all use cases
func generateAuthToken(ctx context.Context, endpoint, region string, values url.Values, signingID string, creds aws.CredentialsProvider, optFns ...func(options *BuildAuthTokenOptions)) (string, error) {
	if len(region) == 0 {
		return "", fmt.Errorf("region is required")
	}
	if len(endpoint) == 0 {
		return "", fmt.Errorf("endpoint is required")
	}

	o := BuildAuthTokenOptions{}

	for _, fn := range optFns {
		fn(&o)
	}

	if o.ExpiresIn == 0 {
		o.ExpiresIn = 15 * time.Minute
	}

	if creds == nil {
		return "", fmt.Errorf("credetials provider must not ne nil")
	}

	// the scheme is arbitrary and is only needed because validation of the URL requires one.
	if !(strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")) {
		endpoint = "https://" + endpoint
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = values.Encode()
	signer := v4.NewSigner()

	credentials, err := creds.Retrieve(ctx)
	if err != nil {
		return "", err
	}

	expires := o.ExpiresIn
	// if creds expire before expiresIn, set that as the expiration time
	if credentials.CanExpire && !credentials.Expires.IsZero() {
		credsExpireIn := credentials.Expires.Sub(sdk.NowTime())
		expires = min(o.ExpiresIn, credsExpireIn)
	}
	query := req.URL.Query()
	query.Set("X-Amz-Expires", strconv.Itoa(int(expires.Seconds())))
	req.URL.RawQuery = query.Encode()

	signedURI, _, err := signer.PresignHTTP(ctx, credentials, req, emptyPayloadHash, signingID, region, sdk.NowTime().UTC())
	if err != nil {
		return "", err
	}

	url := signedURI
	if strings.HasPrefix(url, "http://") {
		url = url[len("http://"):]
	} else if strings.HasPrefix(url, "https://") {
		url = url[len("https://"):]
	}

	return url, nil
}

func validateURL(hostPort string) (host, port string) {
	colon := strings.LastIndexByte(hostPort, ':')
	if colon != -1 {
		host, port = hostPort[:colon], hostPort[colon+1:]
	}
	if !validatePort(port) {
		port = ""
		return
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

func validatePort(port string) bool {
	if _, err := strconv.Atoi(port); err == nil {
		return true
	}
	return false
}
