package enhancedclient

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Schema[T]) AttributeDefinitions() []types.AttributeDefinition {
	return s.attributeDefinitions
}

func (s *Schema[T]) WithAttributeDefinitions(attributeDefinitions []types.AttributeDefinition) *Schema[T] {
	s.attributeDefinitions = attributeDefinitions

	return s
}

func (s *Schema[T]) KeySchema() []types.KeySchemaElement {
	return s.keySchema
}

func (s *Schema[T]) WithKeySchema(keySchema []types.KeySchemaElement) *Schema[T] {
	s.keySchema = keySchema

	return s
}

func (s *Schema[T]) TableName() *string {
	return s.tableName
}

func (s *Schema[T]) WithTableName(tableName *string) *Schema[T] {
	s.tableName = tableName

	return s
}

func (s *Schema[T]) BillingMode() types.BillingMode {
	return s.billingMode
}

func (s *Schema[T]) WithBillingMode(billingMode types.BillingMode) *Schema[T] {
	s.billingMode = billingMode

	return s
}

func (s *Schema[T]) DeletionProtectionEnabled() *bool {
	return s.deletionProtectionEnabled
}

func (s *Schema[T]) WithDeletionProtectionEnabled(deletionProtectionEnabled *bool) *Schema[T] {
	s.deletionProtectionEnabled = deletionProtectionEnabled

	return s
}

func (s *Schema[T]) GlobalSecondaryIndexes() []types.GlobalSecondaryIndex {
	return s.globalSecondaryIndexes
}

// WithGlobalSecondaryIndexes overwrites the global secondary indexes
func (s *Schema[T]) WithGlobalSecondaryIndexes(globalSecondaryIndexes []types.GlobalSecondaryIndex) *Schema[T] {
	s.globalSecondaryIndexes = globalSecondaryIndexes

	return s
}

// WithGlobalSecondaryIndex creates or updates in place a global secondary index
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

func (s *Schema[T]) LocalSecondaryIndexes() []types.LocalSecondaryIndex {
	return s.localSecondaryIndexes
}

func (s *Schema[T]) WithLocalSecondaryIndexes(localSecondaryIndexes []types.LocalSecondaryIndex) *Schema[T] {
	s.localSecondaryIndexes = localSecondaryIndexes

	return s
}

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

func (s *Schema[T]) OnDemandThroughput() *types.OnDemandThroughput {
	return s.onDemandThroughput
}

func (s *Schema[T]) WithOnDemandThroughput(onDemandThroughput *types.OnDemandThroughput) *Schema[T] {
	s.onDemandThroughput = onDemandThroughput

	return s
}

func (s *Schema[T]) ProvisionedThroughput() *types.ProvisionedThroughput {
	return s.provisionedThroughput
}

func (s *Schema[T]) WithProvisionedThroughput(provisionedThroughput *types.ProvisionedThroughput) *Schema[T] {
	s.provisionedThroughput = provisionedThroughput

	return s
}

func (s *Schema[T]) ResourcePolicy() *string {
	return s.resourcePolicy
}

func (s *Schema[T]) WithResourcePolicy(resourcePolicy *string) *Schema[T] {
	s.resourcePolicy = resourcePolicy

	return s
}

func (s *Schema[T]) SSESpecification() *types.SSESpecification {
	return s.sseSpecification
}

func (s *Schema[T]) WithSSESpecification(sseSpecification *types.SSESpecification) *Schema[T] {
	s.sseSpecification = sseSpecification

	return s
}

func (s *Schema[T]) StreamSpecification() *types.StreamSpecification {
	return s.streamSpecification
}

func (s *Schema[T]) WithStreamSpecification(streamSpecification *types.StreamSpecification) *Schema[T] {
	s.streamSpecification = streamSpecification

	return s
}

func (s *Schema[T]) TableClass() types.TableClass {
	return s.tableClass
}

func (s *Schema[T]) WithTableClass(tableClass types.TableClass) *Schema[T] {
	s.tableClass = tableClass

	return s
}

func (s *Schema[T]) Tags() []types.Tag {
	return s.tags
}

func (s *Schema[T]) WithTags(tags []types.Tag) *Schema[T] {
	s.tags = tags

	return s
}

func (s *Schema[T]) WarmThroughput() *types.WarmThroughput {
	return s.warmThroughput
}

func (s *Schema[T]) WithWarmThroughput(warmThroughput *types.WarmThroughput) *Schema[T] {
	s.warmThroughput = warmThroughput

	return s
}

func (s *Schema[T]) MultiRegionConsistency() types.MultiRegionConsistency {
	return s.multiRegionConsistency
}

func (s *Schema[T]) WithMultiRegionConsistency(multiRegionConsistency types.MultiRegionConsistency) *Schema[T] {
	s.multiRegionConsistency = multiRegionConsistency

	return s
}

func (s *Schema[T]) ReplicaUpdates() []types.ReplicationGroupUpdate {
	return s.replicaUpdates
}

func (s *Schema[T]) WithReplicaUpdates(replicaUpdates []types.ReplicationGroupUpdate) *Schema[T] {
	s.replicaUpdates = replicaUpdates

	return s
}
