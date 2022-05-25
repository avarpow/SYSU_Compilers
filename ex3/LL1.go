package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var CHARCAST = map[string]uint8{
	`E'`: 'R',
	`T'`: 'Y',
}

const debug = true

type DebugLevel int

const debugLevel = INFO
const (
	DEBUG DebugLevel = iota
	INFO
	WARNNING
	ERROR
)

func unique(uint8Slice []uint8) []uint8 {
	keys := make(map[uint8]bool)
	list := []uint8{}
	for _, entry := range uint8Slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func debugPrintf(level DebugLevel, format string, a ...interface{}) {
	if level >= debugLevel {
		fmt.Printf(format, a...)
	}
}
func debugPrint(level DebugLevel, format string) {
	if level >= debugLevel {
		fmt.Print(format)
	}
}

type Token []uint8
type GrammarLL1 struct {
	grammar       map[uint8]([]Token)
	terminals     []uint8
	nonTerminals  []uint8
	unfoldGrammar [](map[uint8](Token))
	first         map[uint8]([]uint8)
	follow        map[uint8]([]uint8)
	parseTable    map[uint8](map[uint8]([]uint8))
	QTs           [][4]uint8
	ready         bool
}

func (Grammar *GrammarLL1) buildGrammar(grammar_filename string) {
	Grammar.grammar = make(map[uint8]([]Token))
	Grammar.readGrammarFromFile(grammar_filename)
	Grammar.genTerminalAndNonterminal()
	Grammar.genUnfoldGrammar()
	Grammar.genFirst()
	Grammar.genFollow()
	Grammar.printFirstFollow()
	Grammar.genParseTable()
	Grammar.QTs = make([][4]uint8, 0)
	Grammar.ready = true
}
func (Grammar *GrammarLL1) printFirstFollow() {
	level := INFO
	debugPrintf(level, " %10s   %10s\n", "first", "follow")
	for key, _ := range Grammar.grammar {
		debugPrintf(level, "%c -> ", key)
		firstStr := ""
		followStr := ""
		for _, token := range Grammar.first[key] {
			firstStr += fmt.Sprintf("%c ", token)
		}
		for _, token := range Grammar.follow[key] {
			followStr += fmt.Sprintf("%c ", token)
		}
		debugPrintf(level, " %-10s  %-10s\n", firstStr, followStr)
	}
}

func printGrammar(grammar map[uint8]([]Token)) {
	level := INFO
	for key, value := range grammar {
		debugPrintf(level, "%c -> ", key)
		for _, token := range value {
			debugPrintf(level, "%s | ", token)
		}
		debugPrint(level, "\n")
	}
	debugPrint(level, "\n")
}
func printUnfoldGrammar(unfoldGrammar [](map[uint8](Token))) {
	level := DEBUG
	for _, value := range unfoldGrammar {
		for key, token := range value {
			debugPrintf(level, "%c -> %s", key, token)
			debugPrint(level, "\n")
		}
	}
	debugPrint(level, "\n")
}
func (Grammar *GrammarLL1) readGrammarFromFile(grammar_filename string) {
	//declear a empty map from uint8 to slice uint8
	//read file
	b, err := ioutil.ReadFile(grammar_filename)
	if err != nil {
		log.Print(err)
	}
	//[]byte to string
	s := string(b)
	//remove space
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	for k, v := range CHARCAST {
		s = strings.Replace(s, k, string(v), -1)
	}
	//split by line
	lines := strings.Split(string(s), "\n")
	for _, line := range lines {
		//split by ->
		tokens := strings.Split(line, "->")
		//check tokens length
		if len(tokens) != 2 {
			log.Print("Error: grammar file format error")
			os.Exit(1)
		}
		//split by |
		split_tokens := strings.Split(tokens[1], "|")
		if len(tokens[0]) != 1 {
			log.Print("Error: grammar file format error")
			os.Exit(1)
		}
		//check if key in grammar
		key := tokens[0][0]
		if _, ok := Grammar.grammar[key]; !ok {
			Grammar.grammar[key] = make([]Token, 0)
		}
		for _, token := range split_tokens {
			Grammar.grammar[key] = append(Grammar.grammar[key], Token(token))
		}

	}
	printGrammar(Grammar.grammar)
}
func printSlice(level DebugLevel, slice []uint8) {
	for _, v := range slice {
		debugPrintf(level, "%c ", v)
	}
	debugPrint(level, "\n")
}
func (Grammar *GrammarLL1) genTerminalAndNonterminal() {
	level := INFO
	for key, value := range Grammar.grammar {
		Grammar.nonTerminals = append(Grammar.nonTerminals, key)
		for _, token := range value {
			for _, char := range token {
				if isTerminal(char) {
					Grammar.terminals = append(Grammar.terminals, char)
				}
			}
		}
	}
	Grammar.terminals = append(Grammar.terminals, '#')
	Grammar.terminals = unique(Grammar.terminals)
	Grammar.nonTerminals = unique(Grammar.nonTerminals)
	debugPrint(level, "print terminal\n")
	printSlice(level, Grammar.terminals)
	debugPrint(level, "print nonTerminal\n")
	printSlice(level, Grammar.nonTerminals)
	debugPrint(level, "\n")
}
func (Grammar *GrammarLL1) genUnfoldGrammar() {
	Grammar.unfoldGrammar = make([](map[uint8](Token)), 0)
	for key, value := range Grammar.grammar {
		for _, token := range value {
			Grammar.unfoldGrammar = append(Grammar.unfoldGrammar, map[uint8](Token){key: token})
		}
	}
	printUnfoldGrammar(Grammar.unfoldGrammar)

}

// type Node struct {
// 	in   uint8
// 	out  uint8
// 	next []uint8
// }
// type Graph struct {
// 	nodes map[uint8]Node
// }
func isNumber(token uint8) bool {
	return token >= '0' && token <= '9' || token >= 'a' && token <= 'z' && token != 'e'
}
func isNonTerminal(token uint8) bool {
	return token >= 'A' && token <= 'Z'
}
func isEmptyToken(token uint8) bool {
	return token == 'e'
}
func isTerminal(token uint8) bool {
	return !(isNonTerminal(token) || isEmptyToken(token))
}

func (Grammar *GrammarLL1) __buildFirst(key uint8, buildOK *map[uint8]bool) {
	level := DEBUG
	debugPrintf(level, "build first %c\n", key)
	for _, token := range Grammar.grammar[key] {
		firstChar := token[0]
		if isTerminal(firstChar) {
			//是终结符,添加进入first
			Grammar.first[key] = append(Grammar.first[key], firstChar)
			debugPrintf(level, "[%c] add %c  terminal\n", key, firstChar)
		} else if isEmptyToken(firstChar) {
			//是空,添加进入first
			debugPrintf(level, "[%c] add %c  empty\n", key, firstChar)
			Grammar.first[key] = append(Grammar.first[key], firstChar)
		} else if isNonTerminal(firstChar) {
			if !(*buildOK)[firstChar] {
				Grammar.__buildFirst(firstChar, buildOK)
			}
			debugPrintf(level, "[%c] add %c  firstChar set\n", key, firstChar)
			for _, first := range Grammar.first[firstChar] {
				//将该非终结符的first集合添加到该非终结符的first集合中
				Grammar.first[key] = append(Grammar.first[key], first)
				debugPrintf(level, "\r[%c] add %c  first set\n", key, first)

			}
		} else {
			debugPrintf(ERROR, "Wrong token: %c\n", firstChar)
			os.Exit(1)
		}
	}
	Grammar.first[key] = unique(Grammar.first[key])
}
func printFirst(first map[uint8][]uint8) {
	level := INFO
	debugPrintf(level, "first list\n")
	for key, value := range first {
		debugPrintf(level, "%c -> ", key)
		for _, token := range value {
			debugPrintf(level, "%c ", token)
		}
		debugPrint(level, "\n")
	}
}

func (Grammar *GrammarLL1) genFirst() {
	Grammar.first = make(map[uint8]([]uint8))
	buildOK := make(map[uint8]bool)
	for key, _ := range Grammar.grammar {
		buildOK[key] = false
	}
	for key, _ := range Grammar.grammar {
		if !buildOK[key] {
			Grammar.__buildFirst(key, &buildOK)
			buildOK[key] = true
		}
	}
	printFirst(Grammar.first)
}
func (Grammar *GrammarLL1) __buildFollow(key uint8, buildOK *map[uint8]bool) {
	level := DEBUG
	if key == 'S' {
		Grammar.follow[key] = append(Grammar.follow[key], '#')
		return
	}
	debugPrintf(level, "build follow %c\n", key)
	for k, value := range Grammar.grammar {
		for _, token := range value {
			index := strings.Index(string(token), string(key))
			if index != -1 && index != len(token)-1 {
				debugPrintf(level, "find %c in token %s\n", key, token)
				next := token[index+1]
				debugPrintf(level, "next = %c\n", next)
				if isTerminal(next) {
					Grammar.follow[key] = append(Grammar.follow[key], next)
					debugPrintf(level, "[%c] add %c  terminal\n", key, next)
				} else if isNonTerminal(next) {
					if !(*buildOK)[next] {
						Grammar.__buildFollow(next, buildOK)
					}
					debugPrintf(level, "[%c] add %c  follow set\n", key, next)
					for _, follow := range Grammar.first[next] {
						if follow == 'e' {
							if !(*buildOK)[next] {
								Grammar.__buildFollow(next, buildOK)
							}
							Grammar.follow[key] = append(Grammar.follow[key], Grammar.follow[next]...)
							continue
						}
						Grammar.follow[key] = append(Grammar.follow[key], follow)
						debugPrintf(level, "\r[%c] add %c  follow set\n", key, follow)

					}
				} else {
					debugPrintf(ERROR, "Wrong token: %c\n", next)
					os.Exit(1)
				}
			} else if index == len(token)-1 {
				if k == 'S' {
					Grammar.follow[key] = append(Grammar.follow[key], '#')
				} else if k != token[len(token)-1] {
					if !(*buildOK)[k] {
						Grammar.__buildFollow(k, buildOK)
					}
					debugPrintf(level, "[%c] at last add first [%c]\n", key, k)
					for _, follow := range Grammar.follow[k] {
						Grammar.follow[key] = append(Grammar.follow[key], follow)
						debugPrintf(level, "\r[%c] add %c  first set\n", key, follow)
					}
				}
			}
		}
	}
	Grammar.follow[key] = unique(Grammar.follow[key])
}
func printFollow(follow map[uint8][]uint8) {
	level := INFO
	debugPrint(level, "\nfollow list\n")
	for key, value := range follow {
		debugPrintf(level, "%c -> ", key)
		for _, token := range value {
			debugPrintf(level, "%c ", token)
		}
		debugPrint(level, "\n")
	}
	debugPrint(level, "\n")
}
func (Grammar *GrammarLL1) genFollow() {
	level := DEBUG
	Grammar.follow = make(map[uint8]([]uint8))
	buildOK := make(map[uint8]bool)
	for key, _ := range Grammar.grammar {
		buildOK[key] = false
	}
	for key, _ := range Grammar.grammar {
		if !buildOK[key] {
			Grammar.__buildFollow(key, &buildOK)
			debugPrintf(level, "build follow %c ok\n", key)
			buildOK[key] = true
		}
	}
	printFollow(Grammar.follow)
}
func (Grammar *GrammarLL1) printParseTable() {
	level := INFO
	debugPrint(level, "\n      ParseTable\n")
	title := "      "
	for _, t := range Grammar.terminals {
		title += fmt.Sprintf("%-6c", t)
	}
	debugPrintf(level, "%s\n", title)
	for _, nt := range Grammar.nonTerminals {
		debugPrintf(level, "%c     ", nt)
		printStr := ""
		for _, t := range Grammar.terminals {
			printStr += fmt.Sprintf("%-6s", Grammar.parseTable[nt][t])
		}
		debugPrintf(level, "%s\n", printStr)
	}
}

func (Grammar *GrammarLL1) genParseTable() {
	level := INFO
	debugPrint(level, "\n      genParseTable\n")
	Grammar.parseTable = make(map[uint8](map[uint8][]uint8))
	for key, value := range Grammar.grammar {
		Grammar.parseTable[key] = make(map[uint8][]uint8)
		for _, token := range value {
			if token[0] == 'e' {
				for _, follow := range Grammar.follow[key] {
					Grammar.parseTable[key][follow] = token
				}
			} else if isTerminal(token[0]) {
				Grammar.parseTable[key][token[0]] = token
			} else if isNonTerminal(token[0]) {
				for _, first := range Grammar.first[token[0]] {
					Grammar.parseTable[key][first] = token
				}
			} else {
				debugPrintf(ERROR, "Wrong token: %s\n", token)
				os.Exit(1)
			}
		}
	}
	Grammar.printParseTable()
}
func stringfySYN(SYN []uint8) string {
	str := ""
	for _, v := range SYN {
		if v >= 128 {
			v = v - 128
			if v == '+' || v == '-' || v == '*' || v == '/' {
				str += "GEQ(" + string(v) + ")"
			} else {
				str += "PUSH(" + string(v) + ")"
			}
		} else {
			str += string(v)
		}
	}
	return str
}
func stringfySEM(SEM []uint8) string {
	ret := ""
	for _, v := range SEM {
		if v > 128 {
			ret += "t"
			ret += string(v - 128)
		} else {
			ret += string(v)
		}
		ret += " "
	}
	return ret
}
func printState(stack []uint8, finishStack []uint8, SEM_stack []uint8, QT [4]uint8, expression string, index int) {
	level := INFO
	leftExppression := expression[:index]
	rightExppression := expression[index:]
	debugPrintf(level, "%-10s%5s%-20s%5s%-10s%5s%-10s\n", "matched", "", "matching", "", "SEM_stack", "", "QT")
	debugPrintf(level, "%10s%5s%-20s%5s%-10s%5s%-10s\n", leftExppression, "", rightExppression, "", stringfySEM(SEM_stack), "", stringfyQT(QT))
	copystack := make([]uint8, len(stack))
	copy(copystack, stack)
	//reverse stack
	for i, j := 0, len(copystack)-1; i < j; i, j = i+1, j-1 {
		copystack[i], copystack[j] = copystack[j], copystack[i]
	}

	debugPrintf(level, "%10s%5s%-10s\n\n", finishStack, "", stringfySYN(copystack))
}
func printErrorState(stack []uint8, finishStack []uint8, expression string, index int) {
	level := ERROR
	leftExppression := expression[:index]
	rightExppression := expression[index:]
	debugPrintf(level, "\nError state dump\n")
	debugPrintf(level, "%-10s%5s%-10s\n", "matched", "", "matching")
	debugPrintf(level, "%10s%5s%-10s\n", leftExppression, "", rightExppression)
	copystack := make([]uint8, len(stack))
	copy(copystack, stack)
	//reverse stack
	for i, j := 0, len(copystack)-1; i < j; i, j = i+1, j-1 {
		copystack[i], copystack[j] = copystack[j], copystack[i]
	}
	debugPrintf(level, "%10s%5s%-10s\n", finishStack, "", copystack)
}
func stringfyQT(QT [4]uint8) string {
	ret := ""
	if QT[0] == 0 {
		return "none"
	}
	ret += fmt.Sprintf("%c ", QT[0])
	if QT[2] > 128 {
		ret += fmt.Sprint("t")
		QT[2] -= 128
	}
	ret += fmt.Sprintf("%c ", QT[2])
	if QT[1] > 128 {
		ret += fmt.Sprint("t")
		QT[1] -= 128
	}
	ret += fmt.Sprintf("%c ", QT[1])
	ret += fmt.Sprintf("t%d ", QT[3])
	return ret
}
func printQT(QT [4]uint8) {
	level := INFO
	debugPrintf(level, "QT: ")
	debugPrintf(level, "%c ", QT[0])
	if QT[2] > 128 {
		debugPrint(level, "t")
		QT[2] -= 128
	}
	debugPrintf(level, "%c ", QT[2])
	if QT[1] > 128 {
		debugPrint(level, "t")
		QT[1] -= 128
	}
	debugPrintf(level, "%c ", QT[1])
	debugPrintf(level, "t%d ", QT[3])
	debugPrint(level, "\n")
}
func (Grammar *GrammarLL1) PrintQuaternary() {
	level := INFO
	debugPrintf(level, "Quaternary:\n")
	for _, v := range Grammar.QTs {
		printQT(v)
	}
}
func (Grammar *GrammarLL1) ParseExpression(expression string) error {
	level := INFO
	debugPrintf(level, "\nParseExpression  %s\n", expression)
	if !Grammar.ready {
		debugPrintf(level, "Grammar not builded.\n")
		return errors.New("Grammar not builded.")
	}
	step := 0
	//add end symbol
	expression += "#"
	expression = strings.Replace(expression, " ", "", -1)
	stack := make([]uint8, 0)
	finishStack := make([]uint8, 0)
	stack = append(stack, 'S')
	temp_variable := uint8(0)
	SEM_stack := make([]uint8, 0)
	QT := [4]uint8{}
	Grammar.QTs = make([][4]uint8, 0)
	for index := 0; len(stack) > 0; {
		printState(stack, finishStack, SEM_stack, QT, expression, index)
		QT[0] = 0
		char := expression[index]
		topStack := stack[len(stack)-1]
		if topStack > 128 {
			topStack = topStack - 128
			if topStack == '+' || topStack == '-' || topStack == '*' || topStack == '/' {
				num1 := SEM_stack[len(SEM_stack)-1]
				num2 := SEM_stack[len(SEM_stack)-2]
				QT[0] = topStack
				QT[1] = num1
				QT[2] = num2
				QT[3] = uint8(temp_variable)
				SEM_stack = SEM_stack[:len(SEM_stack)-2]
				SEM_stack = append(SEM_stack, (temp_variable+'0')+128)
				temp_variable++
				Grammar.QTs = append(Grammar.QTs, QT)
			} else if isNumber(topStack) {
				debugPrintf(level, "push %c to SEM_stack\n", topStack)
				SEM_stack = append(SEM_stack, topStack)
			} else {
			}
			stack = stack[:len(stack)-1]
		} else if isTerminal(topStack) {
			if topStack == char || topStack == 'n' && isNumber(char) {
				if topStack == char {
					debugPrintf(level, "step:%d match a operater %c %c\n", step, topStack, char)
				} else {
					debugPrintf(level, "step:%d match a number %c %c\n", step, topStack, char)
				}
				step++
				finishStack = append(finishStack, topStack)
				//pop stack
				stack = stack[:len(stack)-1]
				index++
			} else {
				debugPrintf(ERROR, "Error: %c != %c\n", topStack, char)
				printErrorState(stack, finishStack, expression, index)
				return errors.New("Error: " + string(topStack) + " != " + string(char))
			}
		} else if isNonTerminal(topStack) {
			if isNumber(char) {
				//lookup in ParseTable
				token := Grammar.parseTable[topStack]['n']
				//pop stack
				stack = stack[:len(stack)-1]
				//check token start with operater
				if token[0] == '+' || token[0] == '-' || token[0] == '*' || token[0] == '/' {
					token_copy := make([]uint8, len(token))
					copy(token_copy, token)
					//insert push operater in the token
					token_copy = append(token_copy[:2], token_copy[1:]...)
					token_copy[2] = 128 + token_copy[0]
					for i := len(token_copy) - 1; i >= 0; i-- {
						stack = append(stack, token_copy[i])
					}
				} else if token[0] == 'n' {
					token_copy := make([]uint8, len(token))
					copy(token_copy, token)
					token_copy = append(token_copy, char+128)
					//push token_copy
					for i := len(token_copy) - 1; i >= 0; i-- {
						stack = append(stack, token_copy[i])
					}
				} else {
					//push token
					for i := len(token) - 1; i >= 0; i-- {
						stack = append(stack, token[i])
					}
				}
			} else {
				//lookup in ParseTable
				token := Grammar.parseTable[topStack][char]
				if len(token) == 0 {
					printErrorState(stack, finishStack, expression, index)
					return errors.New("Error: NonTerminal [" + string(topStack) + "] lookup fail")
				}
				//pop stack
				stack = stack[:len(stack)-1]
				//check token start with operater
				if token[0] == '+' || token[0] == '-' || token[0] == '*' || token[0] == '/' {
					token_copy := make([]uint8, len(token))
					copy(token_copy, token)
					//insert push operater in the token
					token_copy = append(token_copy[:2], token_copy[1:]...)
					token_copy[2] = 128 + token_copy[0]
					for i := len(token_copy) - 1; i >= 0; i-- {
						stack = append(stack, token_copy[i])
					}
				} else if token[0] == 'n' {
					token_copy := make([]uint8, len(token))
					copy(token_copy, token)
					token_copy = append(token_copy, char+128)
					//push token_copy
					for i := len(token_copy) - 1; i >= 0; i-- {
						stack = append(stack, token_copy[i])
					}
				} else {
					//push token
					for i := len(token) - 1; i >= 0; i-- {
						stack = append(stack, token[i])
					}
				}
			}
		} else if topStack == 'e' {
			stack = stack[:len(stack)-1]
		} else {
			debugPrintf(ERROR, "Error: %c\n", topStack)
			printErrorState(stack, finishStack, expression, index)
			return errors.New("Error: " + string(topStack) + " != " + string(char))
		}
	}
	Grammar.PrintQuaternary()
	return nil
}

func main() {
	grammar_filename := "grammar.txt"
	Grammar := GrammarLL1{}
	//read grammar
	Grammar.buildGrammar(grammar_filename)
	expression := "a+b*c"
	err := Grammar.ParseExpression(expression)
	if err != nil {
		debugPrintf(ERROR, "Parse fail %s\n", err)
	} else {
		debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	}
	expression = "a+b*(c+d)/f"
	err = Grammar.ParseExpression(expression)
	if err != nil {
		debugPrintf(ERROR, "Parse fail %s\n", err)
	} else {
		debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	}
	//read from stdin
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter expression: ")
		expression, _ := reader.ReadString('\n')
		expression = strings.Replace(expression, " ", "", -1)
		expression = strings.Replace(expression, "\r", "", -1)
		expression = strings.Replace(expression, "\n", "", -1)
		err = Grammar.ParseExpression(expression)
		if err != nil {
			debugPrintf(ERROR, "Parse fail %s\n", err)
		} else {
			debugPrintf(ERROR, "Parse expression %s success.\n", expression)
		}
	}

}
