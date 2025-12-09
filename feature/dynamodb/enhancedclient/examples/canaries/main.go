package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/enhancedclient"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Item struct {
	ID    string `dynamodbav:"id,partition"`
	Email string `dynamodbav:"email,sort"`
	Name  string `dynamodbav:"name"`
	Body  string `dynamodbav:"body"`
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	tableName := "CSM_aws-sdk-go-v2_LL_GetItem_Small"

	ddb := dynamodb.NewFromConfig(cfg)

	sch, err := enhancedclient.NewSchema[Item]()
	if err != nil {
		panic(err)
	}

	sch = sch.WithTableName(aws.String(tableName))

	tbl, err := enhancedclient.NewTable[Item](ddb, func(options *enhancedclient.TableOptions[Item]) {
		options.Schema = sch
	})
	if err != nil {
		panic(err)
	}

	log.Print("GetItem() up to 10")
	for c := range 10 {
		m := enhancedclient.Map{}.
			With("id", fmt.Sprintf("%d", c)).
			With("email", fmt.Sprintf("user-%d@amazon.dev.null", c))

		i, err := tbl.GetItem(
			context.Background(),
			m,
		)
		if err != nil {
			log.Printf("Error getting item %v: %v", m, err)
		}
		if i != nil {
			log.Printf("Got item %#+v", i)
		}
	}

	log.Print("Query()")
	{
		keyCond := expression.Key("id").Equal(expression.Value("1")) //.
		//		And(expression.Key("email").BeginsWith("user"))

		expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
		if err != nil {
			panic(err)
		}

		for res := range tbl.Query(context.Background(), expr) {
			if res.Error() != nil {
				log.Printf("Query() error: %v", res.Error())

				continue
			}

			log.Printf("Got item %#+v", res.Item())
		}
	}

	log.Print("Scan()")
	{
		f := expression.Name("id").Contains("1").
			And(expression.Name("email").Contains("user"))

		expr, err := expression.NewBuilder().WithFilter(f).Build()
		if err != nil {
			panic(err)
		}

		total := 0
		for res := range tbl.Scan(context.Background(), expr) {
			if res.Error() != nil {
				log.Printf("Scan() error: %v", res.Error())

				continue
			}

			id := res.Item().ID
			if strings.Contains(id, "1") {
				log.Printf("Got item %#+v", id)
			} else {
				log.Printf("WTF?: %v", id)
			}

			total++
		}

		log.Printf("Total: %d", total)
	}
}
