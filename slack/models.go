package slack

// InputElement represents the slack input model
type InputElement struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

// ConversationSelectElement represents a slack conversation select model
type ConversationSelectElement struct {
	Type                 string `json:"type,omitempty"`
	SelectedConversation string `json:"selected_conversation,omitempty"`
}

// ViewState represents the current state of each block element
type ViewState struct {
	Values map[string]map[string]interface{} `json:"values,omitempty"`
}

// ViewOptions represents the options for the elements
type ViewOptions struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	Emoji bool   `json:"emoji,omitempty"`
}

// Element represents an element inside the modal blocks
type Element struct {
	ActionID                     string       `json:"action_id,omitempty"`
	Type                         string       `json:"type,omitempty"`
	DefaultToCurrentConversation bool         `json:"default_to_current_conversation,omitempty"`
	ResponseURLEnabled           bool         `json:"response_url_enabled,omitempty"`
	Multiline                    bool         `json:"multiline,omitempty"`
	Placeholder                  *ViewOptions `json:"placeholder,omitempty"`
}

// Block represents the different blocks in the modal
type Block struct {
	BlockID  string       `json:"block_id,omitempty"`
	Type     string       `json:"type,omitempty"`
	Element  *Element     `json:"element,omitempty"`
	Label    *ViewOptions `json:"label,omitempty"`
	Hint     *ViewOptions `json:"hint,omitempty"`
	Optional bool         `json:"optional,omitempty"`
}

// ModalDefinition represents the slack modal data
type ModalDefinition struct {
	Type            string       `json:"type,omitempty"`
	Title           *ViewOptions `json:"title,omitempty"`
	Submit          *ViewOptions `json:"submit,omitempty"`
	Close           *ViewOptions `json:"close,omitempty"`
	PrivateMetadata string       `json:"private_metadata,omitempty"`
	Blocks          []Block      `json:"blocks,omitempty"`
	State           *ViewState   `json:"state,omitempty"`
}

// Response contains the properties necessary to respond to a message
type Response struct {
	ResponseAction string `json:"response_action,omitempty"`
	ResponseType   string `json:"response_type,omitempty"`
	Text           string `json:"text,omitempty"`
}

// User represents a slack user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// ResponseURL represents a slack response url object
type ResponseURL struct {
	ActionID  string `json:"action_id"`
	ChannelID string `json:"channel_id"`
	URL       string `json:"response_url"`
}

// ModalRequest represents the request sent to trigger a modal and received from a modal
type ModalRequest struct {
	TriggerID    string          `json:"trigger_id"`
	View         ModalDefinition `json:"view"`
	User         User            `json:"user"`
	ResponseURLS []ResponseURL   `json:"response_urls"`
}

// Request represents the incoming request body from Slack
type Request struct {
	APIAppID            string `schema:"api_app_id"`
	ChannelID           string `schema:"channel_id"`
	ChannelName         string `schema:"channel_name"`
	AppCommand          string `schema:"command"`
	IsEnterpriseInstall bool   `schema:"is_enterprise_install"`
	ResponseURL         string `schema:"response_url"`
	TeamDomain          string `schema:"team_domain"`
	TeamID              string `schema:"team_id"`
	Text                string `schema:"text"`
	Token               string `schema:"token"`
	TriggerID           string `schema:"trigger_id"`
	UserID              string `schema:"user_id"`
	UserName            string `schema:"user_name"`
	ModalPayload        string `schema:"payload"`
}
