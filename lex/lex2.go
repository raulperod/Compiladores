package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	df    = -1
	diag  = 47 // diagonal
	ast   = 42 // asterisk
	space = 32
)

func isLetter(c int) bool {
	return c > 64 && c < 91 || c > 96 && c < 123
}

func isNumber(c int) bool {
	return c > 47 && c < 58
}

func isSpecial(c int) bool {
	return !(isNumber(c) || isLetter(c))
}

func convertLetterToken(cadena string) {
	switch cadena {
	case "":
		return
	case "break":
		fmt.Println("T_BREAK")
	case "case":
		fmt.Println("T_CASE")
	case "chan":
		fmt.Println("T_CHAN")
	case "const":
		fmt.Println("T_CONST")
	case "continue":
		fmt.Println("T_CONTINUE")
	case "defer":
		fmt.Println("T_DEFER")
	case "default":
		fmt.Println("T_DEFAULT")
	case "else":
		fmt.Println("T_ELSE")
	case "fallthrough":
		fmt.Println("T_FALLTHROUGH")
	case "for":
		fmt.Println("T_FOR")
	case "func":
		fmt.Println("T_FUNC")
	case "go":
		fmt.Println("T_GO")
	case "goto":
		fmt.Println("T_GOTO")
	case "if":
		fmt.Println("T_IF")
	case "import":
		fmt.Println("T_IMPORT")
	case "interface":
		fmt.Println("T_INTERFACE")
	case "map":
		fmt.Println("T_MAP")
	case "package":
		fmt.Println("T_PACKAGE")
	case "range":
		fmt.Println("T_RANGE")
	case "return":
		fmt.Println("T_RETURN")
	case "select":
		fmt.Println("T_SELECT")
	case "struct":
		fmt.Println("T_STRUCT")
	case "switch":
		fmt.Println("T_SWITCH")
	case "type":
		fmt.Println("T_TYPE")
	case "var":
		fmt.Println("T_VAR")
	default:
		fmt.Println("T_ID")
	}
}

func convertSpecialToken(cadena string) {
	switch cadena {
	case "":
		return
	case "+":
		fmt.Println("T_PLUS")
	case "++":
		fmt.Println("T_PLUS_PLUS")
	case "+=":
		fmt.Println("T_PLUS_EQ")
	case "-":
		fmt.Println("T_SUB")
	case "--":
		fmt.Println("T_SUB_SUB")
	case "-=":
		fmt.Println("T_SUB_EQ")
	case "=":
		fmt.Println("T_EQ")
	case "==":
		fmt.Println("T_EQ_EQ")
	case ":=":
		fmt.Println("T_COL_EQ")
	case "!":
		fmt.Println("T_NOT")
	case "!=":
		fmt.Println("T_NOT_EQ")
	case ">":
		fmt.Println("T_GT")
	case ">=":
		fmt.Println("T_GT_EQ")
	case ">>":
		fmt.Println("T_GT_GT")
	case ">>=":
		fmt.Println("T_GT_GT_EQ")
	case "<":
		fmt.Println("T_LT")
	case "<=":
		fmt.Println("T_LT_EQ")
	case "*":
		fmt.Println("T_MULT")
	case "*=":
		fmt.Println("T_MULT_EQ")
	case "/":
		fmt.Println("T_SLASH")
	case "/=":
		fmt.Println("T_SLASH_EQ")
	case "%":
		fmt.Println("T_MOD")
	case "%=":
		fmt.Println("T_MOD_EQ")
	case "&":
		fmt.Println("T_AND")
	case "&&":
		fmt.Println("T_AND_AND")
	case "&=":
		fmt.Println("T_AND_EQ")
	case "|":
		fmt.Println("T_OR")
	case "||":
		fmt.Println("T_OR_OR")
	case "|=":
		fmt.Println("T_OR_EQ")
	case "(":
		fmt.Println("T_PTS_LEFT")
	case ")":
		fmt.Println("T_PTS_RIGHT")
	case "()":
		fmt.Println("T_PTS_LEFT")
		fmt.Println("T_PTS_RIGHT")
	case "[":
		fmt.Println("T_BKT_LEFT")
	case "]":
		fmt.Println("T_BKT_RIGHT")
	case "[]":
		fmt.Println("T_BKT_LEFT")
		fmt.Println("T_BKT_RIGHT")
	case "{":
		fmt.Println("T_CBKT_LEFT")
	case "}":
		fmt.Println("T_CBKT_RIGHT")
	case "{}":
		fmt.Println("T_CBKT_LEFT")
		fmt.Println("T_CBKT_RIGHT")
	case ";":
		fmt.Println("T_SEMICOL")
	case ",":
		fmt.Println("T_COMMA")
	case ".":
		fmt.Println("T_DOT")
	default:
		return
	}
}

func convertToken(code string, is int) {
	switch is {
	case 0:
		convertLetterToken(code)
	case 1:
		fmt.Println("T_INT")
	case 2:
		convertSpecialToken(code)
	default:
		return
	}
}

func main() {

	file, _ := os.Open("fuente.go")
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	isComment := false
	prev := df
	code := ""
	is := 0

	for fileScanner.Scan() {
		prev = df
		for _, c := range fileScanner.Text() {
			if !isComment {
				if c == diag {
					if prev == diag {
						prev, code = df, ""
						break
					}
					prev = diag
					code += "/"
				} else if c == ast && prev == diag {
					isComment, prev = true, df
				} else {
					if c != space {
						if isLetter(int(c)) && is != 0 {
							convertToken(code, is)
							code, is = "", 0
						} else if isNumber(int(c)) && is != 1 {
							convertToken(code, is)
							code, is = "", 1
						} else if isSpecial(int(c)) && is != 2 {
							convertToken(code, is)
							code, is = "", 2
						}
						prev = df
						code += string(c)
					} else if prev != space {
						convertToken(code, is)
						prev, code, is = space, "", 0
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
		convertToken(code, is)
		code, is = "", 0
	}

}
