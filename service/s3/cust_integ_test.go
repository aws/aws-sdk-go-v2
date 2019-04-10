// +build integration

package s3_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestInteg_WriteToObject(t *testing.T) {
	_, err := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    aws.String("key name"),
		Body:   bytes.NewReader([]byte("hello world")),
	}).Send(context.Background())
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	getResp, err := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: bucketName,
		Key:    aws.String("key name"),
	}).Send(context.Background())
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	b, _ := ioutil.ReadAll(getResp.Body)
	if e, a := []byte("hello world"), b; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestInteg_PresignedGetPut(t *testing.T) {
	putReq := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: bucketName,
		Key:    aws.String("presigned-key"),
	})
	var err error

	// Presign a PUT request
	var puturl string
	puturl, err = putReq.Presign(300 * time.Second)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	// PUT to the presigned URL with a body
	var putHTTPReq *http.Request
	buf := bytes.NewReader([]byte("hello world"))
	putHTTPReq, err = http.NewRequest("PUT", puturl, buf)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	var putResp *http.Response
	putResp, err = http.DefaultClient.Do(putHTTPReq)
	if err != nil {
		t.Errorf("expect put with presign url no error, got %v", err)
	}
	if e, a := 200, putResp.StatusCode; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	// Presign a GET on the same URL
	getReq := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: bucketName,
		Key:    aws.String("presigned-key"),
	})

	var getURL string
	getURL, err = getReq.Presign(300 * time.Second)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	// Get the body
	var getResp *http.Response
	getResp, err = http.Get(getURL)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	var b []byte
	defer getResp.Body.Close()
	b, err = ioutil.ReadAll(getResp.Body)
	if e, a := "hello world", string(b); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}
