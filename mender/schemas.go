package mender

type DeploymentStatus int

const (
	Downloading DeploymentStatus = iota + 1
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

type MenderDeploymentArtifact struct {
	ID         string           `json:"id,omitempty"`
	Name       string           `json:"artifact_name"`
	Source     DeploymentSource `json:"source"`
	Compatible []string         `json:"device_types_compatible"`
}

type MenderDeploymentInstructions struct {
	ID       string                   `json:"id"`
	Artifact MenderDeploymentArtifact `json:"artifact"`
}

type Artifact struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Source     DeploymentSource `json:"source"`
	Compatible []string         `json:"compatible"`
}

type Instructions struct {
	ID       string   `json:"id"`
	Artifact Artifact `json:"artifact"`
}

type CurrentDeployment struct {
	State        DeploymentStatus `json:"state"`
	Instructions Instructions     `json:"instructions"`
}

type Artifacts struct {
	Current *CurrentDeployment             `json:"current"`
	Archive []MenderDeploymentInstructions `json:"archive"`
}
