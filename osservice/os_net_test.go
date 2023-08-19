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

func (o *OsServiceSuite) TestNet_ListInterfaces() {
	args := []struct {
		name string
		nets []osservice.NetInterface
		err  error
	}{
		{
			name: "bytes",
			nets: []osservice.NetInterface{
				{
					Name:     "1",
					IPAddrV4: "2",
				},
				{
					Name:     "3",
					IPAddrV4: "4",
				},
			},
			err: nil,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockNet := mocks.NewMockNet(ctrl)
		waitInterfaces := make(chan struct{})
		mockNet.EXPECT().ListInterfaces().DoAndReturn(func() []osservice.NetInterface {
			defer close(waitInterfaces)
			return arg.nets
		})
		opts := []osservice.Option{osservice.WithNet(mockNet)}

		new(TestServer).With(opts, func() {
			client, err := osservice.NewNetClient(SrvTestHost, SrvTestPort, ClientTimeout)
			req.Nil(err, arg.nets)

			data := client.ListInterfaces()
			o.waitFor(waitInterfaces, "waitInterfaces")
			req.ElementsMatch(arg.nets, data, arg.name)

			ctrl.Finish()
		})
	}
}
