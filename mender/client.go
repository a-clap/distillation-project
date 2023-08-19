package mender

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync/atomic"
	"time"

	"mender/pkg/device"

	"github.com/carlmjohnson/requests"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -package mocks -destination mocks/mocks_mender.go . Signer,Device,Downloader,Installer,Rebooter,LoadSaver,Callbacks
type Client struct {
	Timeout time.Duration
	Signer
	Device
	Downloader
	Installer
	Rebooter
	LoadSaver
	Callbacks

	teenantToken string
	jwtToken     string
	paths        *serverPaths
	artifacts    Artifacts
	updating     atomic.Bool
}

const (
	defaultTimeout = 1 * time.Second
	artifactsKey   = "artifacts"
)

func New(options ...Option) (*Client, error) {
	c := &Client{
		Timeout: defaultTimeout,
	}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if err := c.verify(); err != nil {
		return nil, err
	}

	c.loadArtifacts()

	// Should we continue update?
	if c.artifacts.Current != nil {
		// We need token
		if err := c.Connect(); err != nil {
			return nil, err
		}
		c.updateFromState(c.artifacts.Current.State, c.artifacts.Current.DeploymentInstructions)
	}

	return c, nil
}

// verify is responsible for checking if Client is provided with correct options
func (c *Client) verify() error {
	var errs []error
	if c.Signer == nil {
		errs = append(errs, ErrNeedSignerVerifier)
	}

	if c.paths == nil {
		errs = append(errs, ErrNeedServerURLAndToken)
	}

	if c.Device == nil {
		errs = append(errs, ErrNeedDevice)
	}

	if c.Downloader == nil {
		errs = append(errs, ErrNeedDownloader)
	}

	if c.Installer == nil {
		errs = append(errs, ErrNeedInstaller)
	}

	if c.Rebooter == nil {
		errs = append(errs, ErrNeedRebooter)
	}

	if c.LoadSaver == nil {
		errs = append(errs, ErrNeedLoadSaver)
	}

	if c.Callbacks == nil {
		errs = append(errs, ErrNeedCallbacks)
	}

	return errors.Join(errs...)
}

func (c *Client) Connect() error {
	ids, err := c.Device.ID()
	if err != nil {
		return err
	}
	// Encode ids to json
	idsMap := make(map[string]interface{})
	for _, id := range ids {
		idsMap[id.Name] = id.Value
	}

	id, err := json.Marshal(idsMap)
	if err != nil {
		return err
	}

	auth := struct {
		ID     string `json:"id_data"`
		PubKey string `json:"pubkey"`
		Token  string `json:"tenant_token"`
	}{
		ID:     string(id),
		PubKey: c.PublicKeyPEM(),
		Token:  c.teenantToken,
	}

	jsonReq, err := json.Marshal(auth)
	if err != nil {
		panic(err)
	}

	sign, err := c.Sign(jsonReq)
	if err != nil {
		panic(err)
	}

	sig := base64.StdEncoding.EncodeToString(sign)
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	err = requests.
		URL(c.paths.AuthRequest()).
		ContentType("application/json").
		Accept("application/json").
		Header("X-MEN-Signature", sig).
		BodyBytes(jsonReq).
		ToString(&c.jwtToken).
		Fetch(ctx)

	// Check for specific statuses
	if err != nil {
		if requests.HasStatusErr(err, http.StatusUnauthorized) {
			return ErrNeedAuthentication
		}
		return fmt.Errorf("http request error: %w", err)
	}

	return nil
}

func (c *Client) UpdateInventory() error {
	attrs, err := c.Device.Attributes()
	if err != nil {
		return err
	}

	info, err := c.Device.Info()
	if err != nil {
		return err
	}
	attrs = append(attrs,
		device.Attribute{
			Name:  "device_type",
			Value: []string{info.DeviceType},
		},
		device.Attribute{
			Name:  "artifact_name",
			Value: []string{info.ArtifactName},
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	err = requests.
		URL(c.paths.Inventory()).
		ContentType("application/json").
		Accept("application/json").
		BodyJSON(attrs).
		Bearer(c.jwtToken).
		Patch().
		Fetch(ctx)

	return err
}

func (c *Client) PullReleases() (newRelease bool, err error) {
	if c.IsDuringUpdate() {
		return false, ErrDuringUpdate
	}

	info, err := c.Device.Info()
	if err != nil {
		return
	}

	params := map[string][]string{
		"artifact_name": {info.ArtifactName},
		"device_type":   {info.DeviceType},
	}

	var artifact DeploymentInstructions

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	err = requests.
		URL(c.paths.Deployment()).
		ContentType("application/json").
		Params(params).
		Accept("application/json").
		ToJSON(&artifact).
		AddValidator(requests.CheckStatus(http.StatusOK)).
		Bearer(c.jwtToken).
		Fetch(ctx)

	if err != nil {
		// Status 204 - No updates for device
		if requests.HasStatusErr(err, http.StatusNoContent) {
			return false, nil
		}
		return false, fmt.Errorf("http request error: %w", err)
	}

	// Make sure this release is compatible with us
	if idx := slices.Index(artifact.Artifact.Compatible, info.DeviceType); idx == -1 {
		return false, nil
	}

	// Make sure we already don't have such artifact
	if _, err := c.getInstructions(artifact.Artifact.Name); err == nil {
		return false, nil
	}

	c.artifacts.Archive = append(c.artifacts.Archive, artifact)
	if err := c.saveArtifacts(); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) AvailableReleases() []string {
	releases := make([]string, 0, len(c.artifacts.Archive))
	for _, artifact := range c.artifacts.Archive {
		releases = append(releases, artifact.Artifact.Name)
	}
	return releases
}

func (c *Client) NotifyServer(status DeploymentStatus, artifactName string) error {
	ins, err := c.getInstructions(artifactName)
	if err != nil {
		return err
	}

	jsonStatus := struct {
		Status string `json:"status"`
	}{
		Status: toServerStatus(status),
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	err = requests.
		URL(c.paths.DeploymentStatus(ins.ID)).
		ContentType("application/json").
		Accept("application/json").
		BodyJSON(jsonStatus).
		Bearer(c.jwtToken).
		Put().
		Fetch(ctx)

	return err
}

func (c *Client) Update(artifactName string) error {
	if c.IsDuringUpdate() {
		return ErrDuringUpdate
	}

	instructions, err := c.getInstructions(artifactName)
	if err != nil {
		return err
	}

	c.updateFromState(Downloading, instructions)

	return nil
}

func (c *Client) IsDuringUpdate() bool {
	return c.updating.Load()
}

func (c *Client) StopUpdate() error {
	if !c.IsDuringUpdate() {
		return nil
	}

	return nil
}

func (c *Client) updateFromState(state DeploymentStatus, ins *DeploymentInstructions) {
	c.updating.Store(true)
	c.artifacts.Current = &CurrentDeployment{
		State:                  state,
		DeploymentInstructions: ins,
	}
	go c.handleUpdate()
}

func (c *Client) handleUpdate() {
	var (
		artifactName   = c.artifacts.Current.DeploymentInstructions.Artifact.Name
		srcURL         = c.artifacts.Current.DeploymentInstructions.Artifact.Source.URI
		err            error
		downloadedFile string
	)

	shouldContinue := func(status DeploymentStatus) bool {
		if next := c.Callbacks.NextState(status); !next {
			// Cleanup
			c.doCleanup(Failure, srcURL)
			return false
		}
		return true
	}

	for c.updating.Load() {
		switch c.artifacts.Current.State {
		case Downloading:
			// Save c.artifacts.Current.State
			downloadedFile, err = c.handleDownload(artifactName, srcURL)
			if err != nil {
				return
			}
			c.artifacts.Current.State = PauseBeforeInstalling
		case PauseBeforeInstalling:
			if !shouldContinue(Installing) {
				return
			}
			c.artifacts.Current.State = Installing
		case Installing:
			if err := c.handleInstall(artifactName, downloadedFile); err != nil {
				return
			}
			c.artifacts.Current.State = PauseBeforeRebooting
		case PauseBeforeRebooting:
			if !shouldContinue(Rebooting) {
				return
			}
			c.artifacts.Current.State = Rebooting
		case Rebooting:
			// Store state before reboot
			c.artifacts.Current.State = PauseBeforeCommitting
			if err := c.saveArtifacts(); err != nil {
				c.Callbacks.Error(err)
				c.doCleanup(Failure, artifactName)
				return
			}
			// Execute reboot
			if err := c.handleReboot(artifactName); err != nil {
				return
			}
			// Well... we shouldn't be here
			c.updating.Store(false)
		case PauseBeforeCommitting:
			if !shouldContinue(Success) {
				return
			}
			c.artifacts.Current.State = Success
		case Success:
			// Everything went fine
			c.handleSuccess(artifactName)
		case Failure:
			c.updating.Store(false)
		default:
			// Shouldn't achieve this state
			c.Callbacks.Error(fmt.Errorf("unsupported state: %v", c.artifacts.Current.State))
			c.updating.Store(false)
		}
	}
}

func (c *Client) notifyServerDuringUpdate(status DeploymentStatus, artifactName string) error {
	if err := c.NotifyServer(status, artifactName); err != nil {
		c.Callbacks.Error(fmt.Errorf("failed to notify server with status %v: %w", toServerStatus(status), err))
		c.doCleanup(Failure, artifactName)
		return err
	}
	return nil
}

func (c *Client) handleSuccess(artifactName string) {
	// Finished
	c.updating.Store(false)

	// Notify server that we are done
	if err := c.NotifyServer(Success, artifactName); err != nil {
		c.Callbacks.Error(fmt.Errorf("failed to notify server with status %v: %w", toServerStatus(Success), err))
	}

	// Remove just installed artifact from Archive
	c.artifacts.Archive = slices.DeleteFunc(c.artifacts.Archive, func(instructions DeploymentInstructions) bool {
		return c.artifacts.Current.Artifact.Name == instructions.Artifact.Name
	})

	// And current artifact itself
	c.artifacts.Current = nil
	// Store updated artifacts
	if err := c.saveArtifacts(); err != nil {
		c.Callbacks.Error(err)
	}
}

func (c *Client) handleDownload(artifactName, srcURL string) (string, error) {
	if err := c.notifyServerDuringUpdate(Downloading, artifactName); err != nil {
		return "", err
	}

	dst := path.Join(os.TempDir(), artifactName+".mender")
	downloading, errs, err := c.Downloader.Download(dst, srcURL)
	if err != nil {
		c.Callbacks.Error(fmt.Errorf("download %v failed: %w", srcURL, err))
		c.doCleanup(Failure, artifactName)
		return "", err
	}

	progress := 0
	for progress < 100 {
		select {
		case progress = <-downloading:
			c.Callbacks.Update(Downloading, progress)
		case err := <-errs:
			c.Callbacks.Error(fmt.Errorf("download %v failed: %w", srcURL, err))
			c.doCleanup(Failure, artifactName)
			return "", err
		}
	}

	// Download finished - notify server
	if err := c.notifyServerDuringUpdate(PauseBeforeInstalling, artifactName); err != nil {
		return "", err
	}

	// PauseBeforeInstall - notify user
	c.Callbacks.Update(PauseBeforeInstalling, 100)

	return dst, nil
}

func (c *Client) handleReboot(artifactName string) error {
	if err := c.notifyServerDuringUpdate(Rebooting, artifactName); err != nil {
		return err
	}

	// we can notify user, however.. after a second we will reboot so, doest it matter?
	c.Callbacks.Update(Rebooting, 1)

	if err := c.Rebooter.Reboot(); err != nil {
		c.Callbacks.Error(fmt.Errorf("failed to reboot: %w", err))
		c.doCleanup(Failure, artifactName)
		return err
	}
	return nil
}

func (c *Client) handleInstall(artifactName, src string) error {
	if err := c.notifyServerDuringUpdate(Installing, artifactName); err != nil {
		return err
	}

	progressChan, errs, err := c.Installer.Install(src)
	if err != nil {
		c.Callbacks.Error(fmt.Errorf("install %v failed: %w", src, err))
		c.doCleanup(Failure, artifactName)
		return err
	}

	progress := 0
	for progress < 100 {
		select {
		case progress = <-progressChan:
			c.Callbacks.Update(Installing, progress)
		case err := <-errs:
			c.Callbacks.Error(fmt.Errorf("install %v failed: %w", src, err))
			c.doCleanup(Failure, artifactName)
			return err
		}
	}

	// Install finished - notify server
	if err := c.notifyServerDuringUpdate(PauseBeforeRebooting, artifactName); err != nil {
		return err
	}

	// Install finished - notify user
	c.Callbacks.Update(PauseBeforeRebooting, 100)
	return nil

}

func (c *Client) loadArtifacts() {
	maybeData := c.LoadSaver.Load(artifactsKey)
	if maybeData == nil {
		return
	}

	rawBytes, err := json.Marshal(maybeData)
	if err != nil {
		return
	}

	_ = json.Unmarshal(rawBytes, &c.artifacts)

}

func (c *Client) saveArtifacts() error {
	if err := c.LoadSaver.Save(artifactsKey, c.artifacts); err != nil {
		return fmt.Errorf("failed to save artifacts: %w", err)
	}
	return nil
}

func (c *Client) doCleanup(status DeploymentStatus, artifactName string) {
	_ = c.NotifyServer(status, artifactName)
	c.updating.Store(false)
}

// Verify overloads SignerVerifier - it checks first, if sig is base64 encoded
func (c *Client) Verify(data []byte, sig []byte) error {
	sig = maybeDecodeBase64(sig)
	return c.Signer.Verify(data, sig)
}

// getInstructions finds proper instructions in internal DeploymentInstructions
func (c *Client) getInstructions(artifactName string) (*DeploymentInstructions, error) {
	// Maybe we are using current artifact
	if c.artifacts.Current != nil &&
		c.artifacts.Current.DeploymentInstructions.Artifact.Name == artifactName {
		return c.artifacts.Current.DeploymentInstructions, nil
	}

	idx := slices.IndexFunc(c.artifacts.Archive, func(instructions DeploymentInstructions) bool {
		return instructions.Artifact.Name == artifactName
	})

	if idx == -1 {
		return nil, fmt.Errorf("artifact %v not found", artifactName)
	}
	return &c.artifacts.Archive[idx], nil
}

// maybeDecodeBase64 checks if it is possible to decode string with base64 encoding
// if it is - returned decoded bytes, if not - leave buffer as is
func maybeDecodeBase64(sig []byte) []byte {
	decoded, err := base64.StdEncoding.DecodeString(string(sig))
	if err != nil {
		return sig
	}
	return decoded
}
