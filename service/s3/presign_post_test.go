package s3

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

func TestPresignPutObject(t *testing.T) {
	fixedTime := time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC)
	defer mockTime(fixedTime)()

	cases := map[string]struct {
		input            PutObjectInput
		options          []func(*PresignPostOptions)
		expectedExpires  time.Time
		expectedURL      string
		region           string
		pathStyleEnabled bool
		BaseEndpoint     string
	}{
		"sample": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
		},
		"bucket and key have the same value": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("bucket"),
			},
		},
		"expires override": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			expectedExpires: fixedTime.Add(5 * time.Minute),
			options: []func(o *PresignPostOptions){
				func(o *PresignPostOptions) {
					o.Expires = 5 * time.Minute
				},
			},
		},
		"body is ignored": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				// This will be ignored
				Body: strings.NewReader("hello-world"),
			},
		},
		"different region": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			region:      "eu-central-1",
			expectedURL: "https://bucket.s3.eu-central-1.amazonaws.com",
		},
		"mrap endpoint is changed": {
			input: PutObjectInput{
				Bucket: aws.String("arn:aws:s3::123456789012:accesspoint:mfzwi23gnjvgw.mrap"),
				Key:    aws.String("mockkey"),
			},
			expectedURL: "https://mfzwi23gnjvgw.mrap.accesspoint.s3-global.amazonaws.com",
		},
		"use path style bucket hosting pattern": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			expectedURL:      "https://s3.us-west-2.amazonaws.com/bucket",
			pathStyleEnabled: true,
		},
		"use path style bucket and key have the same value ": {
			input: PutObjectInput{
				Bucket: aws.String("value"),
				Key:    aws.String("value"),
			},
			expectedURL:      "https://s3.us-west-2.amazonaws.com/value",
			pathStyleEnabled: true,
		},
		"use path style bucket with custom baseEndpoint": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			expectedURL:      "https://s3.custom-domain.com/bucket",
			pathStyleEnabled: true,
			BaseEndpoint:     "https://s3.custom-domain.com",
		},
		"use path style bucket with custom baseEndpoint with path": {
			input: PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			BaseEndpoint:     "https://my-custom-domain.com/path_my_path",
			pathStyleEnabled: true,
			expectedURL:      "https://my-custom-domain.com/path_my_path/bucket",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			region := "us-west-2"
			if tc.region != "" {
				region = tc.region
			}
			cfg := aws.Config{
				Region:      region,
				Credentials: unit.StubCredentialsProvider{},
				Retryer: func() aws.Retryer {
					return aws.NopRetryer{}
				},
			}
			presignClient := NewPresignClient(NewFromConfig(cfg, func(options *Options) {
				options.UsePathStyle = tc.pathStyleEnabled
				if tc.BaseEndpoint != "" {
					options.BaseEndpoint = aws.String(tc.BaseEndpoint)
				}
			}))
			postObject, err := presignClient.PresignPostObject(ctx, &tc.input, tc.options...)
			if err != nil {
				t.Error(err)
			}
			if postObject == nil {
				t.Error("expected non-nil postObject")
			}
			if tc.expectedURL != "" {
				if tc.expectedURL != postObject.URL {
					t.Errorf("expected URL %q; got %q", tc.expectedURL, postObject.URL)
				}
			} else {
				if "https://bucket.s3.us-west-2.amazonaws.com" != postObject.URL {
					t.Error("expected URL to contain 'https://amazon.com', was: ", postObject.URL)
				}
			}

			if len(postObject.Values) < 1 {
				t.Error("expected non-empty values")
			}
			policy, ok := postObject.Values["policy"]
			if !ok {
				t.Error("expected non-empty policy on postObject")
			}
			decoded, err := base64.StdEncoding.DecodeString(policy)
			if err != nil {
				t.Error("expected base64 encoded policy, got error", err, "policy", policy)
			}
			var policyJSON map[string]interface{}
			err = json.Unmarshal(decoded, &policyJSON)
			if err != nil {
				t.Error("expected valid JSON for policy, got error", err, "with policy", policy)
			}
			actualExpires, ok := policyJSON["expiration"]
			if !ok {
				t.Error("expected non-empty expiration on policy JSON policy", policyJSON)
			}

			if !time.Time.IsZero(tc.expectedExpires) {
				isEqual, err := isTimeEqual(actualExpires.(string), tc.expectedExpires)
				if err != nil {
					t.Error("Error parsing expires", actualExpires, err)
				}
				if !isEqual {
					t.Error("expected expiration to be", tc.expectedExpires, "got", actualExpires)
				}
			} else {
				// Check the default is set. Go serializes JSON values as RFC3339
				expectedExpires := fixedTime.Add(15 * time.Minute).Format(time.RFC3339)
				if actualExpires != expectedExpires {
					t.Error("expected expiration to be", expectedExpires, "got", actualExpires)
				}
			}
		})
	}
}

// Test that comes straight from the docs https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-post-example.html
// Unfortunately it can't be verified with the exact same values
// since the sample in the docs lowercases all headers `x-amzn-header`
// while the SDK does not `X-Amzn-Header`, so the signature and policy are different.
// However, the values have been manually inspected to match the desired output
func TestSampleFromPublicDocs(t *testing.T) {
	accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	bucket := "sigv4examplebucket"
	key := "user/user1"
	testTime := time.Date(2015, time.December, 29, 0, 0, 0, 0, time.UTC)
	defer mockTime(testTime)()
	expiresIn := 36 * time.Hour
	staticCredentials := staticCredentialsProvider{Key: accessKeyID, Secret: secretAccessKey}
	ctx := context.Background()

	cfg := aws.Config{
		Region:      "us-east-1",
		Credentials: staticCredentials,
		Retryer: func() aws.Retryer {
			return aws.NopRetryer{}
		},
	}

	presignClient := NewPresignClient(NewFromConfig(cfg))
	input := PutObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}
	conditions := []interface{}{
		[]interface{}{"starts-with", "$key", "user/user1/"},
		map[string]string{"acl": "public-read"},
		map[string]string{"success_action_redirect": "http://sigv4examplebucket.s3.amazonaws.com/successful_upload.html"},
		[]interface{}{"starts-with", "$Content-Type", "image/"},
		map[string]string{"x-amz-meta-uuid": "14365123651274"},
		[]interface{}{"starts-with", "$x-amz-meta-tag", ""},
	}
	opts := func(o *PresignPostOptions) {
		o.Expires = expiresIn
		o.Conditions = conditions
	}
	postObject, err := presignClient.PresignPostObject(ctx, &input, opts)
	if err != nil {
		t.Error(err)
	}
	if postObject == nil {
		t.Error("expected non-nil postObject")
	}
	values := postObject.Values
	signature, ok := values["X-Amz-Signature"]
	if !ok {
		t.Error("expected non-empty signature on postObject", values)
	}
	// Signature and policy are VERY sensitive to any change in output or order. If these tests fail,
	// it can be due to a change in order for the policy or a change in capitalization
	if signature != "41eb7f468113e77dca133475d38815dbe1f92b073964f4a0575f036e9c02d28a" {
		t.Error("expected signature to equal to be precomputed", signature, "got", values)
	}
	policy, ok := values["policy"]
	if !ok {
		t.Error("expected non-empty policy on values", values)
	}
	expectedPolicy := "eyJjb25kaXRpb25zIjpbeyJYLUFtei1BbGdvcml0aG0iOiJBV1M0LUhNQUMtU0hBMjU2In0seyJidWN" +
		"rZXQiOiJzaWd2NGV4YW1wbGVidWNrZXQifSx7IlgtQW16LUNyZWRlbnRpYWwiOiJBS0lBSU9TRk9ETk" +
		"43RVhBTVBMRS8yMDE1MTIyOS91cy1lYXN0LTEvczMvYXdzNF9yZXF1ZXN0In0seyJYLUFtei1EYXRlI" +
		"joiMjAxNTEyMjlUMDAwMDAwWiJ9LFsic3RhcnRzLXdpdGgiLCIka2V5IiwidXNlci91c2VyMS8iXSx7" +
		"ImFjbCI6InB1YmxpYy1yZWFkIn0seyJzdWNjZXNzX2FjdGlvbl9yZWRpcmVjdCI6Imh0dHA6Ly9zaWd" +
		"2NGV4YW1wbGVidWNrZXQuczMuYW1hem9uYXdzLmNvbS9zdWNjZXNzZnVsX3VwbG9hZC5odG1sIn0sWy" +
		"JzdGFydHMtd2l0aCIsIiRDb250ZW50LVR5cGUiLCJpbWFnZS8iXSx7IngtYW16LW1ldGEtdXVpZCI6I" +
		"jE0MzY1MTIzNjUxMjc0In0sWyJzdGFydHMtd2l0aCIsIiR4LWFtei1tZXRhLXRhZyIsIiJdXSwiZXhw" +
		"aXJhdGlvbiI6IjIwMTUtMTItMzBUMTI6MDA6MDBaIn0="
	if policy != expectedPolicy {
		t.Error("expected policy to equal", expectedPolicy, "got", policy)
	}
}

func TestBuildPresignPostRequest(t *testing.T) {
	cases := map[string]struct {
		credentials       aws.Credentials
		extraConditions   []interface{}
		isKeyConditionSet bool
	}{
		"credentials without access token": {
			credentials:     credentialsNoToken,
			extraConditions: []interface{}{},
		},
		"credentials with access token": {
			credentials:     credentialsWithToken,
			extraConditions: []interface{}{},
		},
		"no extra conditions": {
			credentials:     credentialsWithToken,
			extraConditions: []interface{}{},
		},
		"extra conditions": {
			credentials: credentialsWithToken,
			extraConditions: []interface{}{
				map[string]string{"acl": "public-read"},
				[]string{"starts-with", "$Content-Type", "image/"},
			},
		},
		"extra conditions collision": {
			credentials: credentialsWithToken,
			extraConditions: []interface{}{
				map[string]string{"bucket": "otherBucket"},
			},
		},
		"a key condition is set, no extra one is generated": {
			credentials: credentialsNoToken,
			extraConditions: []interface{}{
				[]interface{}{"starts-with", "$key", "user/user1/"},
			},
			isKeyConditionSet: true,
		},
	}
	requiredFields := []string{
		"X-Amz-Algorithm",
		"X-Amz-Credential",
		"X-Amz-Date",
		"X-Amz-Signature",
		"key",
		"policy",
	}

	requiredConditions := []string{"X-Amz-Algorithm", "bucket", "X-Amz-Credential", "X-Amz-Date"}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			target := postSignAdapter{}
			aBucketKey := "someKey"
			bucket := "someBucket"
			signingTime := sdk.NowTime()
			expiration := signingTime.Add(time.Hour)
			fields, err := target.PresignPost(tc.credentials, bucket, aBucketKey, "region", "service", signingTime, tc.extraConditions, expiration)
			if err != nil {
				t.Errorf("PresignPostHTTP returned unexepected error: %s", err.Error())
			}
			if len(fields) == 0 {
				t.Errorf("PresignPostHTTP returned no fields")
			}

			for _, field := range requiredFields {
				_, ok := fields[field]
				if !ok {
					t.Errorf("Fields response did not contain required key %s. Res %v", field, fields)
				}
			}

			if tc.credentials.SessionToken != "" {
				_, ok := fields["X-Amz-Security-Token"]
				if !ok {
					t.Errorf("Credentials are using a session token, but is not set on the fields response")
				}
			}

			actualKey := fields["key"]
			if actualKey != aBucketKey {
				t.Errorf("PresignPostHTTP did not contain expected \"key\" %s. Has %s", aBucketKey, actualKey)
			}
			policy := fields["policy"]
			decoded, err := base64.StdEncoding.DecodeString(policy)
			if err != nil {
				t.Errorf("Decoding policy document %s failed with error %v", policy, err)
			}
			var doc map[string]interface{}
			err = json.Unmarshal(decoded, &doc)
			if err != nil {
				t.Errorf("Policy document %s failed to parse to JSON with error %v", policy, err)
			}
			_, ok := doc["conditions"]
			if !ok {
				t.Errorf("Conditions field not present in policy document %s", policy)
			}
			exp, ok := doc["expiration"]
			if !ok {
				t.Errorf("Expiration field not present in policy document %s", policy)
			}
			docExpiration, ok := exp.(string)
			if !ok {
				t.Errorf("Expiration field is not a time as expected, is %v", doc["expiration"])
			}
			isEqual, err := isTimeEqual(docExpiration, expiration)
			if err != nil {
				t.Errorf("PresignPost did not parse expiration time %s. Error %v", docExpiration, err)
			}
			if !isEqual {
				t.Errorf("Expected policy expiration to be %v. Got %v", expiration, docExpiration)
			}
			conditions := doc["conditions"].([]interface{})
			if len(conditions) == 0 {
				t.Errorf("Policy document didn't contain any conditions")
			}
			for _, required := range requiredConditions {
				val := findInSlice(conditions, required)
				if val == nil {
					t.Errorf("Policy document didn't contain required conditions %s. Has %v", required, conditions)
				}
			}
			actualBucket := findInSlice(conditions, "bucket")
			if !reflect.DeepEqual(bucket, actualBucket) {
				t.Errorf("Expected bucket to be %v, was %v", bucket, actualBucket)
			}
			actualDate := findInSlice(conditions, "X-Amz-Date")
			signingTimeStr := signingTime.UTC().Format("20060102T150405Z")
			if signingTimeStr != actualDate {
				t.Errorf("Expected date to be %v, was %v", signingTimeStr, actualDate)
			}
			if len(tc.extraConditions) > 0 {
				for _, ec := range tc.extraConditions {
					if !isPresent(ec, conditions) {
						t.Errorf("Expected item %v not found on conditions %v", ec, conditions)
					}
				}
			}
			if !tc.isKeyConditionSet {
				// check the default is set
				conditionKey := findInSlice(conditions, "key")
				if conditionKey == nil {
					t.Errorf("Expected Condition 'key' to be set on policy conditions, none found. Conditions %v", conditions)
				}
				actualVal, ok := conditionKey.(string)
				if !ok {
					t.Errorf("Expected condition key to be a string, was %v", conditionKey)
				}
				if actualVal != aBucketKey {
					t.Errorf("Expected bucket key to be %v, was %v", aBucketKey, actualVal)
				}
			} else {
				// check the key condition is not set
				conditionKey := findInSlice(conditions, "key")
				if conditionKey != nil {
					t.Errorf("Expected condition key to be nil since %v was set, was %v", tc.isKeyConditionSet, conditionKey)
				}

			}
		})
	}
}

func mockTime(t time.Time) func() {
	sdk.NowTime = func() time.Time { return t }
	return func() { sdk.NowTime = time.Now }
}

type staticCredentialsProvider struct {
	Key    string
	Secret string
}

func (p staticCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{AccessKeyID: p.Key, SecretAccessKey: p.Secret}, nil
}

var credentialsNoToken = aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}
var credentialsWithToken = aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION"}

func isPresent(needle interface{}, haystack []interface{}) bool {
	needleValue := reflect.ValueOf(needle)
	for _, item := range haystack {
		itemValue := reflect.ValueOf(item)

		// special checks for slices and maps, since interface{} are not typecasted
		// by reflect.DeepEquals
		isSlice := itemValue.Kind() == reflect.Slice && needleValue.Kind() == reflect.Slice
		if isSlice && areSlicesEqual(needleValue, itemValue) {
			return true
		}
		isMap := itemValue.Kind() == reflect.Map && needleValue.Kind() == reflect.Map
		if isMap && areMapsEqual(needleValue, itemValue) {
			return true
		}

		// else do a regular deep equal check
		if reflect.DeepEqual(item, needle) {
			return true
		}
	}
	return false
}

func areSlicesEqual(a reflect.Value, b reflect.Value) bool {
	if a.Len() != b.Len() {
		return false
	}

	for i := 0; i < a.Len(); i++ {
		aValue := a.Index(i).Interface()
		bValue := b.Index(i).Interface()

		if !reflect.DeepEqual(aValue, bValue) {
			return false
		}
	}

	return true
}

func areMapsEqual(aVal reflect.Value, bVal reflect.Value) bool {
	// Check if 'a' is a map
	if aVal.Kind() != reflect.Map {
		return false
	}

	// Check if both maps have the same number of keys
	if aVal.Len() != bVal.Len() {
		return false
	}

	// Iterate over the keys and values in the first map
	for _, key := range aVal.MapKeys() {
		aValue := aVal.MapIndex(key)
		if !aValue.IsValid() {
			return false
		}
		bValue := bVal.MapIndex(key)
		if !bValue.IsValid() {
			return false
		}

		// Compare values using reflect.DeepEqual
		if !reflect.DeepEqual(aValue.Interface(), bValue.Interface()) {
			return false
		}
	}
	return true
}

// filters items in slice that have a map[string]interface{} and returns
// the first items map that has the key from "key"
func findInSlice(slice []interface{}, key string) interface{} {
	for _, item := range slice {
		// filter only the values with keys. Ignore stuff like arrays
		if v, ok := item.(map[string]interface{}); ok {
			// once in the maps, check if they have the desired key
			if _, ok := v[key]; ok {
				return v[key]
			}
		}
	}
	return nil
}

func isTimeEqual(t1s string, t2 time.Time) (bool, error) {
	t1, err := time.Parse(time.RFC3339, t1s)
	if err != nil {
		return false, err
	}
	areEqual := t1.Format(time.RFC3339) == t2.Format(time.RFC3339)
	return areEqual, nil
}
