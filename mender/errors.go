package mender

import (
	"errors"
)

var (
	ErrNeedAuthentication    = errors.New("device is not authenticated, go to GUI and accept device")
	ErrNeedSignerVerifier    = errors.New("SignerVerifier is mandatory")
	ErrNeedServerURLAndToken = errors.New("server URL and teenantToken are mandatory")
	ErrNeedDevice            = errors.New("Device is mandatory")
	ErrNeedDownloader        = errors.New("Downloader is mandatory")
	ErrNeedInstaller         = errors.New("Installer is mandatory")
	ErrNeedRebooter          = errors.New("Rebooter is mandatory")
	ErrNeedLoadSaver         = errors.New("LoadSaver is mandatory")
	ErrNeedCallbacks         = errors.New("Callbacks are mandatory")
)
