package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func connectToServer(as *appState) error {
	as.startNetworkStatusOverlay()
	defer as.stopNetworkStatusOverlay()

	// Check that the server is active
	_, err := http.Get(serverURL + "/ping")
	if err != nil {
		return err
	}

	// Fetch active guild data for the bot
	resp, err := http.Get(serverURL + "/discord/guilds")
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data := make(map[string]string)
	json.Unmarshal(body, &data)

	as.serverMap = data
	as.updateServerSelectOptions()

	return nil
}
