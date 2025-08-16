package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func SetupSystemTray(a fyne.App, stopCh chan bool) error {
	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayMenu(
			fyne.NewMenu("Clipboard Manager",
				fyne.NewMenuItem("View History", func() {
					ShowHistoryWindow()
				}),
				fyne.NewMenuItem("Clear All", func() {
					ClearAllEntries()
				}),
				fyne.NewMenuItemSeparator(),
				fyne.NewMenuItem("Quit", func() {
					stopCh <- true
					a.Quit()
				}),
			),
		)
	}
	return nil
}
