package main

import (
	"bufio"
	"os"
)

const diag = 47 // diagonal
const ast = 42  // asterisk
const space = 32
const df = -1

func main() {

	file, _ := os.Open("prueba.go")
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	code := ""
	isComment := false
	prev := df

	for fileScanner.Scan() {
		for _, c := range fileScanner.Text() {
			if !isComment {
				if c == diag {
					if prev == diag {
						prev = df
						break
					}
					prev = diag
				} else if c == ast && prev == diag {
					isComment = true
				} else {
					if c != space {
						code += string(c)
					}
					prev = df
				}
			} else {
				if c == ast {
					prev = ast
				} else if c == diag && prev == ast {
					isComment = false
				} else {
					prev = df
				}
			}
		}
	}
}
