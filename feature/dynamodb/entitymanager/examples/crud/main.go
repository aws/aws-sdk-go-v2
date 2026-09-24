package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/entitymanager"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Item example struct
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

	tableName := fmt.Sprintf("table_%s", time.Now().Format("2006_01_02_15_04_05"))

	ddb := dynamodb.NewFromConfig(cfg)

	sch, err := entitymanager.NewSchema[Item]()
	if err != nil {
		panic(err)
	}

	sch = sch.WithTableName(aws.String(tableName))

	tbl, err := entitymanager.NewTable[Item](ddb, func(options *entitymanager.TableOptions[Item]) {
		options.Schema = sch
	})
	if err != nil {
		panic(err)
	}

	if exists, err := tbl.Exists(context.Background()); !exists || err != nil {
		if err != nil {
			panic(err)
		}

		if err := tbl.CreateWithWait(context.Background(), time.Minute*2); err != nil {
			panic(err)
		}

		defer func() {
			if err := tbl.DeleteWithWait(context.Background(), time.Minute*2); err != nil {
				panic(err)
			}
		}()
	}

	log.Print("PutItem() up to 10")
	for c := range 10 {
		i, err := tbl.PutItem(context.Background(), &Item{
			ID:    fmt.Sprintf("%d", c),
			Email: fmt.Sprintf("user-%d@amazon.dev.null", c),
			Name:  fmt.Sprintf("First%d", c),
			Body:  fmt.Sprintf("Last%d", c),
		})
		if err != nil {
			log.Printf("Error putting item: %v", err)
		}
		if i != nil {
			log.Printf("Put item %#+v", i)
		}
	}

	log.Print("GetItem() up to 10")
	for c := range 10 {
		m := entitymanager.Map{}.
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
				log.Printf("error: %v", res.Error())

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
				log.Printf("Error: %v", id)
			}

			total++
		}

		log.Printf("Total: %d", total)
	}
}
