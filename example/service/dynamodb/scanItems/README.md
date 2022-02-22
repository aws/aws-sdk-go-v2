# Amazon DynamoDB Scan Items Example

This is an example using the AWS SDK for Go to list items in a DynamoDB table.

### Usage

The example uses the table name provided, and lists all items in a dynamoDB's table.

```
go run scanItems.go -table <table-name> -region <region-name>

  -table name
        The name of the DynamoDB table to list item from.
  -region name
        The `region` of your AWS project.

```
