package slack

type ViewOptions struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

type Element struct {
	ActionID                     string      `json:"action_id"`
	Type                         string      `json:"type"`
	DefaultToCurrentConversation bool        `json:"default_to_current_conversation"`
	ResponseURLEnabled           bool        `json:"response_url_enabled"`
	Multiline                    bool        `json:"multiline"`
	Placeholder                  ViewOptions `json:"placeholder"`
}

type Block struct {
	BlockID  string      `json:"block_id"`
	Type     string      `json:"type"`
	Element  Element     `json:"element"`
	Label    ViewOptions `json:"label"`
	Hint     ViewOptions `json:"hint"`
	Optional bool        `json:"optional"`
}

type ModalDefinition struct {
	Type            string      `json:"type"`
	Title           ViewOptions `json:"title"`
	Submit          ViewOptions `json:"submit"`
	Close           ViewOptions `json:"close"`
	PrivateMetadata string      `json:"private_metadata"`
	Blocks          []Block     `json:"blocks"`
}

// GenerateRESLModal returns a payload that contains a language-specific
// resl modal
func GenerateRESLModal(languageName, placeholder string) ModalDefinition {
	return ModalDefinition{
		Type: "modal",
		Title: ViewOptions{
			Type:  "plain_text",
			Text:  "RESL",
			Emoji: true,
		},
		Submit: ViewOptions{
			Type:  "plain_text",
			Text:  "Run Code",
			Emoji: true,
		},
		Close: ViewOptions{
			Type:  "plain_text",
			Text:  "Cancel",
			Emoji: true,
		},
		PrivateMetadata: languageName,
		Blocks: []Block{
			Block{
				BlockID: "main_block",
				Type:    "input",
				Element: Element{
					Type:      "plain_text_input",
					ActionID:  "code_input",
					Multiline: true,
					Placeholder: ViewOptions{
						Type: "plain_text",
						Text: placeholder,
					},
				},
				Label: ViewOptions{
					Type: "plain_text",
					Text: "Enter " + languageName + " here",
				},
				Hint: ViewOptions{
					Type: "plain_text",
					Text: "Wrapping your code in backticks is optional",
				},
			},
			Block{
				BlockID:  "response_block",
				Type:     "input",
				Optional: true,
				Label: ViewOptions{
					Type: "plain_text",
					Text: "Select a channel to post the result in",
				},
				Element: Element{
					ActionID:                     "conversation_select_action",
					Type:                         "conversations_select",
					DefaultToCurrentConversation: true,
					ResponseURLEnabled:           true,
				},
			},
		},
	}
}
