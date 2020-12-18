package slack

// GenerateRESLModal returns a payload that contains a language-specific
// resl modal
func GenerateRESLModal(languageName string, languageShortName string, placeholder string) ModalDefinition {
	return ModalDefinition{
		Type: "modal",
		Title: &ViewOptions{
			Type:  "plain_text",
			Text:  "RESL",
			Emoji: true,
		},
		Submit: &ViewOptions{
			Type:  "plain_text",
			Text:  "Run Code",
			Emoji: true,
		},
		Close: &ViewOptions{
			Type:  "plain_text",
			Text:  "Cancel",
			Emoji: true,
		},
		PrivateMetadata: languageShortName,
		Blocks: []Block{
			Block{
				BlockID: "main_code_block",
				Type:    "input",
				Element: &Element{
					Type:      "plain_text_input",
					ActionID:  "code_input",
					Multiline: true,
					Placeholder: &ViewOptions{
						Type: "plain_text",
						Text: placeholder,
					},
				},
				Label: &ViewOptions{
					Type: "plain_text",
					Text: "Enter " + languageName + " here",
				},
				Hint: &ViewOptions{
					Type: "plain_text",
					Text: "Wrapping your code in backticks is optional",
				},
			},
			Block{
				BlockID:  "response_block",
				Type:     "input",
				Optional: true,
				Label: &ViewOptions{
					Type: "plain_text",
					Text: "Select a channel to post the result in",
				},
				Element: &Element{
					ActionID:                     "conversation_select_action",
					Type:                         "conversations_select",
					DefaultToCurrentConversation: true,
					ResponseURLEnabled:           true,
				},
			},
		},
	}
}
