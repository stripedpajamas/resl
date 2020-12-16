package slack

type ViewOptions struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

type Element struct {
	ActionID                     string      `json:"action_id,omitempty"`
	Type                         string      `json:"type,omitempty"`
	DefaultToCurrentConversation bool        `json:"default_to_current_conversation,omitempty"`
	ResponseURLEnabled           bool        `json:"response_url_enabled,omitempty"`
	Multiline                    bool        `json:"multiline,omitempty"`
	Placeholder                  ViewOptions `json:"placeholder,omitempty"`
}

type Block struct {
	BlockID  string      `json:"block_id,omitempty"`
	Type     string      `json:"type,omitempty"`
	Element  Element     `json:"element,omitempty"`
	Label    ViewOptions `json:"label,omitempty"`
	Hint     ViewOptions `json:"hint,omitempty"`
	Optional bool        `json:"optional,omitempty"`
}

type ModalDefinition struct {
	Type            string      `json:"type,omitempty"`
	Title           ViewOptions `json:"title,omitempty"`
	Submit          ViewOptions `json:"submit,omitempty"`
	Close           ViewOptions `json:"close,omitempty"`
	PrivateMetadata string      `json:"private_metadata,omitempty"`
	Blocks          []Block     `json:"blocks,omitempty"`
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
