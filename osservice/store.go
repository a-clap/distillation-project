package osservice

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	_ Store = (*configStoreOs)(nil)
)

type configStoreOs struct {
	v *viper.Viper
}

func newLoadSaver(configFile string) (*configStoreOs, error) {
	dir, file := filepath.Split(configFile)
	ext := filepath.Ext(file)

	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(strings.TrimSuffix(file, ext))

	if len(ext) > 0 {
		// Skip dot
		v.SetConfigType(ext[1:])
	}

	if err := v.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		if err = v.SafeWriteConfig(); err != nil {
			return nil, fmt.Errorf("failed to create config: %w", err)
		}
	}

	return &configStoreOs{v: v}, nil
}

func (l *configStoreOs) Save(key string, data []byte) error {
	l.v.Set(key, data)
	return l.v.WriteConfig()
}

func (l *configStoreOs) Load(key string) []byte {
	return l.v.Get(key).([]byte)
}
