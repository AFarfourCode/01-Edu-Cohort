package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <string> OR <filename>")
		return
	}

	inputArgs := os.Args[1]
	// read arguments or txt file
	if fileContent, err := os.ReadFile(inputArgs); err == nil {
		inputArgs = string(fileContent)
	}

	if inputArgs == "\\n" {
		fmt.Println()
		return
	}
	// read standard.txt
	readinput, err := os.ReadFile("standard.txt")
	if err != nil {
		fmt.Print("error of reading file")
		return
	}

	// standard.txt edit
	editInput := strings.ReplaceAll(string(readinput), "\r\n", "\n")
	lines := strings.Split(editInput, "\n")

	// input edit
	inputArgs = strings.ReplaceAll(inputArgs, "\r\n", "\n") // if this argument is a text file
	inputArgs = strings.ReplaceAll(inputArgs, "\\n", "\n")  // if this argument from terminal direct

	argslines := strings.Split(inputArgs, "\n")

	for _, word := range argslines {
		if word == "" {
			fmt.Println()
			continue
		}
		for i := 0; i < 8; i++ {
			for _, ch := range word {
				startline := int(ch-32)*9 + 1 + i

				if startline >= 0 && startline < len(lines) {
					fmt.Print(lines[startline])
				}
			}
			fmt.Println()
		}
	}
}
