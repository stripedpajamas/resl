package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

type LanguageProperties struct {
	Name           string `json:"langName"`
	Extension      string `json:"extension"`
	Placeholder    string `json:"placeholder"`
	FileName       string `json:"fileName"`
	RunCommand     string `json:"runCmd"`
	CompileCommand string `json:"compileCmd"`
}

func writeCodeFile(fileName string, code []byte) error {
	fmt.Printf("Writing file: %s \n", fileName)

	if err := ioutil.WriteFile(fileName, code, 0755); err != nil {
		return err
	}

	fmt.Printf("Created file to execute code: %s \n", fileName)

	return nil
}

func runCode(languageConfig LanguageProperties) (string, error) {
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

func main() {

	binary, lookErr := exec.LookPath("node")
	if lookErr != nil {
		panic(lookErr)
	}

	fmt.Println(binary)

	if err := writeCodeFile("test.js", []byte("console.log('hello world')")); err != nil {
		panic(err)
	}

	props := LanguageProperties{
		Name:       "JavaScript",
		RunCommand: "node",
		FileName:   "test.js",
	}

	outp, err := runCode(props)
	if err != nil {
		panic(err)
	}

	fmt.Println(outp)
}
