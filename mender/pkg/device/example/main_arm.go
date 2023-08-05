package main

import (
	"log"

	"github.com/a-clap/distillation-ota/pkg/mender/device"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dev := device.New()

	if info, err := dev.Info(); err != nil {
		log.Println("Failed to read info:", err)
	} else {
		log.Printf("Got info: %#v", info)
	}

	if attributes, err := dev.Attributes(); err != nil {
		log.Println("Failed to read attributes:", err)
	} else {
		log.Printf("Got attributes: %#v", attributes)
	}

	if id, err := dev.ID(); err != nil {
		log.Println("Failed to read id:", err)
	} else {
		log.Printf("Got id: %#v", id)
	}
}
