package device

type Option func(d *Device)

func WithPath(path string) Option {
	return func(d *Device) {
		d.path = path
	}
}

func WithInventoryDir(invDir string) Option {
	return func(d *Device) {
		d.inventoryDir = invDir
	}
}

func WithIdentityDir(identityDir string) Option {
	return func(d *Device) {
		d.identityDir = identityDir
	}
}

func WithInfoDir(infoDir string) Option {
	return func(d *Device) {
		d.infoDir = infoDir
	}
}
