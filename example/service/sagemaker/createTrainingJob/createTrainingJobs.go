// +build example

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jviney/aws-sdk-go-v2/aws/external"
	"github.com/jviney/aws-sdk-go-v2/service/sagemaker"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// Usage: go run -tags example createTrainingJobs.go
func main() {

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}

	sagemakerSvc := sagemaker.New(cfg)

	//Define intput variables
	role := "arn:aws:iam::xxxxxx:role/<name of role>"
	name := "k-means-in-sagemaker"
	MaxRuntimeInSeconds := int64(60 * 60)
	S3OutputPath := "s3://<bucket where your model artifact will be saved"
	InstanceCount := int64(2)
	VolumeSizeInGB := int64(75)
	TrainingInstanceType := sagemaker.TrainingInstanceType("ml.c4.8xlarge")
	TrainingImage := "174872318107.dkr.ecr.us-west-2.amazonaws.com/kmeans:1"
	TrainingInputMode := sagemaker.TrainingInputMode("File")

	ChannelName := "train"
	S3DataDistributionType := sagemaker.S3DataDistribution("FullyReplicated")
	S3DataType := sagemaker.S3DataType("S3Prefix")
	S3Uri := "s3://<bucket where the input data is available>"

	HyperParameters := map[string]string{
		"k":               "10",
		"feature_dim":     "784",
		"mini_batch_size": "500",
	}

	params := &sagemaker.CreateTrainingJobInput{
		RoleArn:         &role,
		TrainingJobName: &name,

		StoppingCondition: &sagemaker.StoppingCondition{
			MaxRuntimeInSeconds: &MaxRuntimeInSeconds,
		},

		OutputDataConfig: &sagemaker.OutputDataConfig{
			S3OutputPath: &S3OutputPath,
		},

		ResourceConfig: &sagemaker.ResourceConfig{
			InstanceCount:  &InstanceCount,
			VolumeSizeInGB: &VolumeSizeInGB,
			InstanceType:   TrainingInstanceType,
		},

		AlgorithmSpecification: &sagemaker.AlgorithmSpecification{
			TrainingImage:     &TrainingImage,
			TrainingInputMode: TrainingInputMode,
		},

		InputDataConfig: []sagemaker.Channel{
			{
				ChannelName: &ChannelName,
				DataSource: &sagemaker.DataSource{
					S3DataSource: &sagemaker.S3DataSource{
						S3DataDistributionType: S3DataDistributionType,
						S3DataType:             S3DataType,
						S3Uri:                  &S3Uri,
					},
				},
			},
		},
		HyperParameters: HyperParameters,
	}

	req := sagemakerSvc.CreateTrainingJobRequest(params)

	resp, err := req.Send(context.TODO())
	if err != nil {
		exitErrorf("Error creating TrainingJob, %v", err)
		return
	}

	fmt.Println(resp)
}
