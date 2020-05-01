package lexruntimeservice_test

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/smithyprototype/lexruntimeservice"
	lexruntime "github.com/aws/aws-sdk-go-v2/service/smithyprototype/lexruntimeservice"
)

func ExampleClient_ListSessions_pagination() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("unable to load configuration, %v", err)
	}

	client := lexruntimeservice.NewFromConfig(cfg)

	// Create a paginator with the client and input parameters.
	p := lexruntime.NewListSessionsPaginator(client, &lexruntime.ListSessionsInput{
		BotAlias: aws.String("botAlias"),
		BotName:  aws.String("botName"),
		UserId:   aws.String("userID"),
	})

	for p.HasMorePages() {
		o, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get next page, %v", err)
		}
		fmt.Println("Page:", o)
	}
}
