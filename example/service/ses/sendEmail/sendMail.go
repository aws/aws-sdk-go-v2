package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

const SesRegion = "us-east-1"

// SendEmail attempts to send email on specified email address
func SendEmail(from string, to []string, cc []string, bcc []string) error {
	fmt.Println("Sending email using ses")
	client, err := getSESClient(SesRegion)
	if err != nil {
		fmt.Printf("error while creating session for sending email %v", err)
		return err
	}

	messageHTMLBody := "<html></head><title>This is html body</title></head><body><h>hello there</h><br>this is message body</body></html>"
	messageTextBody := "This is message body"
	messageSubject := "This email is from ses"
	message := &types.Message{
		Body: &types.Body{
			Html: &types.Content{
				Data: &messageHTMLBody,
				// UTF-8, ISO-8859-1,
				Charset: aws.String("UTF-8"),
			},
			Text: &types.Content{
				Data: &messageTextBody,
				// UTF-8, ISO-8859-1,
				Charset: aws.String("UTF-8"),
			},
		},
		Subject: &types.Content{
			Data:    &messageSubject,
			Charset: aws.String("UTF-8"),
		},
	}
	emailInput := ses.SendEmailInput{
		Message: message,
		Destination: &types.Destination{
			ToAddresses:  to,
			CcAddresses:  cc,
			BccAddresses: bcc,
		},
		Source: &from, // you need to put same from name which you have in SES
	}

	emailOutput, err := client.SendEmail(context.TODO(), &emailInput)

	if err != nil {
		return fmt.Errorf("send email failed. Error: %v", err)
	}
	fmt.Printf("Email sent successfully. Response: %v", emailOutput)
	return nil
}

func getSESClient(region string) (client *ses.Client, err error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(region))
	if err != nil {
		return client, fmt.Errorf("issue getting ses aws configuration. Error: %v", err)
	}
	client = ses.NewFromConfig(cfg)
	return client, err
}

func main() {
	to := []string{""}  // email id to set in to
	cc := []string{""}  // email id to set in cc
	bcc := []string{""} // email id to set in bcc
	fromEmail := ""     // from email id, This needs to be same as configured in ses setup
	err := SendEmail(fromEmail, to, cc, bcc)
	if err != nil {
		err = fmt.Errorf("error while sending email. ERROR: %v", err)
		fmt.Println(err)
	}
}
