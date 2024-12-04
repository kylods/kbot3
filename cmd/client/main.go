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
)

const version string = "INDEV"

func main() {
	a := app.New()
	w := a.NewWindow("KBot Media Player " + version)

	// Initialize bottomMediaInfoContainer
	mediaProgressBar := widget.NewSlider(0, 1)
	mediaProgressBar.SetValue(0)
	mediaProgressBar.Disable()
	mediaProgressBar.Resize(fyne.Size{Width: 100, Height: 20})
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

	// initialize topbar
	serverSelectionDropdown := widget.NewSelect([]string{"Midnight Cookout", "KBot Testing Grounds", "The Groovers"}, func(s string) {
		prop := canvas.NewRectangle(color.Transparent)
		prop.SetMinSize(fyne.NewSize(50, 50))

		a3 := widget.NewActivity()
		d := dialog.NewCustomWithoutButtons("Please wait...", container.NewStack(prop, a3), w)
		a3.Start()
		d.Show()

		go func() {
			time.Sleep(time.Second * 3)
			a3.Stop()
			d.Hide()
		}()
	})

	topMainContainer := container.NewBorder(nil, nil, nil, serverSelectionDropdown)

	content := container.NewBorder(
		topMainContainer,
		bottomMediaInfoContainer,
		nil,
		nil,
	)

	w.SetContent(content)
	w.Resize(fyne.Size{Width: 800, Height: 600})
	w.ShowAndRun()
}
