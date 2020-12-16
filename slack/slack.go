package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const slackViewsOpenURL = "https://slack.com/api/views.open"

// SlackResponse contains the properties necessary to respond to a message
type SlackResponse struct {
	ResponseType string `json:"response_type,omitempty"`
	Text         string `json:"text,omitempty"`
}

type SlackModal struct {
	TriggerID string          `json:"trigger_id"`
	View      ModalDefinition `json:"view"`
}

// PublicAcknowledgement shows the originally sent command in the channel
func PublicAcknowledgement() ([]byte, error) {
	return json.Marshal(SlackResponse{
		ResponseType: "in_channel",
	})
}

// PrivateAcknowledgement sends back a "visible to only you" message to the
// command initiator
func PrivateAcknowledgement(text string) ([]byte, error) {
	return json.Marshal(SlackResponse{
		Text: text,
	})
}

// SendChannelResponse sends text to a response url in a channel
func SendChannelResponse(url, text string) error {
	reqBody, err := json.Marshal(SlackResponse{
		ResponseType: "in_channel",
		Text:         text,
	})

	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		return err
	}

	// we don't care about slack's response
	resp.Body.Close()

	return nil
}

// SendModal sends a modal to the user who typed the command. The modal
// has language-specific placeholder code and shows the chosen language name
func SendModal(triggerID, languageName, placeholder string) error {
	client := &http.Client{}

	authToken := os.Getenv("SLACK_TOKEN")
	log.Printf("Auth token: %s\n", authToken)

	reqBody, err := json.Marshal(SlackModal{
		TriggerID: triggerID,
		View:      GenerateRESLModal(languageName, placeholder),
	})
	if err != nil {
		return err
	}
	log.Printf("Modal request body: %s\n", string(reqBody))

	req, err := http.NewRequest("POST", slackViewsOpenURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Slack response: %s\n", string(body))

	return nil
}
