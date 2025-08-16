package main

import (
	"database/sql"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/dev-parvej/clipboard_manager/components"
)

// ShowHistoryWindow displays clipboard items grouped by date
func ShowHistoryWindow() error {
	w := fyne.CurrentApp().NewWindow("Clipboard History")

	// Title
	title := canvas.NewText("Clipboard History", color.White)
	title.TextSize = 24
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Main container for all date sections
	mainContainer := container.NewVBox()

	// Toolbar
	clearAllBtn := widget.NewButtonWithIcon("Clear All", theme.DeleteIcon(), func() {
		ClearAllEntries()
		FetchAndShowTheData(mainContainer, w)
	})
	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		FetchAndShowTheData(mainContainer, w)
	})
	buttons := container.NewHBox(
		layout.NewSpacer(),
		container.NewPadded(refreshBtn),
		container.NewPadded(clearAllBtn),
		layout.NewSpacer(),
	)

	toolbar := container.NewVBox(
		container.NewCenter(title),
		container.NewPadded(buttons),
		widget.NewSeparator(),
	)

	// Status bar
	statusBar := container.NewHBox(
		widget.NewLabel(""),
		layout.NewSpacer(),
		widget.NewLabel("v1.0.0"),
	)

	// OUTER scrollable content
	centerScroll := container.NewVScroll(mainContainer)

	content := container.NewBorder(
		toolbar,
		container.NewPadded(statusBar),
		nil, nil,
		centerScroll,
	)

	FetchAndShowTheData(mainContainer, w)
	w.SetContent(content)
	w.Resize(fyne.NewSize(750, 700))
	w.Show()
	return nil
}

func FetchAndShowTheData(mainContainer *fyne.Container, w fyne.Window) {
	mainContainer.Objects = nil // Clear old content

	type clipboardEntry struct {
		id        int
		content   string
		imagePath sql.NullString
		timestamp time.Time
	}

	// Fetch from DB
	fetchData := func() []clipboardEntry {
		var entries []clipboardEntry
		rows, err := db.Query("SELECT id, content, image_path, timestamp FROM clipboard ORDER BY timestamp DESC")
		if err != nil {
			log.Println(err)
			return entries
		}
		defer rows.Close()

		for rows.Next() {
			var entry clipboardEntry
			if err := rows.Scan(&entry.id, &entry.content, &entry.imagePath, &entry.timestamp); err == nil {
				entries = append(entries, entry)
			}
		}
		return entries
	}

	entries := fetchData()

	// Group by date
	dateGroups := make(map[string][]clipboardEntry)
	for _, entry := range entries {
		dateStr := entry.timestamp.Format("2006-01-02")
		dateGroups[dateStr] = append(dateGroups[dateStr], entry)
	}

	// Build each date section
	for dateStr, dateEntries := range dateGroups {
		// Fix closure capture
		ds := dateStr

		// Header
		dateHeader := container.NewHBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Date: %s", ds),
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Clear Date", theme.DeleteIcon(), func() {
				ClearEntriesByDate(ds)
				FetchAndShowTheData(mainContainer, w)
			}),
		)

		// Row conversion
		convertToRows := func() [][]string {
			rows := make([][]string, len(dateEntries))
			for i, entry := range dateEntries {
				content := entry.content
				if entry.imagePath.Valid && entry.imagePath.String != "" {
					content = "[Image]"
				}
				if len(content) > 50 {
					content = content
				}
				rows[i] = []string{
					content,
					entry.timestamp.Format("15:04:05"),
					"", // Action column
				}
			}
			return rows
		}

		rows := convertToRows()

		// Table
		tableData := components.TableData{
			Headers:   []string{"Content", "Time", "Actions"},
			Rows:      rows,
			Widths:    []float32{750, 150, 100},
			OnRefresh: convertToRows,
			ButtonColumn: map[int]components.ButtonConfig{
				2: {
					Text:   "",
					Icon:   theme.DeleteIcon(),
					Height: 40,
					Width:  50,
					OnClick: func(rowIndex int) {
						DeleteEntry(dateEntries[rowIndex].id)
						FetchAndShowTheData(mainContainer, w)
					},
				},
			},
			TrimConfig: map[int]components.TrimConfig{
				0: {
					MaxChars: 100,
					OnMoreClick: func(rowIndex int, fullText string) {
						// Get current window size
						winSize := w.Canvas().Size()

						// Scrollable content
						scrollContent := container.NewScroll(widget.NewLabel(fullText))
						scrollContent.SetMinSize(fyne.NewSize(winSize.Width*0.8, winSize.Height*0.8))

						// Create copy button
						copyBtn := widget.NewButton("Copy", func() {
							w.Clipboard().SetContent(fullText) // copy full text to clipboard
						})

						// Close button
						closeBtn := widget.NewButton("Close", func() {
							// Will close dialog later, see below
						})

						// Combine buttons in horizontal layout
						buttonBar := container.NewHBox(copyBtn, layout.NewSpacer(), closeBtn)

						// Create a container for content + buttons
						contentWithButtons := container.NewBorder(nil, buttonBar, nil, nil, scrollContent)

						// Show dialog
						d := dialog.NewCustom("Full Description", "", contentWithButtons, w)

						// Close button callback
						closeBtn.OnTapped = func() {
							d.Hide()
						}

						d.Show()
					},
				},
			},
		}
		table := components.CreateTable(tableData)

		dateSection := container.NewVBox(
			container.NewPadded(dateHeader),
			container.NewPadded(table.CanvasObject()),
		)

		// Add images if any
		for _, entry := range dateEntries {
			if entry.imagePath.Valid && entry.imagePath.String != "" {
				if file, err := os.Open(entry.imagePath.String); err == nil {
					if img, err := png.Decode(file); err == nil {
						file.Close()
						cImg := canvas.NewImageFromImage(img)
						cImg.FillMode = canvas.ImageFillContain
						cImg.SetMinSize(fyne.NewSize(200, 150))
						dateSection.Add(cImg)
					}
				}
			}
		}

		mainContainer.Add(dateSection)
	}
}

// DB utility functions

func ClearEntriesByDate(dateStr string) {
	start, _ := time.Parse("2006-01-02", dateStr)
	end := start.Add(24 * time.Hour)
	_, err := db.Exec("DELETE FROM clipboard WHERE timestamp >= ? AND timestamp < ?", start, end)
	if err != nil {
		log.Println("Error clearing date:", err)
	}
}

func DeleteEntry(id int) {
	_, err := db.Exec("DELETE FROM clipboard WHERE id=?", id)
	if err != nil {
		log.Println("Delete error:", err)
	}
}

func ClearAllEntries() {
	_, err := db.Exec("DELETE FROM clipboard")
	if err != nil {
		log.Println("Clear all error:", err)
	}
}
