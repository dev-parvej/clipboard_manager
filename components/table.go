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
	Width   float32
	Height  float32
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

// buildTable constructs the table content
func (t *CustomTable) buildTable() {
	t.Container.Objects = nil // clear existing

	colCount := len(t.data.Headers)
	naturalWidths := make([]float32, colCount)

	// ---------- 1. Measure natural widths ----------
	// Headers
	for colIndex, header := range t.data.Headers {
		lbl := widget.NewLabel(header)
		w := lbl.MinSize().Width
		if w > naturalWidths[colIndex] {
			naturalWidths[colIndex] = w
		}
	}

	// Rows
	for _, row := range t.data.Rows {
		for colIndex, cell := range row {

			if btnConfig, isButton := t.data.ButtonColumn[colIndex]; isButton {
				btn := widget.NewButton(btnConfig.Text, nil)
				if btnConfig.Icon != nil {
					btn.SetIcon(btnConfig.Icon)
				}
				if btnConfig.Width > 0 && btnConfig.Width > naturalWidths[colIndex] {
					naturalWidths[colIndex] = btnConfig.Width
				}
			} else {
				lbl := widget.NewLabel(cell)
				w := lbl.MinSize().Width
				if w > naturalWidths[colIndex] {
					naturalWidths[colIndex] = w
				}
			}
		}
	}

	// ---------- 2. Compute proportions ----------
	var total float32
	for _, w := range naturalWidths {
		total += w
	}
	proportions := make([]float32, colCount)
	for i, w := range naturalWidths {
		if total > 0 {
			proportions[i] = w / total
		} else {
			proportions[i] = 1.0 / float32(colCount)
		}
	}

	// Get available width from container (fallback: sum of naturals)
	availableWidth := t.Container.Size().Width
	if availableWidth <= 0 {
		availableWidth = total
	}

	// Final responsive widths
	colWidths := make([]float32, colCount)
	for i := 0; i < colCount; i++ {
		colWidths[i] = proportions[i] * availableWidth
	}

	// ---------- 3. HEADERS ----------
	headerRow := container.NewHBox()
	for colIndex, header := range t.data.Headers {
		lbl := widget.NewLabel(header)
		lbl.TextStyle = fyne.TextStyle{Bold: true}
		headerRow.Add(fixedSize(lbl, colWidths[colIndex], lbl.MinSize().Height))
	}
	t.Container.Add(headerRow)

	// ---------- 4. ROWS ----------
	for rowIndex, row := range t.data.Rows {
		rowBox := container.NewHBox()
		for colIndex, cell := range row {
			var obj fyne.CanvasObject

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
				obj = btn
			} else {
				if trimConfig, shouldTrim := t.data.TrimConfig[colIndex]; shouldTrim && len(cell) > trimConfig.MaxChars {
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

					obj = container.NewHBox(textLabel, moreBtn)
				} else {
					lbl := widget.NewLabel(cell)
					if t.data.Modifiers != nil {
						if modifier, exists := t.data.Modifiers[colIndex]; exists {
							modifier(lbl, cell)
						}
					}
					obj = lbl
				}
			}

			rowBox.Add(fixedSize(obj, colWidths[colIndex], obj.MinSize().Height))
		}
		t.Container.Add(rowBox)
	}

	t.Container.Refresh()
}

// helper: forces a widget into fixed size
func fixedSize(obj fyne.CanvasObject, w, h float32) fyne.CanvasObject {
	return container.New(&fixedSizeLayout{w, h}, obj)
}

type fixedSizeLayout struct {
	w, h float32
}

func (f *fixedSizeLayout) Layout(objects []fyne.CanvasObject, _ fyne.Size) {
	if len(objects) == 0 {
		return
	}
	objects[0].Resize(fyne.NewSize(f.w, f.h))
	objects[0].Move(fyne.NewPos(0, 0))
}

func (f *fixedSizeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(f.w, f.h)
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
