package backend

import (
	"github.com/a-clap/logging"
	"gui/backend/wifi"
)

type WifiConnected struct {
	Connected bool   `json:"connected"`
	AP        string `json:"AP"`
}

func (b *Backend) WifiAPList() []string {
	logger.Debug("WifiAPList")

	aps, err := wifi.AP()

	if err != nil {
		logger.Error("WifiAPList", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrWIFIAPList)

		return nil
	}

	return aps
}

func (b *Backend) WifiIsConnected() WifiConnected {
	logger.Debug("WifiIsConnected")

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
	logger.Debug("WifiConnect", logging.String("ap", ap))
	_ = wifi.Disconnect()

	if err := wifi.Connect(ap, psk); err != nil {
		b.eventEmitter.OnError(ErrWifiConnect)
	}
}
