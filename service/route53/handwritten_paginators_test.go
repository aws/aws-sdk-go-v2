package route53

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/smithy-go/middleware"
	"testing"
)

var limit int32
var startRecordName *string
var startRecordIdentifier *string
var startRecordType types.RRType

type testListRRSMiddleware struct {
	id int
}

func (m *testListRRSMiddleware) ID() string {
	return fmt.Sprintf("mock middleware %d", m.id)
}

func (m *testListRRSMiddleware) HandleInitialize(ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler) (
	output middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	params := input.Parameters.(*ListResourceRecordSetsInput)
	startRecordName = params.StartRecordName
	limit = *params.MaxItems
	startRecordIdentifier = params.StartRecordIdentifier
	startRecordType = params.StartRecordType
	return middleware.InitializeOutput{Result: &ListResourceRecordSetsOutput{}}, metadata, nil
}

type testCase struct {
	startRecordName       *string
	limit                 int32
	startRecordIdentifier *string
	startRecordType       types.RRType
}

func TestListResourceRecordSetsPaginator(t *testing.T) {
	cases := map[string]testCase{
		"page limit 5 with record name but without record type and identifier": {
			startRecordName: aws.String("testRecord1"),
			limit:           5,
		},
		"page limit 10 with record name and type": {
			startRecordName: aws.String("testRecord2"),
			limit:           10,
			startRecordType: types.RRTypeTxt,
		},
		"page limit 15 with record name, type and identifier": {
			startRecordName:       aws.String("testRecord3"),
			limit:                 15,
			startRecordIdentifier: aws.String("testID1"),
			startRecordType:       types.RRTypeTxt,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := NewFromConfig(aws.Config{})

			paginator := NewListResourceRecordSetsPaginator(client, &ListResourceRecordSetsInput{
				StartRecordName:       c.startRecordName,
				StartRecordType:       c.startRecordType,
				StartRecordIdentifier: c.startRecordIdentifier,
			}, func(options *ListResourceRecordSetsPaginatorOptions) {
				options.Limit = c.limit
			})
			if !paginator.HasMorePages() {
				t.Errorf("Expect paginator has more page, got not")
			}

			paginator.NextPage(context.TODO(), initializeMiddlewareFn(&testListRRSMiddleware{1}))

			testNextPageResult(c, paginator, t)
		})
	}
}

// insert middleware at the beginning of initialize step to see if page limit and other params
// can be passed to API call's stack input
func initializeMiddlewareFn(initializeMiddleware middleware.InitializeMiddleware) func(*Options) {
	return func(options *Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			return stack.Initialize.Add(initializeMiddleware, middleware.Before)
		})
	}
}

// unit test can not control client API call's output, so just check params' default nil value
func testNextPageResult(c testCase, p *ListResourceRecordSetsPaginator, t *testing.T) {
	if c.limit != limit {
		t.Errorf("Expect page limit to be %d, got %d", c.limit, limit)
	}
	if *c.startRecordName != *startRecordName {
		t.Errorf("Expect startRecordName to be %s, got %s", *c.startRecordName, *startRecordName)
	}
	if c.startRecordType != startRecordType {
		t.Errorf("Expect startRecordType to be %s, got %s", c.startRecordType, startRecordType)
	}
	if c.startRecordIdentifier != nil && *c.startRecordIdentifier != *startRecordIdentifier {
		t.Errorf("Expect startRecordIdentifier to be %s, got %s",
			*c.startRecordIdentifier, *startRecordIdentifier)
	}
	if p.startRecordName != nil || p.startRecordType != "" || p.startRecordIdentifier != nil {
		t.Errorf("Expect paginator record params to be zero value, got %s, %s and %s",
			*p.startRecordName, p.startRecordType, *p.startRecordIdentifier)
	}
}
