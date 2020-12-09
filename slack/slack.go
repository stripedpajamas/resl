package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SlackResponse contains the properties necessary to respond to a message
type SlackResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// SendChannelResponse sends text to a response url in a channel
func SendChannelResponse(url, text string) (string, error) {
	reqBody, err := json.Marshal(SlackResponse{
		ResponseType: "in_channel",
		Text:         text,
	})

	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
