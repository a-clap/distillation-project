// MIT License
//
// Copyright (c) 2023 a-clap
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"mender"
	"mender/pkg/signer"

	"osservice"

	"golang.org/x/exp/slog"
)

var (
	port         = flag.Int("port", 50003, "the server port")
	configDir    = flag.String("config", ".", "config file location")
	menderServer = flag.String("mender", "https://eu.hosted.mender.io/", "mender server url")
	menderToken  = flag.String("token", "", "mender server token")
	PEMFile      = flag.String("pem", "", "pem file")
)

func main() {
	flag.Parse()

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})

	slog.SetDefault(slog.New(textHandler))

	pem, err := os.ReadFile(*PEMFile)
	if err != nil {
		log.Fatalln(err)
	}

	menderStore := path.Join(*configDir, "mender.yaml")
	menderSigner, err := signer.New(signer.WithPrivKey(pem))
	if err != nil {
		log.Fatalln(err)
	}

	client, err := mender.NewBuilder().
		WithTimeout(5*time.Second).
		WithSignerVerifier(menderSigner).
		WithServer(*menderServer, *menderToken).
		WithStore(menderStore).
		WithStdIOInterface().
		Build()

	if err != nil {
		log.Fatal(err)
	}

	configFile := path.Join(*configDir, "osservice.yaml")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	opts := []osservice.Option{
		osservice.WithPort(*port),
		osservice.WithConfigFile(configFile),
		osservice.WithMender(client),
	}

	osSrv, err := osservice.New(opts...)
	if err != nil {
		log.Fatal(err)
	}

	wait := make(chan struct{})
	go func() {
		log.Printf(`Running server on :%v`, *port)
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
