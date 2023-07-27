package main

import (
	"fmt"
	"log"

	"github.com/a-clap/distillation-ota/pkg/mender/loadsaver"
)

type Release struct {
	Name  string `mapstructure:"name"`
	State string `mapstructure:"state"`
}

type Releases struct {
	Values []Release `mapstructure:"releases"`
}

func main() {
	ls, err := loadsaver.New("/home/adamclap/repo/.go/ota/config13.yaml")
	if err != nil {
		log.Fatal(err)
	}

	r := Release{
		Name:  "1",
		State: "2",
	}
	err = ls.Save("relesae", r)
	fmt.Println(err)
}
