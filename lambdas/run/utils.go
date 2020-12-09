package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/stripedpajamas/resl/models"
)

func getCodePayloadFromRequestBody(body RequestBody) ([]byte, error) {
	var payload = parseText(body.Text)

	if valid, err := validateLanguage(payload); !valid {
		return nil, err
	}

	parseCode(&payload)

	fmt.Printf("Parsed Code: %s", payload.Code)
	fmt.Printf("Parsed Language: %s", payload.Language)

	return json.Marshal(payload)
}

func parseCode(payload *models.CodeProcessRequest) {
	payload.Code = strings.ReplaceAll(payload.Code, "&amp;", "&")
	payload.Code = strings.ReplaceAll(payload.Code, "&lt;", "<")
	payload.Code = strings.ReplaceAll(payload.Code, "&gt;", ">")

	code := payload.Code

	i := 0
	j := len(payload.Code) - 1

	for i <= j && code[i] == '`' && code[j] == code[i] {
		i++
		j--
	}

	if i == 1 || i == 3 {
		payload.Code = code[i : j+1]
	}
}

func validateLanguage(payload models.CodeProcessRequest) (bool, error) {
	if _, ok := languageConfig[payload.Language]; !ok {
		return false, errors.New("language not supported")
	}

	return true, nil
}

func parseText(text string) models.CodeProcessRequest {
	text = strings.Trim(text, " ")

	for idx, c := range text {
		if c == ' ' {
			return models.CodeProcessRequest{
				Code:     text[idx+1:],
				Language: text[0:idx],
			}
		}
	}

	return models.CodeProcessRequest{
		Code:     "",
		Language: text,
	}
}
