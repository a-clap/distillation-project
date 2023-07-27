package mender

import (
	"time"

	"github.com/a-clap/distillation-ota/pkg/mender/downloader"
	"github.com/a-clap/distillation-ota/pkg/mender/installer"
	"github.com/a-clap/distillation-ota/pkg/mender/loadsaver"
	"github.com/a-clap/distillation-ota/pkg/mender/rebooter"
	"github.com/a-clap/distillation-ota/pkg/mender/signer"
)

type Builder struct {
	url, token     string
	device         Device
	signerVerifier Signer
	timeout        time.Duration
	downloader     Downloader
	installer      Installer
	rebooter       Rebooter
	loadSaver      LoadSaver
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
	}

	return b
}

func (b *Builder) Build() (*Client, error) {
	opts := make([]Option, 0, 8)

	if b.signerVerifier == nil {
		signer, err := signer.New(signer.WithNewKey())
		if err != nil {
			return nil, err
		}
		b.signerVerifier = &builderSignerVerifier{signer}
	}

	opts = append(opts, WithSigner(b.signerVerifier))
	opts = append(opts, WithTimeout(b.timeout))
	opts = append(opts, WithDownloader(b.downloader))
	opts = append(opts, WithInstaller(b.installer))
	opts = append(opts, WithRebooter(b.rebooter))

	if b.device != nil {
		opts = append(opts, WithDevice(b.device))
	}

	if b.url != "" && b.token != "" {
		opts = append(opts, WithServer(b.url, b.token))
	}

	if b.loadSaver != nil {
		opts = append(opts, WithLoadSaver(b.loadSaver))
	}

	return New(opts...)
}

func (b *Builder) WithSignerVerifier(s Signer) *Builder {
	b.signerVerifier = s
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

func (b *Builder) WithStore(configFilePath string) *Builder {
	if saver, err := loadsaver.New(configFilePath); err == nil {
		b.loadSaver = saver
	}
	return b
}

var (
	_ Downloader = (*builderDownloader)(nil)
	_ Installer  = (*builderInstaller)(nil)
	_ Rebooter   = (*builderRebooter)(nil)
	_ Signer     = (*builderSignerVerifier)(nil)
)

type builderSignerVerifier struct {
	*signer.Signer
}

type builderRebooter struct {
}

// Reboot implements Rebooter.
func (*builderRebooter) Reboot() error {
	return rebooter.Reboot()
}

type builderInstaller struct {
	*installer.Installer
}

// Install implements Installer.
func (b *builderInstaller) Install(src string) (progress chan int, errCh chan error, err error) {
	return b.Installer.Install(src)
}

type builderDownloader struct {
}

// Download implements Downloader.
func (*builderDownloader) Download(dst string, srcUrl string) (progress chan int, errCh chan error, err error) {
	return downloader.Download(dst, srcUrl)
}
