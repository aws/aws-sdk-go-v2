package amplifybackend

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/amplifybackend/types"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type mockListBackendJobsClient struct {
	outputs []*ListBackendJobsOutput
	inputs  []*ListBackendJobsInput
	t       *testing.T
	limit   int32
}

func (c *mockListBackendJobsClient) ListBackendJobs(ctx context.Context, input *ListBackendJobsInput, optFns ...func(*Options)) (*ListBackendJobsOutput, error) {
	c.inputs = append(c.inputs, input)
	requestCnt := len(c.inputs)
	testCurRequest(len(c.outputs), requestCnt, c.limit, aws.ToInt32(input.MaxResults), c.t)
	return c.outputs[requestCnt-1], nil
}

type listBackendJobsTestCase struct {
	limit                  int32
	requestCnt             int
	stopOnDuplicationToken bool
	outputs                []*ListBackendJobsOutput
}

func TestListBackendJobsPaginator(t *testing.T) {
	cases := map[string]listBackendJobsTestCase{
		"page limit 3": {
			limit:      3,
			requestCnt: 3,
			outputs: []*ListBackendJobsOutput{
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job1"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job2"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job3"),
						},
					},
					NextToken: aws.String("token1"),
				},
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job4"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job5"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job6"),
						},
					},
					NextToken: aws.String("token2"),
				},
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job7"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job8"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job9"),
						},
					},
				},
			},
		},
		"total count 2 due to nil nextToken": {
			limit:      3,
			requestCnt: 2,
			outputs: []*ListBackendJobsOutput{
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job1"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job2"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job3"),
						},
					},
					NextToken: aws.String("token1"),
				},
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job4"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job5"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job6"),
						},
					},
				},
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job7"),
						},
					},
				},
			},
		},
		"total count 2 due to duplicate nextToken": {
			limit:                  3,
			requestCnt:             2,
			stopOnDuplicationToken: true,
			outputs: []*ListBackendJobsOutput{
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job1"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job2"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job3"),
						},
					},
					NextToken: aws.String("token1"),
				},
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job4"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job5"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job6"),
						},
					},
					NextToken: aws.String("token1"),
				},
				{
					Jobs: []types.BackendJobRespObj{
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job7"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job8"),
						},
						{
							AppId: aws.String("App"),
							JobId: aws.String("Job9"),
						},
					},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := mockListBackendJobsClient{
				t:       t,
				outputs: c.outputs,
				inputs:  []*ListBackendJobsInput{},
				limit:   c.limit,
			}
			paginator := NewListBackendJobsPaginator(&client, &ListBackendJobsInput{}, func(options *ListBackendJobsPaginatorOptions) {
				options.Limit = c.limit
				options.StopOnDuplicateToken = c.stopOnDuplicationToken
			})

			for paginator.HasMorePages() {
				_, err := paginator.NextPage(context.TODO())
				if err != nil {
					t.Errorf("error: %v", err)
				}
			}

			inputLen := len(client.inputs)
			testTotalRequests(c.requestCnt, inputLen, t)
			for i := 1; i < inputLen; i++ {
				if *client.inputs[i].NextToken != *c.outputs[i-1].NextToken {
					t.Errorf("Expect next input's nextToken to be eaqul to %s, got %s",
						*c.outputs[i-1].NextToken, *client.inputs[i].NextToken)
				}
			}
		})
	}
}

func testCurRequest(maxReqCnt, actualReqCnt int, expectLimit, actualLimit int32, t *testing.T) {
	if actualReqCnt > maxReqCnt {
		t.Errorf("Paginator calls client more than expected %d times", maxReqCnt)
	}
	if expectLimit != actualLimit {
		t.Errorf("Expect page limit to be %d, got %d", expectLimit, actualLimit)
	}
}

func testTotalRequests(expect, actual int, t *testing.T) {
	if actual != expect {
		t.Errorf("Expect total request number to be %d, got %d", expect, actual)
	}
}
