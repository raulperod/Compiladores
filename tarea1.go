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
	if err != nil {
		return err
	}
	defer fo.Close()

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
	prev := -1

	for fileScanner.Scan() {
		for _, c := range fileScanner.Text() {
			if !isComment {
				if c == diag {
					if prev == diag {
						prev = -1
						break
					}
					prev = diag
				} else if c == ast && prev == diag {
					isComment = true
				} else {
					if c != space {
						code += string(c)
					}
					prev = -1
				}
			} else {
				if c == ast {
					prev = ast
				} else if c == diag && prev == ast {
					isComment = false
				} else {
					prev = -1
				}
			}
		}
	}

	if err := WriteStringToFile("fuente.go", code); err != nil {
		panic(err)
	}

}
