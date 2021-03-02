package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	Service ServiceConfig
	Slack   SlackConfig
	File    FileConfig
	Log     LogConfig
}

// ServiceConfig ...
type ServiceConfig struct {
	Port string `toml:"Port"`
}

// SlackConfig ...
type SlackConfig struct {
	FileURL  string `toml:"FileURL"`
	BotToken string `toml:"BotToken"`
	Channel  string `toml:"Channel"`
	ChatURL  string `toml:"ChatURL"`
}

// FileConfig ...
type FileConfig struct {
	Path  string `toml:"Path"`
	DLURL string `toml:"DLURL"`
}

// LogConfig ...
type LogConfig struct {
	LogLevel   string `toml:"LogLevel"`
	OutputPath string `toml:"OutputPath"`
	FileExtend string `toml:"FileExtend"`
}

// GetConfig ...
func GetConfig() (config Config) {
	// read config.
	exe, _ := os.Executable()
	confPath := filepath.Join(filepath.Dir(exe), "config.toml")
	_, err := toml.DecodeFile(confPath, &config)
	if err != nil {
		// Error Handling
		fmt.Println(err.Error())
	}

	return config
}
