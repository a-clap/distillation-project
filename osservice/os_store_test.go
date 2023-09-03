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

package osservice_test

import (
	"osservice"
	"osservice/mocks"

	"github.com/golang/mock/gomock"
)

func (o *OsServiceSuite) storeClient(msg string) *osservice.StoreClient {
	client, err := osservice.NewStoreClient(SrvTestHost, SrvTestPort, ClientTimeout)
	o.Require().Nil(err, msg)
	o.Require().NotNil(client, msg)
	return client
}

func (o *OsServiceSuite) TestTime_Save() {
	args := []struct {
		name   string
		key    string
		value  []byte
		err    error
		retErr error
	}{
		{
			name:   "bytes",
			key:    "key",
			value:  []byte("my message"),
			retErr: nil,
			err:    nil,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		waitForSave := make(chan struct{})
		mockStore := mocks.NewMockStore(ctrl)
		mockStore.EXPECT().Save(arg.key, arg.value).DoAndReturn(func(any, any) error {
			defer close(waitForSave)
			return arg.err
		})
		opts := []osservice.Option{osservice.WithStore(mockStore)}

		new(TestServer).With(opts, func() {
			client := o.storeClient(arg.name)

			err := client.Save(arg.key, arg.value)
			o.waitFor(waitForSave, "waitForSave")

			if arg.err != nil {
				req.NotNil(err, arg.name)
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
			}

			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestTime_Load() {
	args := []struct {
		name  string
		key   string
		value []byte
		err   error
	}{
		{
			name:  "bytes",
			key:   "key",
			value: []byte("my message"),
			err:   nil,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		waitForLoad := make(chan struct{})
		mockStore := mocks.NewMockStore(ctrl)
		mockStore.EXPECT().Load(arg.key).DoAndReturn(func(any) []byte {
			defer close(waitForLoad)
			return arg.value
		})

		opts := []osservice.Option{osservice.WithStore(mockStore)}

		new(TestServer).With(opts, func() {
			client := o.storeClient(arg.name)

			data := client.Load(arg.key)

			req.EqualValues(arg.value, data, arg.name)

			ctrl.Finish()
		})
	}
}
