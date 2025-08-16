package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
)

func StartClipboardMonitor(myApp fyne.App) {
	// Use the existing app instead of creating a new one
	myWindow := myApp.NewWindow("Clipboard Monitor")
	myWindow.Hide() // Hide the window since we don't need UI

	var last string

	// Use a ticker for better performance
	ticker := time.NewTicker(500 * time.Millisecond)

	// Run clipboard monitoring in a goroutine
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			// Access clipboard directly - this should work from goroutine in newer Fyne versions
			clipboard := myWindow.Clipboard()
			current := clipboard.Content()

			if current != "" && current != last {
				StoreTextClipboard(current)
				last = current
				fmt.Println("Captured:", current)
			}

			// This can run outside the main thread
			CleanupOldEntries()
		}
	}()
}
