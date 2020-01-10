package rdsutils

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPV4Signer interface is used to presign a request
type HTTPV4Signer interface {
	Presign(ctx context.Context, r *http.Request, body io.ReadSeeker, service, region string, exp time.Duration, signTime time.Time) (http.Header, error)
}

// BuildAuthToken will return an authorization token used as the password for a DB
// connection.
//
// * endpoint - Endpoint consists of the port needed to connect to the DB. <host>:<port>
// * region - Region is the location of where the DB is
// * dbUser - User account within the database to sign in with
// * signer - Signer used to be signed with
//
// The following example shows how to use BuildAuthToken to create an authentication
// token for connecting to a MySQL database in RDS.
//
//   signer := v4.NewSigner(credsProvider)
//   authToken, err := BuildAuthToken(ctx, dbEndpoint, awsRegion, dbUser, signer)
//
//   // Create the MySQL DNS string for the DB connection
//   // user:password@protocol(endpoint)/dbname?<params>
//   connectStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?allowCleartextPasswords=true&tls=rds",
//      dbUser, authToken, dbEndpoint, dbName,
//   )
//
//   // Use db to perform SQL operations on database
//   db, err := sql.Open("mysql", connectStr)
//
// See http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
// for more information on using IAM database authentication with RDS.
func BuildAuthToken(ctx context.Context, endpoint, region, dbUser string, signer HTTPV4Signer) (string, error) {
	// the scheme is arbitrary and is only needed because validation of the URL requires one.
	if !(strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")) {
		endpoint = "https://" + endpoint
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	values := req.URL.Query()
	values.Set("Action", "connect")
	values.Set("DBUser", dbUser)
	req.URL.RawQuery = values.Encode()

	_, err = signer.Presign(ctx, req, nil, "rds-db", region, 15*time.Minute, time.Now())
	if err != nil {
		return "", err
	}

	url := req.URL.String()
	if strings.HasPrefix(url, "http://") {
		url = url[len("http://"):]
	} else if strings.HasPrefix(url, "https://") {
		url = url[len("https://"):]
	}

	return url, nil
}
