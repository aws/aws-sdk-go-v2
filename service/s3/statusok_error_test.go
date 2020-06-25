package s3_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/awserr"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/service/s3"
)

const (
	errMsg         = `<Error><Code>ErrorCode</Code><Message>message body</Message><RequestId>requestID</RequestId><HostId>hostID=</HostId></Error>`
	xmlPreambleMsg = `<?xml version="1.0" encoding="UTF-8"?>`
)

var lastModifiedTime = time.Date(2009, 11, 23, 0, 0, 0, 0, time.UTC)

func TestCopyObjectNoError(t *testing.T) {
	const successMsg = `
<?xml version="1.0" encoding="UTF-8"?>
<CopyObjectResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><LastModified>2009-11-23T0:00:00Z</LastModified><ETag>&quot;1da64c7f13d1e8dbeaea40b905fd586c&quot;</ETag></CopyObjectResult>`

	req := newCopyTestSvc(successMsg).CopyObjectRequest(&s3.CopyObjectInput{
		Bucket:     aws.String("bucketname"),
		CopySource: aws.String("bucketname/exists.txt"),
		Key:        aws.String("destination.txt"),
	})
	res, err := req.Send(context.Background())

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if e, a := fmt.Sprintf(`%q`, "1da64c7f13d1e8dbeaea40b905fd586c"), *res.CopyObjectResult.ETag; e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := lastModifiedTime, *res.CopyObjectResult.LastModified; !e.Equal(a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
}

func TestCopyObjectError(t *testing.T) {
	req := newCopyTestSvc(errMsg).CopyObjectRequest(&s3.CopyObjectInput{
		Bucket:     aws.String("bucketname"),
		CopySource: aws.String("bucketname/doesnotexist.txt"),
		Key:        aws.String("destination.txt"),
	})
	_, err := req.Send(context.Background())

	if err == nil {
		t.Error("expected error, but received none")
	}
	e := err.(awserr.Error)

	if e, a := "ErrorCode", e.Code(); e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := "message body", e.Message(); e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
}

func TestUploadPartCopySuccess(t *testing.T) {
	const successMsg = `
<?xml version="1.0" encoding="UTF-8"?>
<UploadPartCopyResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><LastModified>2009-11-23T0:00:00Z</LastModified><ETag>&quot;1da64c7f13d1e8dbeaea40b905fd586c&quot;</ETag></UploadPartCopyResult>`

	req := newCopyTestSvc(successMsg).UploadPartCopyRequest(&s3.UploadPartCopyInput{
		Bucket:     aws.String("bucketname"),
		CopySource: aws.String("bucketname/doesnotexist.txt"),
		Key:        aws.String("destination.txt"),
		PartNumber: aws.Int64(0),
		UploadId:   aws.String("uploadID"),
	})
	res, err := req.Send(context.Background())

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := fmt.Sprintf(`%q`, "1da64c7f13d1e8dbeaea40b905fd586c"), *res.CopyPartResult.ETag; e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := lastModifiedTime, *res.CopyPartResult.LastModified; !e.Equal(a) {
		t.Errorf("expected %v, but received %v", e, a)
	}
}

func TestUploadPartCopyError(t *testing.T) {
	req := newCopyTestSvc(errMsg).UploadPartCopyRequest(&s3.UploadPartCopyInput{
		Bucket:     aws.String("bucketname"),
		CopySource: aws.String("bucketname/doesnotexist.txt"),
		Key:        aws.String("destination.txt"),
		PartNumber: aws.Int64(0),
		UploadId:   aws.String("uploadID"),
	})
	_, err := req.Send(context.Background())

	if err == nil {
		t.Error("expected an error")
	}
	e := err.(awserr.Error)

	if e, a := "ErrorCode", e.Code(); e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := "message body", e.Message(); e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
}

func TestCompleteMultipartUploadSuccess(t *testing.T) {
	const successMsg = `
<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Location>locationName</Location><Bucket>bucketName</Bucket><Key>keyName</Key><ETag>"etagVal"</ETag></CompleteMultipartUploadResult>`
	req := newCopyTestSvc(successMsg).CompleteMultipartUploadRequest(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String("bucketname"),
		Key:      aws.String("key"),
		UploadId: aws.String("uploadID"),
	})
	res, err := req.Send(context.Background())

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := `"etagVal"`, *res.ETag; e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := "bucketName", *res.Bucket; e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := "keyName", *res.Key; e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := "locationName", *res.Location; e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
}

func TestCompleteMultipartUploadError(t *testing.T) {
	req := newCopyTestSvc(errMsg).CompleteMultipartUploadRequest(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String("bucketname"),
		Key:      aws.String("key"),
		UploadId: aws.String("uploadID"),
	})
	_, err := req.Send(context.Background())

	if err == nil {
		t.Error("expected an error")
	}
	e := err.(awserr.Error)

	if e, a := "ErrorCode", e.Code(); e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
	if e, a := "message body", e.Message(); e != a {
		t.Errorf("expected %s, but received %s", e, a)
	}
}

func newCopyTestSvc(errMsg string) *s3.Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, errMsg, http.StatusOK)
	}))
	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)
	cfg.Retryer = aws.NoOpRetryer{}
	svc := s3.New(cfg)
	svc.ForcePathStyle = true

	return svc
}

func TestStatusOKPayloadHandling(t *testing.T) {
	cases := map[string]struct {
		Header   http.Header
		Payloads [][]byte
		OpCall   func(api *s3.Client) error
		Err      string
	}{
		"200 error": {
			Header: http.Header{
				"Content-Length": []string{strconv.Itoa(len(errMsg))},
			},
			Payloads: [][]byte{[]byte(errMsg)},
			OpCall: func(c *s3.Client) error {
				req := c.CopyObjectRequest(&s3.CopyObjectInput{
					Bucket:     aws.String("bucketname"),
					CopySource: aws.String("bucketname/doesnotexist.txt"),
					Key:        aws.String("destination.txt"),
				})
				_, err := req.Send(context.Background())
				return err
			},
			Err: "ErrorCode: message body",
		},
		"200 error partial response": {
			Header: http.Header{
				"Content-Length": []string{strconv.Itoa(len(errMsg))},
			},
			Payloads: [][]byte{
				[]byte(errMsg[:20]),
			},
			OpCall: func(c *s3.Client) error {
				req := c.CopyObjectRequest(&s3.CopyObjectInput{
					Bucket:     aws.String("bucketname"),
					CopySource: aws.String("bucketname/doesnotexist.txt"),
					Key:        aws.String("destination.txt"),
				})
				_, err := req.Send(context.Background())
				return err
			},
			Err: "unexpected EOF",
		},
		"200 error multipart": {
			Header: http.Header{
				"Transfer-Encoding": []string{"chunked"},
			},
			Payloads: [][]byte{
				[]byte(errMsg[:20]),
				[]byte(errMsg[20:]),
			},
			OpCall: func(c *s3.Client) error {
				req := c.CopyObjectRequest(&s3.CopyObjectInput{
					Bucket:     aws.String("bucketname"),
					CopySource: aws.String("bucketname/doesnotexist.txt"),
					Key:        aws.String("destination.txt"),
				})
				_, err := req.Send(context.Background())
				return err
			},
			Err: "ErrorCode: message body",
		},
		"200 error multipart partial response": {
			Header: http.Header{
				"Transfer-Encoding": []string{"chunked"},
			},
			Payloads: [][]byte{
				[]byte(errMsg[:20]),
			},
			OpCall: func(c *s3.Client) error {
				req := c.CopyObjectRequest(&s3.CopyObjectInput{
					Bucket:     aws.String("bucketname"),
					CopySource: aws.String("bucketname/doesnotexist.txt"),
					Key:        aws.String("destination.txt"),
				})
				_, err := req.Send(context.Background())
				return err
			},
			Err: "XML syntax error",
		},
		"200 error multipart no payload": {
			Header: http.Header{
				"Transfer-Encoding": []string{"chunked"},
			},
			Payloads: [][]byte{},
			OpCall: func(c *s3.Client) error {
				req := c.CopyObjectRequest(&s3.CopyObjectInput{
					Bucket:     aws.String("bucketname"),
					CopySource: aws.String("bucketname/doesnotexist.txt"),
					Key:        aws.String("destination.txt"),
				})
				_, err := req.Send(context.Background())
				return err
			},
			Err: "empty response payload",
		},
		"response with only xml preamble": {
			Header: http.Header{
				"Transfer-Encoding": []string{"chunked"},
			},
			Payloads: [][]byte{
				[]byte(xmlPreambleMsg),
			},
			OpCall: func(c *s3.Client) error {
				req := c.CopyObjectRequest(&s3.CopyObjectInput{
					Bucket:     aws.String("bucketname"),
					CopySource: aws.String("bucketname/doesnotexist.txt"),
					Key:        aws.String("destination.txt"),
				})
				_, err := req.Send(context.Background())
				return err
			},
			Err: "empty response payload",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ww := w.(interface {
					http.ResponseWriter
					http.Flusher
				})

				for k, vs := range c.Header {
					for _, v := range vs {
						ww.Header().Add(k, v)
					}
				}
				ww.WriteHeader(http.StatusOK)
				ww.Flush()

				for _, p := range c.Payloads {
					ww.Write(p)
					ww.Flush()
				}
			}))
			defer srv.Close()

			cfg := unit.Config()
			cfg.EndpointResolver = aws.ResolveWithEndpointURL(srv.URL)
			svc := s3.New(cfg)
			svc.ForcePathStyle = true
			err := c.OpCall(svc)
			if len(c.Err) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.Err, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %v error in, %v", e, a)
				}
			} else if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
