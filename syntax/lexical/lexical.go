package lexical

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	df           = -1
	diag         = 47 // diagonal
	ast          = 42 // asterisk
	space        = 32
	tab          = 9
	badState     = -1
	badState2    = -2
	mediumState  = 118
	mediumState2 = 119
	initState    = 0
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

func NewQueue() *Queue {
	var q Queue
	return &q
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
	var t = q.init

	for t != nil {
		fmt.Printf("| %s ", t.t_type)
		t = t.next
	}
	fmt.Println("|")
}

func (q *Queue) LastLine() int {
	return q.init.line
}

func (q *Queue) TOQ() string {
	return q.init.t_type
}

func GetTT() [193][68]int {

	var file, _ = os.Open("lexical_files/transitions_table_raw.csv")
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

func GetTokens() map[int]string {

	var file, _ = os.Open("lexical_files/words.csv")
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

func GetSymbol(index int) int {
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

func GetToken(tokens map[int]string, state, line int) *Token {
	if state != initState {
		var newToken Token
		newToken.line = line
		newToken.t_type = tokens[state]
		return &newToken
	}
	return nil
}

func LexicalAnalysis(archivo string) *Queue {
	var file, _ = os.Open(archivo)
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var isComment = false
	var prev = df
	var tt = GetTT()
	var tokens = GetTokens()
	var prevState = 0
	var state = 0
	var line = 1
	var inputTokens = NewQueue()

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
					state = tt[state][GetSymbol(int(c))]
					prev = diag
				} else if c == ast && prev == diag {
					isComment, prev = true, df
					state, prevState = 0, 0
				} else {
					if c != space && c != tab {
						prevState = state
						state = tt[state][GetSymbol(int(c))]
						if state == badState {
							inputTokens.Append(GetToken(tokens, prevState, line))
							state, prevState = 0, 0
							state = tt[state][GetSymbol(int(c))]
							if state == badState {
								fmt.Println("Error: Se leyo un cadena irreconocible en la linea:", line)
								os.Exit(1)
							}
						} else if state == badState2 {
							fmt.Println("Error: Se leyo un cadena irreconocible en la linea:", line)
							os.Exit(1)
						}
						prev = df
					} else { // si es espacio o tab
						prev = df
						if state != mediumState && state != mediumState2 {
							inputTokens.Append(GetToken(tokens, state, line))
							state, prevState = 0, 0
						}

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
		if state == mediumState || state == mediumState2 {
			fmt.Println("Error: Se leyo un cadena irreconocible en la linea:", line)
			os.Exit(1)
		}
		inputTokens.Append(GetToken(tokens, state, line))
		state, prevState = 0, 0
		prev = df
		line++
	}
	// add dollar
	inputTokens.Append(&Token{line: line, t_type: "DOLLAR"})

	return inputTokens
}
