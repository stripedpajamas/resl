package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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

// LanguageProperties represents properties for running each supported language
type LanguageProperties struct {
	Name           string `json:"langName"`
	Extension      string `json:"extension"`
	Placeholder    string `json:"placeholder"`
	FileName       string `json:"fileName"`
	RunCommand     string `json:"runCmd"`
	CompileCommand string `json:"compileCmd"`
}

// LanguageConfig represents the model matching the languages.json file
type LanguageConfig map[string]LanguageProperties

// ParseLanguageConfig parses the languages file into a LanguageConfig model
func ParseLanguageConfig() (LanguageConfig, error) {
	data, err := ioutil.ReadFile("../../languages.json")
	if err != nil {
		fmt.Println(err)
		return LanguageConfig{}, err
	}

	var config LanguageConfig

	if err = json.Unmarshal(data, &config); err != nil {
		return LanguageConfig{}, err
	}

	return config, nil
}
