package slack

// CodeBlockName represents the name of the modal code block
const CodeBlockName = "main_code_block"

// ConversationSelectBlockName represents the name of the modal convo select block
const ConversationSelectBlockName = "response_block"

const plainTextType = "plain_text"
const inputType = "input"

// GenerateRESLModal returns a payload that contains a language-specific resl modal
func GenerateRESLModal(languageName string, languageShortName string, placeholder string) ModalDefinition {
	return ModalDefinition{
		Type: "modal",
		Title: &ViewOptions{
			Type:  plainTextType,
			Text:  "RESL",
			Emoji: true,
		},
		Submit: &ViewOptions{
			Type:  plainTextType,
			Text:  "Run Code",
			Emoji: true,
		},
		Close: &ViewOptions{
			Type:  plainTextType,
			Text:  "Cancel",
			Emoji: true,
		},
		PrivateMetadata: languageShortName,
		Blocks: []Block{
			Block{
				BlockID: CodeBlockName,
				Type:    inputType,
				Element: &Element{
					Type:      "plain_text_input",
					ActionID:  "code_input",
					Multiline: true,
					Placeholder: &ViewOptions{
						Type: plainTextType,
						Text: placeholder,
					},
				},
				Label: &ViewOptions{
					Type: plainTextType,
					Text: "Enter " + languageName + " here",
				},
				Hint: &ViewOptions{
					Type: plainTextType,
					Text: "Wrapping your code in backticks is optional",
				},
			},
			Block{
				BlockID:  ConversationSelectBlockName,
				Type:     inputType,
				Optional: true,
				Label: &ViewOptions{
					Type: plainTextType,
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
