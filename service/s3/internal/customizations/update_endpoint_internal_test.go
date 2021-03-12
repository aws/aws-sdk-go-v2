package customizations

import (
	"net/url"
	"strconv"
	"testing"
)

func TestRemoveBucketFromPath(t *testing.T) {
	cases := []struct {
		url      url.URL
		bucket   string
		expected string
	}{
		{
			url: url.URL{
				Scheme:  "https",
				Host:    "amazonaws.com",
				Path:    "/bucket-name/key/path",
				RawPath: "/bucket-name/key/path",
			},
			bucket:   "bucket-name",
			expected: "https://amazonaws.com/key/path",
		},
		{
			url: url.URL{
				Scheme:  "https",
				Host:    "amazonaws.com",
				Path:    "/bucket-name/key/path/with/bucket-name",
				RawPath: "/bucket-name/key/path/with/bucket-name",
			},
			bucket:   "bucket-name",
			expected: "https://amazonaws.com/key/path/with/bucket-name",
		},
		{
			url: url.URL{
				Scheme:  "https",
				Host:    "amazonaws.com",
				Path:    "/arn:aws:s3:us-east-1:012345678901:accesspoint:myap/key/path?isEscaped=true",
				RawPath: "/arn%3Aaws%3As3%3Aus-east-1%3A012345678901%3Aaccesspoint%3Amyap/key/path%3FisEscaped%3Dtrue",
			},
			bucket:   "arn:aws:s3:us-east-1:012345678901:accesspoint:myap",
			expected: "https://amazonaws.com/key/path%3FisEscaped%3Dtrue",
		},
		{
			url: url.URL{
				Scheme:  "https",
				Host:    "amazonaws.com",
				Path:    "/path/to/key",
				RawPath: "/path/to/key",
			},
			bucket:   "not-a-match",
			expected: "https://amazonaws.com/path/to/key",
		},
		{
			url: url.URL{
				Scheme:  "https",
				Host:    "amazonaws.com",
				Path:    "",
				RawPath: "",
			},
			bucket:   "not-a-match",
			expected: "https://amazonaws.com/",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			removeBucketFromPath(&tt.url, tt.bucket)

			if e, a := tt.expected, tt.url.String(); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
