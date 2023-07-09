package mender

import (
	"fmt"
	"strings"
)

func ProgressBarUpdate(v int) {
	if v == 0 {
		fmt.Printf("\r[>%s] %v %%", strings.Repeat(" ", 100), v)
	} else if v == 100 {
		fmt.Printf("\r[%s] %v %%", strings.Repeat("=", 100), v)
	} else {
		fmt.Printf("\r[%s>%s] %v %%", strings.Repeat("=", v), strings.Repeat(" ", 100-v), v)
	}
}
