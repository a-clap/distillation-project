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

package osservice

import (
	"mender"

	"osservice/osproto"
)

//go:generate protoc --experimental_allow_proto3_optional --go_out=osproto --go_opt=paths=source_relative --go-grpc_out=osproto --go-grpc_opt=paths=source_relative --proto_path osproto time.proto store.proto net.proto wifi.proto update.proto

var (
	toProtoUpdateState = func() func(s mender.DeploymentStatus) osproto.UpdateState {
		toStatus := map[mender.DeploymentStatus]osproto.UpdateState{
			mender.Downloading:           osproto.UpdateState_Downloading,
			mender.PauseBeforeInstalling: osproto.UpdateState_PauseBeforeInstalling,
			mender.Installing:            osproto.UpdateState_Installing,
			mender.PauseBeforeRebooting:  osproto.UpdateState_PauseBeforeRebooting,
			mender.Rebooting:             osproto.UpdateState_Rebooting,
			mender.PauseBeforeCommitting: osproto.UpdateState_PauseBeforeCommitting,
			mender.Success:               osproto.UpdateState_Success,
			mender.Failure:               osproto.UpdateState_Failure,
			mender.AlreadyInstalled:      osproto.UpdateState_AlreadyInstalled,
		}

		return func(s mender.DeploymentStatus) osproto.UpdateState {
			return toStatus[s]
		}
	}()

	fromProtoUpdateState = func() func(s osproto.UpdateState) mender.DeploymentStatus {
		toStatus := map[osproto.UpdateState]mender.DeploymentStatus{
			osproto.UpdateState_Downloading:           mender.Downloading,
			osproto.UpdateState_PauseBeforeInstalling: mender.PauseBeforeInstalling,
			osproto.UpdateState_Installing:            mender.Installing,
			osproto.UpdateState_PauseBeforeRebooting:  mender.PauseBeforeRebooting,
			osproto.UpdateState_Rebooting:             mender.Rebooting,
			osproto.UpdateState_PauseBeforeCommitting: mender.PauseBeforeCommitting,
			osproto.UpdateState_Success:               mender.Success,
			osproto.UpdateState_Failure:               mender.Failure,
			osproto.UpdateState_AlreadyInstalled:      mender.AlreadyInstalled,
		}

		return func(s osproto.UpdateState) mender.DeploymentStatus {
			return toStatus[s]
		}
	}()
)
