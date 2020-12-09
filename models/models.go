package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

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

// CodeProcessRequest represents the payload sent to the code runner lambda
type CodeProcessRequest struct {
	ResponseURL string             `json:"responseUrl"`
	Code        string             `json:"code"`
	Props       LanguageProperties `json:"props"`
}

func ImportLanguageConfig(filePath string) (LanguageConfig, error) {
	var config LanguageConfig

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path.Join(dir, filePath))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}
