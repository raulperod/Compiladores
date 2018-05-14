package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	df        = -1
	diag      = 47 // diagonal
	ast       = 42 // asterisk
	space     = 32
	tab       = 9
	badState  = -1
	initState = 0
)

func getTT() [193][68]int {

	var file, _ = os.Open("tabla_de_transiciones_2.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)

	var stateTransitions = [193][68]int{}
	var state = 0

	for fileScanner.Scan() {
		var line = strings.Split(fileScanner.Text(), ",")
		var symbol = 0
		for _, e := range line {
			stateTransitions[state][symbol], _ = strconv.Atoi(e)
			symbol++
		}
		state++
	}
	return stateTransitions

}

func getTokens() map[int]string {

	var file, _ = os.Open("words.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var tokens = make(map[int]string)
	// coloco los estados de los tokens identificadores
	for i := 0; i < 193; i++ {
		tokens[i] = "T_IDENT"
	}
	// coloco los demas remplazando los de ident si es necesario
	for fileScanner.Scan() {
		var line = strings.Split(fileScanner.Text(), " ")
		var state, _ = strconv.Atoi(line[1])
		var token = line[0]
		tokens[state] = token
	}
	// obteniendo los estados que son identificadores
	return tokens
}

func getSymbol(index int) int {
	switch {
	case index > 32 && index < 65: // caracteres especiales y numeros 0-9
		return index - 33 // 0 - 31
	case index > 64 && index < 91: // letras mayusculas
		return 32
	case index > 90 && index < 126: // caracteres especiales y letras minusculas
		return index - 58 // 33 - 67
	default:
		return -1
	}
}

func printToken(tokens map[int]string, state int) {
	if state != initState {
		fmt.Println(tokens[state])
	}
}

func main() {

	var file, _ = os.Open("prueba.go")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var isComment = false
	var prev = df
	var tt = getTT()
	var tokens = getTokens()
	var prevState = 0
	var state = 0

	for fileScanner.Scan() {
		for _, c := range fileScanner.Text() {
			if !isComment {
				if c == diag {
					if prev == diag {
						prev = df
						state, prevState = 0, 0
						break
					}
					prevState = state
					state = tt[state][getSymbol(int(c))]
					prev = diag
				} else if c == ast && prev == diag {
					isComment, prev = true, df
					state, prevState = 0, 0
				} else {
					if c != space && c != tab {
						prevState = state
						state = tt[state][getSymbol(int(c))]
						if state == badState {
							printToken(tokens, prevState)
							state, prevState = 0, 0
							state = tt[state][getSymbol(int(c))]
							if state == badState {
								os.Exit(1)
							}
						}
						prev = df
					} else if prev != space || prev != tab {
						prev = space
						printToken(tokens, state)
						state, prevState = 0, 0
					}
				}
			} else {
				if c == ast {
					prev = ast
				} else if c == diag && prev == ast {
					isComment, prev = false, df
				} else {
					prev = df
				}
			}
		}
		printToken(tokens, state)
		state, prevState = 0, 0
		prev = df
	}

}
