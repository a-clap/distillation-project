package main

import (
	"fmt"

	"github.com/a-clap/distillation-ota/pkg/mender/rebooter"
)

func main() {
	fmt.Println(rebooter.Reboot())
}
