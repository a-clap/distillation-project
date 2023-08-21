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

func WithCommitter(comm Committer) Option {
	return func(c *Client) error {
		c.Committer = comm
		return nil
	}
}

func WithCallbacks(callbacks Callbacks) Option {
	return func(c *Client) error {
		c.Callbacks = callbacks
		return nil
	}
}
