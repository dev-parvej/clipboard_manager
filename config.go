package main

import (
	"os"
	"path/filepath"
	"time"
)

var (
	AppFolder       string
	DBPath          string
	ImageFolder     string
	CleanupInterval = 1 * time.Hour
)

func init() {
	home, _ := os.UserHomeDir()
	AppFolder = filepath.Join(home, ".clipboard_manager")
	os.MkdirAll(AppFolder, 0755)

	DBPath = filepath.Join(AppFolder, "clipboard.db")
	ImageFolder = filepath.Join(AppFolder, "images")
	os.MkdirAll(ImageFolder, 0755)
}
