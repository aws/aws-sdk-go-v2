package lexruntimeservice_test

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	lexruntime "github.com/aws/aws-sdk-go-v2/service/smithyprototype/lexruntimeservice"
)

func ExampleNew() {
	client := lexruntime.New(lexruntime.Options{
		RegionID:    "us-west-2",
		Credentials: customCredProvider,
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

func ExampleNewFromConfig_customOptions() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	client := lexruntime.NewFromConfig(cfg, func(o *lexruntime.Options) {
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

var customCredProvider = unit.Config().Credentials
