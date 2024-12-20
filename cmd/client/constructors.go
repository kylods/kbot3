package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func mediaBarConstructor() *fyne.Container {
	// Initialize bottomMediaInfoContainer
	mediaProgressBar := widget.NewSlider(0, 1)
	mediaProgressBar.SetValue(0)
	mediaProgressBar.Disable()

	// placeholder values
	mediaLengthSeconds := 90
	mediaNowSeconds := 0

	mediaDurationLabel := widget.NewLabel("00:00")
	mediaName := widget.NewLabel("No Media Playing")

	// placeholder code for media timer
	go func() {
		for {
			time.Sleep(time.Second / 10)
			mediaProgressBar.Max = float64(mediaLengthSeconds)
			if mediaLengthSeconds > mediaNowSeconds {
				mediaNowSeconds++
			} else {
				mediaNowSeconds = 0
			}

			if mediaNowSeconds == 0 || mediaLengthSeconds == 0 {
				mediaProgressBar.SetValue(0)
			} else {
				mediaProgressBar.SetValue(float64(mediaNowSeconds))
			}

			var formattedText string = fmt.Sprintf("%02d:%02d / %02d:%02d", mediaNowSeconds/60, mediaNowSeconds%60, mediaLengthSeconds/60, mediaLengthSeconds%60)

			mediaDurationLabel.SetText(formattedText)
		}
	}()

	bottomMediaInfoContainer := container.NewVBox(container.NewBorder(nil, nil, mediaDurationLabel, mediaName), mediaProgressBar)

	return bottomMediaInfoContainer
}

func headerBarConstructor(w fyne.Window, as *appState) *fyne.Container {
	// initialize topbar
	serverSelectionDropdown := widget.NewSelect([]string{}, func(s string) {

	})

	as.serverSelectDropdown = serverSelectionDropdown

	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		if reader == nil {
			return
		}

		uploadFileLogic(reader)

	}, w)

	uploadFileButton := widget.NewButton("Upload File", func() {
		fileDialog.Resize(fyne.NewSize(w.Content().Size().Width*.8, w.Content().Size().Height*.8))
		fileDialog.Show()
	})

	connectEntry := widget.NewEntry()
	connectEntry.Text = serverURL

	connectDialog := dialog.NewForm(
		"Connect to KBot",
		"Connect",
		"Cancel",
		[]*widget.FormItem{{Text: "IP Address:Port", Widget: connectEntry}},
		func(b bool) {
			if b { // Runs when pressing 'Connect'
				serverURL = connectEntry.Text
				connectToServer(as)
			} else { // Runs when dismissed
				return
			}
		},
		w)

	connectButton := widget.NewButton("Connect to Server", func() {
		connectDialog.Resize(fyne.NewSize(430, 100))
		connectDialog.Show()
	})

	buttonsBox := container.NewHBox(uploadFileButton, connectButton)

	topMainContainer := container.NewBorder(nil, nil, buttonsBox, serverSelectionDropdown)

	return topMainContainer
}

func queueBarConstructor() *fyne.Container {

	playerControlsBar := widget.NewToolbar()

	queueList := widget.NewList(
		func() int {
			return len(queueDataArray)
		},
		func() fyne.CanvasObject {

			entry := newQueueEntry()

			return entry
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*queueEntryWidget).UUID = queueDataArray[lii].UUID
			co.(*queueEntryWidget).SetText(queueDataArray[lii].Title)
		})

	queueBar := container.NewBorder(nil, playerControlsBar, nil, nil, queueList)

	return queueBar
}

func contentConstructor(w fyne.Window, as *appState) *fyne.Container {
	queueDataArray = []queueEntryData{}
	mediaBar := mediaBarConstructor()
	headerBar := headerBarConstructor(w, as)
	queueBar := queueBarConstructor()

	content := container.NewBorder(
		headerBar,
		mediaBar,
		nil,
		queueBar,
	)

	return content
}
