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
