package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
)

func StartClipboardMonitor(myApp fyne.App, stopCh <-chan bool) {
	myWindow := myApp.NewWindow("Clipboard Monitor")
	myWindow.Hide()

	var last string
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-stopCh:
				fmt.Println("Clipboard monitor stopped")
				return
			case <-ticker.C:
				clipboard := myApp.Clipboard()
				current := clipboard.Content()
				if current != "" && current != last {
					StoreTextClipboard(current)
					last = current
					fmt.Println("Captured:", current)
				}
				CleanupOldEntries()
			}
		}
	}()
}
