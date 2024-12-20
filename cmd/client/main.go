package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"mime/multipart"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type appState struct {
	serverMap             map[string]string
	serverSelectDropdown  *widget.Select
	networkStatusLabel    *widget.Label
	networkStatusActivity *widget.Activity
	networkStatusDialog   *dialog.CustomDialog
}

func (s *appState) updateServerSelectOptions() {
	newOptions := []string{}
	for _, name := range s.serverMap {
		newOptions = append(newOptions, name)
	}

	s.serverSelectDropdown.Options = newOptions
	s.serverSelectDropdown.Refresh()
}

func (s *appState) startNetworkStatusOverlay() {
	s.networkStatusActivity.Start()
	s.networkStatusDialog.Show()
}

func (s *appState) stopNetworkStatusOverlay() {
	s.networkStatusActivity.Stop()
	s.networkStatusDialog.Hide()
}

const version string = "INDEV"

var serverURL string = "http://localhost:8080"

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

func uploadFileLogic(reader fyne.URIReadCloser) error {
	defer reader.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create form file field
	formFile, err := writer.CreateFormFile("audioFile", reader.URI().Name())
	if err != nil {
		return fmt.Errorf("error creating form file: %w", err)
	}

	// Copy filedata to form field
	_, err = io.Copy(formFile, reader)
	if err != nil {
		return fmt.Errorf("error copying file data: %w", err)
	}

	// Close writer & finalize multipart form
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error closing multipart writer: %w", err)
	}

	// Send request to API
	apiURL := serverURL + "/upload"
	req, err := http.NewRequest(http.MethodPost, apiURL, &requestBody)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s", string(respBody))
	}

	return nil
}

func main() {

	a := app.New()
	w := a.NewWindow("KBot Media Player " + version)

	clientAppState := &appState{}

	content := contentConstructor(w, clientAppState)

	prop := canvas.NewRectangle(color.Transparent)
	prop.SetMinSize(fyne.NewSize(50, 50))

	networkActivityWidget := widget.NewActivity()
	networkActivityDialog := dialog.NewCustomWithoutButtons("Requesting Server Data...", container.NewStack(prop, networkActivityWidget), w)

	clientAppState.networkStatusActivity = networkActivityWidget
	clientAppState.networkStatusDialog = networkActivityDialog

	w.SetContent(content)
	w.Resize(fyne.Size{Width: 1024, Height: 768})
	w.ShowAndRun()
}
