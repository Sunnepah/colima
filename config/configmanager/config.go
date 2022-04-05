package configmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/abiosoft/colima/cli"
	"github.com/abiosoft/colima/config"
	"github.com/abiosoft/colima/util/yamlutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Save saves the config.
func Save(c config.Config) error {
	return yamlutil.Save(c, config.File())
}

// oldConfigFile returns the path to config file of versions <0.4.0.
// TODO: remove later, only for backward compatibility
func oldConfigFile() string {
	_, configFileName := filepath.Split(config.File())
	return filepath.Join(os.Getenv("HOME"), "."+config.Profile().ID, configFileName)
}

// Load loads the config.
// Error is only returned if the config file exists but could not be loaded.
// No error is returned if the config file does not exist.
func Load() (config.Config, error) {
	cFile := config.File()
	if _, err := os.Stat(cFile); err != nil {
		oldCFile := oldConfigFile()

		// config file does not exist, check older version for backward compatibility
		if _, err := os.Stat(oldCFile); err != nil {
			return config.Config{}, nil
		}

		// older version exists
		logrus.Infof("settings from older %s version detected and copied", config.AppName)
		if err := cli.Command("cp", oldCFile, cFile).Run(); err != nil {
			logrus.Warn(fmt.Errorf("error copying config: %w, proceeding with defaults", err))
			return config.Config{}, nil
		}
	}

	var c config.Config
	b, err := os.ReadFile(cFile)
	if err != nil {
		return c, fmt.Errorf("could not load previous settings: %w", err)
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return c, fmt.Errorf("could not load previous settings: %w", err)
	}
	return c, nil
}

// Teardown deletes the config.
func Teardown() error {
	if _, err := os.Stat(config.Dir()); err == nil {
		return os.RemoveAll(config.Dir())
	}
	return nil
}