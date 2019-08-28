package config

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

func Select(configpaths ...string) (string, error) {
	if len(configpaths) == 0 {
		return "", errors.New("nil config files")
	}

	for _, configpath := range configpaths {
		if _, err := os.Stat(configpath); err == nil {
			return configpath, nil
		}
	}

	return "", errors.New("no such file in the root directory")
}

func Load(config *ProblemConfig, configpath string) error {
	for _, configX := range []struct {
		Type      string
		Unmarshal func([]byte, interface{}) error
	}{
		{".json", json.Unmarshal}, {".yml", yaml.Unmarshal},
		{".toml", toml.Unmarshal}, {".xml", xml.Unmarshal}} {
		if _, err := os.Stat(configpath + configX.Type); err == nil {
			f, err := os.Open(configpath + configX.Type)
			if err != nil {
				return err
			}

			b, err := ioutil.ReadAll(f)
			f.Close()
			if err != nil {
				return err
			}
			err = configX.Unmarshal(b, config)
			if err != nil {
				return err
			}
			config.LoadType = configX.Type
			return nil
		}
	}

	return errors.New("no such file in the root directory")
}

func Save(config *ProblemConfig, configpath string) error {
	var b []byte
	var err error
	switch config.LoadType {
	case ".json":
		b, err = json.MarshalIndent(config, "", "    ")
		if err != nil {
			return err
		}
	case ".yml":
		b, err = yaml.Marshal(config)
		if err != nil {
			return err
		}
	case ".toml":
		b, err = toml.Marshal(config)
		if err != nil {
			return err
		}
	case ".xml":
		b, err = xml.MarshalIndent(config, "", "    ")
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(configpath + config.LoadType); err == nil {
		f, err := os.OpenFile(configpath+config.LoadType, os.O_WRONLY|os.O_TRUNC, 0333)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		f.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("no such file in the root directory")
}
