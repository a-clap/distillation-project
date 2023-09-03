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
	"errors"
)

var (
	ErrNeedAuthentication    = errors.New("device is not authenticated, go to GUI and accept device")
	ErrNeedSignerVerifier    = errors.New("SignerVerifier is mandatory")
	ErrNeedServerURLAndToken = errors.New("server URL and teenantToken are mandatory")
	ErrNeedDevice            = errors.New("device is mandatory")
	ErrNeedDownloader        = errors.New("downloader is mandatory")
	ErrNeedInstaller         = errors.New("installer is mandatory")
	ErrNeedRebooter          = errors.New("rebooter is mandatory")
	ErrNeedLoadSaver         = errors.New("LoadSaver is mandatory")
	ErrNeedCallbacks         = errors.New("callbacks are mandatory")
	ErrNeedCommitter         = errors.New("committer is mandatory")
	ErrDuringUpdate          = errors.New("currently during update")
)
