/*
Package dynamodbbatchwriter exposes a wrapper to a dynamodb client which can
be used to do automatic batch writes of PutRequests and DeleteRequests.

BatchWriter would normally be used to simplify the task of making
performant pipelines that read from some other service like ElastiCache and
save on DynamoDB.

Code using BatchWriter will typically look like this (assuming there is a
goroutine writing to itemChannel):

	batchWriter := dynamodbbatchwriter.New("table-name", dynamoClient)
	defer func() {
		for !batchwriter.Empty() {
			batchWriter.Flush()
		}
	}
	for item := range(itemChannel) {
		// Where itemChannel is of type `chan map[string]interface{}`
		dynamoItem := dynamodbattribute.MarshalMap(item)
		batchWriter.PutItem(dynamoItem) // or DeleteItem
	}

Or, in the common case where the pipeline runs forever and it needs to be
flushed periodically, one could use time.NewTicker and the select{} statement
instead of iterating through the itemChannel to do the flushing safely.

The input to DeleteItem should be the Key of a normal DeleteItem request.
e.g. {"word": {"S": "potato"}, "position": {"N": 5}}.

Note: batchWriter is NOT thread safe. Use it only from within a single thread.
*/
package dynamodbbatchwriter
