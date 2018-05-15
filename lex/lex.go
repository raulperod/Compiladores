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

type Token struct {
	line   int
	t_type string
	next   *Token
}

type Queue struct {
	init, final *Token
	length      int
}

func newQueue() Queue {
	var q Queue
	q.init = nil
	q.final = nil
	q.length = 0
	return q
}

func (q *Queue) Append(token *Token) {
	if token == nil {
		return
	}

	var t *Token
	t = token

	if q.init != nil {
		t.next = nil
		q.final.next = t
		q.final = t
	} else {
		t.next = nil
		q.init = t
		q.final = t
	}

	q.length++
}

func (q *Queue) Pop() *Token {
	var t *Token

	if q.init != nil {
		t = q.init
		q.init = t.next
		q.length--
		if q.length == 0 {
			q.final = nil
			return nil
		}
		return t
	}

	return nil
}

func (q *Queue) PrintQueue() {
	var t *Token
	t = q.init

	for t != nil {
		fmt.Println("Token:", t.t_type, "Linea:", t.line)
		t = t.next
	}
}

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

func getToken(tokens map[int]string, state, line int) *Token {
	if state != initState {
		var newToken Token
		newToken.line = line
		newToken.t_type = tokens[state]
		return &newToken
	}
	return nil
}

func main() {

	var file, _ = os.Open("test.go")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var isComment = false
	var prev = df
	var tt = getTT()
	var tokens = getTokens()
	var prevState = 0
	var state = 0
	var line = 1
	var inputTokens = newQueue()

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
							inputTokens.Append(getToken(tokens, prevState, line))
							state, prevState = 0, 0
							state = tt[state][getSymbol(int(c))]
							if state == badState {
								os.Exit(1)
							}
						}
						prev = df
					} else if prev != space || prev != tab {
						prev = space
						inputTokens.Append(getToken(tokens, state, line))
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
		inputTokens.Append(getToken(tokens, state, line))
		state, prevState = 0, 0
		prev = df
		line++
	}

	inputTokens.PrintQueue()
}
