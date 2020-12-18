package slack

// CodeBlockName represents the name of the modal code block
const CodeBlockName = "main_code_block"

// LanguageBlockName represents the language selector name
const LanguageBlockName = "language_block"

// ConversationSelectBlockName represents the name of the modal convo select block
const ConversationSelectBlockName = "response_block"

// CodeActionID represents the name of the code element action
const CodeActionID = "code_input"

// LanguageActionID represents the action of the language selector input
const LanguageActionID = "select_language"

const plainTextType = "plain_text"
const inputType = "input"

// GenerateRESLModal returns a payload that contains a language-specific resl modal
func GenerateRESLModal(languageName string, languageShortName string, placeholder string) ModalDefinition {
	blocks := []Block{
		Block{
			BlockID: CodeBlockName,
			Type:    inputType,
			Element: &Element{
				Type:      "plain_text_input",
				ActionID:  CodeActionID,
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
	}

	// Todo make options dynamic based on supported lanugages in languages.json
	if languageShortName == "" {
		blocks = append(blocks, Block{
			BlockID: LanguageBlockName,
			Type:    inputType,
			Label: &ViewOptions{
				Type: plainTextType,
				Text: "Select a coding language",
			},
			Element: &Element{
				Type:     "static_select",
				ActionID: LanguageActionID,
				Options: []SelectOption{
					SelectOption{
						Text: ViewOptions{
							Type: plainTextType,
							Text: "JavaScript",
						},
						Value: "js",
					},
					SelectOption{
						Text: ViewOptions{
							Type: plainTextType,
							Text: "Python",
						},
						Value: "py",
					},
					SelectOption{
						Text: ViewOptions{
							Type: plainTextType,
							Text: "Python3",
						},
						Value: "py3",
					},
				},
			},
		})
	}

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
		Blocks:          blocks,
	}
}
