package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

const version string = "INDEV"

type queueEntryWidget struct {
	widget.Label
	UUID uuid.UUID
}

type queueEntryData struct {
	Title    string
	Duration uint
	QueuedBy string
	UUID     uuid.UUID
}

var queueDataArray []queueEntryData

func fetchQueueData() {
	return
}

func newQueueEntry() *queueEntryWidget {
	qe := &queueEntryWidget{}
	qe.ExtendBaseWidget(qe)
	return qe
}

func (qe *queueEntryWidget) TappedSecondary(pe *fyne.PointEvent) {
	menu := fyne.NewMenu(
		"Right Menu",
		fyne.NewMenuItem("Remove from queue", func() {
			for i := range len(queueDataArray) {
				if qe.UUID == queueDataArray[i].UUID {
					queueDataArray = append(queueDataArray[:i], queueDataArray[i+1:]...)

					fetchQueueData()

					break
				}
			}
		}),
	)

	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(qe)
	popUpPos := entryPos.Add(pe.Position)
	c := fyne.CurrentApp().Driver().CanvasForObject(qe)
	widget.ShowPopUpMenuAtPosition(menu, c, popUpPos)
}

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

func headerBarConstructor(w fyne.Window) *fyne.Container {
	// initialize topbar
	serverSelectionDropdown := widget.NewSelect([]string{"Midnight Cookout", "KBot Testing Grounds", "The Groovers"}, func(s string) {
		prop := canvas.NewRectangle(color.Transparent)
		prop.SetMinSize(fyne.NewSize(50, 50))

		a3 := widget.NewActivity()
		d := dialog.NewCustomWithoutButtons("Requesting Server Data...", container.NewStack(prop, a3), w)
		a3.Start()
		d.Show()

		go func() {
			time.Sleep(time.Second * 3)
			a3.Stop()
			d.Hide()
		}()

	})

	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		if reader == nil {
			return
		}
	}, w)

	uploadFileButton := widget.NewButton("Upload File", func() {
		fileDialog.Resize(fyne.NewSize(w.Content().Size().Width*.8, w.Content().Size().Height*.8))
		fileDialog.Show()
	})

	topMainContainer := container.NewBorder(nil, nil, uploadFileButton, serverSelectionDropdown)

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

func contentConstructor(w fyne.Window) *fyne.Container {
	queueDataArray = []queueEntryData{}
	mediaBar := mediaBarConstructor()
	headerBar := headerBarConstructor(w)
	queueBar := queueBarConstructor()

	content := container.NewBorder(
		headerBar,
		mediaBar,
		nil,
		queueBar,
	)

	return content
}

func main() {
	a := app.New()
	w := a.NewWindow("KBot Media Player " + version)

	content := contentConstructor(w)

	w.SetContent(content)
	w.Resize(fyne.Size{Width: 1024, Height: 768})
	w.ShowAndRun()
}
