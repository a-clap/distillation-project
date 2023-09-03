package backend

import (
	"gui/backend/wifi"

	"github.com/a-clap/logging"
)

type WifiConnected struct {
	Connected bool   `json:"connected"`
	AP        string `json:"AP"`
}

func (b *Backend) WifiAPList() []string {
	aps, err := wifi.AP()
	if err != nil {
		logger.Error("WifiAPList", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrWIFIAPList)

		return nil
	}

	return aps
}

func (b *Backend) WifiIsConnected() WifiConnected {
	conn, ap, err := wifi.Connected()
	if err != nil {
		logger.Error("WifiIsConnected", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrWifiIsConnected)
		return WifiConnected{}
	}

	return WifiConnected{
		Connected: conn,
		AP:        ap,
	}
}

func (b *Backend) WifiConnect(ap, psk string) {
	_ = wifi.Disconnect()

	if err := wifi.Connect(ap, psk); err != nil {
		b.eventEmitter.OnError(ErrWifiConnect)
	}
}
