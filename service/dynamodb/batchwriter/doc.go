/*
Package batchwriter exposes a wrapper to a dynamodb client which can
be used to do automatic batch writes of PutRequests and DeleteRequests.

BatchWriter would normally be used to simplify the task of making
performant pipelines that read from some other service like ElastiCache and
save on DynamoDB.

Code using BatchWriter will typically look like this (assuming there is a
goroutine writing to itemChannel):

	batchWriter := batchwriter.New("table-name", dynamoClient)
	defer func() {
		for !batchwriter.Empty() {
			batchWriter.Flush()
		}
	}
	for item := range itemChannel {
		// Where itemChannel is of type `chan map[string]interface{}`
		dynamoItem := dynamodbattribute.MarshalMap(item)
		batchWriter.PutItem(&dynamodb.PutRequest{
			Item: dynamoItem,
		}) // or DeleteItem(&dynamodb.DeleteRequest{Key: dynamoItem})
	}

Or, in the common case where the pipeline runs forever and it needs to be
flushed periodically, one could use time.NewTicker and the select{} statement
instead of iterating through the itemChannel to do the flushing safely.

Note: batchWriter is NOT thread safe. Use it only from within a single thread.
*/
package batchwriter
