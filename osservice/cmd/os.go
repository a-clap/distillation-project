package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"osservice"
)

var (
	host       = flag.String("host", "", "the server host")
	port       = flag.Int("port", 50003, "the server port")
	configFile = flag.String("config", "./config.yaml", "config file location")
)

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	opts := []osservice.Option{
		osservice.WithHost(*host),
		osservice.WithPort(*port),
		osservice.WithConfigFile(*configFile),
	}

	osSrv, err := osservice.New(opts...)
	if err != nil {
		log.Fatal(err)
	}

	wait := make(chan struct{})
	go func() {
		log.Printf(`Running server on %v:%v`, *host, *port)
		err := osSrv.Run()
		if err != nil {
			log.Println("Server failed with ", err)
		}
		log.Println("Server stopped")
		close(wait)
	}()

	handledSignals := []os.Signal{
		os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGPIPE, syscall.SIGTERM,
	}

	sigs := make(chan os.Signal, len(handledSignals))
	signal.Notify(sigs, handledSignals...)
	// For any handled signal, just do cleanup
	<-sigs
	// Stop server gracefully
	osSrv.Stop()
	// Wait until closed
	<-wait
}
