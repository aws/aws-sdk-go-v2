package converters

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestTimePtrConverter_FromAttributeValue(t *testing.T) {
	cases := []struct {
		input          types.AttributeValue
		opts           []string
		expectedOutput any
		expectedError  bool
	}{
		{input: &types.AttributeValueMemberN{Value: "1136214245"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		{input: &types.AttributeValueMemberN{Value: "1136214245.113621424"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 113621424).In(time.UTC)), expectedError: false},
		{input: &types.AttributeValueMemberS{Value: "01/02 03:04:05PM #06 +0000"}, opts: []string{"TZ=UTC", "format=01/02 03:04:05PM #06 -0700"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.Layout
		{input: &types.AttributeValueMemberS{Value: "01/02 03:04:05PM '06 +0000"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.ANSIC
		{input: &types.AttributeValueMemberS{Value: "Mon Jan  2 15:04:05 2006"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.UnixDate
		{input: &types.AttributeValueMemberS{Value: "Mon Jan 02 15:04:05 UTC 2006"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.RubyDate
		{input: &types.AttributeValueMemberS{Value: "Mon Jan 02 15:04:05 +0000 2006"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.RFC822
		{input: &types.AttributeValueMemberS{Value: "02 Jan 06 15:04 UTC"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214240, 0).In(time.UTC)), expectedError: false},
		//time.RFC822Z
		{input: &types.AttributeValueMemberS{Value: "02 Jan 06 15:04 +0000"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214240, 0).In(time.UTC)), expectedError: false},
		//time.RFC850
		{input: &types.AttributeValueMemberS{Value: "Monday, 02-Jan-06 15:04:05 UTC"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.RFC1123
		{input: &types.AttributeValueMemberS{Value: "Mon, 02 Jan 2006 15:04:05 UTC"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.RFC1123Z
		{input: &types.AttributeValueMemberS{Value: "Mon, 02 Jan 2006 15:04:05 +0000"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.RFC3339
		{input: &types.AttributeValueMemberS{Value: "2006-01-02T15:04:05+00:00"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.RFC3339Nano
		{input: &types.AttributeValueMemberS{Value: "2006-01-02T15:04:05.999999999+00:00"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 999999999).In(time.UTC)), expectedError: false},
		//time.Kitchen
		{input: &types.AttributeValueMemberS{Value: "3:04PM"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(-62167164960, 0).In(time.UTC)), expectedError: false},
		//time.Stamp
		{input: &types.AttributeValueMemberS{Value: "Jan 02 15:04:05"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(-62167078555, 0).In(time.UTC)), expectedError: false},
		//time.StampMilli
		{input: &types.AttributeValueMemberS{Value: "Jan 02 15:04:05.000"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(-62167078555, 0).In(time.UTC)), expectedError: false},
		//time.StampMicro
		{input: &types.AttributeValueMemberS{Value: "Jan 02 15:04:05.000000"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(-62167078555, 0).In(time.UTC)), expectedError: false},
		//time.StampNano
		{input: &types.AttributeValueMemberS{Value: "Jan 02 15:04:05.000000000"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(-62167078555, 0).In(time.UTC)), expectedError: false},
		//time.DateTime
		{input: &types.AttributeValueMemberS{Value: "2006-01-02 15:04:05"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), expectedError: false},
		//time.DateOnly
		{input: &types.AttributeValueMemberS{Value: "2006-01-02"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(1136160000, 0).In(time.UTC)), expectedError: false},
		//time.TimeOnly
		{input: &types.AttributeValueMemberS{Value: "15:04:05"}, opts: []string{"TZ=UTC"}, expectedOutput: aws.Time(time.Unix(-62167164955, 0).In(time.UTC)), expectedError: false},
		// errors
		{input: nil, opts: nil, expectedOutput: nil, expectedError: true},
		{input: (types.AttributeValue)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: (*types.AttributeValueMemberN)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: (*types.AttributeValueMemberS)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: (*types.AttributeValueMemberBOOL)(nil), opts: nil, expectedOutput: nil, expectedError: true},
		{input: &types.AttributeValueMemberBOOL{Value: true}, opts: nil, expectedOutput: nil, expectedError: true},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc := TimePtrConverter{}

			actualOutput, actualError := tc.FromAttributeValue(c.input, c.opts)

			if actualError == nil && c.expectedError {
				t.Fatalf("expected error, got none")
			}

			if actualError != nil && !c.expectedError {
				t.Fatalf("unexpected error, got: %v", actualError)
			}

			if actualError != nil && c.expectedError {
				return
			}

			if !reflect.DeepEqual(c.expectedOutput, actualOutput) {
				t.Fatalf("%#+v != %#+v", c.expectedOutput, actualOutput)
			}
		})
	}
}

func TestTimePtrConverter_ToAttributeValue(t *testing.T) {
	cases := []struct {
		input          *time.Time
		opts           []string
		expectedOutput types.AttributeValue
		expectedError  bool
	}{
		{input: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), opts: []string{"TZ=UTC"}, expectedOutput: &types.AttributeValueMemberS{Value: time.Unix(1136214245, 0).In(time.UTC).Format(time.RFC3339)}, expectedError: false},
		{input: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), opts: []string{"TZ=UTC", "as=number"}, expectedOutput: &types.AttributeValueMemberN{Value: "1136214245"}, expectedError: false},
		{input: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), opts: []string{"TZ=UTC", "as=string"}, expectedOutput: &types.AttributeValueMemberS{Value: time.Unix(1136214245, 0).In(time.UTC).Format(time.RFC3339)}, expectedError: false},
		{input: aws.Time(time.Unix(1136214245, 113621424).In(time.UTC)), opts: []string{"TZ=UTC", "as=number"}, expectedOutput: &types.AttributeValueMemberN{Value: "1136214245.113621424"}, expectedError: false},
		{input: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), opts: []string{"TZ=UTC", "format=01/02 03:04:05PM #06 -0700"}, expectedOutput: &types.AttributeValueMemberS{Value: "01/02 03:04:05PM #06 +0000"}, expectedError: false},
		{input: aws.Time(time.Unix(1136214245, 0).In(time.UTC)), opts: []string{"TZ=UTC", "format=01/02 03:04:05PM 06 -0700"}, expectedOutput: &types.AttributeValueMemberS{Value: "01/02 03:04:05PM 06 +0000"}, expectedError: false},
		{input: &time.Time{}, opts: nil, expectedOutput: &types.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"}, expectedError: false},
		// errors
		{input: nil, opts: nil, expectedOutput: nil, expectedError: true},
		{input: (*time.Time)(nil), opts: nil, expectedOutput: (types.AttributeValue)(nil), expectedError: true},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tc := TimePtrConverter{}

			actualOutput, actualError := tc.ToAttributeValue(c.input, c.opts)

			if actualError == nil && c.expectedError {
				t.Fatalf("expected error, got none")
			}

			if actualError != nil && !c.expectedError {
				t.Fatalf("unexpected error, got: %v", actualError)
			}

			if actualError != nil && c.expectedError {
				return
			}

			if !reflect.DeepEqual(c.expectedOutput, actualOutput) {
				t.Fatalf("%#+v != %#+v", c.expectedOutput, actualOutput)
			}
		})
	}
}
