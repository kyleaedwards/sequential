package main

import (
	"encoding/json"
	"io/ioutil"
)

// Sequential configuration settings
type Config struct {
	Styles ConfigStyles `json:"styles"`
}

// Configurable TUI text colors
type ConfigStyles struct {
	SelectedColor  string `json:"selectedColor"`
	CompletedColor string `json:"completedColor"`
	DisabledColor  string `json:"disabledColor"`
}

var defaultStyles = ConfigStyles{
	SelectedColor:  "#40c997",
	CompletedColor: "#907af0",
	DisabledColor:  "#777777",
}

// Fills empty unmarshalled config data with default values.
func setDefaults(config *Config) {
	if config.Styles.SelectedColor == "" {
		config.Styles.SelectedColor = defaultStyles.SelectedColor
	}
	if config.Styles.CompletedColor == "" {
		config.Styles.CompletedColor = defaultStyles.CompletedColor
	}
	if config.Styles.DisabledColor == "" {
		config.Styles.DisabledColor = defaultStyles.DisabledColor
	}
}

// Reads configuration file from ~/.sequential/config.json. If the
// file is not present or config fields aren't set, it falls back
// to defaults defined in this file.
func LoadConfig() Config {
	fh := EnsureFileExists("config.json")
	defer fh.Close()

	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return Config{Styles: defaultStyles}
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{Styles: defaultStyles}
	}

	setDefaults(&config)
	return config
}
