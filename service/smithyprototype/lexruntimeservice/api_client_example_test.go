package lexruntimeservice_test

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	lexruntime "github.com/aws/aws-sdk-go-v2/service/smithyprototype/lexruntimeservice"
)

func ExampleNewClient() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	client := lexruntime.NewClient(cfg)
	res, err := client.GetSession(context.TODO(), &lexruntime.GetSessionInput{
		BotAlias: aws.String("botAlias"),
		BotName:  aws.String("botName"),
		UserId:   aws.String("userID"),
	})
	if err != nil {
		log.Fatalf("failed to get session, %v", err)
	}

	fmt.Println("session:", res.SessionId)
}

func ExampleNewClient_customOptions() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	client := lexruntime.NewClient(cfg, func(o *lexruntime.ClientOptions) {
		o.RegionID = "us-west-2"
	})
	res, err := client.GetSession(context.TODO(), &lexruntime.GetSessionInput{
		BotAlias: aws.String("botAlias"),
		BotName:  aws.String("botName"),
		UserId:   aws.String("userID"),
	})
	if err != nil {
		log.Fatalf("failed to get session, %v", err)
	}

	fmt.Println("session:", res.SessionId)
}

var external = mockExternal{}

type mockExternal struct {
}

func (mockExternal) LoadDefaultAWSConfig() (aws.Config, error) {
	return aws.Config{}, nil
}
