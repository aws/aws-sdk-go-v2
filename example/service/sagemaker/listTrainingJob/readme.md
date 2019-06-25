## usage
`go run listTrainingJobs.go <number of jobs to be displayed>`

`go run listTrainingJobs.go 1`

## output
`
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
`
