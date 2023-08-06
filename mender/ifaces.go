package mender

import (
	"mender/pkg/device"
)

var (
	toServerStatus = func() func(s DeploymentStatus) string {
		toStatus := map[DeploymentStatus]string{
			Downloading:           "downloading",
			PauseBeforeInstalling: "pause_before_installing",
			Installing:            "installing",
			PauseBeforeRebooting:  "pause_before_rebooting",
			Rebooting:             "rebooting",
			PauseBeforeCommitting: "pause_before_committing",
			Success:               "success",
			Failure:               "failure",
			AlreadyInstalled:      "already-installed",
		}
		return func(s DeploymentStatus) string {
			return toStatus[s]
		}
	}()

	toReadableStatus = func() func(s DeploymentStatus) string {
		toStatus := map[DeploymentStatus]string{
			Downloading:           "Downloading",
			PauseBeforeInstalling: "Pause before installing",
			Installing:            "Installing",
			PauseBeforeRebooting:  "Pause before rebooting",
			Rebooting:             "Rebooting",
			PauseBeforeCommitting: "Pause before committing",
			Success:               "Success",
			Failure:               "Failure",
			AlreadyInstalled:      "Already installed",
		}
		return func(s DeploymentStatus) string {
			return toStatus[s]
		}
	}()
)

type Signer interface {
	Sign(data []byte) ([]byte, error)
	Verify(data []byte, sig []byte) error // I think it is not necessary
	PublicKeyPEM() string
}

type Device interface {
	Info() (device.Info, error)
	ID() ([]device.Attribute, error)
	Attributes() ([]device.Attribute, error)
}

// Installer installs artifact on device
type Installer interface {
	Install(src string) (progress chan int, errCh chan error, err error) // progress should return int in range <0, 100> (%)
}

// Downloader downloads file from url
type Downloader interface {
	Download(dst, srcURL string) (progress chan int, errCh chan error, err error) // progress should return int in range <0, 100> (%)
}

// Rebooter allows to reboot device
type Rebooter interface {
	Reboot() error // Reboot should reboot device but also store information somewhere, if device is after reboot
}

// LoadSaver allows to save and load data in persistent storage
type LoadSaver interface {
	Save(key string, data interface{}) error
	Load(key string) interface{}
}

type Callbacks interface {
	Update(status DeploymentStatus, progress int)
	NextState(status DeploymentStatus) bool
	Error(err error)
}
