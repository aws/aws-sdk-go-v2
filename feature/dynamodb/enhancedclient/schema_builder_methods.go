package enhancedclient

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// AttributeDefinitions returns the attribute definitions for the DynamoDB table schema.
func (s *Schema[T]) AttributeDefinitions() []types.AttributeDefinition {
	return s.attributeDefinitions
}

// WithAttributeDefinitions sets the attribute definitions for the schema and returns the updated schema.
func (s *Schema[T]) WithAttributeDefinitions(attributeDefinitions []types.AttributeDefinition) *Schema[T] {
	s.attributeDefinitions = attributeDefinitions

	return s
}

// KeySchema returns the key schema elements for the DynamoDB table.
func (s *Schema[T]) KeySchema() []types.KeySchemaElement {
	return s.keySchema
}

// WithKeySchema sets the key schema for the table and returns the updated schema.
func (s *Schema[T]) WithKeySchema(keySchema []types.KeySchemaElement) *Schema[T] {
	s.keySchema = keySchema

	return s
}

// TableName returns the name of the DynamoDB table.
func (s *Schema[T]) TableName() *string {
	return s.tableName
}

// WithTableName sets the table name and returns the updated schema.
func (s *Schema[T]) WithTableName(tableName *string) *Schema[T] {
	s.tableName = tableName

	return s
}

// BillingMode returns the billing mode for the DynamoDB table.
func (s *Schema[T]) BillingMode() types.BillingMode {
	return s.billingMode
}

// WithBillingMode sets the billing mode for the table and returns the updated schema.
func (s *Schema[T]) WithBillingMode(billingMode types.BillingMode) *Schema[T] {
	s.billingMode = billingMode

	return s
}

// DeletionProtectionEnabled returns whether deletion protection is enabled for the table.
func (s *Schema[T]) DeletionProtectionEnabled() *bool {
	return s.deletionProtectionEnabled
}

// WithDeletionProtectionEnabled sets deletion protection for the table and returns the updated schema.
func (s *Schema[T]) WithDeletionProtectionEnabled(deletionProtectionEnabled *bool) *Schema[T] {
	s.deletionProtectionEnabled = deletionProtectionEnabled

	return s
}

// GlobalSecondaryIndexes returns the global secondary indexes for the table.
func (s *Schema[T]) GlobalSecondaryIndexes() []types.GlobalSecondaryIndex {
	return s.globalSecondaryIndexes
}

// WithGlobalSecondaryIndexes overwrites the global secondary indexes and returns the updated schema.
func (s *Schema[T]) WithGlobalSecondaryIndexes(globalSecondaryIndexes []types.GlobalSecondaryIndex) *Schema[T] {
	s.globalSecondaryIndexes = globalSecondaryIndexes

	return s
}

// WithGlobalSecondaryIndex creates or updates a global secondary index by name using the provided function.
// If the index does not exist, it is created. Returns the updated schema.
func (s *Schema[T]) WithGlobalSecondaryIndex(name string, fn func(gsi *types.GlobalSecondaryIndex)) *Schema[T] {
	var gsi *types.GlobalSecondaryIndex
	for idx := range s.globalSecondaryIndexes {
		gsi = &s.globalSecondaryIndexes[idx]
		if gsi.IndexName != nil && *gsi.IndexName == name {
			fn(gsi)
			break
		}
	}

	if gsi == nil {
		gsi = &types.GlobalSecondaryIndex{
			IndexName: pointer(name),
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
		}
		fn(gsi)
		s.globalSecondaryIndexes = append(s.globalSecondaryIndexes, *gsi)
	}

	attrs := map[string]bool{}
	for _, ad := range s.attributeDefinitions {
		attrs[*ad.AttributeName] = true
	}

	for _, ks := range gsi.KeySchema {
		if _, ok := attrs[*ks.AttributeName]; !ok {
			f, _ := s.cachedFields.FieldByName(*ks.AttributeName)
			at, _ := typeToScalarAttributeType(f.Type)
			s.attributeDefinitions = append(s.attributeDefinitions, types.AttributeDefinition{
				AttributeName: ks.AttributeName,
				AttributeType: at,
			})
		}
	}

	return s
}

// LocalSecondaryIndexes returns the local secondary indexes for the table.
func (s *Schema[T]) LocalSecondaryIndexes() []types.LocalSecondaryIndex {
	return s.localSecondaryIndexes
}

// WithLocalSecondaryIndexes overwrites the local secondary indexes and returns the updated schema.
func (s *Schema[T]) WithLocalSecondaryIndexes(localSecondaryIndexes []types.LocalSecondaryIndex) *Schema[T] {
	s.localSecondaryIndexes = localSecondaryIndexes

	return s
}

// WithLocalSecondaryIndex creates or updates a local secondary index by name using the provided function.
// If the index does not exist, it is created. Returns the updated schema.
func (s *Schema[T]) WithLocalSecondaryIndex(name string, fn func(gsi *types.LocalSecondaryIndex)) *Schema[T] {
	existing := false
	for idx := range s.localSecondaryIndexes {
		lsi := s.localSecondaryIndexes[idx]
		if lsi.IndexName != nil && *lsi.IndexName == name {
			fn(&lsi)
			existing = true
		}
	}

	if !existing {
		lsi := types.LocalSecondaryIndex{
			IndexName: pointer(name),
		}
		fn(&lsi)
		s.localSecondaryIndexes = append(s.localSecondaryIndexes, lsi)
	}

	return s
}

// OnDemandThroughput returns the on-demand throughput settings for the table.
func (s *Schema[T]) OnDemandThroughput() *types.OnDemandThroughput {
	return s.onDemandThroughput
}

// WithOnDemandThroughput sets the on-demand throughput and returns the updated schema.
func (s *Schema[T]) WithOnDemandThroughput(onDemandThroughput *types.OnDemandThroughput) *Schema[T] {
	s.onDemandThroughput = onDemandThroughput

	return s
}

// ProvisionedThroughput returns the provisioned throughput settings for the table.
func (s *Schema[T]) ProvisionedThroughput() *types.ProvisionedThroughput {
	return s.provisionedThroughput
}

// WithProvisionedThroughput sets the provisioned throughput and returns the updated schema.
func (s *Schema[T]) WithProvisionedThroughput(provisionedThroughput *types.ProvisionedThroughput) *Schema[T] {
	s.provisionedThroughput = provisionedThroughput

	return s
}

// ResourcePolicy returns the resource policy for the table.
func (s *Schema[T]) ResourcePolicy() *string {
	return s.resourcePolicy
}

// WithResourcePolicy sets the resource policy and returns the updated schema.
func (s *Schema[T]) WithResourcePolicy(resourcePolicy *string) *Schema[T] {
	s.resourcePolicy = resourcePolicy

	return s
}

// SSESpecification returns the server-side encryption specification for the table.
func (s *Schema[T]) SSESpecification() *types.SSESpecification {
	return s.sseSpecification
}

// WithSSESpecification sets the server-side encryption specification and returns the updated schema.
func (s *Schema[T]) WithSSESpecification(sseSpecification *types.SSESpecification) *Schema[T] {
	s.sseSpecification = sseSpecification

	return s
}

// StreamSpecification returns the stream specification for the table.
func (s *Schema[T]) StreamSpecification() *types.StreamSpecification {
	return s.streamSpecification
}

// WithStreamSpecification sets the stream specification and returns the updated schema.
func (s *Schema[T]) WithStreamSpecification(streamSpecification *types.StreamSpecification) *Schema[T] {
	s.streamSpecification = streamSpecification

	return s
}

// TableClass returns the table class for the DynamoDB table.
func (s *Schema[T]) TableClass() types.TableClass {
	return s.tableClass
}

// WithTableClass sets the table class and returns the updated schema.
func (s *Schema[T]) WithTableClass(tableClass types.TableClass) *Schema[T] {
	s.tableClass = tableClass

	return s
}

// Tags returns the tags associated with the table.
func (s *Schema[T]) Tags() []types.Tag {
	return s.tags
}

// WithTags sets the tags for the table and returns the updated schema.
func (s *Schema[T]) WithTags(tags []types.Tag) *Schema[T] {
	s.tags = tags

	return s
}

// WarmThroughput returns the warm throughput settings for the table.
func (s *Schema[T]) WarmThroughput() *types.WarmThroughput {
	return s.warmThroughput
}

// WithWarmThroughput sets the warm throughput and returns the updated schema.
func (s *Schema[T]) WithWarmThroughput(warmThroughput *types.WarmThroughput) *Schema[T] {
	s.warmThroughput = warmThroughput

	return s
}

// MultiRegionConsistency returns the multi-region consistency setting for the table.
func (s *Schema[T]) MultiRegionConsistency() types.MultiRegionConsistency {
	return s.multiRegionConsistency
}

// WithMultiRegionConsistency sets the multi-region consistency and returns the updated schema.
func (s *Schema[T]) WithMultiRegionConsistency(multiRegionConsistency types.MultiRegionConsistency) *Schema[T] {
	s.multiRegionConsistency = multiRegionConsistency

	return s
}

// ReplicaUpdates returns the replication group updates for the table.
func (s *Schema[T]) ReplicaUpdates() []types.ReplicationGroupUpdate {
	return s.replicaUpdates
}

// WithReplicaUpdates sets the replication group updates and returns the updated schema.
func (s *Schema[T]) WithReplicaUpdates(replicaUpdates []types.ReplicationGroupUpdate) *Schema[T] {
	s.replicaUpdates = replicaUpdates

	return s
}
