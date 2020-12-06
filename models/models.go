package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
func ParseLanguageConfig(filePath string) (LanguageConfig, error) {
	var config LanguageConfig

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return config, err
	}

	data, err := ioutil.ReadFile(path.Join(dir, filePath))
	if err != nil {
		fmt.Println(err)
		return config, err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
