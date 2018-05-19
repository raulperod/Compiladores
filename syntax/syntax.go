package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"./lexical"
)

type Stack struct {
	length     int
	topOfStack *Element
}

type Element struct {
	state int
	value string
	next  *Element
}

func NewStack() *Stack {
	var s Stack
	return &s
}

func (s *Stack) Push(state int, value string) {
	var e Element
	e.state = state
	e.value = value

	if s.topOfStack != nil {
		e.next = s.topOfStack
		s.topOfStack = &e
	} else {
		e.next = nil
		s.topOfStack = &e
	}
	s.length++
}

func (s *Stack) Pop() *Element {
	var e *Element

	if s.topOfStack != nil {
		e = s.topOfStack
		s.topOfStack = e.next
		s.length--
		return e
	}

	return nil
}

func (s *Stack) PopN(n int) {
	for i := 0; i < n; i++ {
		s.Pop()
	}
}

func (s *Stack) PrintStack() {
	var e *Element
	e = s.topOfStack

	for e != nil {
		fmt.Println("Estado:", e.state, "Valor:", e.value)
		e = e.next
	}
}

func (s *Stack) TOS() int {
	return s.topOfStack.state
}

func GetTerminals() map[string]int {
	var file, _ = os.Open("syntax_files/terminales.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var terminals = make(map[string]int)
	var c = 0
	// se agregan los terminales en el orden en
	// que estaran en la tabla de accion
	for fileScanner.Scan() {
		var terminal = fileScanner.Text()
		terminals[terminal] = c
		c++
	}
	// obteniendo los estados que son identificadores
	return terminals
}

func GetNoTerminals() map[string]int {
	var file, _ = os.Open("syntax_files/no_terminales.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var noTerminals = make(map[string]int)
	var c = 0
	// se agregan los no terminales en el orden en
	// que estaran en la tabla de goto
	for fileScanner.Scan() {
		var noTerminal = fileScanner.Text()
		noTerminals[noTerminal] = c
		c++
	}
	// obteniendo los estados que son identificadores
	return noTerminals
}

type Rule struct {
	lenght     int
	noTerminal string
}

func GetRules() [9]Rule {
	var file, _ = os.Open("syntax_files/rules_length.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var rules = [9]Rule{}
	var i = 0
	// se agrega el tamaño de la regla en cuanto a
	// elementos que la contiene
	for fileScanner.Scan() {
		var line = strings.Split(fileScanner.Text(), ",")
		rules[i].noTerminal = line[0]
		rules[i].lenght, _ = strconv.Atoi(line[1])
		i++
	}
	// obteniendo los estados que son identificadores
	return rules
}

func GetActionTable() [228][83]int {
	var file, _ = os.Open("syntax_files/action_2.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var actionTable = [228][83]int{}
	var state = 0

	for fileScanner.Scan() {
		var line = strings.Split(fileScanner.Text(), ",")
		var terminal = 0
		for _, e := range line {
			actionTable[state][terminal], _ = strconv.Atoi(e)
			terminal++
		}
		state++
	}
	return actionTable
}

func GetGotoTable() [228][59]int {
	var file, _ = os.Open("syntax_files/goto_2.csv")
	defer file.Close()
	var fileScanner = bufio.NewScanner(file)
	var gotoTable = [228][59]int{}
	var state = 0

	for fileScanner.Scan() {
		var line = strings.Split(fileScanner.Text(), ",")
		var noTerminal = 0
		for _, e := range line {
			gotoTable[state][noTerminal], _ = strconv.Atoi(e)
			noTerminal++
		}
		state++
	}
	return gotoTable
}

func IsAccepted(action [228][83]int, stack *Stack, input *lexical.Queue, terminales map[string]int) bool {
	return action[stack.TOS()][terminales[input.TOQ()]] == 1000
}

func IsShift(action [228][83]int, stack *Stack, input *lexical.Queue, terminales map[string]int) bool {
	var tos = stack.TOS()
	var terminal = terminales[input.TOQ()]
	fmt.Println("shift:", action[tos][terminal])
	return action[tos][terminal] > 0 && action[tos][terminal] < 1000
}

func GetShift(action [228][83]int, stack *Stack, input *lexical.Queue, terminales map[string]int) (int, string) {
	var tos = stack.TOS()
	var terminal = terminales[input.TOQ()]
	return action[tos][terminal], input.TOQ()
}

func IsReduce(action [228][83]int, stack *Stack, input *lexical.Queue, terminales map[string]int) bool {
	var tos = stack.TOS()
	var terminal = terminales[input.TOQ()]
	fmt.Println("reduce:", action[tos][terminal])
	return action[tos][terminal] > 1000 && action[tos][terminal] < 2000
}

func GetReduce(action [228][83]int, stack *Stack, input *lexical.Queue, terminales map[string]int) int {
	var tos = stack.TOS()
	var terminal = terminales[input.TOQ()]
	return action[tos][terminal] - 1000
}

func SyntacticAnalysis(inputTokens *lexical.Queue) (bool, string) {
	// obtengo las tablas necesarias
	var terminales = GetTerminals()     // map
	var noTerminales = GetNoTerminals() // map
	var rules = GetRules()              // Rule
	var actionTable = GetActionTable()  // int
	var gotoTable = GetGotoTable()      // int
	var stack = NewStack()
	var input = inputTokens
	// analisis sintactico
	stack.Push(0, rules[0].noTerminal)

	for !IsAccepted(actionTable, stack, input, terminales) {
		if IsShift(actionTable, stack, input, terminales) {
			var shifti, terminal = GetShift(actionTable, stack, input, terminales)
			stack.Push(shifti, terminal) // ingresa el estado y el terminal a la pila
			input.Pop()                  // recorre al siguiente token
			stack.PrintStack()
			input.PrintQueue()
			time.Sleep(5 * time.Second)
		} else if IsReduce(actionTable, stack, input, terminales) {
			var reducei = GetReduce(actionTable, stack, input, terminales)
			var ruleiLength = rules[reducei].lenght
			var noTerminalRulei = rules[reducei].noTerminal
			var noTerminalGoto = noTerminales[noTerminalRulei]
			stack.PopN(ruleiLength)
			stack.Push(gotoTable[stack.TOS()][noTerminalGoto], noTerminalRulei)
			stack.PrintStack()
			input.PrintQueue()
			time.Sleep(5 * time.Second)
		} else {
			return false, "Error"
		}
	}

	return true, "Sin errores"
}

func main() {

	var input = lexical.LexicalAnalysis("test.txt")
	var valid, err = SyntacticAnalysis(input)
	if !valid {
		fmt.Println(err)
	} else {
		fmt.Println("La cadena es valida")
	}

}
