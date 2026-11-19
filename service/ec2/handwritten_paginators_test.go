package ec2

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// --- DescribeVpcEncryptionControls ---

type mockDescribeVpcEncryptionControlsClient struct {
	outputs []*DescribeVpcEncryptionControlsOutput
	inputs  []*DescribeVpcEncryptionControlsInput
	t       *testing.T
}

func (c *mockDescribeVpcEncryptionControlsClient) DescribeVpcEncryptionControls(_ context.Context, input *DescribeVpcEncryptionControlsInput, _ ...func(*Options)) (*DescribeVpcEncryptionControlsOutput, error) {
	c.inputs = append(c.inputs, input)
	requestCnt := len(c.inputs)
	if requestCnt > len(c.outputs) {
		c.t.Fatalf("paginator called client more times than expected (%d)", len(c.outputs))
	}
	return c.outputs[requestCnt-1], nil
}

func TestDescribeVpcEncryptionControlsPaginator(t *testing.T) {
	cases := map[string]struct {
		limit                int32
		requestCnt           int
		stopOnDuplicateToken bool
		outputs              []*DescribeVpcEncryptionControlsOutput
	}{
		"multi-page": {
			limit:      5,
			requestCnt: 3,
			outputs: []*DescribeVpcEncryptionControlsOutput{
				{NextToken: aws.String("token1")},
				{NextToken: aws.String("token2")},
				{NextToken: nil},
			},
		},
		"single page": {
			limit:      10,
			requestCnt: 1,
			outputs: []*DescribeVpcEncryptionControlsOutput{
				{NextToken: nil},
			},
		},
		"stop on duplicate token": {
			limit:                10,
			requestCnt:           2,
			stopOnDuplicateToken: true,
			outputs: []*DescribeVpcEncryptionControlsOutput{
				{NextToken: aws.String("token1")},
				{NextToken: aws.String("token1")},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := &mockDescribeVpcEncryptionControlsClient{
				t:       t,
				outputs: c.outputs,
				inputs:  []*DescribeVpcEncryptionControlsInput{},
			}
			paginator := NewDescribeVpcEncryptionControlsPaginator(client, &DescribeVpcEncryptionControlsInput{}, func(o *DescribeVpcEncryptionControlsPaginatorOptions) {
				o.Limit = c.limit
				o.StopOnDuplicateToken = c.stopOnDuplicateToken
			})

			for paginator.HasMorePages() {
				_, err := paginator.NextPage(context.TODO())
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			inputLen := len(client.inputs)
			testTotalRequestsEC2(c.requestCnt, inputLen, t)

			for i := 1; i < inputLen; i++ {
				if aws.ToInt32(client.inputs[i].MaxResults) != c.limit {
					t.Errorf("page %d: expected MaxResults %d, got %d", i, c.limit, aws.ToInt32(client.inputs[i].MaxResults))
				}
				expectedToken := aws.ToString(c.outputs[i-1].NextToken)
				actualToken := aws.ToString(client.inputs[i].NextToken)
				if expectedToken != actualToken {
					t.Errorf("page %d: expected NextToken %q, got %q", i, expectedToken, actualToken)
				}
			}
		})
	}
}

// --- GetVpcResourcesBlockingEncryptionEnforcement ---

type mockGetVpcResourcesBlockingEncryptionEnforcementClient struct {
	outputs []*GetVpcResourcesBlockingEncryptionEnforcementOutput
	inputs  []*GetVpcResourcesBlockingEncryptionEnforcementInput
	t       *testing.T
}

func (c *mockGetVpcResourcesBlockingEncryptionEnforcementClient) GetVpcResourcesBlockingEncryptionEnforcement(_ context.Context, input *GetVpcResourcesBlockingEncryptionEnforcementInput, _ ...func(*Options)) (*GetVpcResourcesBlockingEncryptionEnforcementOutput, error) {
	c.inputs = append(c.inputs, input)
	requestCnt := len(c.inputs)
	if requestCnt > len(c.outputs) {
		c.t.Fatalf("paginator called client more times than expected (%d)", len(c.outputs))
	}
	return c.outputs[requestCnt-1], nil
}

func TestGetVpcResourcesBlockingEncryptionEnforcementPaginator(t *testing.T) {
	cases := map[string]struct {
		vpcID                string
		limit                int32
		requestCnt           int
		stopOnDuplicateToken bool
		outputs              []*GetVpcResourcesBlockingEncryptionEnforcementOutput
	}{
		"multi-page": {
			vpcID:      "vpc-12345",
			limit:      5,
			requestCnt: 3,
			outputs: []*GetVpcResourcesBlockingEncryptionEnforcementOutput{
				{NextToken: aws.String("token1")},
				{NextToken: aws.String("token2")},
				{NextToken: nil},
			},
		},
		"single page": {
			vpcID:      "vpc-67890",
			limit:      10,
			requestCnt: 1,
			outputs: []*GetVpcResourcesBlockingEncryptionEnforcementOutput{
				{NextToken: nil},
			},
		},
		"stop on duplicate token": {
			vpcID:                "vpc-abcde",
			limit:                10,
			requestCnt:           2,
			stopOnDuplicateToken: true,
			outputs: []*GetVpcResourcesBlockingEncryptionEnforcementOutput{
				{NextToken: aws.String("token1")},
				{NextToken: aws.String("token1")},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := &mockGetVpcResourcesBlockingEncryptionEnforcementClient{
				t:       t,
				outputs: c.outputs,
				inputs:  []*GetVpcResourcesBlockingEncryptionEnforcementInput{},
			}
			paginator := NewGetVpcResourcesBlockingEncryptionEnforcementPaginator(client, &GetVpcResourcesBlockingEncryptionEnforcementInput{
				VpcId: aws.String(c.vpcID),
			}, func(o *GetVpcResourcesBlockingEncryptionEnforcementPaginatorOptions) {
				o.Limit = c.limit
				o.StopOnDuplicateToken = c.stopOnDuplicateToken
			})

			for paginator.HasMorePages() {
				_, err := paginator.NextPage(context.TODO())
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			inputLen := len(client.inputs)
			testTotalRequestsEC2(c.requestCnt, inputLen, t)

			for i := 0; i < inputLen; i++ {
				if aws.ToString(client.inputs[i].VpcId) != c.vpcID {
					t.Errorf("page %d: expected VpcId %q, got %q", i, c.vpcID, aws.ToString(client.inputs[i].VpcId))
				}
			}
			for i := 1; i < inputLen; i++ {
				if aws.ToInt32(client.inputs[i].MaxResults) != c.limit {
					t.Errorf("page %d: expected MaxResults %d, got %d", i, c.limit, aws.ToInt32(client.inputs[i].MaxResults))
				}
				expectedToken := aws.ToString(c.outputs[i-1].NextToken)
				actualToken := aws.ToString(client.inputs[i].NextToken)
				if expectedToken != actualToken {
					t.Errorf("page %d: expected NextToken %q, got %q", i, expectedToken, actualToken)
				}
			}
		})
	}
}

func testTotalRequestsEC2(expect, actual int, t *testing.T) {
	t.Helper()
	if actual != expect {
		t.Errorf("expected %d total requests, got %d", expect, actual)
	}
}