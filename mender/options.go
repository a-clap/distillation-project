package mender

import (
	"time"
)

type Option func(c *Client) error

func WithSigner(sv Signer) Option {
	return func(c *Client) error {
		c.Signer = sv
		return nil
	}
}

func WithServer(url, teenantToken string) Option {
	return func(c *Client) error {
		c.paths = newServerPaths(url)
		c.teenantToken = teenantToken
		return nil
	}
}

func WithDevice(dev Device) Option {
	return func(c *Client) error {
		c.Device = dev
		return nil
	}
}

func WithTimeout(t time.Duration) Option {
	return func(c *Client) error {
		c.Timeout = t
		return nil
	}
}

func WithDownloader(d Downloader) Option {
	return func(c *Client) error {
		c.Downloader = d
		return nil
	}
}

func WithInstaller(i Installer) Option {
	return func(c *Client) error {
		c.Installer = i
		return nil
	}
}

func WithRebooter(r Rebooter) Option {
	return func(c *Client) error {
		c.Rebooter = r
		return nil
	}
}

func WithLoadSaver(saver LoadSaver) Option {
	return func(c *Client) error {
		c.LoadSaver = saver
		return nil
	}
}

func WithCallbacks(callbacks Callbacks) Option {
	return func(c *Client) error {
		c.Callbacks = callbacks
		return nil
	}
}
