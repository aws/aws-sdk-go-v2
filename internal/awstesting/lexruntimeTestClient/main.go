package main

import (
	"context"
	cryptorand "crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/lexruntimeservice"
	"github.com/awslabs/smithy-go/ptr"
	smithyrand "github.com/awslabs/smithy-go/rand"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

var (
	botAlias string
	botName  string
	userID   string
)

func init() {
	uuid, _ := smithyrand.NewUUID(cryptorand.Reader).GetUUID()

	flag.StringVar(&botAlias, "bot-alias", "TEST", "alias of the bot to use")
	flag.StringVar(&botName, "bot-name", "OrderFlowers", "name of the bot to use")
	flag.StringVar(&userID, "user-id", uuid, "unique application user ID")
}

func main() {
	flag.Parse()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to get config, %v", err)
	}
	cfg.Retryer = aws.NoOpRetryer{}
	ctx := context.Background()

	client := lexruntimeservice.NewFromConfig(cfg, func(o *lexruntimeservice.Options) {
		o.HTTPSigner = v4.NewSigner(o.Credentials)

		origEndpointResolver := o.EndpointResolver
		o.EndpointResolver = aws.EndpointResolverFunc(func(prefix, region string) (aws.Endpoint, error) {
			// Need to specify endpoint prefix because Smithy does not have it.
			return origEndpointResolver.ResolveEndpoint("runtime.lex", region)
		})

		o.HTTPClient = smithyhttp.WrapLogClient(logger{}, o.HTTPClient, false)
	})

	if err = putSession(ctx, client); err != nil {
		log.Fatalf("failed to call PutSession, %v", err)
	}

	if err = getSession(ctx, client); err != nil {
		log.Fatalf("failed to call GetSession, %v", err)
	}

	// TODO fails because signer missing middleware
	if err = postContent(ctx, client); err != nil {
		log.Fatalf("failed to call PostContent, %v", err)
	}

	if err = postText(ctx, client); err != nil {
		log.Fatalf("failed to call PostText, %v", err)
	}

	if err = deleteSession(ctx, client); err != nil {
		log.Fatalf("failed to call DeleteSession, %v", err)
	}
}

func putSession(ctx context.Context, client *lexruntimeservice.Client) error {
	log.Println("Putting Session", userID)
	resp, err := client.PutSession(ctx, &lexruntimeservice.PutSessionInput{
		Accept:   ptr.String("audio/mpeg"),
		BotAlias: ptr.String(botAlias),
		BotName:  ptr.String(botName),
		UserId:   ptr.String(userID),
	})
	if err != nil {
		return err
	}

	log.Println("PutSession successful")
	log.Println(resp)
	log.Println()
	return nil
}

func getSession(ctx context.Context, client *lexruntimeservice.Client) error {
	log.Println("Getting Session", userID)
	resp, err := client.GetSession(ctx, &lexruntimeservice.GetSessionInput{
		BotAlias: ptr.String(botAlias),
		BotName:  ptr.String(botName),
		UserId:   ptr.String(userID),
	})
	if err != nil {
		return err
	}

	log.Println("GetSession successful")
	log.Println(resp)
	log.Println()
	return nil
}

func postContent(ctx context.Context, client *lexruntimeservice.Client) error {
	log.Println("Posting Content", userID)
	resp, err := client.PostContent(ctx, &lexruntimeservice.PostContentInput{
		Accept:      ptr.String("audio/mpeg"),
		BotAlias:    ptr.String(botAlias),
		BotName:     ptr.String(botName),
		UserId:      ptr.String(userID),
		ContentType: ptr.String("text/plain; charset=utf-8"),
		InputStream: strings.NewReader("hello there how are you doing today"),
	})
	if err != nil {
		return err
	}

	log.Println("PostContent successful")
	log.Println(resp)
	defer resp.AudioStream.Close()

	b, err := ioutil.ReadAll(resp.AudioStream)
	if err != nil {
		return fmt.Errorf("failed to read response, %w", err)
	}
	log.Println("Response body length,", len(b))
	log.Println()
	return nil
}

func postText(ctx context.Context, client *lexruntimeservice.Client) error {
	log.Println("Posting Text", userID)
	resp, err := client.PostText(ctx, &lexruntimeservice.PostTextInput{
		BotAlias:  ptr.String(botAlias),
		BotName:   ptr.String(botName),
		UserId:    ptr.String(userID),
		InputText: ptr.String("hello there how are you"),
	})
	if err != nil {
		return err
	}

	log.Println("PostText successful")
	log.Println(resp)
	log.Println()
	return nil
}

func deleteSession(ctx context.Context, client *lexruntimeservice.Client) error {
	log.Println("Deleting Session", userID)
	resp, err := client.DeleteSession(ctx, &lexruntimeservice.DeleteSessionInput{
		BotAlias: ptr.String(botAlias),
		BotName:  ptr.String(botName),
		UserId:   ptr.String(userID),
	})
	if err != nil {
		return err
	}

	log.Println("DeleteSession successful")
	log.Println(resp)
	log.Println()
	return nil
}

type logger struct{}

func (logger) Logf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
