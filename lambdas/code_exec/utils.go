package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/stripedpajamas/resl/models"
)

func writeCodeFile(fileName string, code []byte) (string, error) {
	filePath := path.Join("/tmp", fileName)
	fmt.Printf("Writing file: %s", filePath)

	err := ioutil.WriteFile(filePath, code, 0755)
	if err != nil {
		return "", err
	}

	fmt.Printf("Created file to execute code: %s", filePath)

	return filePath, nil
}

func runCode(fileName string, languageConfig models.LanguageProperties) (string, error) {
	cmd := fmt.Sprintf(languageConfig.RunCommand, fileName)
	fmt.Printf("Running command %s", cmd)

	runCmd := exec.Command(cmd)

	runOut, err := runCmd.Output()
	if err != nil {
		return "", err
	}

	return string(runOut), nil
}
