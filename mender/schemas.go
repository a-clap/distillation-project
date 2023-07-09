package mender

type DeploymentStatus int

const (
	Downloading DeploymentStatus = iota
	PauseBeforeInstalling
	Installing
	PauseBeforeRebooting
	Rebooting
	PauseBeforeCommiting
	Success
	Failure
	AlreadyInstalled
)

type AuthRequest struct {
	ID     string `json:"id_data"`
	PubKey string `json:"pubkey"`
	Token  string `json:"tenant_token,omitempty"`
}

type UpdateStatus struct {
	Progress int
	Status   DeploymentStatus
	Error    error
}

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

type ArtifactDeploymentStatus struct {
	State DeploymentStatus `json:"state"`
	DeploymentArtifact
}

type Artifacts struct {
	Current *ArtifactDeploymentStatus `json:"current"`
	Archive []DeploymentInstructions  `json:"archive"`
}
