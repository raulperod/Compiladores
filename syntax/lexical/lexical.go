package lexical

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	df                  = -1
	slash               = 47 // diagonal
	backslach           = 92 // diagonal inverted
	ast                 = 42 // asterisk
	space               = 32
	tab                 = 9
	badState            = -1
	badState2           = -2
	mediumState         = 118
	mediumState2        = 119
	mediumState3        = 194
	commentState        = 195
	commentState2       = 198
	commentMediumState  = 196
	commentMediumState2 = 197
	lineState           = 2
	initState           = 0
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
		fmt.Printf("| %s, %d ", t.t_type, t.line)
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

func GetTT() [199][68]int {
	var file, _ = os.Open("./syntax/lexical_files/transitions_table_raw.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var stateTransitions = [199][68]int{}
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

func GetTL() [3][68]int {
	var file, _ = os.Open("./syntax/lexical_files/line_raw.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var stateTransitions = [3][68]int{}
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
	var file, _ = os.Open("./syntax/lexical_files/words.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var tokens = make(map[int]string)
	// coloco los estados de los tokens identificadores
	for i := 0; i < 114; i++ {
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
	if state != initState && state != commentState2 {
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
	var tt = GetTT()
	var tokens = GetTokens()
	var prevState = 0
	var state = 0
	var line = 1
	var inputTokens = NewQueue()

	for fileScanner.Scan() {
		for _, c := range fileScanner.Text() {
			// si es un comentario en linea se pasa a la siguiente
			if state == commentState {
				break
			}
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
			} else { // si es espacio o tab
				if state != mediumState && state != mediumState2 && state != commentMediumState && state != commentMediumState2 {
					inputTokens.Append(GetToken(tokens, state, line))
					state, prevState = 0, 0
				}
			}
		}
		// cuando se termina de leer una linea
		if state == mediumState || state == mediumState2 { // si estaba leyendo un string y no se termino
			fmt.Println("Error: Se leyo un cadena irreconocible en la linea:", line)
			os.Exit(1)
		} else if state == commentState { // si es comentario de una linea
			state, prevState = 0, 0
			line++
		} else if state == commentMediumState || state == commentMediumState2 { // si esta en un comentario multilinea
			line++
		} else { // caso "normal"
			inputTokens.Append(GetToken(tokens, state, line))
			state, prevState = 0, 0
			line++
		}
	}
	// add dollar
	inputTokens.Append(&Token{line: line, t_type: "DOLLAR"})
	return inputTokens
}

func LexicalAnalysisForWeb(text string) (*Queue, string, bool) {
	var tt = GetTT()
	var tl = GetTL()
	var tokens = GetTokens()
	var prevState = 0
	var state = 0
	var stateLine = 0
	var line = 1
	var inputTokens = NewQueue()

	for idx, c := range text {
		if idx == 0 || idx == len(text)-1 {
			continue
		}
		if c != space && c != tab {
			prevState = state
			state = tt[state][GetSymbol(int(c))]
			stateLine = tl[stateLine][GetSymbol(int(c))]
			//fmt.Println(string(c), state, stateLine)
			if stateLine == lineState {
				// cuando se termina de leer una linea
				if prevState == mediumState3 { // si estaba leyendo un string y no se termino
					return inputTokens, "Error: Se leyo un cadena irreconocible en la linea: " + strconv.Itoa(line), true
				} else if state == commentMediumState || state == commentMediumState2 { // si esta en un comentario multilinea
					stateLine = 0
					line++
				} else {
					state, prevState, stateLine = 0, 0, 0
					line++
				}
			} else if state == badState {
				inputTokens.Append(GetToken(tokens, prevState, line))
				state, prevState = 0, 0
				state = tt[state][GetSymbol(int(c))]
				if state == badState {
					return inputTokens, "Error: Se leyo un cadena irreconocible en la linea: " + strconv.Itoa(line), true
				}
			} else if state == badState2 {
				return inputTokens, "Error: Se leyo un cadena irreconocible en la linea: " + strconv.Itoa(line), true
			}
		} else { // si es espacio o tab
			if state != mediumState && state != mediumState2 && state != commentState && state != commentMediumState && state != commentMediumState2 {
				inputTokens.Append(GetToken(tokens, state, line))
				state, prevState = 0, 0
			}
		}
		if idx == len(text)-2 {
			inputTokens.Append(GetToken(tokens, state, line))
		}
	}
	// add dollar
	inputTokens.Append(&Token{line: line, t_type: "DOLLAR"})
	return inputTokens, "", false
}
