package config

import (
	"errors"
	"fmt"

	"github.com/k81/kate/utils"
	"gopkg.in/ini.v1"
)

type iniConfiger struct {
	*ini.File
}

// Load load all configs
func (c *iniConfiger) Load(file string) error {
	var err error
	if c.File, err = ini.Load(file); err != nil {
		return fmt.Errorf("load config: %v", err)
	}

	configs := []Config{
		Profiling,
		MySQL,
		Redis,
		HTTP,
	}

	for _, config := range configs {
		section := c.File.Section(config.SectionName())
		if err = config.Load(section); err != nil {
			return fmt.Errorf("load config: section=%v, error=%v", config.SectionName(), err)
		}
	}
	return nil
}

func (c *iniConfiger) Get(name string) (string, error) {
	key, err := c.getKey(name)
	if err != nil {
		return "", err
	}

	return key.String(), nil
}

func (c *iniConfiger) MustGet(name string, defaultValue string) string {
	key, err := c.getKey(name)
	if err != nil {
		return defaultValue
	}

	return key.MustString(defaultValue)
}

func (c *iniConfiger) getKey(name string) (key *ini.Key, err error) {
	parts := utils.Split(name, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid key: %v, should have format: section.key", name)
	}

	if c.File == nil {
		return nil, errors.New("config not loaded")
	}
	return c.Section(parts[0]).Key(parts[1]), nil
}
