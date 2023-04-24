package route53

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"testing"
)

type mockListResourceRecordSetsClient struct {
	outputs []*ListResourceRecordSetsOutput
	inputs  []*ListResourceRecordSetsInput
	t       *testing.T
}

func (c *mockListResourceRecordSetsClient) ListResourceRecordSets(ctx context.Context, input *ListResourceRecordSetsInput, optFns ...func(*Options)) (*ListResourceRecordSetsOutput, error) {
	c.inputs = append(c.inputs, input)
	requestCnt := len(c.inputs)
	if *input.MaxItems != *c.outputs[requestCnt-1].MaxItems {
		c.t.Errorf("Expect page limit to be %d, got %d", *c.outputs[requestCnt-1].MaxItems, *input.MaxItems)
	}
	if outputLen := len(c.outputs); requestCnt > outputLen {
		c.t.Errorf("Paginator calls client more than expected %d times", outputLen)
	}
	return c.outputs[requestCnt-1], nil
}

type listRRSTestCase struct {
	limit                  int32
	requestCnt             int
	stopOnDuplicationToken bool
	outputs                []*ListResourceRecordSetsOutput
}

func TestListResourceRecordSetsPaginator(t *testing.T) {
	cases := map[string]listRRSTestCase{
		"page limit 5": {
			limit:      5,
			requestCnt: 3,
			outputs: []*ListResourceRecordSetsOutput{
				{
					MaxItems:             aws.Int32(5),
					NextRecordName:       aws.String("testRecord1"),
					NextRecordIdentifier: aws.String("testID1"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          true,
				},
				{
					MaxItems:             aws.Int32(5),
					NextRecordName:       aws.String("testRecord2"),
					NextRecordIdentifier: aws.String("testID2"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          true,
				},
				{
					MaxItems:             aws.Int32(5),
					NextRecordName:       aws.String("testRecord3"),
					NextRecordIdentifier: aws.String("testID3"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          false,
				},
			},
		},
		"page limit 10 with duplicate record": {
			limit:                  10,
			requestCnt:             4,
			stopOnDuplicationToken: true,
			outputs: []*ListResourceRecordSetsOutput{
				{
					MaxItems:             aws.Int32(10),
					NextRecordName:       aws.String("testRecord1"),
					NextRecordIdentifier: aws.String("testID1"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          true,
				},
				{
					MaxItems:             aws.Int32(10),
					NextRecordName:       aws.String("testRecord2"),
					NextRecordIdentifier: aws.String("testID2"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          true,
				},
				{
					MaxItems:             aws.Int32(10),
					NextRecordName:       aws.String("testRecord3"),
					NextRecordIdentifier: aws.String("testID3"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          true,
				},
				{
					MaxItems:             aws.Int32(10),
					NextRecordName:       aws.String("testRecord3"),
					NextRecordIdentifier: aws.String("testID3"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          true,
				},
				{
					MaxItems:             aws.Int32(10),
					NextRecordName:       aws.String("testRecord5"),
					NextRecordIdentifier: aws.String("testID5"),
					NextRecordType:       types.RRTypeA,
					IsTruncated:          false,
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := mockListResourceRecordSetsClient{
				inputs:  []*ListResourceRecordSetsInput{},
				outputs: c.outputs,
				t:       t,
			}

			paginator := NewListResourceRecordSetsPaginator(&client, &ListResourceRecordSetsInput{}, func(options *ListResourceRecordSetsPaginatorOptions) {
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
			if inputLen != c.requestCnt {
				t.Errorf("Expect total request number to be %d, got %d", c.requestCnt, inputLen)
			}
			for i := 1; i < inputLen; i++ {
				if *client.inputs[i].StartRecordName != *c.outputs[i-1].NextRecordName {
					t.Errorf("Expect next input's RecordName to be eaqul to %s, got %s",
						*c.outputs[i-1].NextRecordName, *client.inputs[i].StartRecordName)
				}
				if *client.inputs[i].StartRecordIdentifier != *c.outputs[i-1].NextRecordIdentifier {
					t.Errorf("Expect next input's RecordIdentifier to be eaqul to %s, got %s",
						*c.outputs[i-1].NextRecordIdentifier, *client.inputs[i].StartRecordIdentifier)
				}
				if client.inputs[i].StartRecordType != c.outputs[i-1].NextRecordType {
					t.Errorf("Expect next input's RecordType to be eaqul to %s, got %s",
						c.outputs[i-1].NextRecordType, client.inputs[i].StartRecordType)
				}
			}
		})
	}
}
