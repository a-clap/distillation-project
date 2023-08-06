package main

import (
	"fmt"

	"mender/pkg/rebooter"
)

func main() {
	fmt.Println(rebooter.Reboot())
}
