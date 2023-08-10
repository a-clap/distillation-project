package osservice

type Option func(os *Os)

func WithTime(t Time) Option {
	return func(o *Os) {
		o.time = t
	}
}

func WithConfigFile(path string) Option {
	return func(o *Os) {
		var err error
		o.store, err = newLoadSaver(path)
		if err != nil {
			panic(err)
		}
	}
}

func WithWifi(wifi Wifi) Option {
	return func(o *Os) {
		o.wifi = wifi
	}
}

func WithStore(store Store) Option {
	return func(o *Os) {
		o.store = store
	}
}

func WithNet(net Net) Option {
	return func(o *Os) {
		o.net = net
	}
}

func WithPort(port int) Option {
	return func(o *Os) {
		o.port = port
	}
}

func WithHost(host string) Option {
	return func(o *Os) {
		o.host = host
	}
}
