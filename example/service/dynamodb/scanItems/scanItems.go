package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

var (
	tableName string
	region    string
)

func init() {
	flag.StringVar(&tableName, "table", "", "The `name` of the DynamoDB table to list item from.")
	flag.StringVar(&region, "region", "", "The `region` of your AWS project.")
}

// Record holds info about the records returned by Scan
type Record struct {
	ID   string
	URLs []string
}

func main() {
	flag.Parse()
	if len(tableName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, table name required")
	}
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region name required")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	// Using the Config value, create the DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	param := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	scan, err := client.Scan(context.TODO(), param)
	if err != nil {
		log.Fatalf("Query API call failed: %s", err)
	}
	var records []Record
	err = attributevalue.UnmarshalListOfMaps(scan.Items, &records)
	if err != nil {
		log.Fatalf("unable to unmarshal records: %v", err)
	}
	for _, record := range records {
		log.Printf("Record : %v", record)
	}

}
