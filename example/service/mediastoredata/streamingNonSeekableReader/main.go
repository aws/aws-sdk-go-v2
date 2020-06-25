// +build example

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/service/mediastore"
	"github.com/jviney/aws-sdk-go-v2/service/mediastoredata"
)

func main() {
	containerName := os.Args[1]
	objectPath := os.Args[2]

	// Create an AWS Elemental MediaStore Data client using default config.
	config := aws.Config{}
	dataSvc, err := getMediaStoreDataClient(containerName, config)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	// Create a random reader to simulate a unseekable reader, wrap the reader
	// in an io.LimitReader to prevent uploading forever.
	randReader := rand.New(rand.NewSource(0))
	reader := io.LimitReader(randReader, 1024*1024 /* 1MB */)

	// Wrap the unseekable reader with the SDK's RandSeekCloser. This type will
	// allow the SDK's to use the nonseekable reader.
	body := aws.ReadSeekCloser(reader)

	// Make the PutObject API call with the nonseekable reader, causing the SDK
	// to send the request body payload a chunked transfer encoding.
	dataSvc.PutObjectRequest(&mediastoredata.PutObjectInput{
		Path: &objectPath,
		Body: body,
	})

	fmt.Println("object uploaded")
}

// getMediaStoreDataClient uses the AWS Elemental MediaStore API to get the
// endpoint for a container. If the container endpoint can be retrieved a AWS
// Elemental MediaStore Data client will be created and returned. Otherwise
// error is returned.
func getMediaStoreDataClient(containerName string, config aws.Config) (*mediastoredata.Client, error) {
	endpoint, err := containerEndpoint(containerName, config)
	if err != nil {
		return nil, err
	}
	config.EndpointResolver = aws.ResolveWithEndpointURL(aws.StringValue(endpoint))
	dataSvc := mediastoredata.New(config)

	return dataSvc, nil
}

// ContainerEndpoint will attempt to get the endpoint for a container,
// returning error if the container doesn't exist, or is not active within a
// timeout.
func containerEndpoint(name string, config aws.Config) (*string, error) {
	for i := 0; i < 3; i++ {
		ctrlSvc := mediastore.New(config)
		descContainerRequest := ctrlSvc.DescribeContainerRequest(&mediastore.DescribeContainerInput{
			ContainerName: &name,
		})

		descResp, err := descContainerRequest.Send(descContainerRequest.Context())
		if err != nil {
			return nil, err
		}

		if status := descResp.Container.Status; status != "ACTIVE" {
			log.Println("waiting for container to be active, ", status)
			time.Sleep(10 * time.Second)
			continue
		}

		return descResp.Container.Endpoint, nil
	}

	return nil, fmt.Errorf("container is not active")
}
