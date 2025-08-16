package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
)

func main() {
	fmt.Println("🔧 Initializing Clipboard Manager...")

	// Initialize database
	InitDB()
	fmt.Println("✅ Database initialized")

	// Use the simple approach that avoids threading issues
	SetupSystemTraySimple()
}

func SetupSystemTraySimple() {
	myApp := app.New()

	stopCh := make(chan bool)

	StartClipboardMonitor(myApp, stopCh)
	SetupSystemTray(myApp, stopCh)

	fmt.Println("✨ Clipboard Manager is running!")
	fmt.Println("💡 Copy some text to see it captured.")
	fmt.Println("⏹️  Press Ctrl+C to stop the application.")

	// This will keep the app running
	myApp.Run()
}
