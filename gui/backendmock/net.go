package backendmock

import (
	"osservice"
)

type NetMock struct{}

var _ osservice.Net = (*NetMock)(nil)

func (n NetMock) ListInterfaces() []osservice.NetInterface {
	return nil
}
