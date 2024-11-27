package main

import (
	"fmt"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Microphone represents a single microphone
type Microphone struct {
	Number int
	Type   string
	Owner  string
}

// String returns the current state of the microphone as a string
func (m *Microphone) String() string {
	if m.Owner != "" {
		return fmt.Sprintf("%d (%s) live with %s", m.Number, m.Type, m.Owner)
	}
	return fmt.Sprintf("%d (%s) is available", m.Number, m.Type)
}

var microphones []*Microphone

func main() {
	// Initialize the Fyne app
	a := app.New()
	w := a.NewWindow("Microphone Manager")
	w.Resize(fyne.NewSize(600, 400))

	// List widget to display microphone statuses
	micList := widget.NewList(
		func() int { return len(microphones) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(microphones[i].String())
		},
	)

	// Input fields
	numberEntry := widget.NewEntry()
	numberEntry.SetPlaceHolder("Microphone Number")
	typeEntry := widget.NewEntry()
	typeEntry.SetPlaceHolder("Microphone Type")
	ownerEntry := widget.NewEntry()
	ownerEntry.SetPlaceHolder("Owner Name")

	// Function to refresh the microphone list display
	refreshList := func() {
		micList.Refresh()
	}

	// Function to autofill details based on the selected microphone
	micList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(microphones) {
			selectedMic := microphones[id]
			numberEntry.SetText(strconv.Itoa(selectedMic.Number))
			typeEntry.SetText(selectedMic.Type)
			ownerEntry.SetText(selectedMic.Owner)
		}
	}

	// Buttons
	createButton := widget.NewButton("Create Microphone", func() {
		number, err := strconv.Atoi(numberEntry.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Invalid microphone number"), w)
			return
		}
		micType := typeEntry.Text
		if micType == "" {
			dialog.ShowError(fmt.Errorf("Microphone type cannot be empty"), w)
			return
		}
		microphones = append(microphones, &Microphone{Number: number, Type: micType})
		numberEntry.SetText("")
		typeEntry.SetText("")
		ownerEntry.SetText("")
		refreshList()
	})

	liveButton := widget.NewButton("Set Live", func() {
		number, err := strconv.Atoi(numberEntry.Text)
		if err != nil || number <= 0 || number > len(microphones) {
			dialog.ShowError(fmt.Errorf("Invalid microphone number"), w)
			return
		}
		owner := ownerEntry.Text
		if owner == "" {
			dialog.ShowError(fmt.Errorf("Owner name cannot be empty"), w)
			return
		}
		microphones[number-1].Owner = owner
		refreshList()
	})

	dieButton := widget.NewButton("Set Returned", func() {
		number, err := strconv.Atoi(numberEntry.Text)
		if err != nil || number <= 0 || number > len(microphones) {
			dialog.ShowError(fmt.Errorf("Invalid microphone number"), w)
			return
		}
		microphones[number-1].Owner = ""
		refreshList()
	})

	// Layout the application
	inputForm := container.NewVBox(
		widget.NewLabel("Manage Microphones"),
		numberEntry,
		typeEntry,
		ownerEntry,
		createButton,
		liveButton,
		dieButton,
	)
	content := container.NewHSplit(
		micList,
		inputForm,
	)
	content.SetOffset(0.4)

	w.SetContent(content)
	w.ShowAndRun()
}
