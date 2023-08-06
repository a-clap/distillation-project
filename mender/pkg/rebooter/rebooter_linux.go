package rebooter

import (
	"golang.org/x/sys/unix"
)

func Reboot() error {
	unix.Sync()
	return unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
}
