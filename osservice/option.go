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

package osservice

import (
	"mender"
)

type Option func(os *Os)

func WithTime(t Time) Option {
	return func(o *Os) {
		o.time = t
	}
}

func WithConfigFile(path string) Option {
	return func(o *Os) {
		var err error
		o.store, err = newLoadSaver(path)
		if err != nil {
			panic(err)
		}
	}
}

func WithWifi(wifi Wifi) Option {
	return func(o *Os) {
		o.wifi = wifi
	}
}

func WithStore(store Store) Option {
	return func(o *Os) {
		o.store = store
	}
}

func WithNet(net Net) Option {
	return func(o *Os) {
		o.net = net
	}
}

func WithPort(port int) Option {
	return func(o *Os) {
		o.port = port
	}
}

func WithUpdate(u Update) Option {
	return func(o *Os) {
		o.update = u
	}
}

func WithMender(c *mender.Client) Option {
	return func(o *Os) {
		o.update = newUpdateOs(c)
	}
}
