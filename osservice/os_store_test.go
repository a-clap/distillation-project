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

		mockStore := mocks.NewMockStore(ctrl)
		mockStore.EXPECT().Save(arg.key, arg.value).Return(arg.retErr)
		opts := []osservice.Option{osservice.WithStore(mockStore)}

		new(TestServer).With(opts, func() {
			client := o.storeClient(arg.name)

			err := client.Save(arg.key, arg.value)
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

		mockStore := mocks.NewMockStore(ctrl)
		mockStore.EXPECT().Load(arg.key).Return(arg.value)
		opts := []osservice.Option{osservice.WithStore(mockStore)}

		new(TestServer).With(opts, func() {
			client := o.storeClient(arg.name)

			data := client.Load(arg.key)

			req.EqualValues(arg.value, data, arg.name)

			ctrl.Finish()
		})
	}
}
