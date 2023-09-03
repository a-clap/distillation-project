// MIT License
//
// Copyright (c) 2023 a-clap
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mender

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"mender/pkg/committer"
	"mender/pkg/device"
	"mender/pkg/downloader"
	"mender/pkg/installer"
	"mender/pkg/loadsaver"
	"mender/pkg/rebooter"
	"mender/pkg/signer"
)

type Builder struct {
	url, token     string
	callbacks      Callbacks
	device         Device
	signerVerifier Signer
	timeout        time.Duration
	downloader     Downloader
	installer      Installer
	rebooter       Rebooter
	loadSaver      LoadSaver
	committer      Committer
}

func NewBuilder() *Builder {
	b := &Builder{
		url:            "",  // Must be provided
		token:          "",  // Must be provided
		device:         nil, // Must be provided
		signerVerifier: nil, // This will be created in Build method, as it takes a lot of time to create new key
		loadSaver:      nil,
		timeout:        defaultTimeout,
		downloader:     &builderDownloader{},
		installer:      &builderInstaller{installer.New()},
		rebooter:       &builderRebooter{},
		committer:      &builderCommiter{committer.New()},
	}

	return b
}

func (b *Builder) Build() (*Client, error) {
	opts := make([]Option, 0, 8)

	if b.signerVerifier == nil {
		sign, err := signer.New(signer.WithNewKey())
		if err != nil {
			return nil, fmt.Errorf("failed to create signer: %w", err)
		}
		b.signerVerifier = &builderSignerVerifier{sign}
	}

	if b.loadSaver == nil {
		saver, err := loadsaver.New(path.Join(os.TempDir(), "mender.json"))
		if err != nil {
			return nil, fmt.Errorf("failed to create loadSaver: %w", err)
		}
		b.loadSaver = &builderLoadSaver{saver}
	}

	if b.device == nil {
		b.device = device.New()
	}

	if b.url != "" && b.token != "" {
		opts = append(opts, WithServer(b.url, b.token))
	}

	if b.loadSaver != nil {
		opts = append(opts, WithLoadSaver(b.loadSaver))
	}

	opts = append(opts, WithSigner(b.signerVerifier))
	opts = append(opts, WithTimeout(b.timeout))
	opts = append(opts, WithDownloader(b.downloader))
	opts = append(opts, WithInstaller(b.installer))
	opts = append(opts, WithRebooter(b.rebooter))
	opts = append(opts, WithCallbacks(b.callbacks))
	opts = append(opts, WithDevice(b.device))
	opts = append(opts, WithCommitter(b.committer))

	return New(opts...)
}

func (b *Builder) WithSignerVerifier(s Signer) *Builder {
	b.signerVerifier = s
	return b
}

func (b *Builder) WithCommitter(c Committer) *Builder {
	b.committer = c
	return b
}

func (b *Builder) WithServer(url, token string) *Builder {
	b.url, b.token = url, token
	return b
}

func (b *Builder) WithDevice(d Device) *Builder {
	b.device = d
	return b
}

func (b *Builder) WithCallbacks(cb Callbacks) *Builder {
	b.callbacks = cb
	return b
}

func (b *Builder) WithTimeout(t time.Duration) *Builder {
	b.timeout = t
	return b
}

func (b *Builder) WithDownloader(d Downloader) *Builder {
	b.downloader = d
	return b
}

func (b *Builder) WithInstaller(i Installer) *Builder {
	b.installer = i
	return b
}

func (b *Builder) WithRebooter(r Rebooter) *Builder {
	b.rebooter = r
	return b
}

func (b *Builder) WithLoadSaver(saver LoadSaver) *Builder {
	b.loadSaver = saver
	return b
}

func (b *Builder) WithStore(file string) *Builder {
	saver, err := loadsaver.New(file)
	if err != nil {
		// Can't do anything about it, fail fast
		panic(err)
	}
	b.loadSaver = saver
	return b
}

func (b *Builder) WithStdIOInterface() *Builder {
	b.callbacks = &builderStdIOCallbacks{}
	return b
}

var (
	_ Downloader = (*builderDownloader)(nil)
	_ Installer  = (*builderInstaller)(nil)
	_ Rebooter   = (*builderRebooter)(nil)
	_ Signer     = (*builderSignerVerifier)(nil)
	_ LoadSaver  = (*builderLoadSaver)(nil)
	_ Callbacks  = (*builderStdIOCallbacks)(nil)
	_ Committer  = (*builderCommiter)(nil)
)

type builderCommiter struct {
	*committer.Committer
}

type builderSignerVerifier struct {
	*signer.Signer
}

type builderRebooter struct{}

type builderLoadSaver struct {
	*loadsaver.LoadSaver
}

func (*builderRebooter) Reboot() error {
	return rebooter.Reboot()
}

type builderInstaller struct {
	*installer.Installer
}

func (b *builderInstaller) Install(src string) (progress chan int, errCh chan error, err error) {
	return b.Installer.Install(src)
}

type builderDownloader struct{}

func (*builderDownloader) Download(dst string, srcURL string) (progress chan int, errCh chan error, err error) {
	return downloader.Download(dst, srcURL)
}

type builderStdIOCallbacks struct {
	lastStatus DeploymentStatus
}

func (b *builderStdIOCallbacks) Update(status DeploymentStatus, progress int) {
	if b.lastStatus != status {
		fmt.Println(toReadableStatus(status), "...")
		b.lastStatus = status
	}

	b.updateProgressBar(progress)
	if progress == 100 {
		// Add new line
		fmt.Printf("\n")
		fmt.Println(toReadableStatus(status), "finished!")
	}
}

func (b *builderStdIOCallbacks) NextState(status DeploymentStatus) bool {
	validResponse := false
	continueUpdate := false
	scanner := bufio.NewScanner(os.Stdin)

	for !validResponse {
		fmt.Println(`Updater want to move to next state:`, toServerStatus(status))
		fmt.Println(`Would you like to proceed? [Y\n]`)

		scanner.Scan()
		fromUser := scanner.Text()

		validResponse = fromUser == "Y" || fromUser == "n"
		continueUpdate = fromUser == "Y"

		if !validResponse {
			fmt.Println(`I don't understand, please try again.`)
		}
	}

	if continueUpdate {
		fmt.Println("I will continue update")
	} else {
		fmt.Println("Update aborted by user")
	}

	return continueUpdate
}

func (b *builderStdIOCallbacks) Error(err error) {
	fmt.Println("Updater error:", err)
}

func (b *builderStdIOCallbacks) updateProgressBar(v int) {
	if v == 0 {
		fmt.Printf("\r[>%s] %v %%", strings.Repeat(" ", 100), v)
	} else if v == 100 {
		fmt.Printf("\r[%s] %v %%", strings.Repeat("=", 100), v)
	} else {
		fmt.Printf("\r[%s>%s] %v %%", strings.Repeat("=", v), strings.Repeat(" ", 100-v), v)
	}
}
