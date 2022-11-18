// Package config_new defines new config handling implementation
package config_new

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config represents the structure that contains configurations
type Config struct {
	bM       *sync.Mutex
	aM       *sync.Mutex
	Basic    map[string]Entry `yaml:"basic"`
	Advanced map[string]Entry `yaml:"advanced"`
}

// NewConfig creates new Config
func NewConfig() *Config {
	return &Config{
		bM:       &sync.Mutex{},
		aM:       &sync.Mutex{},
		Basic:    make(map[string]Entry),
		Advanced: make(map[string]Entry),
	}
}

// GetBasic fetches config entry from basic config
func (c *Config) GetBasic(key string) Entry {
	c.bM.Lock()
	defer c.bM.Unlock()

	tmp := c.Basic[key]
	return tmp
}

// Get fetches config entry from advanced config
func (c *Config) Get(key string) Entry {
	c.aM.Lock()
	defer c.aM.Unlock()

	tmp := c.Advanced[key]
	return tmp
}

// Set sets new config from given file path
func (c *Config) Set(filePath string) error {
	bytes, err := loadFileData(filePath)
	if err != nil {
		return err
	}

	conf, err := unmarshalConfig(bytes)
	if err != nil {
		return err
	}

	c.bM.Lock()
	c.Basic = conf.Basic
	c.bM.Unlock()

	c.aM.Lock()
	c.Advanced = conf.Advanced
	c.aM.Unlock()

	return nil
}

// unmarshalConfig unmatshalls YAML file bytes into Config
func unmarshalConfig(configData []byte) (*Config, error) {
	conf := &Config{}
	err := yaml.Unmarshal(configData, conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}
	return conf, nil
}

// loadFileData loads YAML config data from file loader int bytes
func loadFileData(configFile string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Clean(configFile))
}
