package device_test

import (
	"os"
	"path"
	"testing"

	"github.com/a-clap/distillation-ota/pkg/mender/device"
	"github.com/stretchr/testify/suite"
)

type DeviceTestSuite struct {
	suite.Suite
}

func TestDevice(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}

func (d *DeviceTestSuite) TestInventory() {
	cwd, _ := os.Getwd()

	req := d.Require()
	args := []struct {
		name    string
		options []device.Option
		attrs   []device.Attribute
		err     error
	}{
		{
			name: "basic",
			options: []device.Option{
				device.WithPath(path.Join(cwd, "testdata")),
				device.WithInventoryDir("inventory"),
			},
			attrs: []device.Attribute{
				{
					Name:  "single",
					Value: []string{"value"},
				},
				{
					Name:  "multiple",
					Value: []string{"value1", "value2"},
				},
				{
					Name:  "another",
					Value: []string{"file"},
				},
			},
		},
	}
	for _, arg := range args {
		dev := device.New(arg.options...)

		attrs, err := dev.Attributes()
		if arg.err != nil {
			req.Nil(attrs, arg.name)
			req.NotNil(err, arg.name)
			continue
		}

		req.Nil(err, arg.name)
		req.NotNil(attrs, arg.name)
		req.ElementsMatch(arg.attrs, attrs, arg.name)

	}

}
