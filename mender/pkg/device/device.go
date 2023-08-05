package device

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"

	"golang.org/x/exp/slices"
)

type Device struct {
	path         string
	inventoryDir string
	identityDir  string
	infoDir      string
}

type Attribute struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

type Info struct {
	DeviceType   string `json:"device_type"`
	ArtifactName string `json:"artifact_name"`
}

const (
	DefaultPath         = "/usr/share/mender"
	DefaultInventoryDir = "inventory"
	DefaultIdentityDir  = "identity"
	DefaultInfoDir      = "info"
)

func New(opts ...Option) *Device {
	d := &Device{
		path:         DefaultPath,
		inventoryDir: DefaultInventoryDir,
		identityDir:  DefaultIdentityDir,
		infoDir:      DefaultInfoDir,
	}

	for _, opt := range opts {
		opt(d)
	}

	d.inventoryDir = path.Join(d.path, d.inventoryDir)
	d.identityDir = path.Join(d.path, d.identityDir)
	d.infoDir = path.Join(d.path, d.infoDir)

	return d
}

func (d *Device) Info() (i Info, err error) {
	const (
		deviceTypeKey   = "device_type"
		artifactNameKey = "artifact_name"
	)

	attrs, err := parseAttributes(d.infoDir)
	if err != nil {
		err = fmt.Errorf("failed to read info: %w", err)
		return
	}

	getInfoValue := func(key string) (*string, error) {
		idx := slices.IndexFunc(attrs, func(attribute Attribute) bool {
			return attribute.Name == key
		})
		if idx == -1 {
			return nil, fmt.Errorf("key %v not found", key)
		}

		if len(attrs[idx].Value) != 1 {
			return nil, fmt.Errorf("key %v not found", key)
		}

		return &attrs[idx].Value[0], nil
	}

	maybeDeviceType, err := getInfoValue(deviceTypeKey)
	if err != nil {
		return
	}
	maybeArtifactName, err := getInfoValue(artifactNameKey)
	if err != nil {
		return
	}

	i.DeviceType = *maybeDeviceType
	i.ArtifactName = *maybeArtifactName

	return
}

func (d *Device) ID() ([]Attribute, error) {
	return parseAttributes(d.identityDir)

}
func (d *Device) Attributes() ([]Attribute, error) {
	return parseAttributes(d.inventoryDir)
}

func parseAttributes(dir string) ([]Attribute, error) {
	runnables, err := listRunnables(dir)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for _, run := range runnables {
		cmd := exec.Command(run)
		cmd.Stdout = &buf

		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}
	attrsMap := make(map[string][]string)

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		key, value, got := bytes.Cut(scanner.Bytes(), []byte("="))
		if !got {
			continue
		}
		attrsMap[string(key)] = append(attrsMap[string(key)], string(value))
	}

	ret := make([]Attribute, 0, len(attrsMap))
	for k, v := range attrsMap {
		ret = append(ret, Attribute{
			Name:  k,
			Value: v,
		})
	}

	return ret, nil
}

func listRunnables(dir string) ([]string, error) {
	isRunnable := func(info fs.FileInfo) bool {
		const (
			runnableBits = os.FileMode(syscall.S_IXUSR | syscall.S_IXGRP | syscall.S_IXOTH)
		)
		return info.Mode()&runnableBits != 0
	}

	var runnable []string

	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if path == dir {
				return nil
			}
			return fs.SkipDir
		}
		if isRunnable(info) {
			runnable = append(runnable, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return runnable, nil
}
