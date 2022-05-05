package main

import (
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

const debugLevel = DEBUG
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

// func uniqueToken(RunTokenSlice []RunToken) []RunToken {
// 	keys := make(map[RunToken]bool)
// 	list := []RunToken{}
// 	for _, entry := range RunTokenSlice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }
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
	first         map[uint8]([]uint8)
	follow        map[uint8]([]uint8)
	parseTable    map[uint8](map[uint8]([]uint8))
	nfa           []Node
	ready         bool
}

func (Grammar *GrammarLR0) buildGrammar(grammar_filename string) {
	Grammar.grammar = make(map[uint8]([]Token))
	Grammar.readGrammarFromFile(grammar_filename)
	Grammar.genTerminalAndNonterminal()
	Grammar.genUnfoldGrammar()
	// Grammar.genFirst()
	// Grammar.genFollow()
	// Grammar.printFirstFollow()
	Grammar.genNFA()
	Grammar.ready = true
}
func (Grammar *GrammarLR0) printFirstFollow() {
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

// type Node struct {
// 	in   uint8
// 	out  uint8
// 	next []uint8
// }
// type Graph struct {
// 	nodes map[uint8]Node
// }
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

func (Grammar *GrammarLR0) __buildFirst(key uint8, buildOK *map[uint8]bool) {
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

func (Grammar *GrammarLR0) genFirst() {
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
func (Grammar *GrammarLR0) __buildFollow(key uint8, buildOK *map[uint8]bool) {
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
func (Grammar *GrammarLR0) genFollow() {
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
func (Grammar *GrammarLR0) printParseTable() {
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
func (Grammar *GrammarLR0) __buildNFA(nfaIndex int) {
	level := DEBUG
	// nfaNode := Grammar.nfa[nfaIndex]
	debugPrintf(level, "build nfa %d\n", nfaIndex)

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
	level := DEBUG
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
func (Grammar *GrammarLR0) __expandNFA(stateSet *[]RunToken) {
	level := DEBUG
	debugPrint(level, "expandNFA\n")
	printStateSet(*stateSet)
	size := len(*stateSet)
	for i := 0; i < size; i++ {
		*stateSet = append(*stateSet, Grammar.__addRunToken((*stateSet)[i])...)
	}
	//remove duplicate
	*stateSet = uniqueToken(*stateSet)
	debugPrint(level, "after expandNFA\n")
	printStateSet(*stateSet)
}
func (Grammar *GrammarLR0) __checkStateSet(stateSet []RunToken) (int, bool) {
	level := DEBUG
	//check if stateSet is in nfa status
	for index, node := range Grammar.nfa {
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

func (Grammar *GrammarLR0) __makeJump(nfaNode *Node, token uint8) {
	level := DEBUG
	debugPrintf(level, "make jump from %d token %c\n", nfaNode.id, token)
	newStateSet := make([]RunToken, 0)
	for _, runtoken := range nfaNode.stateSet {
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
			printStateSet(newStateSet)
		}
	}
	Grammar.__expandNFA(&newStateSet)
	index, flag := Grammar.__checkStateSet(newStateSet)
	if flag {
		debugPrintf(level, "exist token %c from %d jump to %d\n", token, nfaNode.id, index)
		nfaNode.jumpTable[token] = index
	} else {
		//add a new state
		debugPrintf(level, "add new state id:%d\n", len(Grammar.nfa))
		nfaNode.jumpTable[token] = len(Grammar.nfa)
		Grammar.nfa = append(Grammar.nfa, Node{
			id:        len(Grammar.nfa),
			stateSet:  newStateSet,
			jumpTable: make(map[uint8]int),
		})
		debugPrintf(level, "not exist token %c from %d jump to %d\n", token, nfaNode.id, nfaNode.jumpTable[token])

	}
}
func (Grammar *GrammarLR0) __buildJumptable(nfaNode *Node) {
	level := DEBUG
	// nfaNode := Grammar.nfa[nfaIndex]
	debugPrintf(level, "build jump table %d\n", nfaNode.id)
	buildOk := make(map[uint8]bool)
	for _, key := range Grammar.nonTerminals {
		buildOk[key] = false
	}
	for _, key := range Grammar.terminals {
		buildOk[key] = false
	}
	for _, runtoken := range nfaNode.stateSet {
		debugPrint(level, "runToken ")
		printRunToken(runtoken)
		if runtoken.index == len(runtoken.token) {
			//reach the end
			debugPrintf(level, "reach end state %d can be reduceAble\n", nfaNode.id)
			nfaNode.reduceAble = true
		} else if !buildOk[runtoken.token[runtoken.index]] {
			Grammar.__makeJump(nfaNode, runtoken.token[runtoken.index])
			printRunToken(runtoken)

			buildOk[runtoken.token[runtoken.index]] = true
		}
	}

}
func (Grammar *GrammarLR0) genNFA() {
	level := INFO
	Grammar.nfa = make([]Node, 0)
	Grammar.nfa = append(Grammar.nfa, Node{
		id:        0,
		jumpTable: make(map[uint8]int),
		stateSet:  make([]RunToken, 0),
	})
	Grammar.nfa[0].stateSet = append(Grammar.nfa[0].stateSet, RunToken{
		token: Grammar.grammar['S'][0],
		index: 0,
		left:  'S',
	})
	for i := 0; i < len(Grammar.nfa); i++ {
		Grammar.__expandNFA(&Grammar.nfa[i].stateSet)
		Grammar.__buildJumptable(&Grammar.nfa[i])
	}
	for i := 0; i < len(Grammar.nfa); i++ {
		debugPrintf(level, "state id:%d\n", i)
		printStateSet(Grammar.nfa[i].stateSet)
	}

}
func printState(stack []uint8, finishStack []uint8, expression string, index int) {
	level := DEBUG
	leftExppression := expression[:index]
	rightExppression := expression[index:]
	debugPrintf(level, "%-10s%5s%-10s\n", "matched", "", "matching")
	debugPrintf(level, "%10s%5s%-10s\n", leftExppression, "", rightExppression)
	copystack := make([]uint8, len(stack))
	copy(copystack, stack)
	//reverse stack
	for i, j := 0, len(copystack)-1; i < j; i, j = i+1, j-1 {
		copystack[i], copystack[j] = copystack[j], copystack[i]
	}
	debugPrintf(level, "%10s%5s%-10s\n\n", finishStack, "", copystack)
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
func (Grammar *GrammarLR0) ParseExpression(expression string) error {
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
	for index := 0; len(stack) > 0; {
		printState(stack, finishStack, expression, index)
		char := expression[index]
		topStack := stack[len(stack)-1]
		if isTerminal(topStack) {
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
				//push token
				for i := len(token) - 1; i >= 0; i-- {
					stack = append(stack, token[i])
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
				//push token
				for i := len(token) - 1; i >= 0; i-- {
					stack = append(stack, token[i])
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
	return nil
}

func main() {
	grammar_filename := "../grammarlr0.txt"
	Grammar := GrammarLR0{}
	//read grammar
	Grammar.buildGrammar(grammar_filename)
	// expression := "3+1*(5+6)/7"
	// err := Grammar.ParseExpression(expression)
	// if err != nil {
	// 	debugPrintf(ERROR, "Parse fail %s\n", err)
	// } else {
	// 	debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	// }
	// expression = "3+1*(5+6)/7+"
	// err = Grammar.ParseExpression(expression)
	// if err != nil {
	// 	debugPrintf(ERROR, "Parse fail %s\n", err)
	// } else {
	// 	debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	// }
	// expression = "3+1*(5++6"
	// err = Grammar.ParseExpression(expression)
	// if err != nil {
	// 	debugPrintf(ERROR, "Parse fail %s\n", err)
	// } else {
	// 	debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	// }
	// //read from stdin
	// reader := bufio.NewReader(os.Stdin)
	// for {
	// 	fmt.Print("Enter expression: ")
	// 	expression, _ := reader.ReadString('\n')
	// 	expression = strings.Replace(expression, " ", "", -1)
	// 	expression = strings.Replace(expression, "\r", "", -1)
	// 	expression = strings.Replace(expression, "\n", "", -1)
	// 	err = Grammar.ParseExpression(expression)
	// 	if err != nil {
	// 		debugPrintf(ERROR, "Parse fail %s\n", err)
	// 	} else {
	// 		debugPrintf(ERROR, "Parse expression %s success.\n", expression)
	// 	}
	// }

}