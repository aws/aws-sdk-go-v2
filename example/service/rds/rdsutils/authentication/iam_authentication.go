// +build example,skip

package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/aws/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/rds/rdsutils"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-sql-driver/mysql"
)

// Usage ./iam_authentication <region> <db user> <db name> <endpoint to database> <iam arn>
func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "USAGE ERROR: go run concatenateObjects.go <region> <endpoint to database> <iam arn>\n")
		os.Exit(1)
	}

	awsRegion := os.Args[1]
	dbUser := os.Args[2]
	dbName := os.Args[3]
	dbEndpoint := os.Args[4]

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration, %v", err)
		os.Exit(1)
	}
	cfg.Region = awsRegion

	credProvider := stscreds.NewAssumeRoleProvider(sts.New(cfg), os.Args[5])
	signer := v4.NewSigner(credProvider)
	authToken, err := rdsutils.BuildAuthToken(context.Background(), dbEndpoint, awsRegion, dbUser, signer)

	// Create the MySQL DNS string for the DB connection
	// user:password@protocol(endpoint)/dbname?<params>
	dnsStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true",
		dbUser, authToken, dbEndpoint, dbName,
	)

	driver := mysql.MySQLDriver{}
	_ = driver
	// Use db to perform SQL operations on database
	if _, err = sql.Open("mysql", dnsStr); err != nil {
		panic(err)
	}

	fmt.Println("Successfully opened connection to database")
}
