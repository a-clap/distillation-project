package loadsaver

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type LoadSaver struct {
	v *viper.Viper
}

func New(configFile string) (*LoadSaver, error) {
	// Create file if it doesn't exist
	if err := touchFile(configFile); err != nil {
		return nil, err
	}

	dir, file := filepath.Split(configFile)
	ext := filepath.Ext(file)

	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(file)

	if len(ext) > 0 {
		// Skip dot
		v.SetConfigType(ext[1:])
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &LoadSaver{v: v}, nil
}

func (l *LoadSaver) Save(key string, data interface{}) error {
	l.v.Set(key, data)
	return l.v.WriteConfig()
}

func (l *LoadSaver) Load(key string) interface{} {
	return l.v.Get(key)
}

func touchFile(configFile string) error {
	file, err := os.OpenFile(configFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print("oj")
		return err
	}

	file.Close()
	return nil
}
