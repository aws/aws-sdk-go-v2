package aws_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

func TestPagination(t *testing.T) {
	type mockInput struct {
		Foo *string
	}

	type mockOutput struct {
		Bar       *string
		NextToken *string
	}

	cases := []struct {
		input *mockInput
		resps []*mockOutput
	}{
		{
			input: &mockInput{
				Foo: aws.String("foo"),
			},
			resps: []*mockOutput{
				{aws.String("1"), aws.String("token")},
				{aws.String("2"), aws.String("token")},
				{aws.String("3"), aws.String("")},
				{aws.String("4"), aws.String("token")},
			},
		},
		{
			input: &mockInput{
				Foo: aws.String("foo"),
			},
			resps: []*mockOutput{
				{aws.String("1"), aws.String("token")},
				{aws.String("2"), aws.String("token")},
				{aws.String("3"), nil},
				{aws.String("4"), aws.String("token")},
			},
		},
		{
			input: nil,
			resps: []*mockOutput{
				{aws.String("1"), aws.String("token")},
				{aws.String("2"), aws.String("token")},
				{aws.String("3"), nil},
				{aws.String("4"), aws.String("token")},
			},
		},
	}

	retryer := aws.DefaultRetryer{NumMaxRetries: 2}
	op := aws.Operation{
		Name: "Operation",
		Paginator: &aws.Paginator{
			InputTokens:  []string{"NextToken"},
			OutputTokens: []string{"NextToken"},
		},
	}

	for _, c := range cases {
		resps := c.resps
		input := c.input

		inValues := []string{}
		reqNum := 0
		p := aws.Pager{
			NewRequest: func() (*aws.Request, error) {
				h := defaults.Handlers()

				var inCpy *mockInput
				var tmp mockInput
				if input != nil {
					tmp = *input
					inCpy = &tmp
				}
				var output mockOutput
				req := aws.New(unit.Config(), aws.Metadata{}, h, retryer, &op, inCpy, &output)
				req.Handlers.Send.Clear()
				req.Handlers.Unmarshal.Clear()
				req.Handlers.UnmarshalMeta.Clear()
				req.Handlers.ValidateResponse.Clear()
				req.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
					r.Data = resps[reqNum]
					reqNum++
				})
				req.Handlers.Build.PushBack(func(r *aws.Request) {
					in := r.Params.(*mockInput)
					if in == nil {
						inValues = append(inValues, "")
					} else if in.Foo != nil {
						inValues = append(inValues, *in.Foo)
					}
				})
				req.SetContext(aws.BackgroundContext())

				return req, nil
			},
		}

		results := []*string{}
		for p.Next() {
			page := p.CurrentPage()
			output := page.(*mockOutput)
			results = append(results, output.Bar)
		}

		if err := p.Err(); err != nil {
			t.Error("unexpected error", err)
		}

		expected := []*string{
			aws.String("1"),
			aws.String("2"),
			aws.String("3"),
		}

		if e, a := expected, results; !awsutil.DeepEqual(e, a) {
			t.Log("\n-------------------------\nexpected\n-------------------------\n")
			for i, v := range e {
				t.Errorf("\t%d: %v", i, *v)
			}

			t.Log("\n\n-------------------------\nactual\n-------------------------\n")
			for i, v := range a {
				t.Errorf("\t%d: %v", i, *v)
			}
		}
	}
}

func TestPaginationTruncation(t *testing.T) {
	type mockInput struct {
		Foo *string
	}
	input := mockInput{
		Foo: aws.String("foo"),
	}

	type mockOutput struct {
		Bar         *string
		NextToken   *string
		IsTruncated *bool
	}

	resps := []*mockOutput{
		{aws.String("1"), aws.String("token"), aws.Bool(true)},
		{aws.String("2"), aws.String("token"), aws.Bool(true)},
		{aws.String("3"), aws.String(""), aws.Bool(false)},
		{aws.String("4"), aws.String(""), aws.Bool(true)},
	}

	reqNum := 0
	retryer := aws.DefaultRetryer{NumMaxRetries: 2}
	ops := []aws.Operation{
		{
			Name: "Operation",
			Paginator: &aws.Paginator{
				InputTokens:     []string{"NextToken"},
				OutputTokens:    []string{"NextToken"},
				TruncationToken: "IsTruncated",
			},
		},
		{
			Name: "Operation",
			Paginator: &aws.Paginator{
				InputTokens:     []string{"NextToken"},
				OutputTokens:    []string{"NextToken"},
				TruncationToken: "IsTruncated",
			},
		},
		{
			Name: "Operation",
			Paginator: &aws.Paginator{
				InputTokens:     []string{"NextToken"},
				OutputTokens:    []string{"NextToken"},
				TruncationToken: "IsTruncated",
			},
		},
		{
			Name: "Operation",
			Paginator: &aws.Paginator{
				InputTokens:     []string{"NextToken"},
				OutputTokens:    []string{"NextToken"},
				TruncationToken: "IsTruncated",
			},
		},
	}

	p := aws.Pager{
		NewRequest: func() (*aws.Request, error) {
			h := defaults.Handlers()

			tmp := input
			inCpy := &tmp
			op := ops[reqNum]

			var output mockOutput
			req := aws.New(unit.Config(), aws.Metadata{}, h, retryer, &op, inCpy, &output)
			req.Handlers.Send.Clear()
			req.Handlers.Unmarshal.Clear()
			req.Handlers.UnmarshalMeta.Clear()
			req.Handlers.ValidateResponse.Clear()
			req.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
				output := resps[reqNum]
				r.Data = output
				reqNum++
			})
			req.SetContext(aws.BackgroundContext())

			return req, nil
		},
	}

	results := []string{}
	for p.Next() {
		page := p.CurrentPage()
		output := page.(*mockOutput)
		results = append(results, *output.Bar)
	}

	if err := p.Err(); err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := 3, len(results); e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}

	if e, a := []string{"1", "2", "3"}, results; !awsutil.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
}

func BenchmarkPagination(b *testing.B) {
	type mockInput struct {
		Foo *string
	}

	type mockOutput struct {
		Bar       *string
		NextToken *string
	}

	input := &mockInput{
		Foo: aws.String("foo"),
	}
	resps := []*mockOutput{
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("1"), aws.String("token")},
		{aws.String("3"), aws.String("")},
	}

	retryer := aws.DefaultRetryer{NumMaxRetries: 2}
	op := aws.Operation{
		Name: "Operation",
		Paginator: &aws.Paginator{
			InputTokens:  []string{"NextToken"},
			OutputTokens: []string{"NextToken"},
		},
	}

	reqNum := 0
	p := aws.Pager{
		NewRequest: func() (*aws.Request, error) {
			h := defaults.Handlers()

			var inCpy *mockInput
			var tmp mockInput
			if input != nil {
				tmp = *input
				inCpy = &tmp
			}
			var output mockOutput
			req := aws.New(unit.Config(), aws.Metadata{}, h, retryer, &op, inCpy, &output)
			req.Handlers.Send.Clear()
			req.Handlers.Unmarshal.Clear()
			req.Handlers.UnmarshalMeta.Clear()
			req.Handlers.ValidateResponse.Clear()
			req.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
				r.Data = resps[reqNum]
				reqNum++
			})
			req.SetContext(aws.BackgroundContext())

			return req, nil
		},
	}

	for p.Next() {
		p.CurrentPage()
	}
}
