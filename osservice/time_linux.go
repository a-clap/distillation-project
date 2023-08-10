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
type timeOs struct {
}

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
