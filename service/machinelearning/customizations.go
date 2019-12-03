package machinelearning

import (
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/machinelearning/types"
)

func init() {
	initRequest = func(c *Client, r *aws.Request) {
		switch r.Operation.Name {
		case opPredict:
			r.Handlers.Build.PushBack(updatePredictEndpoint)
		}
	}
}

// updatePredictEndpoint rewrites the request endpoint to use the
// "PredictEndpoint" parameter of the Predict operation.
func updatePredictEndpoint(r *aws.Request) {
	if !r.ParamsFilled() {
		return
	}

	r.Metadata.Endpoint = *r.Params.(*types.PredictInput).PredictEndpoint

	uri, err := url.Parse(r.Metadata.Endpoint)
	if err != nil {
		r.Error = err
		return
	}
	r.HTTPRequest.URL = uri
}
