Example using [Amazon SageMaker](https://aws.amazon.com/sagemaker/) with the
AWS SDK for Go to list training jobs that have been created, and their status.

## Usage
Use the example to list the created training jobs and their status, passing in
the number of jobs to show.

```sh
go run listTrainingJobs.go <number of jobs to be displayed>

# E.g.
go run listTrainingJobs.go 1
```

## Output
Example response of a training job and its status.

```
{
  NextToken: "xcskfskdfksdffksdhfjhjghjshdfgjhfjgh"
  TrainingJobSummaries: [{
      CreationTime: 2019-07-03 18:17:34 +0000 UTC,
      LastModifiedTime: 2019-07-03 18:22:15 +0000 UTC,
      TrainingEndTime: 2019-07-03 18:22:15 +0000 UTC,
      TrainingJobArn: "arn:aws:sagemaker:us-west-2:<account_number>:training-job/hpojob-20190703181725-leaj-001-d345f443",
      TrainingJobName: "HPOJob-20190703181725-LEAJ-001-d345f443",
      TrainingJobStatus: Completed
    }]
}
```
