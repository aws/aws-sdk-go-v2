package sts

import (
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
)

func init() {
	initClient = func(c *Client) {
		c.Retryer = retry.AddWithErrorCodes(c.Retryer, ErrCodeIDPCommunicationErrorException)
	}
}
