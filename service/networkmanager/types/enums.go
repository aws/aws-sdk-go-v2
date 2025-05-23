// Code generated by smithy-go-codegen DO NOT EDIT.

package types

type AttachmentErrorCode string

// Enum values for AttachmentErrorCode
const (
	AttachmentErrorCodeVpcNotFound                             AttachmentErrorCode = "VPC_NOT_FOUND"
	AttachmentErrorCodeSubnetNotFound                          AttachmentErrorCode = "SUBNET_NOT_FOUND"
	AttachmentErrorCodeSubnetDuplicatedInAvailabilityZone      AttachmentErrorCode = "SUBNET_DUPLICATED_IN_AVAILABILITY_ZONE"
	AttachmentErrorCodeSubnetNoFreeAddresses                   AttachmentErrorCode = "SUBNET_NO_FREE_ADDRESSES"
	AttachmentErrorCodeSubnetUnsupportedAvailabilityZone       AttachmentErrorCode = "SUBNET_UNSUPPORTED_AVAILABILITY_ZONE"
	AttachmentErrorCodeSubnetNoIpv6Cidrs                       AttachmentErrorCode = "SUBNET_NO_IPV6_CIDRS"
	AttachmentErrorCodeVpnConnectionNotFound                   AttachmentErrorCode = "VPN_CONNECTION_NOT_FOUND"
	AttachmentErrorCodeMaximumNoEncapLimitExceeded             AttachmentErrorCode = "MAXIMUM_NO_ENCAP_LIMIT_EXCEEDED"
	AttachmentErrorCodeDirectConnectGatewayNotFound            AttachmentErrorCode = "DIRECT_CONNECT_GATEWAY_NOT_FOUND"
	AttachmentErrorCodeDirectConnectGatewayExistingAttachments AttachmentErrorCode = "DIRECT_CONNECT_GATEWAY_EXISTING_ATTACHMENTS"
	AttachmentErrorCodeDirectConnectGatewayNoPrivateVif        AttachmentErrorCode = "DIRECT_CONNECT_GATEWAY_NO_PRIVATE_VIF"
)

// Values returns all known values for AttachmentErrorCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AttachmentErrorCode) Values() []AttachmentErrorCode {
	return []AttachmentErrorCode{
		"VPC_NOT_FOUND",
		"SUBNET_NOT_FOUND",
		"SUBNET_DUPLICATED_IN_AVAILABILITY_ZONE",
		"SUBNET_NO_FREE_ADDRESSES",
		"SUBNET_UNSUPPORTED_AVAILABILITY_ZONE",
		"SUBNET_NO_IPV6_CIDRS",
		"VPN_CONNECTION_NOT_FOUND",
		"MAXIMUM_NO_ENCAP_LIMIT_EXCEEDED",
		"DIRECT_CONNECT_GATEWAY_NOT_FOUND",
		"DIRECT_CONNECT_GATEWAY_EXISTING_ATTACHMENTS",
		"DIRECT_CONNECT_GATEWAY_NO_PRIVATE_VIF",
	}
}

type AttachmentState string

// Enum values for AttachmentState
const (
	AttachmentStateRejected                    AttachmentState = "REJECTED"
	AttachmentStatePendingAttachmentAcceptance AttachmentState = "PENDING_ATTACHMENT_ACCEPTANCE"
	AttachmentStateCreating                    AttachmentState = "CREATING"
	AttachmentStateFailed                      AttachmentState = "FAILED"
	AttachmentStateAvailable                   AttachmentState = "AVAILABLE"
	AttachmentStateUpdating                    AttachmentState = "UPDATING"
	AttachmentStatePendingNetworkUpdate        AttachmentState = "PENDING_NETWORK_UPDATE"
	AttachmentStatePendingTagAcceptance        AttachmentState = "PENDING_TAG_ACCEPTANCE"
	AttachmentStateDeleting                    AttachmentState = "DELETING"
)

// Values returns all known values for AttachmentState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AttachmentState) Values() []AttachmentState {
	return []AttachmentState{
		"REJECTED",
		"PENDING_ATTACHMENT_ACCEPTANCE",
		"CREATING",
		"FAILED",
		"AVAILABLE",
		"UPDATING",
		"PENDING_NETWORK_UPDATE",
		"PENDING_TAG_ACCEPTANCE",
		"DELETING",
	}
}

type AttachmentType string

// Enum values for AttachmentType
const (
	AttachmentTypeConnect                  AttachmentType = "CONNECT"
	AttachmentTypeSiteToSiteVpn            AttachmentType = "SITE_TO_SITE_VPN"
	AttachmentTypeVpc                      AttachmentType = "VPC"
	AttachmentTypeDirectConnectGateway     AttachmentType = "DIRECT_CONNECT_GATEWAY"
	AttachmentTypeTransitGatewayRouteTable AttachmentType = "TRANSIT_GATEWAY_ROUTE_TABLE"
)

// Values returns all known values for AttachmentType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AttachmentType) Values() []AttachmentType {
	return []AttachmentType{
		"CONNECT",
		"SITE_TO_SITE_VPN",
		"VPC",
		"DIRECT_CONNECT_GATEWAY",
		"TRANSIT_GATEWAY_ROUTE_TABLE",
	}
}

type ChangeAction string

// Enum values for ChangeAction
const (
	ChangeActionAdd    ChangeAction = "ADD"
	ChangeActionModify ChangeAction = "MODIFY"
	ChangeActionRemove ChangeAction = "REMOVE"
)

// Values returns all known values for ChangeAction. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ChangeAction) Values() []ChangeAction {
	return []ChangeAction{
		"ADD",
		"MODIFY",
		"REMOVE",
	}
}

type ChangeSetState string

// Enum values for ChangeSetState
const (
	ChangeSetStatePendingGeneration  ChangeSetState = "PENDING_GENERATION"
	ChangeSetStateFailedGeneration   ChangeSetState = "FAILED_GENERATION"
	ChangeSetStateReadyToExecute     ChangeSetState = "READY_TO_EXECUTE"
	ChangeSetStateExecuting          ChangeSetState = "EXECUTING"
	ChangeSetStateExecutionSucceeded ChangeSetState = "EXECUTION_SUCCEEDED"
	ChangeSetStateOutOfDate          ChangeSetState = "OUT_OF_DATE"
)

// Values returns all known values for ChangeSetState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ChangeSetState) Values() []ChangeSetState {
	return []ChangeSetState{
		"PENDING_GENERATION",
		"FAILED_GENERATION",
		"READY_TO_EXECUTE",
		"EXECUTING",
		"EXECUTION_SUCCEEDED",
		"OUT_OF_DATE",
	}
}

type ChangeStatus string

// Enum values for ChangeStatus
const (
	ChangeStatusNotStarted ChangeStatus = "NOT_STARTED"
	ChangeStatusInProgress ChangeStatus = "IN_PROGRESS"
	ChangeStatusComplete   ChangeStatus = "COMPLETE"
	ChangeStatusFailed     ChangeStatus = "FAILED"
)

// Values returns all known values for ChangeStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ChangeStatus) Values() []ChangeStatus {
	return []ChangeStatus{
		"NOT_STARTED",
		"IN_PROGRESS",
		"COMPLETE",
		"FAILED",
	}
}

type ChangeType string

// Enum values for ChangeType
const (
	ChangeTypeCoreNetworkSegment              ChangeType = "CORE_NETWORK_SEGMENT"
	ChangeTypeNetworkFunctionGroup            ChangeType = "NETWORK_FUNCTION_GROUP"
	ChangeTypeCoreNetworkEdge                 ChangeType = "CORE_NETWORK_EDGE"
	ChangeTypeAttachmentMapping               ChangeType = "ATTACHMENT_MAPPING"
	ChangeTypeAttachmentRoutePropagation      ChangeType = "ATTACHMENT_ROUTE_PROPAGATION"
	ChangeTypeAttachmentRouteStatic           ChangeType = "ATTACHMENT_ROUTE_STATIC"
	ChangeTypeCoreNetworkConfiguration        ChangeType = "CORE_NETWORK_CONFIGURATION"
	ChangeTypeSegmentsConfiguration           ChangeType = "SEGMENTS_CONFIGURATION"
	ChangeTypeSegmentActionsConfiguration     ChangeType = "SEGMENT_ACTIONS_CONFIGURATION"
	ChangeTypeAttachmentPoliciesConfiguration ChangeType = "ATTACHMENT_POLICIES_CONFIGURATION"
)

// Values returns all known values for ChangeType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ChangeType) Values() []ChangeType {
	return []ChangeType{
		"CORE_NETWORK_SEGMENT",
		"NETWORK_FUNCTION_GROUP",
		"CORE_NETWORK_EDGE",
		"ATTACHMENT_MAPPING",
		"ATTACHMENT_ROUTE_PROPAGATION",
		"ATTACHMENT_ROUTE_STATIC",
		"CORE_NETWORK_CONFIGURATION",
		"SEGMENTS_CONFIGURATION",
		"SEGMENT_ACTIONS_CONFIGURATION",
		"ATTACHMENT_POLICIES_CONFIGURATION",
	}
}

type ConnectionState string

// Enum values for ConnectionState
const (
	ConnectionStatePending   ConnectionState = "PENDING"
	ConnectionStateAvailable ConnectionState = "AVAILABLE"
	ConnectionStateDeleting  ConnectionState = "DELETING"
	ConnectionStateUpdating  ConnectionState = "UPDATING"
)

// Values returns all known values for ConnectionState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConnectionState) Values() []ConnectionState {
	return []ConnectionState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"UPDATING",
	}
}

type ConnectionStatus string

// Enum values for ConnectionStatus
const (
	ConnectionStatusUp   ConnectionStatus = "UP"
	ConnectionStatusDown ConnectionStatus = "DOWN"
)

// Values returns all known values for ConnectionStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConnectionStatus) Values() []ConnectionStatus {
	return []ConnectionStatus{
		"UP",
		"DOWN",
	}
}

type ConnectionType string

// Enum values for ConnectionType
const (
	ConnectionTypeBgp   ConnectionType = "BGP"
	ConnectionTypeIpsec ConnectionType = "IPSEC"
)

// Values returns all known values for ConnectionType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConnectionType) Values() []ConnectionType {
	return []ConnectionType{
		"BGP",
		"IPSEC",
	}
}

type ConnectPeerAssociationState string

// Enum values for ConnectPeerAssociationState
const (
	ConnectPeerAssociationStatePending   ConnectPeerAssociationState = "PENDING"
	ConnectPeerAssociationStateAvailable ConnectPeerAssociationState = "AVAILABLE"
	ConnectPeerAssociationStateDeleting  ConnectPeerAssociationState = "DELETING"
	ConnectPeerAssociationStateDeleted   ConnectPeerAssociationState = "DELETED"
)

// Values returns all known values for ConnectPeerAssociationState. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConnectPeerAssociationState) Values() []ConnectPeerAssociationState {
	return []ConnectPeerAssociationState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"DELETED",
	}
}

type ConnectPeerErrorCode string

// Enum values for ConnectPeerErrorCode
const (
	ConnectPeerErrorCodeEdgeLocationNoFreeIps     ConnectPeerErrorCode = "EDGE_LOCATION_NO_FREE_IPS"
	ConnectPeerErrorCodeEdgeLocationPeerDuplicate ConnectPeerErrorCode = "EDGE_LOCATION_PEER_DUPLICATE"
	ConnectPeerErrorCodeSubnetNotFound            ConnectPeerErrorCode = "SUBNET_NOT_FOUND"
	ConnectPeerErrorCodeIpOutsideSubnetCidrRange  ConnectPeerErrorCode = "IP_OUTSIDE_SUBNET_CIDR_RANGE"
	ConnectPeerErrorCodeInvalidInsideCidrBlock    ConnectPeerErrorCode = "INVALID_INSIDE_CIDR_BLOCK"
	ConnectPeerErrorCodeNoAssociatedCidrBlock     ConnectPeerErrorCode = "NO_ASSOCIATED_CIDR_BLOCK"
)

// Values returns all known values for ConnectPeerErrorCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConnectPeerErrorCode) Values() []ConnectPeerErrorCode {
	return []ConnectPeerErrorCode{
		"EDGE_LOCATION_NO_FREE_IPS",
		"EDGE_LOCATION_PEER_DUPLICATE",
		"SUBNET_NOT_FOUND",
		"IP_OUTSIDE_SUBNET_CIDR_RANGE",
		"INVALID_INSIDE_CIDR_BLOCK",
		"NO_ASSOCIATED_CIDR_BLOCK",
	}
}

type ConnectPeerState string

// Enum values for ConnectPeerState
const (
	ConnectPeerStateCreating  ConnectPeerState = "CREATING"
	ConnectPeerStateFailed    ConnectPeerState = "FAILED"
	ConnectPeerStateAvailable ConnectPeerState = "AVAILABLE"
	ConnectPeerStateDeleting  ConnectPeerState = "DELETING"
)

// Values returns all known values for ConnectPeerState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ConnectPeerState) Values() []ConnectPeerState {
	return []ConnectPeerState{
		"CREATING",
		"FAILED",
		"AVAILABLE",
		"DELETING",
	}
}

type CoreNetworkPolicyAlias string

// Enum values for CoreNetworkPolicyAlias
const (
	CoreNetworkPolicyAliasLive   CoreNetworkPolicyAlias = "LIVE"
	CoreNetworkPolicyAliasLatest CoreNetworkPolicyAlias = "LATEST"
)

// Values returns all known values for CoreNetworkPolicyAlias. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CoreNetworkPolicyAlias) Values() []CoreNetworkPolicyAlias {
	return []CoreNetworkPolicyAlias{
		"LIVE",
		"LATEST",
	}
}

type CoreNetworkState string

// Enum values for CoreNetworkState
const (
	CoreNetworkStateCreating  CoreNetworkState = "CREATING"
	CoreNetworkStateUpdating  CoreNetworkState = "UPDATING"
	CoreNetworkStateAvailable CoreNetworkState = "AVAILABLE"
	CoreNetworkStateDeleting  CoreNetworkState = "DELETING"
)

// Values returns all known values for CoreNetworkState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CoreNetworkState) Values() []CoreNetworkState {
	return []CoreNetworkState{
		"CREATING",
		"UPDATING",
		"AVAILABLE",
		"DELETING",
	}
}

type CustomerGatewayAssociationState string

// Enum values for CustomerGatewayAssociationState
const (
	CustomerGatewayAssociationStatePending   CustomerGatewayAssociationState = "PENDING"
	CustomerGatewayAssociationStateAvailable CustomerGatewayAssociationState = "AVAILABLE"
	CustomerGatewayAssociationStateDeleting  CustomerGatewayAssociationState = "DELETING"
	CustomerGatewayAssociationStateDeleted   CustomerGatewayAssociationState = "DELETED"
)

// Values returns all known values for CustomerGatewayAssociationState. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CustomerGatewayAssociationState) Values() []CustomerGatewayAssociationState {
	return []CustomerGatewayAssociationState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"DELETED",
	}
}

type DeviceState string

// Enum values for DeviceState
const (
	DeviceStatePending   DeviceState = "PENDING"
	DeviceStateAvailable DeviceState = "AVAILABLE"
	DeviceStateDeleting  DeviceState = "DELETING"
	DeviceStateUpdating  DeviceState = "UPDATING"
)

// Values returns all known values for DeviceState. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (DeviceState) Values() []DeviceState {
	return []DeviceState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"UPDATING",
	}
}

type GlobalNetworkState string

// Enum values for GlobalNetworkState
const (
	GlobalNetworkStatePending   GlobalNetworkState = "PENDING"
	GlobalNetworkStateAvailable GlobalNetworkState = "AVAILABLE"
	GlobalNetworkStateDeleting  GlobalNetworkState = "DELETING"
	GlobalNetworkStateUpdating  GlobalNetworkState = "UPDATING"
)

// Values returns all known values for GlobalNetworkState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (GlobalNetworkState) Values() []GlobalNetworkState {
	return []GlobalNetworkState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"UPDATING",
	}
}

type LinkAssociationState string

// Enum values for LinkAssociationState
const (
	LinkAssociationStatePending   LinkAssociationState = "PENDING"
	LinkAssociationStateAvailable LinkAssociationState = "AVAILABLE"
	LinkAssociationStateDeleting  LinkAssociationState = "DELETING"
	LinkAssociationStateDeleted   LinkAssociationState = "DELETED"
)

// Values returns all known values for LinkAssociationState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (LinkAssociationState) Values() []LinkAssociationState {
	return []LinkAssociationState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"DELETED",
	}
}

type LinkState string

// Enum values for LinkState
const (
	LinkStatePending   LinkState = "PENDING"
	LinkStateAvailable LinkState = "AVAILABLE"
	LinkStateDeleting  LinkState = "DELETING"
	LinkStateUpdating  LinkState = "UPDATING"
)

// Values returns all known values for LinkState. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (LinkState) Values() []LinkState {
	return []LinkState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"UPDATING",
	}
}

type PeeringErrorCode string

// Enum values for PeeringErrorCode
const (
	PeeringErrorCodeTransitGatewayNotFound           PeeringErrorCode = "TRANSIT_GATEWAY_NOT_FOUND"
	PeeringErrorCodeTransitGatewayPeersLimitExceeded PeeringErrorCode = "TRANSIT_GATEWAY_PEERS_LIMIT_EXCEEDED"
	PeeringErrorCodeMissingRequiredPermissions       PeeringErrorCode = "MISSING_PERMISSIONS"
	PeeringErrorCodeInternalError                    PeeringErrorCode = "INTERNAL_ERROR"
	PeeringErrorCodeEdgeLocationPeerDuplicate        PeeringErrorCode = "EDGE_LOCATION_PEER_DUPLICATE"
	PeeringErrorCodeInvalidTransitGatewayState       PeeringErrorCode = "INVALID_TRANSIT_GATEWAY_STATE"
)

// Values returns all known values for PeeringErrorCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PeeringErrorCode) Values() []PeeringErrorCode {
	return []PeeringErrorCode{
		"TRANSIT_GATEWAY_NOT_FOUND",
		"TRANSIT_GATEWAY_PEERS_LIMIT_EXCEEDED",
		"MISSING_PERMISSIONS",
		"INTERNAL_ERROR",
		"EDGE_LOCATION_PEER_DUPLICATE",
		"INVALID_TRANSIT_GATEWAY_STATE",
	}
}

type PeeringState string

// Enum values for PeeringState
const (
	PeeringStateCreating  PeeringState = "CREATING"
	PeeringStateFailed    PeeringState = "FAILED"
	PeeringStateAvailable PeeringState = "AVAILABLE"
	PeeringStateDeleting  PeeringState = "DELETING"
)

// Values returns all known values for PeeringState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PeeringState) Values() []PeeringState {
	return []PeeringState{
		"CREATING",
		"FAILED",
		"AVAILABLE",
		"DELETING",
	}
}

type PeeringType string

// Enum values for PeeringType
const (
	PeeringTypeTransitGateway PeeringType = "TRANSIT_GATEWAY"
)

// Values returns all known values for PeeringType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PeeringType) Values() []PeeringType {
	return []PeeringType{
		"TRANSIT_GATEWAY",
	}
}

type RouteAnalysisCompletionReasonCode string

// Enum values for RouteAnalysisCompletionReasonCode
const (
	RouteAnalysisCompletionReasonCodeTransitGatewayAttachmentNotFound                 RouteAnalysisCompletionReasonCode = "TRANSIT_GATEWAY_ATTACHMENT_NOT_FOUND"
	RouteAnalysisCompletionReasonCodeTransitGatewayAttachmentNotInTransitGateway      RouteAnalysisCompletionReasonCode = "TRANSIT_GATEWAY_ATTACHMENT_NOT_IN_TRANSIT_GATEWAY"
	RouteAnalysisCompletionReasonCodeCyclicPathDetected                               RouteAnalysisCompletionReasonCode = "CYCLIC_PATH_DETECTED"
	RouteAnalysisCompletionReasonCodeTransitGatewayAttachmentStableRouteTableNotFound RouteAnalysisCompletionReasonCode = "TRANSIT_GATEWAY_ATTACHMENT_STABLE_ROUTE_TABLE_NOT_FOUND"
	RouteAnalysisCompletionReasonCodeRouteNotFound                                    RouteAnalysisCompletionReasonCode = "ROUTE_NOT_FOUND"
	RouteAnalysisCompletionReasonCodeBlackholeRouteForDestinationFound                RouteAnalysisCompletionReasonCode = "BLACKHOLE_ROUTE_FOR_DESTINATION_FOUND"
	RouteAnalysisCompletionReasonCodeInactiveRouteForDestinationFound                 RouteAnalysisCompletionReasonCode = "INACTIVE_ROUTE_FOR_DESTINATION_FOUND"
	RouteAnalysisCompletionReasonCodeTransitGatewayAttachment                         RouteAnalysisCompletionReasonCode = "TRANSIT_GATEWAY_ATTACHMENT_ATTACH_ARN_NO_MATCH"
	RouteAnalysisCompletionReasonCodeMaxHopsExceeded                                  RouteAnalysisCompletionReasonCode = "MAX_HOPS_EXCEEDED"
	RouteAnalysisCompletionReasonCodePossibleMiddlebox                                RouteAnalysisCompletionReasonCode = "POSSIBLE_MIDDLEBOX"
	RouteAnalysisCompletionReasonCodeNoDestinationArnProvided                         RouteAnalysisCompletionReasonCode = "NO_DESTINATION_ARN_PROVIDED"
)

// Values returns all known values for RouteAnalysisCompletionReasonCode. Note
// that this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (RouteAnalysisCompletionReasonCode) Values() []RouteAnalysisCompletionReasonCode {
	return []RouteAnalysisCompletionReasonCode{
		"TRANSIT_GATEWAY_ATTACHMENT_NOT_FOUND",
		"TRANSIT_GATEWAY_ATTACHMENT_NOT_IN_TRANSIT_GATEWAY",
		"CYCLIC_PATH_DETECTED",
		"TRANSIT_GATEWAY_ATTACHMENT_STABLE_ROUTE_TABLE_NOT_FOUND",
		"ROUTE_NOT_FOUND",
		"BLACKHOLE_ROUTE_FOR_DESTINATION_FOUND",
		"INACTIVE_ROUTE_FOR_DESTINATION_FOUND",
		"TRANSIT_GATEWAY_ATTACHMENT_ATTACH_ARN_NO_MATCH",
		"MAX_HOPS_EXCEEDED",
		"POSSIBLE_MIDDLEBOX",
		"NO_DESTINATION_ARN_PROVIDED",
	}
}

type RouteAnalysisCompletionResultCode string

// Enum values for RouteAnalysisCompletionResultCode
const (
	RouteAnalysisCompletionResultCodeConnected    RouteAnalysisCompletionResultCode = "CONNECTED"
	RouteAnalysisCompletionResultCodeNotConnected RouteAnalysisCompletionResultCode = "NOT_CONNECTED"
)

// Values returns all known values for RouteAnalysisCompletionResultCode. Note
// that this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (RouteAnalysisCompletionResultCode) Values() []RouteAnalysisCompletionResultCode {
	return []RouteAnalysisCompletionResultCode{
		"CONNECTED",
		"NOT_CONNECTED",
	}
}

type RouteAnalysisStatus string

// Enum values for RouteAnalysisStatus
const (
	RouteAnalysisStatusRunning   RouteAnalysisStatus = "RUNNING"
	RouteAnalysisStatusCompleted RouteAnalysisStatus = "COMPLETED"
	RouteAnalysisStatusFailed    RouteAnalysisStatus = "FAILED"
)

// Values returns all known values for RouteAnalysisStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (RouteAnalysisStatus) Values() []RouteAnalysisStatus {
	return []RouteAnalysisStatus{
		"RUNNING",
		"COMPLETED",
		"FAILED",
	}
}

type RouteState string

// Enum values for RouteState
const (
	RouteStateActive    RouteState = "ACTIVE"
	RouteStateBlackhole RouteState = "BLACKHOLE"
)

// Values returns all known values for RouteState. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (RouteState) Values() []RouteState {
	return []RouteState{
		"ACTIVE",
		"BLACKHOLE",
	}
}

type RouteTableType string

// Enum values for RouteTableType
const (
	RouteTableTypeTransitGatewayRouteTable RouteTableType = "TRANSIT_GATEWAY_ROUTE_TABLE"
	RouteTableTypeCoreNetworkSegment       RouteTableType = "CORE_NETWORK_SEGMENT"
	RouteTableTypeNetworkFunctionGroup     RouteTableType = "NETWORK_FUNCTION_GROUP"
)

// Values returns all known values for RouteTableType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (RouteTableType) Values() []RouteTableType {
	return []RouteTableType{
		"TRANSIT_GATEWAY_ROUTE_TABLE",
		"CORE_NETWORK_SEGMENT",
		"NETWORK_FUNCTION_GROUP",
	}
}

type RouteType string

// Enum values for RouteType
const (
	RouteTypePropagated RouteType = "PROPAGATED"
	RouteTypeStatic     RouteType = "STATIC"
)

// Values returns all known values for RouteType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (RouteType) Values() []RouteType {
	return []RouteType{
		"PROPAGATED",
		"STATIC",
	}
}

type SegmentActionServiceInsertion string

// Enum values for SegmentActionServiceInsertion
const (
	SegmentActionServiceInsertionSendVia SegmentActionServiceInsertion = "send-via"
	SegmentActionServiceInsertionSendTo  SegmentActionServiceInsertion = "send-to"
)

// Values returns all known values for SegmentActionServiceInsertion. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (SegmentActionServiceInsertion) Values() []SegmentActionServiceInsertion {
	return []SegmentActionServiceInsertion{
		"send-via",
		"send-to",
	}
}

type SendViaMode string

// Enum values for SendViaMode
const (
	SendViaModeDualHop   SendViaMode = "dual-hop"
	SendViaModeSingleHop SendViaMode = "single-hop"
)

// Values returns all known values for SendViaMode. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (SendViaMode) Values() []SendViaMode {
	return []SendViaMode{
		"dual-hop",
		"single-hop",
	}
}

type SiteState string

// Enum values for SiteState
const (
	SiteStatePending   SiteState = "PENDING"
	SiteStateAvailable SiteState = "AVAILABLE"
	SiteStateDeleting  SiteState = "DELETING"
	SiteStateUpdating  SiteState = "UPDATING"
)

// Values returns all known values for SiteState. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (SiteState) Values() []SiteState {
	return []SiteState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"UPDATING",
	}
}

type TransitGatewayConnectPeerAssociationState string

// Enum values for TransitGatewayConnectPeerAssociationState
const (
	TransitGatewayConnectPeerAssociationStatePending   TransitGatewayConnectPeerAssociationState = "PENDING"
	TransitGatewayConnectPeerAssociationStateAvailable TransitGatewayConnectPeerAssociationState = "AVAILABLE"
	TransitGatewayConnectPeerAssociationStateDeleting  TransitGatewayConnectPeerAssociationState = "DELETING"
	TransitGatewayConnectPeerAssociationStateDeleted   TransitGatewayConnectPeerAssociationState = "DELETED"
)

// Values returns all known values for TransitGatewayConnectPeerAssociationState.
// Note that this can be expanded in the future, and so it is only as up to date as
// the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (TransitGatewayConnectPeerAssociationState) Values() []TransitGatewayConnectPeerAssociationState {
	return []TransitGatewayConnectPeerAssociationState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"DELETED",
	}
}

type TransitGatewayRegistrationState string

// Enum values for TransitGatewayRegistrationState
const (
	TransitGatewayRegistrationStatePending   TransitGatewayRegistrationState = "PENDING"
	TransitGatewayRegistrationStateAvailable TransitGatewayRegistrationState = "AVAILABLE"
	TransitGatewayRegistrationStateDeleting  TransitGatewayRegistrationState = "DELETING"
	TransitGatewayRegistrationStateDeleted   TransitGatewayRegistrationState = "DELETED"
	TransitGatewayRegistrationStateFailed    TransitGatewayRegistrationState = "FAILED"
)

// Values returns all known values for TransitGatewayRegistrationState. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (TransitGatewayRegistrationState) Values() []TransitGatewayRegistrationState {
	return []TransitGatewayRegistrationState{
		"PENDING",
		"AVAILABLE",
		"DELETING",
		"DELETED",
		"FAILED",
	}
}

type TunnelProtocol string

// Enum values for TunnelProtocol
const (
	TunnelProtocolGre     TunnelProtocol = "GRE"
	TunnelProtocolNoEncap TunnelProtocol = "NO_ENCAP"
)

// Values returns all known values for TunnelProtocol. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (TunnelProtocol) Values() []TunnelProtocol {
	return []TunnelProtocol{
		"GRE",
		"NO_ENCAP",
	}
}

type ValidationExceptionReason string

// Enum values for ValidationExceptionReason
const (
	ValidationExceptionReasonUnknownOperation      ValidationExceptionReason = "UnknownOperation"
	ValidationExceptionReasonCannotParse           ValidationExceptionReason = "CannotParse"
	ValidationExceptionReasonFieldValidationFailed ValidationExceptionReason = "FieldValidationFailed"
	ValidationExceptionReasonOther                 ValidationExceptionReason = "Other"
)

// Values returns all known values for ValidationExceptionReason. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ValidationExceptionReason) Values() []ValidationExceptionReason {
	return []ValidationExceptionReason{
		"UnknownOperation",
		"CannotParse",
		"FieldValidationFailed",
		"Other",
	}
}
