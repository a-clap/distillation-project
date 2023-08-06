package mender

import (
	"fmt"
	"strings"
)

const (
	apiAuthRequestv1    = "api/devices/v1/authentication/auth_requests"
	apiInventory        = "api/devices/v1/inventory/device/attributes"
	apiDeployment       = "api/devices/v1/deployments/device/deployments/next"
	apiDeploymentStatus = "api/devices/v1/deployments/device/deployments/%v/status"
)

type serverPaths struct {
	base string
}

func newServerPaths(base string) *serverPaths {
	// Check, if server starts with httpX - http:// or https://
	if !strings.HasPrefix(base, "http") {
		base = "https://" + base
	}

	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return &serverPaths{base: base}
}

func (s *serverPaths) AuthRequest() string {
	return s.base + apiAuthRequestv1
}

func (s *serverPaths) Inventory() string {
	return s.base + apiInventory
}

func (s *serverPaths) Deployment() string {
	return s.base + apiDeployment
}

func (s *serverPaths) DeploymentStatus(id string) string {
	return fmt.Sprintf(s.base+apiDeploymentStatus, id)
}
