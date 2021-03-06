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
type RunToken struct {
	token Token
	left  uint8
	index int
}
type Node struct {
	id         int
	jumpTable  map[uint8]int
	stateSet   [](RunToken)
	reduceAble bool
}
type GrammarLR0 struct {
	grammar       map[uint8]([]Token)
	terminals     []uint8
	nonTerminals  []uint8
	unfoldGrammar [](map[uint8](Token))
	closure       []Node
	ready         bool
}

func (Grammar *GrammarLR0) buildGrammar(grammar_filename string) {
	Grammar.grammar = make(map[uint8]([]Token))
	Grammar.readGrammarFromFile(grammar_filename)
	Grammar.genTerminalAndNonterminal()
	Grammar.genUnfoldGrammar()
	Grammar.genClosure()
	Grammar.printGrammarJumpTable()
	Grammar.ready = true
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
func (Grammar *GrammarLR0) readGrammarFromFile(grammar_filename string) {
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
func (Grammar *GrammarLR0) genTerminalAndNonterminal() {
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
func (Grammar *GrammarLR0) genUnfoldGrammar() {
	Grammar.unfoldGrammar = make([](map[uint8](Token)), 0)
	for key, value := range Grammar.grammar {
		for _, token := range value {
			Grammar.unfoldGrammar = append(Grammar.unfoldGrammar, map[uint8](Token){key: token})
		}
	}
	printUnfoldGrammar(Grammar.unfoldGrammar)

}
func isNumber(token uint8) bool {
	return token >= '0' && token <= '9'
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

func (Grammar *GrammarLR0) __addRunToken(runToken RunToken) []RunToken {
	level := DEBUG
	result := make([]RunToken, 0)
	if runToken.index < len(runToken.token) {
		next := runToken.token[runToken.index]
		if isNonTerminal(next) && next != runToken.left {
			debugPrintf(level, "[%c] is non terminal add\n", next)
			for _, token := range Grammar.grammar[next] {
				r := RunToken{
					token: Token(token),
					index: 0,
					left:  next,
				}
				result = append(result, r)
				debugPrintf(level, "\radd %s\n", token)
				result = append(result, Grammar.__addRunToken(r)...)
			}
		}
	}
	// result = uniqueToken(result)
	return result
}
func printRunToken(runtoken RunToken) {
	level := INFO
	debugPrintf(level, "%c ", runtoken.left)
	debugPrintf(level, "-> ")
	debugPrintf(level, "%s.", runtoken.token[:runtoken.index])
	debugPrintf(level, "%s\n", runtoken.token[runtoken.index:])
}
func printStateSet(stateSet []RunToken) {
	// level := DEBUG
	for _, runtoken := range stateSet {
		printRunToken(runtoken)
	}
}
func containToken(token RunToken, stateSet []RunToken) bool {
	for _, runtoken := range stateSet {
		if runtoken.index == token.index &&
			runtoken.left == token.left &&
			string(runtoken.token) == string(token.token) {
			return true
		}
	}
	return false
}
func uniqueToken(runtokens []RunToken) []RunToken {
	result := make([]RunToken, 0)
	for _, token := range runtokens {
		if !containToken(token, result) {
			result = append(result, token)
		}
	}
	return result
}
func (Grammar *GrammarLR0) __expandClosure(stateSet *[]RunToken) {
	level := DEBUG
	debugPrint(level, "expandClosure\n")
	//printStateSet(*stateSet)
	size := len(*stateSet)
	for i := 0; i < size; i++ {
		*stateSet = append(*stateSet, Grammar.__addRunToken((*stateSet)[i])...)
	}
	//remove duplicate
	*stateSet = uniqueToken(*stateSet)
	debugPrint(level, "after expandClosure\n")
	//printStateSet(*stateSet)
}
func (Grammar *GrammarLR0) __checkStateSet(stateSet []RunToken) (int, bool) {
	level := DEBUG
	//check if stateSet is in closure status
	for index, node := range Grammar.closure {
		flag := true
		if len(stateSet) != len(node.stateSet) {
			flag = false
			continue
		}
		for _, runtoken := range stateSet {
			if !containToken(runtoken, node.stateSet) {
				flag = false
				break
			}
		}
		if flag {
			//current state exist
			debugPrintf(level, "stateSet exist state %d\n", index)
			return index, true
		}
	}
	return -1, false
}

func (Grammar *GrammarLR0) __makeJump(closureNode *Node, token uint8) {
	level := DEBUG
	debugPrintf(level, "make jump from %d token %c\n", closureNode.id, token)
	newStateSet := make([]RunToken, 0)
	for _, runtoken := range closureNode.stateSet {
		if runtoken.index < len(runtoken.token) && runtoken.token[runtoken.index] == token {
			debugPrintf(level, "match token %c at Token %s\n", token, runtoken.token)
			newStateSet = append(newStateSet, RunToken{
				token: make(Token, len(runtoken.token)),
				index: runtoken.index + 1,
				left:  runtoken.left,
			})
			copy(newStateSet[len(newStateSet)-1].token, runtoken.token)
			newStateSet[len(newStateSet)-1].index = runtoken.index + 1
			newStateSet[len(newStateSet)-1].left = runtoken.left
			//printStateSet(newStateSet)
		}
	}
	Grammar.__expandClosure(&newStateSet)
	index, flag := Grammar.__checkStateSet(newStateSet)
	if flag {
		debugPrintf(level, "exist token %c from %d jump to %d\n", token, closureNode.id, index)
		closureNode.jumpTable[token] = index
	} else {
		//add a new state
		debugPrintf(level, "add new state id:%d\n", len(Grammar.closure))
		closureNode.jumpTable[token] = len(Grammar.closure)
		Grammar.closure = append(Grammar.closure, Node{
			id:        len(Grammar.closure),
			stateSet:  newStateSet,
			jumpTable: make(map[uint8]int),
		})
		debugPrintf(level, "not exist token %c from %d jump to %d\n", token, closureNode.id, closureNode.jumpTable[token])

	}
}
func (Grammar *GrammarLR0) __buildJumptable(closureNode *Node) {
	level := DEBUG
	// closureNode := Grammar.closure[closureIndex]
	debugPrintf(level, "build jump table %d\n", closureNode.id)
	buildOk := make(map[uint8]bool)
	for _, key := range Grammar.nonTerminals {
		buildOk[key] = false
	}
	for _, key := range Grammar.terminals {
		buildOk[key] = false
	}
	for _, runtoken := range closureNode.stateSet {
		debugPrint(level, "runToken ")
		//printRunToken(runtoken)
		if runtoken.index == len(runtoken.token) {
			//reach the end
			debugPrintf(level, "reach end state %d can be reduceAble\n", closureNode.id)
			closureNode.reduceAble = true
		} else if !buildOk[runtoken.token[runtoken.index]] {
			Grammar.__makeJump(closureNode, runtoken.token[runtoken.index])
			//printRunToken(runtoken)

			buildOk[runtoken.token[runtoken.index]] = true
		}
	}

}
func printJumpTable(jumpTable map[uint8]int) {
	level := INFO
	for key, value := range jumpTable {
		debugPrintf(level, "token %c to %d\n", key, value)
	}
	debugPrint(level, "\n")
}
func (Grammar *GrammarLR0) printGrammarJumpTable() {
	level := INFO
	title := "state id Reducable?"
	for _, key := range Grammar.terminals {
		title += fmt.Sprintf("%5c", key)
	}
	for _, key := range Grammar.nonTerminals {
		title += fmt.Sprintf("%5c", key)
	}
	debugPrintf(level, "%s\n", title)
	for index := range Grammar.closure {
		txt := fmt.Sprintf("%6s%-9d", "", index)
		node := Grammar.closure[index]
		if node.reduceAble {
			txt += fmt.Sprintf("%4s", "yes")
		} else {
			txt += fmt.Sprintf("%4s", "no")
		}
		for _, key := range Grammar.terminals {
			if value, ok := node.jumpTable[key]; ok {
				txt += fmt.Sprintf("%5d", value)
			} else {
				txt += fmt.Sprintf("%5s", "")
			}
		}
		for _, key := range Grammar.nonTerminals {
			if value, ok := node.jumpTable[key]; ok {
				txt += fmt.Sprintf("%5d", value)
			} else {
				txt += fmt.Sprintf("%5s", "")
			}
		}
		debugPrintf(level, "%s\n", txt)
	}
	debugPrint(level, "\n")
}

func (Grammar *GrammarLR0) genClosure() {
	level := INFO
	Grammar.closure = make([]Node, 0)
	Grammar.closure = append(Grammar.closure, Node{
		id:        0,
		jumpTable: make(map[uint8]int),
		stateSet:  make([]RunToken, 0),
	})
	Grammar.closure[0].stateSet = append(Grammar.closure[0].stateSet, RunToken{
		token: Grammar.grammar['Y'][0],
		index: 0,
		left:  'Y',
	})
	for i := 0; i < len(Grammar.closure); i++ {
		Grammar.__expandClosure(&Grammar.closure[i].stateSet)
		Grammar.__buildJumptable(&Grammar.closure[i])
	}
	for i := 0; i < len(Grammar.closure); i++ {
		debugPrintf(level, "state id:%d\n", i)
		printStateSet(Grammar.closure[i].stateSet)
		// debugPrintf(level, "jumptable state %d\n", i)
		// printJumpTable(Grammar.closure[i].jumpTable)
	}
}
func printLR0State(stack []uint8, state_Stack []int, expression string, index int) {
	level := INFO
	leftExppressionstr := "   "
	rightExppressionstr := "   "
	leftExppression := expression[:index]
	rightExppression := expression[index:]
	for _, value := range leftExppression {
		leftExppressionstr += fmt.Sprintf("%3c", value)
	}
	for _, value := range rightExppression {
		rightExppressionstr += fmt.Sprintf("%3c", value)
	}
	stackstr := ""
	state_stackstr := ""
	for _, value := range stack {
		stackstr += fmt.Sprintf("%3c", value)
	}
	for _, value := range state_Stack {
		state_stackstr += fmt.Sprintf("%3d", value)
	}
	// debugPrintf(level, "%s\n", title)
	debugPrintf(level, "%-30s%5s%-20s\n", "matched", "", "matching")
	debugPrintf(level, "%-30s%5s%-20s\n", leftExppressionstr, "", rightExppression)
	debugPrintf(level, "%-30s%5s%-20s\n", stackstr, "", "")
	debugPrintf(level, "%-30s%5s%-20s\n", state_stackstr, "", "")

}
func (Grammar *GrammarLR0) ParseExpression(expression string) error {
	level := INFO
	debugPrintf(level, "\nParseExpression  %s\n", expression)
	if !Grammar.ready {
		debugPrintf(level, "Grammar not builded.\n")
		return errors.New("grammar not builded")
	}
	step := 0
	index := 0
	//add end symbol
	expression += "#"
	expression = strings.Replace(expression, " ", "", -1)
	stack := make([]uint8, 0)
	state_stack := make([]int, 0)
	// finishStack := make([]uint8, 0)
	stack = append([]uint8{'#'}, stack...)
	state_stack = append(state_stack, 0)
	for {
		printLR0State(stack, state_stack, expression, index)
		if stack[len(stack)-1] == 'Y' {
			return nil
		}
		stateTop := state_stack[len(state_stack)-1]
		if len(state_stack) < len(stack) {
			next := stack[len(stack)-1]
			if value, ok := Grammar.closure[stateTop].jumpTable[next]; ok {
				//state change
				state_stack = append(state_stack, value)
			}
		} else {
			next := expression[index]
			if isNumber(next) {
				if value, ok := Grammar.closure[stateTop].jumpTable['n']; ok {
					//state change
					state_stack = append(state_stack, value)
					stack = append(stack, 'n')
					index++
				}
			} else if next == '(' {
				if value, ok := Grammar.closure[stateTop].jumpTable['(']; ok {
					//state change
					state_stack = append(state_stack, value)
					stack = append(stack, '(')
					index++
				}
			} else if value, ok := Grammar.closure[stateTop].jumpTable[next]; ok {
				//state change
				state_stack = append(state_stack, value)
				stack = append(stack, next)
				index++
			} else if Grammar.closure[stateTop].reduceAble {
				//reduce
				//find left
				for _, runtoken := range Grammar.closure[stateTop].stateSet {
					if runtoken.index == len(runtoken.token) {
						if runtoken.left == 'Y' {
							stack = stack[:len(stack)-len(runtoken.token)]
							stack = append(stack, runtoken.left)
							state_stack = state_stack[:len(state_stack)-len(runtoken.token)]
						} else {
							debugPrintf(level, "next step reduce use: ")
							printRunToken(runtoken)
							debugPrintf(level, "\n")
							stack = stack[:len(stack)-len(runtoken.token)]
							stack = append(stack, runtoken.left)
							state_stack = state_stack[:len(state_stack)-len(runtoken.token)]
						}
					}
				}
			} else {
				//error
				return fmt.Errorf("can not accept next %c at state :%d", next, stateTop)
			}
		}

		step++
	}
}

func main() {
	grammar_filename := "../grammarlr0.txt"
	Grammar := GrammarLR0{}
	//read grammar
	Grammar.buildGrammar(grammar_filename)
	expression := "3*(2-1)"
	err := Grammar.ParseExpression(expression)
	if err != nil {
		debugPrintf(ERROR, "Parse expression %s fail: %s\n", expression, err)
	} else {
		debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	}
	expression = "3+1/7+"
	err = Grammar.ParseExpression(expression)
	if err != nil {
		debugPrintf(ERROR, "Parse expression %s fail: %s\n", expression, err)
	} else {
		debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	}
	expression = "3+(*6"
	err = Grammar.ParseExpression(expression)
	if err != nil {
		debugPrintf(ERROR, "Parse expression %s fail: %s\n", expression, err)
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
			debugPrintf(ERROR, "Parse expression %s fail: %s\n", expression, err)
		} else {
			debugPrintf(ERROR, "Parse expression %s success.\n", expression)
		}
	}

}
