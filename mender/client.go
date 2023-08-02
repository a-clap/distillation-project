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

	"github.com/a-clap/distillation-ota/pkg/mender/device"
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

	if maybeData := c.LoadSaver.Load(artifactsKey); maybeData != nil {
		if rawBytes, err := json.Marshal(maybeData); err == nil {
			_ = json.Unmarshal(rawBytes, &c.artifacts)
		}
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

	// Make sure we already don't have such artifacts
	if idx := slices.IndexFunc(c.artifacts.Archive, func(instructions DeploymentInstructions) bool {
		return instructions.ID == artifact.ID
	}); idx != -1 {
		return false, nil
	}

	c.artifacts.Archive = append(c.artifacts.Archive, artifact)
	if err := c.LoadSaver.Save(artifactsKey, c.artifacts); err != nil {
		return false, fmt.Errorf("failed to save artifacts: %w", err)
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
		Status: getDeploymentStatus(status),
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
		return fmt.Errorf("already during update")
	}

	instructions, err := c.getInstructions(artifactName)
	if err != nil {
		return err
	}

	c.updating.Store(true)
	go c.handleUpdate(instructions)

	return nil
}

func (c *Client) IsDuringUpdate() bool {
	return c.updating.Load()
}

func (c *Client) StopUpdate() error {
	// If not during update process - nothing to do
	if !c.IsDuringUpdate() {
		return nil
	}

	return nil
}

func (c *Client) handleDownload(artifactName, srcUri string) (string, error) {
	if err := c.NotifyServer(Downloading, artifactName); err != nil {
		c.Callbacks.Error(fmt.Errorf("failed to send Downloading status: %w", err))
		c.doCleanup(Failure, artifactName)
		return "", err
	}

	dst := path.Join(os.TempDir(), artifactName, ".tmp")
	downloading, errs, err := c.Downloader.Download(dst, srcUri)
	if err != nil {
		c.Callbacks.Error(fmt.Errorf("download %v failed: %w", srcUri, err))
		c.doCleanup(Failure, artifactName)
		return "", err
	}

	progress := 0
	for progress < 100 {
		select {
		case progress = <-downloading:
			c.Callbacks.Update(Downloading, progress)
		case err := <-errs:
			c.Callbacks.Error(err)
		}
	}

	// Download finished - notify server
	if err := c.NotifyServer(PauseBeforeInstalling, artifactName); err != nil {
		c.Callbacks.Error(fmt.Errorf("failed to notify server: %w", err))
		c.doCleanup(Failure, artifactName)
		return "", err
	}

	// PauseBeforeInstall - notify user
	c.Callbacks.Update(PauseBeforeInstalling, 100)

	return dst, nil
}

func (c *Client) handleReboot(artifactName string) error {
	if err := c.NotifyServer(Rebooting, artifactName); err != nil {
		c.Error(fmt.Errorf("failed to send Rebooting status: %w", err))
		c.doCleanup(Failure, artifactName)
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
	if err := c.NotifyServer(Installing, artifactName); err != nil {
		c.Callbacks.Error(fmt.Errorf("failed to send Installing status: %w", err))
		c.doCleanup(Failure, artifactName)
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
			c.Callbacks.Error(err)
		}
	}

	// Install finished - notify server
	if err := c.NotifyServer(PauseBeforeRebooting, artifactName); err != nil {
		c.Callbacks.Error(fmt.Errorf("failed to notify server: %w", err))
		c.doCleanup(Failure, artifactName)
		return err
	}

	// Install finished - notify user
	c.Callbacks.Update(PauseBeforeRebooting, 100)
	return nil

}

func (c *Client) doCleanup(status DeploymentStatus, artifactName string) {
	_ = c.NotifyServer(status, artifactName)
	c.updating.Store(false)
}

func (c *Client) handleUpdate(ins *DeploymentInstructions) {
	// TODO: maybe it should be more like state machine
	// What if we will want to start from specific stage?
	shouldContinue := func(status DeploymentStatus) bool {
		if next := c.Callbacks.NextState(status); !next {
			// Cleanup
			c.doCleanup(Failure, ins.Artifact.Name)
			return false
		}
		return true
	}

	// Download
	dst, err := c.handleDownload(ins.Artifact.Name, ins.Artifact.Source.URI)
	if err != nil {
		return
	}

	// Wait for confirmation
	if !shouldContinue(Installing) {
		return
	}

	// Install
	if err := c.handleInstall(ins.Artifact.Name, dst); err != nil {
		return
	}

	// Wait for confirmation
	if !shouldContinue(Rebooting) {
		return
	}

	// Reboot
	if err := c.handleReboot(ins.Artifact.Name); err != nil {
		return
	}

}

// Verify overloads SignerVerifier - it checks first, if sig is base64 encoded
func (c *Client) Verify(data []byte, sig []byte) error {
	sig = maybeDecodeBase64(sig)
	return c.Signer.Verify(data, sig)
}

// getInstructions finds proper instructions in internal DeploymentInstructions
func (c *Client) getInstructions(artifactName string) (*DeploymentInstructions, error) {
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
