# Example

This example demonstrates how you can use the AWS SDK for Go V2's Amazon SES client 
to send an email.

# Usage

## How to send email

Input Object:
```go
    ses.SendEmailInput{
        Message: message,
        Destination: &types.Destination{
            ToAddresses:  to,
            CcAddresses:  cc,
            BccAddresses: bcc,
        },
        Source: &from, // you need to put same from name which you have in SES
    }

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
```
SES support both html and text based email, So input object contains both the details as part of body.
```sh
 AWS_REGION=<region> go run sendMail.go
 ```
Note: Please have aws profile added(aws configure)