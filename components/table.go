package components

import (
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ButtonConfig defines button properties for a specific column
type ButtonConfig struct {
	Text    string
	Icon    fyne.Resource
	OnClick func(rowIndex int)
}

// TrimConfig defines trimming properties for a specific column
type TrimConfig struct {
	MaxChars    int
	OnMoreClick func(rowIndex int, fullText string)
}

type CustomTable struct {
	Container   *fyne.Container
	refreshData func()
	data        *TableData
}

type TableData struct {
	Headers      []string
	Rows         [][]string
	Widths       []float32
	Modifiers    map[int]func(label *widget.Label, value string)
	OnRefresh    func() [][]string
	ButtonColumn map[int]ButtonConfig
	TrimConfig   map[int]TrimConfig // Maps column index to trim configuration
}

// Refresh redraws the table content
func (t *CustomTable) Refresh() {
	if t.refreshData != nil {
		t.refreshData()
	}
	t.buildTable()
}

func (t *CustomTable) CanvasObject() fyne.CanvasObject {
	return t.Container
}

// buildTable constructs the grid content
func (t *CustomTable) buildTable() {
	t.Container.Objects = nil // clear existing
	grid := container.NewGridWithColumns(len(t.data.Headers))

	// Add headers
	for _, header := range t.data.Headers {
		lbl := widget.NewLabel(header)
		lbl.TextStyle = fyne.TextStyle{Bold: true}
		grid.Add(lbl)
	}

	// Add rows
	for rowIndex, row := range t.data.Rows {
		for colIndex, cell := range row {
			// Check if this column should display a button
			if btnConfig, isButton := t.data.ButtonColumn[colIndex]; isButton {
				btn := widget.NewButton(btnConfig.Text, func(rIdx int) func() {
					return func() {
						if btnConfig.OnClick != nil {
							btnConfig.OnClick(rIdx)
						}
					}
				}(rowIndex))
				if btnConfig.Icon != nil {
					btn.SetIcon(btnConfig.Icon)
				}
				grid.Add(btn)
			} else {
				// Check for trim configuration
				if trimConfig, shouldTrim := t.data.TrimConfig[colIndex]; shouldTrim && len(cell) > trimConfig.MaxChars {
					// Create container for trimmed text and more button
					trimmedText := cell[:trimConfig.MaxChars]

					re := regexp.MustCompile(`\s+`)

					cleanText := re.ReplaceAllString(
						strings.ReplaceAll(strings.ReplaceAll(trimmedText, "\n", ""), "\r", ""),
						" ",
					)

					textLabel := widget.NewLabel(cleanText)

					moreBtn := widget.NewButton("->", func(rIdx int, fullText string) func() {
						return func() {
							if trimConfig.OnMoreClick != nil {
								trimConfig.OnMoreClick(rIdx, fullText)
							}
						}
					}(rowIndex, cell))

					hBox := container.NewHBox(textLabel, moreBtn)
					grid.Add(hBox)
				} else {
					lbl := widget.NewLabel(cell)
					if t.data.Modifiers != nil {
						if modifier, exists := t.data.Modifiers[colIndex]; exists {
							modifier(lbl, cell)
						}
					}
					grid.Add(lbl)
				}
			}
		}
	}

	t.Container.Add(grid)
	t.Container.Refresh()
}

// CreateTable returns a backward-compatible custom table
func CreateTable(data TableData) *CustomTable {
	c := &CustomTable{
		Container: container.NewVBox(),
		data:      &data,
		refreshData: func() {
			if data.OnRefresh != nil {
				data.Rows = data.OnRefresh()
			}
		},
	}

	c.buildTable()
	return c
}
