package osservice

import (
	"osservice/pkg/wifi"
)

var _ Wifi = (*wifiOs)(nil)

type wifiOs struct {
	conn *wifi.Wifi
}

func newWifiOs() *wifiOs {
	conn, err := wifi.New()
	if err != nil {
		return nil
	}
	return &wifiOs{conn: conn}
}

func (w *wifiOs) APs() ([]wifi.AP, error) {
	return w.conn.APs()
}

func (w *wifiOs) Connected() (wifi.Status, error) {
	return w.conn.Connected()
}

func (w *wifiOs) Disconnect() error {
	return w.conn.Disconnect()
}

func (w *wifiOs) Connect(n wifi.Network) error {
	return w.conn.Connect(n)
}
