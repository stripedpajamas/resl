package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/stripedpajamas/resl/models"
)

func writeCodeFile(fileName string, code []byte) (string, error) {
	file := path.Join("/tmp", fileName)

	fmt.Printf("Writing file: %s \n", file)

	if err := ioutil.WriteFile(file, code, 0755); err != nil {
		return "", err
	}

	fmt.Printf("Created file to execute code: %s \n", file)

	return file, nil
}

func runCode(languageConfig models.LanguageProperties) (string, error) {
	binary, err := exec.LookPath(languageConfig.RunCommand)
	if err != nil {
		return "", err
	}

	fmt.Printf("Running command `%s %s` \n", binary, languageConfig.FileName)

	runCmd := exec.Command(binary, languageConfig.FileName)

	runOut, err := runCmd.Output()
	if err != nil {
		return "", err
	}

	return string(runOut), nil
}
