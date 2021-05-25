module github.com/aws/aws-sdk-go-v2/service/commander

go 1.15

retract [v1.0.0, v1.1.0] // API client incorrectly named, Use AWS Systems Manager Incident Manager (ssmincidents) instead.
