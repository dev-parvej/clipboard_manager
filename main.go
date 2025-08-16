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

	// myApp.SetMetadata(&app.Metadata{
	// 	Name: "",
	// 	ID:   "com.example.clipboardmanager",
	// })

	// Start clipboard monitoring with the app instance
	StartClipboardMonitor(myApp)
	SetupSystemTray(myApp)

	fmt.Println("✨ Clipboard Manager is running!")
	fmt.Println("💡 Copy some text to see it captured.")
	fmt.Println("⏹️  Press Ctrl+C to stop the application.")

	// This will keep the app running
	myApp.Run()
}
