// +build example

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	cfg := Config{}
	if err := cfg.Load(); err != nil {
		exitErrorf("failed to load config, %v", err)
	}

	// Create the config that the DynamoDB service will use.
	awscfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}
	if len(cfg.Region) > 0 {
		// The Region for the DynamoDB table. If Config.Region is not set
		// the region must come from the shared config or AWS_REGION
		// environment variable.
		awscfg.Region = cfg.Region
	}

	// Create the DynamoDB service client to make the query request with.
	svc := dynamodb.New(awscfg)

	// Create the Expression to fill the input struct with.
	filt := expression.Name("Artist").Equal(expression.Value("No One You Know"))
	proj := expression.NamesList(expression.Name("SongTitle"), expression.Name("AlbumTitle"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		exitErrorf("failed to create the Expression, %v", err)
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(cfg.Table),
	}
	if cfg.Limit > 0 {
		params.Limit = aws.Int64(cfg.Limit)
	}

	// Make the DynamoDB Query API call
	req := svc.ScanRequest(params)
	result, err := req.Send()
	if err != nil {
		exitErrorf("failed to make Query API call, %v", err)
	}

	fmt.Println(result)
}

type Config struct {
	Table  string // required
	Region string // optional
	Limit  int64  // optional
}

func (c *Config) Load() error {
	flag.Int64Var(&c.Limit, "limit", 0, "Limit is the max items to be returned, 0 is no limit")
	flag.StringVar(&c.Table, "table", "", "Table to Query on")
	flag.StringVar(&c.Region, "region", "", "AWS Region the table is in")
	flag.Parse()

	if len(c.Table) == 0 {
		flag.PrintDefaults()
		return fmt.Errorf("table name is required.")
	}

	return nil
}
