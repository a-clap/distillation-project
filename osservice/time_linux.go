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
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

const dateExecutable = "timedatectl"

var (
	ntpServiceRE      = regexp.MustCompile(`NTP service:\s?(.+)\n`)
	_            Time = timeOs{}
)

// Linux timeOs implementation
type timeOs struct{}

func (timeOs) Now() (time.Time, error) {
	return time.Now(), nil
}

func (timeOs) NTP() (bool, error) {
	out, err := exec.Command(dateExecutable, "status").Output()
	if err != nil {
		return false, fmt.Errorf("ntp: %w", err)
	}

	got := ntpServiceRE.FindStringSubmatch(string(out))
	if len(got) != 2 {
		return false, fmt.Errorf("ntp parse output: %v", got)
	}

	switch got[1] {
	case "inactive":
		return false, nil
	case "active":
		return true, nil
	}

	return false, fmt.Errorf("unknown output %v", got[1])
}

func (timeOs) SetNTP(enable bool) error {
	args := []string{"set-ntp", "false"}
	if enable {
		args[1] = "true"
	}

	if err := exec.Command(dateExecutable, args...).Run(); err != nil {
		return fmt.Errorf("configure ntp: %w", err)
	}

	return nil
}

func (timeOs) SetNow(now time.Time) error {
	dateString := now.Format("2006-01-02 15:04:05")
	args := []string{"set-time", dateString}

	if err := exec.Command(dateExecutable, args...).Run(); err != nil {
		return fmt.Errorf("setting time: %w", err)
	}

	return nil
}
