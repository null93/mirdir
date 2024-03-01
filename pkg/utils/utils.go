package utils

import (
	"fmt"
	"os"
	"strings"
)

func IsDirectory(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}
	mode := fileInfo.Mode()
	return mode.IsDir()
}

func Exists(path string) bool {
	_, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return true
}

func PromptOverwrite(path string, isDir bool) bool {
	response := ""
	postfix := ""
	if isDir {
		postfix = "/"
	}
	for response != "y" && response != "n" {
		fmt.Printf("%q already exists. Overwrite? (y/n): ", path+postfix)
		fmt.Scanln(&response)
	}
	return response == "y"
}

func GetEnvironmentalVars() map[string]string {
	envVars := map[string]string{}
	for _, envVar := range os.Environ() {
		parts := strings.Split(envVar, "=")
		envVars[parts[0]] = parts[1]
	}
	return envVars
}

func ExitWithError(code int, msg string, err error, printErr bool) {
	fmt.Printf("Error: %s\n", msg)
	if printErr {
		fmt.Println(err)
	}
	os.Exit(code)
}
