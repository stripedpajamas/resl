package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/stripedpajamas/resl/models"
)

func writeCodeFile(fileName string, code []byte) (string, error) {
	fmt.Println("Writing Code File")

	path := path.Join("/tmp", fileName)

	fmt.Printf("Writing file: %s \n", path)

	if err := ioutil.WriteFile(path, code, 0755); err != nil {
		return "", err
	}

	fmt.Printf("Created file to execute code: %s \n", path)

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	fmt.Println(string(dat))

	return path, nil
}

func runCode(languageConfig models.LanguageProperties) (string, error) {
	fmt.Println("Running Code")

	binary, err := exec.LookPath(languageConfig.RunCommand)
	if err != nil {
		return "", err
	}

	debugCmd := exec.Command(binary, "-v")
	debugOut, err := debugCmd.Output()
	if err != nil {
		return "", err
	}
	fmt.Println(debugOut)

	fmt.Printf("Running command `%s %s` \n", binary, languageConfig.FileName)
	runCmd := exec.Command(binary, languageConfig.FileName)
	runOut, err := runCmd.Output()
	if err != nil {
		return "", err
	}

	fmt.Printf("Command output: %s \n", runOut)

	return string(runOut), nil
}
