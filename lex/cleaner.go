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
const df = -1

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

	file, _ := os.Open("prueba.go")
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	code := ""
	isComment := false
	areSpaces := true
	prev := df

	for fileScanner.Scan() {
		areSpaces, prev = true, df
		for _, c := range fileScanner.Text() {
			if !isComment {
				if c == diag {
					if prev == diag {
						prev = df
						break
					}
					prev, areSpaces = diag, false
				} else if c == ast && prev == diag {
					isComment = true
				} else {
					if c != space {
						code += string(c)
						areSpaces = false
					} else if prev != space {
						code += " "
						prev = space
					}
				}
			} else {
				if c == ast {
					prev, areSpaces = ast, false
				} else if c == diag && prev == ast {
					isComment = false
				} else {
					prev = df
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
