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
		mockNet.EXPECT().ListInterfaces().Return(arg.nets)
		opts := []osservice.Option{osservice.WithNet(mockNet)}

		new(TestServer).With(opts, func() {
			client, err := osservice.NewNetClient(SrvTestHost, SrvTestPort, ClientTimeout)
			req.Nil(err, arg.nets)

			data := client.ListInterfaces()

			req.ElementsMatch(arg.nets, data, arg.name)

			ctrl.Finish()
		})
	}
}
