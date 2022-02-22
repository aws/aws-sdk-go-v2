package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

func main() {

	// Initialize config that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(defaultConfig)

	// Create table Movies
	tableName := "Movies"

	param := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Year"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("Title"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Year"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("Title"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}

	_, err = client.CreateTable(context.TODO(), param)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}

	log.Printf("Table is created")
}
