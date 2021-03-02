package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anabebe1020/logger"
)

var (
	config  Config
	logConf logger.LogConf
)

func initialize() bool {
	exe, _ := os.Executable()
	// Get Config.
	config = GetConfig()
	if _, err := os.Stat(config.File.Path); err != nil {
		fmt.Println("not config.", err)
		return false
	}
	// Logger Config.
	logConf.PrgName = "FileUploader"
	logConf.LogLevel = config.Log.LogLevel
	logConf.LogPath = filepath.Join(filepath.Join(filepath.Dir(exe),
		config.Log.OutputPath), logConf.PrgName+config.Log.FileExtend)
	return true
}

func main() {
	if !initialize() {
		return
	}
	service()
}
