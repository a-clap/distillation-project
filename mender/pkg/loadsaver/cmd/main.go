package main

import (
	"log"
	"os"
	"path"

	"github.com/a-clap/distillation-ota/pkg/mender/loadsaver"
)

type Release struct {
	Name  string `mapstructure:"name"`
	State string `mapstructure:"state"`
}

func main() {
	p, err := os.Getwd()
	if err != nil {
		log.Fatal(p)
	}

	ls, err := loadsaver.New(path.Join(p, "config.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	r := Release{
		Name:  "1",
		State: "2",
	}

	if err = ls.Save("release", r); err != nil {
		log.Fatal(err)
	}
}
