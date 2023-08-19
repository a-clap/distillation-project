package osservice

//go:generate mockgen -destination=mocks/mocks.mockgen.go -package=mocks . Time,Store,Net,Wifi,Update,UpdateCallbacks
