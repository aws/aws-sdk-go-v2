// Code generated by smithy-go-codegen DO NOT EDIT.

package types

type ApplicationState string

// Enum values for ApplicationState
const (
	ApplicationStateCreating   ApplicationState = "CREATING"
	ApplicationStateCreated    ApplicationState = "CREATED"
	ApplicationStateStarting   ApplicationState = "STARTING"
	ApplicationStateStarted    ApplicationState = "STARTED"
	ApplicationStateStopping   ApplicationState = "STOPPING"
	ApplicationStateStopped    ApplicationState = "STOPPED"
	ApplicationStateTerminated ApplicationState = "TERMINATED"
)

// Values returns all known values for ApplicationState. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ApplicationState) Values() []ApplicationState {
	return []ApplicationState{
		"CREATING",
		"CREATED",
		"STARTING",
		"STARTED",
		"STOPPING",
		"STOPPED",
		"TERMINATED",
	}
}

type Architecture string

// Enum values for Architecture
const (
	ArchitectureArm64 Architecture = "ARM64"
	ArchitectureX8664 Architecture = "X86_64"
)

// Values returns all known values for Architecture. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (Architecture) Values() []Architecture {
	return []Architecture{
		"ARM64",
		"X86_64",
	}
}

type JobRunMode string

// Enum values for JobRunMode
const (
	JobRunModeBatch     JobRunMode = "BATCH"
	JobRunModeStreaming JobRunMode = "STREAMING"
)

// Values returns all known values for JobRunMode. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (JobRunMode) Values() []JobRunMode {
	return []JobRunMode{
		"BATCH",
		"STREAMING",
	}
}

type JobRunState string

// Enum values for JobRunState
const (
	JobRunStateSubmitted  JobRunState = "SUBMITTED"
	JobRunStatePending    JobRunState = "PENDING"
	JobRunStateScheduled  JobRunState = "SCHEDULED"
	JobRunStateRunning    JobRunState = "RUNNING"
	JobRunStateSuccess    JobRunState = "SUCCESS"
	JobRunStateFailed     JobRunState = "FAILED"
	JobRunStateCancelling JobRunState = "CANCELLING"
	JobRunStateCancelled  JobRunState = "CANCELLED"
	JobRunStateQueued     JobRunState = "QUEUED"
)

// Values returns all known values for JobRunState. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (JobRunState) Values() []JobRunState {
	return []JobRunState{
		"SUBMITTED",
		"PENDING",
		"SCHEDULED",
		"RUNNING",
		"SUCCESS",
		"FAILED",
		"CANCELLING",
		"CANCELLED",
		"QUEUED",
	}
}
