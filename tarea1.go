package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const diag = 47 // diagonal
const ast = 42  // asterisk
const space = 32

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	defer fo.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(fo, strings.NewReader(s))

	if err != nil {
		return err
	}

	return nil
}

func main() {

	file, _ := os.Open("tarea1.cpp")
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	code := ""
	isComment := false
	areSpaces := true
	prev := -1

	for fileScanner.Scan() {
		areSpaces = true
		for _, c := range fileScanner.Text() {
			if !isComment {
				if c == diag {
					if prev == diag {
						prev = -1
						break
					}
					prev, areSpaces = diag, false
				} else if c == ast && prev == diag {
					isComment = true
				} else {
					if c != space {
						code += string(c)
						areSpaces = false
					}
					prev = -1
				}
			} else {
				if c == ast {
					prev, areSpaces = ast, false
				} else if c == diag && prev == ast {
					isComment = false
				} else {
					prev = -1
				}
			}
		}
		if !areSpaces {
			code += " "
		}
	}

	if err := WriteStringToFile("fuente", code); err != nil {
		panic(err)
	}

}
