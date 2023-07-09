package installer_test

import (
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a-clap/distillation-ota/pkg/mender/installer"
	"github.com/a-clap/distillation-ota/pkg/mender/installer/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type InstallerSuite struct {
	suite.Suite
}

func TestInstaller(t *testing.T) {
	suite.Run(t, new(InstallerSuite))
}

var (
	sizeMsg     = `Installing Artifact of size 162508800...`
	percentsMsg = []string{
		"", // empty msg, to match idx with %
		"                                                                    1%\n",
		".                                                                   2%\n",
		".                                                                   3%\n",
		"..                                                                  4%\n",
		"...                                                                 5%\n",
		"...                                                                 6%\n",
		"....                                                                7%\n",
		".....                                                               8%\n",
		".....                                                               9%\n",
		"......                                                             10%\n",
		".......                                                            11%\n",
		".......                                                            12%\n",
		"........                                                           13%\n",
		".........                                                          14%\n",
		".........                                                          15%\n",
		"..........                                                         16%\n",
		"...........                                                        17%\n",
		"...........                                                        18%\n",
		"............                                                       19%\n",
		".............                                                      20%\n",
		".............                                                      21%\n",
		"..............                                                     22%\n",
		"..............                                                     23%\n",
		"...............                                                    24%\n",
		"................                                                   25%\n",
		"................                                                   26%\n",
		".................                                                  27%\n",
		"..................                                                 28%\n",
		"..................                                                 29%\n",
		"...................                                                30%\n",
		"....................                                               31%\n",
		"....................                                               32%\n",
		".....................                                              33%\n",
		"......................                                             34%\n",
		"......................                                             35%\n",
		".......................                                            36%\n",
		"........................                                           37%\n",
		"........................                                           38%\n",
		".........................                                          39%\n",
		"..........................                                         40%\n",
		"..........................                                         41%\n",
		"...........................                                        42%\n",
		"...........................                                        43%\n",
		"............................                                       44%\n",
		".............................                                      45%\n",
		".............................                                      46%\n",
		"..............................                                     47%\n",
		"...............................                                    48%\n",
		"...............................                                    49%\n",
		"................................                                   50%\n",
		".................................                                  51%\n",
		".................................                                  52%\n",
		"..................................                                 53%\n",
		"...................................                                54%\n",
		"...................................                                55%\n",
		"....................................                               56%\n",
		".....................................                              57%\n",
		".....................................                              58%\n",
		"......................................                             59%\n",
		".......................................                            60%\n",
		".......................................                            61%\n",
		"........................................                           62%\n",
		"........................................                           63%\n",
		".........................................                          64%\n",
		"..........................................                         65%\n",
		"..........................................                         66%\n",
		"...........................................                        67%\n",
		"............................................                       68%\n",
		"............................................                       69%\n",
		".............................................                      70%\n",
		"..............................................                     71%\n",
		"..............................................                     72%\n",
		"...............................................                    73%\n",
		"................................................                   74%\n",
		"................................................                   75%\n",
		".................................................                  76%\n",
		"..................................................                 77%\n",
		"..................................................                 78%\n",
		"...................................................                79%\n",
		"....................................................               80%\n",
		"....................................................               81%\n",
		".....................................................              82%\n",
		".....................................................              83%\n",
		"......................................................             84%\n",
		".......................................................            85%\n",
		".......................................................            86%\n",
		"........................................................           87%\n",
		".........................................................          88%\n",
		".........................................................          89%\n",
		"..........................................................         90%\n",
		"...........................................................        91%\n",
		"...........................................................        92%\n",
		"............................................................       93%\n",
		".............................................................      94%\n",
		".............................................................      95%\n",
		"..............................................................     96%\n",
		"...............................................................    97%\n",
		"...............................................................    98%\n",
		"................................................................   99%\n",
		"................................................................. 100%\n",
	}
)

func (i *InstallerSuite) TestInstallInitialErrors() {

	i.Run("stdOutPipe err", func() {
		r := i.Require()
		ctrl := gomock.NewController(i.T())
		defer ctrl.Finish()

		runner := mocks.NewMockRunner(ctrl)

		pipeErr := errors.New("pipe err")
		runner.EXPECT().StdoutPipe().Return(nil, pipeErr)
		runner.EXPECT().StderrPipe().Return(nil, nil).AnyTimes()

		commandRunner := mocks.NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().Command(gomock.Any(), gomock.Any()).Return(runner)

		// Run actual test
		ins := installer.New(installer.WithCommandRunner(commandRunner))
		r.NotNil(ins)

		progress, errs, err := ins.Install("hello")
		r.Nil(progress)
		r.Nil(errs)
		r.ErrorIs(err, pipeErr)
	})

	i.Run("stdErrPipe err", func() {
		r := i.Require()
		ctrl := gomock.NewController(i.T())
		defer ctrl.Finish()

		runner := mocks.NewMockRunner(ctrl)

		pipeErr := errors.New("stdout err")
		runner.EXPECT().StdoutPipe().Return(nil, nil).AnyTimes()
		runner.EXPECT().StderrPipe().Return(nil, pipeErr)

		commandRunner := mocks.NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().Command(gomock.Any(), gomock.Any()).Return(runner)

		// Run actual test
		ins := installer.New(installer.WithCommandRunner(commandRunner))
		r.NotNil(ins)

		progress, errs, err := ins.Install("hello")
		r.Nil(progress)
		r.Nil(errs)
		r.ErrorIs(err, pipeErr)
	})

	i.Run("start err", func() {
		r := i.Require()
		ctrl := gomock.NewController(i.T())
		defer ctrl.Finish()
		errStart := errors.New("start err")
		// Empty pipes
		stdoutPipe := mocks.NewMockReadCloser(ctrl)
		stdoutPipe.EXPECT().Read(gomock.Any()).Return(0, io.EOF).AnyTimes()

		stderrPipe := mocks.NewMockReadCloser(ctrl)
		stderrPipe.EXPECT().Read(gomock.Any()).Return(0, io.EOF).AnyTimes()

		runner := mocks.NewMockRunner(ctrl)

		runner.EXPECT().StdoutPipe().Return(stdoutPipe, nil)
		runner.EXPECT().StderrPipe().Return(stderrPipe, nil)
		runner.EXPECT().Start().Return(errStart)
		runner.EXPECT().Kill().Return(nil)

		commandRunner := mocks.NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().Command(gomock.Any(), gomock.Any()).Return(runner)

		// Run actual test
		ins := installer.New(installer.WithCommandRunner(commandRunner))
		r.NotNil(ins)

		progress, errs, err := ins.Install("hello")
		r.Nil(progress)
		r.Nil(errs)
		r.ErrorIs(err, errStart)
	})

	i.Run("err line", func() {
		r := i.Require()
		ctrl := gomock.NewController(i.T())
		defer ctrl.Finish()

		// Empty pipes
		stdoutPipe := mocks.NewMockReadCloser(ctrl)
		stdoutPipe.EXPECT().Read(gomock.Any()).Return(0, io.EOF)
		stdoutPipe.EXPECT().Close().Return(nil)

		errLine := `ERRO[0005] Error while installing Artifact from command line\n`
		stderrPipe := mocks.NewMockReadCloser(ctrl)

		callErrLine := stderrPipe.EXPECT().Read(gomock.Any()).SetArg(0, []byte(errLine)).Return(len(errLine), nil)
		stderrPipe.EXPECT().Read(gomock.Any()).Return(0, io.EOF).After(callErrLine)
		stderrPipe.EXPECT().Close().Return(nil)

		runner := mocks.NewMockRunner(ctrl)

		runner.EXPECT().StdoutPipe().Return(stdoutPipe, nil)
		runner.EXPECT().StderrPipe().Return(stderrPipe, nil)
		runner.EXPECT().Start().Return(nil)

		wait := make(chan struct{})
		runner.EXPECT().Wait().Return(nil).DoAndReturn(func() error {
			for range wait {
			}
			return nil
		})

		// runner.EXPECT().Kill().Return(nil).AnyTimes()
		// runner.EXPECT().Wait().AnyTimes()

		commandRunner := mocks.NewMockCommandRunner(ctrl)
		commandRunner.EXPECT().Command(gomock.Any(), gomock.Any()).Return(runner)

		// Run actual test
		ins := installer.New(installer.WithCommandRunner(commandRunner))
		r.NotNil(ins)

		progress, errs, err := ins.Install("hello")
		r.NotNil(progress)
		r.NotNil(errs)
		r.Nil(err)

		var running atomic.Bool
		running.Store(true)
		for running.Load() {
			select {
			case <-progress:
				r.FailNow("Shouldn't receive progress")
			case <-time.After(10 * time.Millisecond):
				r.Fail("Shouldn't be here")
			case err = <-errs:
				close(wait)
				running.Store(false)
			}
		}
		<-wait
		// We should receive error with that line
		r.ErrorContains(err, errLine)

	})
}
func (i *InstallerSuite) TestInstallSuccess() {
	r := i.Require()
	ctrl := gomock.NewController(i.T())
	defer ctrl.Finish()

	commandRunner := mocks.NewMockCommandRunner(ctrl)

	stdoutPipe := mocks.NewMockReadCloser(ctrl)
	var call *gomock.Call
	for _, data := range percentsMsg[1:] {
		if call == nil {
			// First call
			call = stdoutPipe.EXPECT().Read(gomock.Any()).SetArg(0, []byte(data)).Return(len(data), nil)
		} else {
			// Make them ordered
			call = stdoutPipe.EXPECT().Read(gomock.Any()).SetArg(0, []byte(data)).Return(len(data), nil).After(call)
		}
	}
	// Make sure that is the last call of read
	stdoutPipe.EXPECT().Read(gomock.Any()).Return(0, io.EOF).After(call)
	stdoutPipe.EXPECT().Close().Return(nil)

	stderrPipe := mocks.NewMockReadCloser(ctrl)
	// errPipe doesn't return anything
	stderrPipe.EXPECT().Read(gomock.Any()).Return(0, io.EOF)
	stderrPipe.EXPECT().Close().Return(nil)

	// should call commandRunner.Command()
	runner := mocks.NewMockRunner(ctrl)

	runner.EXPECT().StdoutPipe().Return(stdoutPipe, nil)
	runner.EXPECT().StderrPipe().Return(stderrPipe, nil)

	wait := make(chan struct{})

	runner.EXPECT().Start().Return(nil)
	runner.EXPECT().Wait().Return(nil).DoAndReturn(func() error {
		for range wait {
		}
		return nil
	})

	commandRunner.EXPECT().Command(gomock.Any(), gomock.Any()).Return(runner)

	// Run actual test
	ins := installer.New(installer.WithCommandRunner(commandRunner))
	r.NotNil(ins)

	progress, errs, err := ins.Install("install")

	r.Nil(err)
	r.NotNil(progress)
	r.NotNil(errs)

	responses := make([]int, 0, 100)
	var running atomic.Bool
	running.Store(true)
	for running.Load() {
		select {
		case resp := <-progress:
			responses = append(responses, resp)
			if resp == 100 {
				// Sounds like success
				close(wait)
				running.Store(false)
			}
		case <-time.After(10 * time.Millisecond):
			r.Fail("Shouldn't be here")
		case err := <-errs:
			r.Fail("unexpected err", err)
		}
	}
	// No need to check anything
}
