package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/a-clap/distillation/pkg/distillation"
)

const addr = "localhost:50002"

func testds() {
	client, err := distillation.NewDSRPCClient(addr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	
	// Contact the server and print out its response.
	for {
		<-time.After(1 * time.Second)
		r, err := client.Get()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, elem := range r {
			log.Println(elem)
			elem.Enabled = true
			_, err := client.Configure(elem)
			if err != nil {
				log.Println(err)
			}
		}
		
		t, err := client.Temperatures()
		log.Println(t, err)
		
	}
	
}

func testGpio() {
	client, err := distillation.NewGPIORPCClient(addr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	
	// Contact the server and print out its response.
	for {
		<-time.After(1 * time.Second)
		r, err := client.Get()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, elem := range r {
			fmt.Println("ID: ", elem.ID)
			fmt.Println("Dir: ", elem.Direction)
			fmt.Println("Active: ", elem.ActiveLevel)
			fmt.Println("Value: ", elem.Value)
		}
		n := r[0]
		n.Value = !n.Value
		
		c, err := client.Configure(n)
		log.Println(c, err)
		
	}
	
}
func testHeaters() {
	client, err := distillation.NewHeaterRPCCLient(addr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	
	// Contact the server and print out its response.
	for {
		<-time.After(1 * time.Second)
		r, err := client.Get()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, elem := range r {
			fmt.Printf("ID: %v, Enabled: %v\n", elem.ID, elem.Enabled)
		}
		n := r[0]
		n.Enabled = !n.Enabled
		
		c, err := client.Configure(n)
		log.Println(c, err)
		
	}
	
}
func testpt() {
	client, err := distillation.NewPTRPCClient(addr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	
	// Contact the server and print out its response.
	for {
		<-time.After(1 * time.Second)
		r, err := client.Get()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, elem := range r {
			log.Println(elem)
			elem.Enabled = true
			_, err := client.Configure(elem)
			if err != nil {
				log.Println(err)
			}
		}
		
		t, err := client.Temperatures()
		log.Println(t, err)
		
	}
	
}

func main() {
	go testpt()
	// testds()
	// testGpio()
	// testHeaters()
	for {
	
	}
}
