package aws_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Use DynamoDB methods for simplicity
func TestPaginationQueryPage(t *testing.T) {
	db := dynamodb.New(unit.Config())
	tokens, pages, numPages, gotToEnd := []map[string]dynamodb.AttributeValue{}, []map[string]dynamodb.AttributeValue{}, 0, false

	reqNum := 0
	resps := []*dynamodb.QueryOutput{
		{
			LastEvaluatedKey: map[string]dynamodb.AttributeValue{"key": {S: aws.String("key1")}},
			Count:            aws.Int64(1),
			Items: []map[string]dynamodb.AttributeValue{
				{
					"key": {S: aws.String("key1")},
				},
			},
		},
		{
			LastEvaluatedKey: map[string]dynamodb.AttributeValue{"key": {S: aws.String("key2")}},
			Count:            aws.Int64(1),
			Items: []map[string]dynamodb.AttributeValue{
				{
					"key": {S: aws.String("key2")},
				},
			},
		},
		{
			LastEvaluatedKey: map[string]dynamodb.AttributeValue{},
			Count:            aws.Int64(1),
			Items: []map[string]dynamodb.AttributeValue{
				{
					"key": {S: aws.String("key3")},
				},
			},
		},
	}

	db.Handlers.Send.Clear() // mock sending
	db.Handlers.Unmarshal.Clear()
	db.Handlers.UnmarshalMeta.Clear()
	db.Handlers.ValidateResponse.Clear()
	db.Handlers.Build.PushBack(func(r *aws.Request) {
		in := r.Params.(*dynamodb.QueryInput)
		if in == nil {
			tokens = append(tokens, nil)
		} else if len(in.ExclusiveStartKey) != 0 {
			tokens = append(tokens, in.ExclusiveStartKey)
		}
	})
	db.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = resps[reqNum]
		reqNum++
	})

	params := &dynamodb.QueryInput{
		Limit:     aws.Int64(2),
		TableName: aws.String("tablename"),
	}
	err := db.QueryPages(params, func(p *dynamodb.QueryOutput, last bool) bool {
		numPages++
		for _, item := range p.Items {
			pages = append(pages, item)
		}
		if last {
			if gotToEnd {
				t.Fatal("last == true, happened twice")
			}
			gotToEnd = true
		}
		return true
	})
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := []map[string]dynamodb.AttributeValue{
		{"key": {S: aws.String("key1")}},
		{"key": {S: aws.String("key2")}},
	}, tokens; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}

	if e, a := []map[string]dynamodb.AttributeValue{
		{"key": {S: aws.String("key1")}},
		{"key": {S: aws.String("key2")}},
		{"key": {S: aws.String("key3")}},
	}, pages; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}

	if e, a := 3, numPages; e != a {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if !gotToEnd {
		t.Error("did not reach of the end of the pagination")
	}
	if params.ExclusiveStartKey != nil {
		t.Errorf("expected a nil value, but received %v", params.ExclusiveStartKey)
	}
}

// Use DynamoDB methods for simplicity
func TestPagination(t *testing.T) {
	db := dynamodb.New(unit.Config())
	tokens, pages, numPages, gotToEnd := []string{}, []string{}, 0, false

	reqNum := 0
	resps := []*dynamodb.ListTablesOutput{
		{TableNames: []string{"Table1", "Table2"}, LastEvaluatedTableName: aws.String("Table2")},
		{TableNames: []string{"Table3", "Table4"}, LastEvaluatedTableName: aws.String("Table4")},
		{TableNames: []string{"Table5"}},
	}

	db.Handlers.Send.Clear() // mock sending
	db.Handlers.Unmarshal.Clear()
	db.Handlers.UnmarshalMeta.Clear()
	db.Handlers.ValidateResponse.Clear()
	db.Handlers.Build.PushBack(func(r *aws.Request) {
		in := r.Params.(*dynamodb.ListTablesInput)
		if in == nil {
			tokens = append(tokens, "")
		} else if in.ExclusiveStartTableName != nil {
			tokens = append(tokens, *in.ExclusiveStartTableName)
		}
	})
	db.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = resps[reqNum]
		reqNum++
	})

	params := &dynamodb.ListTablesInput{Limit: aws.Int64(2)}
	err := db.ListTablesPages(params, func(p *dynamodb.ListTablesOutput, last bool) bool {
		numPages++
		for _, t := range p.TableNames {
			pages = append(pages, t)
		}
		if last {
			if gotToEnd {
				t.Fatal("last == true, happened twice")
			}
			gotToEnd = true
		}
		return true
	})

	if e, a := []string{"Table2", "Table4"}, tokens; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if e, a := []string{"Table1", "Table2", "Table3", "Table4", "Table5"}, pages; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}

	if e, a := 3, numPages; e != a {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if !gotToEnd {
		t.Error("did not reach of the end of the pagination")
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if params.ExclusiveStartTableName != nil {
		t.Errorf("expected nil value, but received %v", params.ExclusiveStartTableName)
	}
}

// Use DynamoDB methods for simplicity
func TestPaginationEachPage(t *testing.T) {
	db := dynamodb.New(unit.Config())
	tokens, pages, numPages, gotToEnd := []string{}, []string{}, 0, false

	reqNum := 0
	resps := []*dynamodb.ListTablesOutput{
		{TableNames: []string{"Table1", "Table2"}, LastEvaluatedTableName: aws.String("Table2")},
		{TableNames: []string{"Table3", "Table4"}, LastEvaluatedTableName: aws.String("Table4")},
		{TableNames: []string{"Table5"}},
	}

	db.Handlers.Send.Clear() // mock sending
	db.Handlers.Unmarshal.Clear()
	db.Handlers.UnmarshalMeta.Clear()
	db.Handlers.ValidateResponse.Clear()
	db.Handlers.Build.PushBack(func(r *aws.Request) {
		in := r.Params.(*dynamodb.ListTablesInput)
		if in == nil {
			tokens = append(tokens, "")
		} else if in.ExclusiveStartTableName != nil {
			tokens = append(tokens, *in.ExclusiveStartTableName)
		}
	})
	db.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = resps[reqNum]
		reqNum++
	})

	params := &dynamodb.ListTablesInput{Limit: aws.Int64(2)}
	req := db.ListTablesRequest(params)
	err := req.EachPage(func(p interface{}, last bool) bool {
		numPages++
		for _, t := range p.(*dynamodb.ListTablesOutput).TableNames {
			pages = append(pages, t)
		}
		if last {
			if gotToEnd {
				t.Fatal("last == true, happened twice")
			}
			gotToEnd = true
		}

		return true
	})

	if e, a := []string{"Table2", "Table4"}, tokens; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if e, a := []string{"Table1", "Table2", "Table3", "Table4", "Table5"}, pages; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if e, a := 3, numPages; e != a {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if !gotToEnd {
		t.Error("did not reach of the end of the pagination")
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
}

// Use DynamoDB methods for simplicity
func TestPaginationEarlyExit(t *testing.T) {
	db := dynamodb.New(unit.Config())
	numPages, gotToEnd := 0, false

	reqNum := 0
	resps := []*dynamodb.ListTablesOutput{
		{TableNames: []string{"Table1", "Table2"}, LastEvaluatedTableName: aws.String("Table2")},
		{TableNames: []string{"Table3", "Table4"}, LastEvaluatedTableName: aws.String("Table4")},
		{TableNames: []string{"Table5"}},
	}

	db.Handlers.Send.Clear() // mock sending
	db.Handlers.Unmarshal.Clear()
	db.Handlers.UnmarshalMeta.Clear()
	db.Handlers.ValidateResponse.Clear()
	db.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = resps[reqNum]
		reqNum++
	})

	params := &dynamodb.ListTablesInput{Limit: aws.Int64(2)}
	err := db.ListTablesPages(params, func(p *dynamodb.ListTablesOutput, last bool) bool {
		numPages++
		if numPages == 2 {
			return false
		}
		if last {
			if gotToEnd {
				t.Fatal("last == true, happened twice")
			}
			gotToEnd = true
		}
		return true
	})

	if e, a := 2, numPages; e != a {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if gotToEnd {
		t.Error("should not have reached the end of pagination")
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
}

func TestSkipPagination(t *testing.T) {
	client := s3.New(unit.Config())
	client.Handlers.Send.Clear() // mock sending
	client.Handlers.Unmarshal.Clear()
	client.Handlers.UnmarshalMeta.Clear()
	client.Handlers.ValidateResponse.Clear()
	client.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = &s3.HeadBucketOutput{}
	})

	req := client.HeadBucketRequest(&s3.HeadBucketInput{Bucket: aws.String("bucket")})

	numPages, gotToEnd := 0, false
	req.EachPage(func(p interface{}, last bool) bool {
		numPages++
		if last {
			gotToEnd = true
		}
		return true
	})
	if e, a := 1, numPages; e != a {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if !gotToEnd {
		t.Error("did not reach of the end of the pagination")
	}
}

// Use S3 for simplicity
func TestPaginationTruncation(t *testing.T) {
	client := s3.New(unit.Config())

	reqNum := 0
	resps := []*s3.ListObjectsOutput{
		{IsTruncated: aws.Bool(true), Contents: []s3.Object{{Key: aws.String("Key1")}}},
		{IsTruncated: aws.Bool(true), Contents: []s3.Object{{Key: aws.String("Key2")}}},
		{IsTruncated: aws.Bool(false), Contents: []s3.Object{{Key: aws.String("Key3")}}},
		{IsTruncated: aws.Bool(true), Contents: []s3.Object{{Key: aws.String("Key4")}}},
	}

	client.Handlers.Send.Clear() // mock sending
	client.Handlers.Unmarshal.Clear()
	client.Handlers.UnmarshalMeta.Clear()
	client.Handlers.ValidateResponse.Clear()
	client.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = resps[reqNum]
		reqNum++
	})

	params := &s3.ListObjectsInput{Bucket: aws.String("bucket")}

	results := []string{}
	err := client.ListObjectsPages(params, func(p *s3.ListObjectsOutput, last bool) bool {
		results = append(results, *p.Contents[0].Key)
		return true
	})

	if e, a := []string{"Key1", "Key2", "Key3"}, results; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	// Try again without truncation token at all
	reqNum = 0
	resps[1].IsTruncated = nil
	resps[2].IsTruncated = aws.Bool(true)
	results = []string{}
	err = client.ListObjectsPages(params, func(p *s3.ListObjectsOutput, last bool) bool {
		results = append(results, *p.Contents[0].Key)
		return true
	})

	if e, a := []string{"Key1", "Key2"}, results; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
}

func TestPaginationNilToken(t *testing.T) {
	client := route53.New(unit.Config())

	reqNum := 0
	resps := []*route53.ListResourceRecordSetsOutput{
		{
			ResourceRecordSets: []route53.ResourceRecordSet{
				{Name: aws.String("first.example.com.")},
			},
			IsTruncated:          aws.Bool(true),
			NextRecordName:       aws.String("second.example.com."),
			NextRecordType:       route53.RRTypeMx,
			NextRecordIdentifier: aws.String("second"),
			MaxItems:             aws.String("1"),
		},
		{
			ResourceRecordSets: []route53.ResourceRecordSet{
				{Name: aws.String("second.example.com.")},
			},
			IsTruncated:    aws.Bool(true),
			NextRecordName: aws.String("third.example.com."),
			NextRecordType: route53.RRTypeMx,
			MaxItems:       aws.String("1"),
		},
		{
			ResourceRecordSets: []route53.ResourceRecordSet{
				{Name: aws.String("third.example.com.")},
			},
			IsTruncated: aws.Bool(false),
			MaxItems:    aws.String("1"),
		},
	}
	client.Handlers.Send.Clear() // mock sending
	client.Handlers.Unmarshal.Clear()
	client.Handlers.UnmarshalMeta.Clear()
	client.Handlers.ValidateResponse.Clear()

	idents := []string{}
	client.Handlers.Build.PushBack(func(r *aws.Request) {
		p := r.Params.(*route53.ListResourceRecordSetsInput)
		idents = append(idents, aws.StringValue(p.StartRecordIdentifier))

	})
	client.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = resps[reqNum]
		reqNum++
	})

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String("id-zone"),
	}

	results := []string{}
	err := client.ListResourceRecordSetsPages(params, func(p *route53.ListResourceRecordSetsOutput, last bool) bool {
		results = append(results, *p.ResourceRecordSets[0].Name)
		return true
	})

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if e, a := []string{"", "second", ""}, idents; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
	if e, a := []string{"first.example.com.", "second.example.com.", "third.example.com."}, results; !reflect.DeepEqual(e, a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
}

func TestPaginationNilInput(t *testing.T) {
	// Code generation doesn't have a great way to verify the code is correct
	// other than being run via unit tests in the SDK. This should be fixed
	// So code generation can be validated independently.

	client := s3.New(unit.Config())
	client.Handlers.Validate.Clear()
	client.Handlers.Send.Clear() // mock sending
	client.Handlers.Unmarshal.Clear()
	client.Handlers.UnmarshalMeta.Clear()
	client.Handlers.ValidateResponse.Clear()
	client.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = &s3.ListObjectsOutput{}
	})

	gotToEnd := false
	numPages := 0
	err := client.ListObjectsPages(nil, func(p *s3.ListObjectsOutput, last bool) bool {
		numPages++
		if last {
			gotToEnd = true
		}
		return true
	})

	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
	if e, a := 1, numPages; e != a {
		t.Errorf("expect %d number pages but got %d", e, a)
	}
	if !gotToEnd {
		t.Errorf("expect to of gotten to end, did not")
	}
}

func TestPaginationWithContextNilInput(t *testing.T) {
	// Code generation doesn't have a great way to verify the code is correct
	// other than being run via unit tests in the SDK. This should be fixed
	// So code generation can be validated independently.

	client := s3.New(unit.Config())
	client.Handlers.Validate.Clear()
	client.Handlers.Send.Clear() // mock sending
	client.Handlers.Unmarshal.Clear()
	client.Handlers.UnmarshalMeta.Clear()
	client.Handlers.ValidateResponse.Clear()
	client.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = &s3.ListObjectsOutput{}
	})

	gotToEnd := false
	numPages := 0
	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	err := client.ListObjectsPagesWithContext(ctx, nil, func(p *s3.ListObjectsOutput, last bool) bool {
		numPages++
		if last {
			gotToEnd = true
		}
		return true
	})

	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
	if e, a := 1, numPages; e != a {
		t.Errorf("expect %d number pages but got %d", e, a)
	}
	if !gotToEnd {
		t.Errorf("expect to of gotten to end, did not")
	}
}

func TestPagination_Standalone(t *testing.T) {
	type testPageInput struct {
		NextToken *string
	}
	type testPageOutput struct {
		Value     *string
		NextToken *string
	}
	type testCase struct {
		Value, PrevToken, NextToken *string
	}

	cases := map[string][]testCase{
		"nil NextToken": {
			testCase{aws.String("FirstValue"), aws.String("InitalToken"), aws.String("FirstToken")},
			testCase{aws.String("SecondValue"), aws.String("FirstToken"), aws.String("SecondToken")},
			testCase{aws.String("ThirdValue"), aws.String("SecondToken"), nil},
		},
		"zero NextToken": {
			testCase{aws.String("FirstValue"), aws.String("InitalToken"), aws.String("FirstToken")},
			testCase{aws.String("SecondValue"), aws.String("FirstToken"), aws.String("SecondToken")},
			testCase{aws.String("ThirdValue"), aws.String("SecondToken"), aws.String("")},
		},
	}

	for cName, c := range cases {
		t.Run(cName, func(t *testing.T) {
			input := testPageInput{
				NextToken: c[0].PrevToken,
			}

			svc := awstesting.NewClient(unit.Config())
			i := 0
			p := aws.Pagination{
				NewRequest: func() (*aws.Request, error) {
					r := svc.NewRequest(
						&aws.Operation{
							Name: "Operation",
							Paginator: &aws.Paginator{
								InputTokens:  []string{"NextToken"},
								OutputTokens: []string{"NextToken"},
							},
						},
						&input, &testPageOutput{},
					)
					// Setup handlers for testing
					r.Handlers.Clear()
					r.Handlers.Build.PushBack(func(req *aws.Request) {
						if e, a := len(c), i; a >= e {
							t.Fatalf("expect no more than %d requests, got %d", e, a)
						}
						in := req.Params.(*testPageInput)
						if e, a := aws.StringValue(c[i].PrevToken), aws.StringValue(in.NextToken); e != a {
							t.Errorf("%d, expect NextToken input %q, got %q", i, e, a)
						}
					})
					r.Handlers.Unmarshal.PushBack(func(req *aws.Request) {
						out := &testPageOutput{
							Value: c[i].Value,
						}
						if c[i].NextToken != nil {
							next := *c[i].NextToken
							out.NextToken = aws.String(next)
						}
						req.Data = out
					})
					return r, nil
				},
			}

			for p.Next() {
				data := p.Page().(*testPageOutput)

				if e, a := aws.StringValue(c[i].Value), aws.StringValue(data.Value); e != a {
					t.Errorf("%d, expect Value to be %q, got %q", i, e, a)
				}
				if e, a := aws.StringValue(c[i].NextToken), aws.StringValue(data.NextToken); e != a {
					t.Errorf("%d, expect NextToken to be %q, got %q", i, e, a)
				}

				i++
			}
			if e, a := len(c), i; e != a {
				t.Errorf("expected to process %d pages, did %d", e, a)
			}
			if err := p.Err(); err != nil {
				t.Fatalf("%d, expected no error, got %v", i, err)
			}
		})
	}
}

// Benchmarks
var benchResps = []dynamodb.ListTablesOutput{
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE", "NXT"}, LastEvaluatedTableName: aws.String("NXT")},
	{TableNames: []string{"TABLE"}},
}

var benchDb = func() *dynamodb.DynamoDB {
	db := dynamodb.New(unit.Config())
	db.Handlers.Send.Clear() // mock sending
	db.Handlers.Unmarshal.Clear()
	db.Handlers.UnmarshalMeta.Clear()
	db.Handlers.ValidateResponse.Clear()
	return db
}

func BenchmarkCodegenIterator(b *testing.B) {
	reqNum := 0
	db := benchDb()
	db.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = benchResps[reqNum]
		reqNum++
	})

	input := &dynamodb.ListTablesInput{Limit: aws.Int64(2)}
	iter := func(fn func(*dynamodb.ListTablesOutput, bool) bool) error {
		page := db.ListTablesRequest(input)
		for ; page.Request != nil; page.Request = page.NextPage() {
			page.Send()
			out := page.Data.(*dynamodb.ListTablesOutput)
			if result := fn(out, !page.HasNextPage()); page.Error != nil || !result {
				return page.Error
			}
		}
		return nil
	}

	for i := 0; i < b.N; i++ {
		reqNum = 0
		iter(func(p *dynamodb.ListTablesOutput, last bool) bool {
			return true
		})
	}
}

func BenchmarkEachPageIterator(b *testing.B) {
	reqNum := 0
	db := benchDb()
	db.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		r.Data = benchResps[reqNum]
		reqNum++
	})

	input := &dynamodb.ListTablesInput{Limit: aws.Int64(2)}
	for i := 0; i < b.N; i++ {
		reqNum = 0
		req := db.ListTablesRequest(input)
		req.EachPage(func(p interface{}, last bool) bool {
			return true
		})
	}
}
