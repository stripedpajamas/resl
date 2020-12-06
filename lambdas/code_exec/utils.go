package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/stripedpajamas/resl/models"
)

func writeCodeFile(fileName string, code []byte) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filePath := path.Join(dir, "tmp", fileName)
	err = ioutil.WriteFile(filePath, code, 0755)
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
