package models

// CodeProcessRequest represents the payload sent to the code runner lambda
type CodeProcessRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

// CodeProcessResponse represents the payload returned by the code runner lambda
type CodeProcessResponse struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exitCode"`
}

// LanguageConfig represents properties for running each supported language
type LanguageConfig struct {
	Name        string `json:"langName"`
	Extension   string `json:"extension"`
	Placeholder string `json:"placeholder"`
	FileName    string `json:"fileName"`
	Command     string `json:"runCmd"`
}
