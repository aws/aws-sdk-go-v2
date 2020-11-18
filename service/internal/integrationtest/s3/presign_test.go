// +build integration

package s3

func TestInteg_PresignURL_PutObject(t *testing.T) {

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	bucketName := "mockbucket-01"
	key := "random"

	params := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   bytes.NewReader([]byte(`Hello-world`)),
	}

	presignerClient := s3.NewPresignClient(client, func(options *s3.PresignOptions) {
		options.Expires = 600 * time.Second
	})

	presignRequest, err := presignerClient.PresignPutObject(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	// Putobject
	t.Logf("url : %v \n", presignRequest.URL)
	t.Logf("method : %v \n", presignRequest.Method)
	t.Logf("signed headers : %v \n", presignRequest.SignedHeader)

	t.Logf("attempting to put request")

	req, err := http.NewRequest(presignRequest.Method, presignRequest.URL, nil)
	if err != nil {
		t.Fatalf("failed to build presigned request, %v", err)
	}

	for k, vs := range presignRequest.SignedHeader {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	// Need to ensure that the content length member is set of the HTTP Request
	// or the request will not be transmitted correctly with a content length
	// value across the wire.
	if contLen := req.Header.Get("Content-Length"); len(contLen) > 0 {
		req.ContentLength, _ = strconv.ParseInt(contLen, 10, 64)
	}

	req.Body = ioutil.NopCloser(params.Body)

	// Upload the file contents to S3.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to do PUT request, %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to put S3 object, %d:%s", resp.StatusCode, resp.Status)
	}
}
