package mender

type DeploymentStatus int

const (
	Downloading DeploymentStatus = iota
	PauseBeforeInstalling
	Installing
	PauseBeforeRebooting
	Rebooting
	PauseBeforeCommitting
	Success
	Failure
	AlreadyInstalled
)

type DeploymentSource struct {
	URI    string `json:"uri,omitempty"`
	Expire string `json:"expire,omitempty"`
}

type DeploymentArtifact struct {
	ID         string           `json:"id,omitempty"`
	Name       string           `json:"artifact_name"`
	Source     DeploymentSource `json:"source"`
	Compatible []string         `json:"device_types_compatible"`
}

type DeploymentInstructions struct {
	ID       string             `json:"id"`
	Artifact DeploymentArtifact `json:"artifact"`
}

type CurrentDeployment struct {
	State DeploymentStatus `json:"state"`
	*DeploymentInstructions
}

type Artifacts struct {
	Current *CurrentDeployment       `json:"current"`
	Archive []DeploymentInstructions `json:"archive"`
}
