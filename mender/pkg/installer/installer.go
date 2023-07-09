package installer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

//go:generate mockgen -package mocks -destination mocks/mocks_installer.go . CommandRunner,Runner,ReadCloser

type ReadCloser interface {
	io.ReadCloser
}

type Installer struct {
	CommandRunner
	runner Runner
	finish chan struct{}
}

var (
	findPercentRE = regexp.MustCompile(`\.*(\d+)%`)
	findSizeRE    = regexp.MustCompile(`\.* size (\d+)\s?\.*`)
	findErrorRE   = regexp.MustCompile(`^ERR`)
	ErrNotFound   = errors.New("string doesn't match")
)

func New(options ...Option) *Installer {
	i := &Installer{}

	for _, option := range options {
		option(i)
	}

	if i.CommandRunner == nil {
		i.CommandRunner = &cmdRunner{}
	}

	return i
}

func (i *Installer) Install(artifactName string) (progress chan int, errs chan error, err error) {
	i.runner = i.CommandRunner.Command("mender", "-l", "error", "install", artifactName)

	outPipe, err := i.runner.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get StdoutPipe: %w", err)
	}
	errPipe, err := i.runner.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get StdoutPipe: %w", err)
	}

	return i.handle(outPipe, errPipe)
}

func (i *Installer) Wait() {
	for range i.finish {

	}
}

func (i *Installer) Kill() error {
	if i.runner == nil {
		return nil
	}

	err := i.runner.Kill()
	close(i.finish)

	i.runner = nil

	return err

}
func (i *Installer) handle(outPipe, errPipe io.ReadCloser) (chan int, chan error, error) {
	progress := make(chan int, 100)
	errs := make(chan error, 1)

	i.finish = make(chan struct{})

	if err := i.runner.Start(); err != nil {
		_ = i.Kill()
		return nil, nil, err
	}

	reader := handlePipes(i.finish, outPipe, errPipe)
	go i.handleRunner(progress, errs, reader)

	return progress, errs, nil
}

func (i *Installer) handleRunner(progress chan int, errs chan error, reader chan string) {
	running := true
	go func() {
		_ = i.runner.Wait()
		close(i.finish)
	}()

	for running {
		select {
		case <-i.finish:
			running = false
		case line := <-reader:
			if _, err := findSize(line); err == nil {
				// TODO: what to do with it?
			} else if v, err := findPercent(line); err == nil {
				progress <- v
			} else {
				errs <- fmt.Errorf("unknown line %v", line)
			}
		}
	}
}

func handlePipes(finish chan struct{}, pipes ...io.ReadCloser) chan string {
	reader := make(chan string, len(pipes))

	pipeReader := func(closer io.ReadCloser) {
		scanner := bufio.NewScanner(closer)
		defer closer.Close()

		running := true
		for running && scanner.Scan() {
			reader <- scanner.Text()
			// Check if we should finish
			select {
			case <-finish:
				running = false
				break
			default:
			}
		}
	}

	for _, pipe := range pipes {
		go pipeReader(pipe)
	}

	return reader
}

func findErr(line string) (string, error) {
	if findErrorRE.MatchString(line) {
		return line, nil
	}
	return "", ErrNotFound
}
func findSize(line string) (int, error) {
	match := findSizeRE.FindStringSubmatch(line)
	if match == nil || len(match) != 2 {
		return 0, ErrNotFound
	}

	v, err := strconv.ParseInt(match[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int from %s:%w", match[1], err)
	}

	return int(v), nil
}

func findPercent(line string) (int, error) {
	match := findPercentRE.FindStringSubmatch(line)
	if match == nil || len(match) != 2 {
		return 0, ErrNotFound
	}

	v, err := strconv.ParseInt(match[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int from %s:%w", match[1], err)
	}

	return int(v), nil
}
