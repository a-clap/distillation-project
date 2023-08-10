package osservice_test

import (
	"os"

	"osservice"
	"osservice/mocks"
	"osservice/pkg/wifi"

	"github.com/golang/mock/gomock"
)

func (o *OsServiceSuite) wifiClient(msg string) *osservice.WifiClient {
	client, err := osservice.NewWifiClient(SrvTestHost, SrvTestPort, ClientTimeout)
	o.Require().Nil(err, msg)
	o.Require().NotNil(client, msg)
	return client

}

func (o *OsServiceSuite) TestWifi_APs() {
	args := []struct {
		name string
		aps  []wifi.AP
		err  error
	}{
		{
			name: "good",
			aps: []wifi.AP{
				{
					ID:   1,
					SSID: "123",
				},
				{
					ID:   2,
					SSID: "455",
				},
			},
			err: nil,
		},
		{
			name: "nil aps",
			aps:  nil,
			err:  nil,
		},
		{
			name: "empty aps",
			aps:  []wifi.AP{},
			err:  nil,
		},
		{
			name: "err",
			err:  os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)
		mockWifi.EXPECT().APs().Return(arg.aps, arg.err)
		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {
			client := o.wifiClient(arg.name)
			aps, err := client.APs()
			if arg.err != nil {
				req.NotNil(err, arg.name)
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
				req.ElementsMatch(arg.aps, aps, err)
			}

			ctrl.Finish()
		})
	}
}
func (o *OsServiceSuite) TestWifi_Disconnect() {
	args := []struct {
		name string
		err  error
	}{
		{
			name: "good",
			err:  nil,
		},
		{
			name: "err",
			err:  os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)
		mockWifi.EXPECT().Disconnect().Return(arg.err)
		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {
			client := o.wifiClient(arg.name)
			err := client.Disconnect()
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
func (o *OsServiceSuite) TestWifi_Connect() {
	args := []struct {
		name string
		net  wifi.Network
		err  error
	}{
		{
			name: "good",
			net: wifi.Network{
				AP: wifi.AP{
					ID:   1,
					SSID: "123",
				},
				Password: "456",
			},
			err: nil,
		},
		{
			name: "err",
			err:  os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)
		mockWifi.EXPECT().Connect(arg.net).Return(arg.err)
		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {
			client := o.wifiClient(arg.name)
			err := client.Connect(arg.net)
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
func (o *OsServiceSuite) TestWifi_Connected() {
	args := []struct {
		name      string
		retStatus wifi.Status
		err       error
	}{
		{
			name: "good",
			retStatus: wifi.Status{
				Connected: true,
				SSID:      "my ap",
			},
			err: nil,
		},
		{
			name: "good #2",
			retStatus: wifi.Status{
				Connected: false,
				SSID:      "",
			},
			err: nil,
		},
		{
			name: "err",
			retStatus: wifi.Status{
				Connected: false,
				SSID:      "",
			},
			err: os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)
		mockWifi.EXPECT().Connected().Return(arg.retStatus, arg.err)
		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {

			client := o.wifiClient(arg.name)
			status, err := client.Connected()
			if arg.err != nil {
				req.NotNil(err, arg.name)
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
				req.EqualValues(arg.retStatus, status, arg.name)
			}

			ctrl.Finish()
		})
	}
}
