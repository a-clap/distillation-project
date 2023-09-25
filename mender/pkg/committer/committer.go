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

package committer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"

	"go.uber.org/atomic"
)

//go:generate mockgen -package mocks -destination mocks/mocks_installer.go . CommandRunner,Runner,ReadCloser

type ReadCloser interface {
	io.ReadCloser
}

type Committer struct {
	cmd    CommandRunner
	runner Runner
	finish chan struct{}
}

var (
	findErrorRE = regexp.MustCompile(`^ERR`)
	ErrFound    = errors.New("errors match")
)

func New(options ...Option) *Committer {
	i := &Committer{}

	for _, option := range options {
		option(i)
	}

	if i.cmd == nil {
		i.cmd = &cmdRunner{}
	}

	return i
}

func (c *Committer) Commit() error {
	c.runner = c.cmd.Command("mender", "-l", "error", "commit")

	outPipe, err := c.runner.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdoutPipe: %w", err)
	}

	errPipe, err := c.runner.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderrPipe: %w", err)
	}

	return c.handle(outPipe, errPipe)
}

func (c *Committer) kill() error {
	if c.runner == nil {
		return nil
	}

	err := c.runner.Kill()

	c.runner = nil

	return err
}

func (c *Committer) handle(outPipe, errPipe io.ReadCloser) error {
	if err := c.runner.Start(); err != nil {
		return err
	}

	c.finish = make(chan struct{})
	reader := handlePipes(c.finish, outPipe, errPipe)

	go func() {
		_ = c.runner.Wait()
		close(c.finish)
	}()

	var (
		err       error
		errString string
	)

	running := atomic.NewBool(true)
	for running.Load() {
		select {
		case <-c.finish:
			running.Store(false)
		case line := <-reader:
			errString, err = findErr(line)
			if errors.Is(err, ErrFound) {
				err = fmt.Errorf("commit :%w", errors.New(errString))
				running.Store(false)
			}
		}
	}

	// In case of error
	if err != nil {
		_ = c.kill()
	}

	return err
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
		return line, ErrFound
	}

	return "", nil
}
