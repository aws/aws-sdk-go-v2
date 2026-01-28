package enhancedclient

// Index represents metadata for a DynamoDB index, including its name and type (global/local, partition/sort).
type Index struct {
	Name      string // Index name
	Global    bool   // True if the index is a global secondary index
	Local     bool   // True if the index is a local secondary index
	Partition bool   // True if the index is a partition key
	Sort      bool   // True if the index is a sort key
}
