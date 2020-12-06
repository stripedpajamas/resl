package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/stripedpajamas/resl/models"
)

func writeCodeFile(fileName string, code []byte) error {
	fmt.Printf("Writing file: %s", fileName)

	if err := ioutil.WriteFile(fileName, code, 0755); err != nil {
		return err
	}

	fmt.Printf("Created file to execute code: %s", fileName)

	return nil
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
