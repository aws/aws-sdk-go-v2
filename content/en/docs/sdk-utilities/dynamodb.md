---
title: "Amazon DynamoDB Utilities"
linkTitle: "Amazon DynamoDB"
---

{{% pageinfo color="warning" %}}
The content on this page does not reflect the current status or feature set of {{% alias sdk-go %}}.
{{% /pageinfo %}}

## {{% alias service="DDBlong" %}} Attributes Converter

The attributes converter simplifies converting {{% alias service="DDBlong" %}} attribute values to and from concrete Go
types. Conversions make it easy to work with attribute values in Go and to write values to
{{% alias service="DDBlong" %}} tables. For example, you can create records in Go and then use the converter when you
want to write those records as attribute values to a {{% alias service="DDB" %}} table.

The following example converts a structure to an {{% alias service="DDBlong" %}}
`AttributeValues` map and then puts the data to the `exampleTable`.

```go
type Record struct {
    MyField string
    Letters []string
    A2Num   map[string]int
}
r := Record{
    MyField: "dynamodbattribute.ConvertToX example",
    Letters: []string{"a", "b", "c", "d"},
    A2Num:   map[string]int{"a": 1, "b": 2, "c": 3},
}

//...

svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2")}))
item, err := dynamodbattribute.ConvertToMap(r)
if err != nil {
    fmt.Println("Failed to convert", err)
    return
}
result, err := svc.PutItem(&dynamodb.PutItemInput{
    Item:      item,
    TableName: aws.String("exampleTable"),
})
fmt.Println("Item put to dynamodb", result, err)
```

For more information about the converter utility, see the [attribute]({{< apiref "feature/dynamodb/attribute" >}})
package in the {{% alias sdk-api %}}.

