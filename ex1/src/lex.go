package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type TokenType int

const DEBUG = false
const (
	// TokenType
	INCLUDE = iota
	DEFINE
	HEAD_FILE
	BLOCKCOMMENT
	LINECOMMENT
	IDENTIFIER
	STRING
	CHAR
	INT
	FLOAT
	DOUBLE
	VOID
	IF
	ELSE
	FOR
	WHILE
	RETURN
	BREAK
	CONTINUE
	LBRACE
	RBRACE
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	SEMICOLON
	COMMA
	ASSIGN
	PLUS
	MINUS
	MUL
	DIV
	MOD
	EQ
	NEQ
	LT
	GT
	LEQ
	GEQ
	AND
	OR
	NOT
	EOF
	DO
	CONST
	STRUCT
	UNION
	ENUM
	TYPEDEF
	EXTERN
	STATIC
	AUTO
	REGISTER
	SIGNED
	UNSIGNED
	SHORT
	LONG
	PLUSASSIGN
	MINUSASSIGN
	MULASSIGN
	DIVASSIGN
	MODASSIGN
	ANDASSIGN
	ORASSIGN
	XORASSIGN
	LSHIFTASSIGN
	RSHIFTASSIGN
	BITAND
	BITOR
	BITXOR
	BITNOT
	BITLSHIFT
	BITRSHIFT
	MACRO
	NUMBER
	UNKNOWN
	BADTOKEN
)

var TokenTypeStrings = map[TokenType]string{
	INCLUDE:      "INCLUDE",
	DEFINE:       "DEFINE",
	HEAD_FILE:    "HEAD_FILE",
	BLOCKCOMMENT: "BLOCKCOMMENT",
	LINECOMMENT:  "LINECOMMENT",
	IDENTIFIER:   "IDENTIFIER",
	STRING:       "STRING",
	CHAR:         "CHAR",
	INT:          "INT",
	FLOAT:        "FLOAT",
	DOUBLE:       "DOUBLE",
	VOID:         "VOID",
	IF:           "IF",
	ELSE:         "ELSE",
	FOR:          "FOR",
	WHILE:        "WHILE",
	RETURN:       "RETURN",
	BREAK:        "BREAK",
	CONTINUE:     "CONTINUE",
	LBRACE:       "LBRACE",
	RBRACE:       "RBRACE",
	LPAREN:       "LPAREN",
	RPAREN:       "RPAREN",
	LBRACKET:     "LBRACKET",
	RBRACKET:     "RBRACKET",
	SEMICOLON:    "SEMICOLON",
	COMMA:        "COMMA",
	ASSIGN:       "ASSIGN",
	PLUS:         "PLUS",
	MINUS:        "MINUS",
	MUL:          "MUL",
	DIV:          "DIV",
	MOD:          "MOD",
	EQ:           "EQ",
	NEQ:          "NEQ",
	LT:           "LT",
	GT:           "GT",
	LEQ:          "LEQ",
	GEQ:          "GEQ",
	AND:          "AND",
	OR:           "OR",
	NOT:          "NOT",
	EOF:          "EOF",
	DO:           "DO",
	CONST:        "CONST",
	STRUCT:       "STRUCT",
	UNION:        "UNION",
	ENUM:         "ENUM",
	TYPEDEF:      "TYPEDEF",
	EXTERN:       "EXTERN",
	STATIC:       "STATIC",
	AUTO:         "AUTO",
	REGISTER:     "REGISTER",
	SIGNED:       "SIGNED",
	UNSIGNED:     "UNSIGNED",
	SHORT:        "SHORT",
	LONG:         "LONG",
	PLUSASSIGN:   "PLUSASSIGN",
	MINUSASSIGN:  "MINUSASSIGN",
	MULASSIGN:    "MULASSIGN",
	DIVASSIGN:    "DIVASSIGN",
	MODASSIGN:    "MODASSIGN",
	ANDASSIGN:    "ANDASSIGN",
	ORASSIGN:     "ORASSIGN",
	XORASSIGN:    "XORASSIGN",
	LSHIFTASSIGN: "LSHIFTASSIGN",
	RSHIFTASSIGN: "RSHIFTASSIGN",
	BITAND:       "BITAND",
	BITOR:        "BITOR",
	BITXOR:       "BITXOR",
	BITNOT:       "BITNOT",
	BITLSHIFT:    "BITLSHIFT",
	BITRSHIFT:    "BITRSHIFT",
	MACRO:        "MACRO",
	NUMBER:       "NUMBER",
	UNKNOWN:      "UNKNOWN",
	BADTOKEN:     "BADTOKEN",
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"include":    INCLUDE,
	"define":     DEFINE,
	"head_file":  HEAD_FILE,
	"identifier": IDENTIFIER,
	"string":     STRING,
	"char":       CHAR,
	"int":        INT,
	"float":      FLOAT,
	"double":     DOUBLE,
	"void":       VOID,
	"if":         IF,
	"else":       ELSE,
	"for":        FOR,
	"while":      WHILE,
	"return":     RETURN,
	"break":      BREAK,
	"continue":   CONTINUE,
	"do":         DO,
	"const":      CONST,
	"struct":     STRUCT,
	"union":      UNION,
	"enum":       ENUM,
	"typedef":    TYPEDEF,
	"extern":     EXTERN,
	"static":     STATIC,
	"auto":       AUTO,
	"register":   REGISTER,
	"signed":     SIGNED,
	"unsigned":   UNSIGNED,
	"short":      SHORT,
	"long":       LONG,
}
var operaters = map[string]TokenType{
	"+=":  PLUSASSIGN,
	"-=":  MINUSASSIGN,
	"*=":  MULASSIGN,
	"/=":  DIVASSIGN,
	"%=":  MODASSIGN,
	"&=":  ANDASSIGN,
	"|=":  ORASSIGN,
	"^=":  XORASSIGN,
	"<<=": LSHIFTASSIGN,
	">>=": RSHIFTASSIGN,
	"&":   BITAND,
	"|":   BITOR,
	"^":   BITXOR,
	"~":   BITNOT,
	"<<":  BITLSHIFT,
	">>":  BITRSHIFT,
	"\"":  STRING,
	"'":   CHAR,
	"#":   MACRO,
	"{":   LBRACE,
	"}":   RBRACE,
	"(":   LPAREN,
	")":   RPAREN,
	"[":   LBRACKET,
	"]":   RBRACKET,
	";":   SEMICOLON,
	",":   COMMA,
	"=":   ASSIGN,
	"+":   PLUS,
	"-":   MINUS,
	"*":   MUL,
	"/":   DIV,
	"%":   MOD,
	"==":  EQ,
	"!=":  NEQ,
	"<":   LT,
	">":   GT,
	"<=":  LEQ,
	">=":  GEQ,
	"&&":  AND,
	"||":  OR,
	"!":   NOT,
	"/*":  BLOCKCOMMENT,
	"//":  LINECOMMENT,
}
var operatorStrings = map[TokenType]string{
	PLUSASSIGN:   "PLUSASSIGN",
	MINUSASSIGN:  "MINUSASSIGN",
	MULASSIGN:    "MULASSIGN",
	DIVASSIGN:    "DIVASSIGN",
	MODASSIGN:    "MODASSIGN",
	ANDASSIGN:    "ANDASSIGN",
	ORASSIGN:     "ORASSIGN",
	XORASSIGN:    "XORASSIGN",
	LSHIFTASSIGN: "LSHIFTASSIGN",
	RSHIFTASSIGN: "RSHIFTASSIGN",
	BITAND:       "BITAND",
	BITOR:        "BITOR",
	BITXOR:       "BITXOR",
	BITNOT:       "BITNOT",
	BITLSHIFT:    "BITLSHIFT",
	BITRSHIFT:    "BITRSHIFT",
	STRING:       "STRING",
	CHAR:         "CHAR",
	MACRO:        "MACRO",
	LBRACE:       "LBRACE",
	RBRACE:       "RBRACE",
	LPAREN:       "LPAREN",
	RPAREN:       "RPAREN",
	LBRACKET:     "LBRACKET",
	RBRACKET:     "RBRACKET",
	SEMICOLON:    "SEMICOLON",
	COMMA:        "COMMA",
	ASSIGN:       "ASSIGN",
	PLUS:         "PLUS",
	MINUS:        "MINUS",
	MUL:          "MUL",
	DIV:          "DIV",
	MOD:          "MOD",
	EQ:           "EQ",
	NEQ:          "NEQ",
	LT:           "LT",
	GT:           "GT",
	LEQ:          "LEQ",
	GEQ:          "GEQ",
	AND:          "AND",
	OR:           "OR",
	NOT:          "NOT",
	BLOCKCOMMENT: "BLOCKCOMMENT",
	LINECOMMENT:  "LINECOMMENT",
}

var operatorChars = []byte{'+', '-', '*', '/', '%', '&', '|', '^', '~', '<', '>', '=', '!', '?', ':', ',', '#', '"', '\'', '\\', '.'}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
func isEscape(ch byte) bool {
	return ch == '\\'
}
func isAlphaLodashNum(ch byte) bool {
	return isAlphaLodash(ch) || isDigit(ch)
}
func isAlphaLodash(ch byte) bool {
	return isAlpha(ch) || ch == '_'
}

func isStringQuote(ch byte) bool {
	return ch == '"'
}
func isCharQuote(ch byte) bool {
	return ch == '\''
}

func isOperaterChar(ch byte) bool {
	for _, c := range operatorChars {
		if c == ch {
			return true
		}
	}
	return false
}

func isOperaterString(s string) bool {
	_, ok := operaters[s]
	return ok
}

func iskeyword(s string) bool {
	_, ok := keywords[s]
	return ok
}
func isStop(ch byte) bool {
	return ch == ',' || ch == ';' || ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == '[' || ch == ']'
}

func (token Token) String() string {
	return fmt.Sprintf("<Type : %13s %10s Line : %3d\tColumn : %3d>\n", TokenTypeStrings[token.Type], token.Literal, token.Line, token.Column)
}

//states
const (
	START = iota
	LETTER_STATE
	DIGIT_STATE_NO_POINT_NO_E
	DIGIT_STATE_WITH_POINT_NO_E
	DIGIT_STATE_WITH_POINT_WITH_E
	DIGIT_STATE_NO_POINT_WITH_E
	CHAR_STATE
	CHAR_STATE_ESCAPE
	CHAR_STATE_LETTER
	CHAR_STATE_END
	STRING_STATE
	STRING_STATE_END
	OPERATER_STATE
	ERROR
	STOP
)

var stateStrings = map[int]string{
	START:                         "START",
	LETTER_STATE:                  "LETTER_STATE",
	DIGIT_STATE_NO_POINT_NO_E:     "DIGIT_STATE_NO_POINT_NO_E",
	DIGIT_STATE_WITH_POINT_NO_E:   "DIGIT_STATE_WITH_POINT_NO_E",
	DIGIT_STATE_WITH_POINT_WITH_E: "DIGIT_STATE_WITH_POINT_WITH_E",
	DIGIT_STATE_NO_POINT_WITH_E:   "DIGIT_STATE_NO_POINT_WITH_E",
	CHAR_STATE:                    "CHAR_STATE",
	CHAR_STATE_ESCAPE:             "CHAR_STATE_ESCAPE",
	CHAR_STATE_END:                "CHAR_STATE_END",
	STRING_STATE:                  "STRING_STATE",
	OPERATER_STATE:                "OPERATER_STATE",
	ERROR:                         "ERROR",
	STOP:                          "STOP",
}

func debugPrint(ch byte, state int, cur_string string) {
	if DEBUG {
		fmt.Print("now char is ", string(ch), "  ")
		fmt.Print("cur_string is ", cur_string, "  ")
		fmt.Println("state is ", stateStrings[state])
	}
}
func restart(state *int, cur_string *string, i *int, ch byte, cur *int) {
	if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
		*cur_string = ""
	} else {
		*cur_string = ""
		*i = *i - 1
		*cur = *cur - 1
	}
	*state = START
}
func main() {
	//read a cpp file
	data, err := ioutil.ReadFile("../demo.c")
	if err != nil {
		panic(err)
	}
	//init
	var result []Token
	fmt.Println(string(data))
	line, cur := 1, -1
	state := START
	cur_string := ""
	//scan the cpp file by byte
	for i := 0; i < len(data); i++ {
		ch := data[i]
		cur++
		if ch == '\n' {
			line++
			cur = -1
		}
		if DEBUG {
			fmt.Print("now char is ", string(ch), "\t")
			fmt.Print("cur_string is ", cur_string, "\t")
			fmt.Println("state is ", stateStrings[state], "\t")
		}
		//state machine
		if state == START {
			if isSpace(ch) {
				continue
			}
			if isDigit(ch) {
				cur_string += string(ch)
				state = DIGIT_STATE_NO_POINT_NO_E
				continue
			}
			if isAlpha(ch) {
				cur_string += string(ch)
				state = LETTER_STATE
				continue
			}
			if isStringQuote(ch) {
				cur_string += string(ch)
				state = STRING_STATE
				continue
			}
			if isCharQuote(ch) {
				cur_string += string(ch)
				state = CHAR_STATE
				continue
			}
			if isOperaterChar(ch) {
				cur_string += string(ch)
				state = OPERATER_STATE
				continue
			}
			if isStop(ch) {
				cur_string += string(ch)
				state = STOP
				continue
			}

		}
		if state == DIGIT_STATE_NO_POINT_NO_E {
			if isDigit(ch) {
				cur_string += string(ch)
				continue
			}
			if ch == 'E' || ch == 'e' {
				cur_string += string(ch)
				state = DIGIT_STATE_NO_POINT_WITH_E
				continue
			}
			if ch == '.' {
				cur_string += string(ch)
				state = DIGIT_STATE_WITH_POINT_NO_E
				continue
			}
			result = append(result, Token{NUMBER, cur_string, line, cur})
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}
		if state == DIGIT_STATE_WITH_POINT_NO_E {
			if isDigit(ch) {
				cur_string += string(ch)
				continue
			}
			if ch == 'E' || ch == 'e' {
				cur_string += string(ch)
				state = DIGIT_STATE_WITH_POINT_WITH_E
				continue
			}
			result = append(result, Token{NUMBER, cur_string, line, cur})
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}
		if state == DIGIT_STATE_WITH_POINT_WITH_E {
			if isDigit(ch) {
				cur_string += string(ch)
				continue
			}
			result = append(result, Token{NUMBER, cur_string, line, cur})
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}
		if state == LETTER_STATE {
			if isAlphaLodashNum(ch) {
				cur_string += string(ch)
				continue
			}
			if iskeyword(cur_string) {
				result = append(result, Token{keywords[cur_string], cur_string, line, cur})
			} else {
				result = append(result, Token{IDENTIFIER, cur_string, line, cur})
			}
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}
		if state == STRING_STATE {
			if isStringQuote(ch) {
				cur_string += string(ch)
				state = STRING_STATE_END
				continue
			}
			cur_string += string(ch)
			continue
		}
		if state == CHAR_STATE {
			if isEscape(ch) {
				cur_string += string(ch)
				state = CHAR_STATE_ESCAPE
				continue
			}
			cur_string += string(ch)
			state = CHAR_STATE_LETTER
			continue
		}
		if state == CHAR_STATE_ESCAPE {
			cur_string += string(ch)
			state = CHAR_STATE_LETTER
			continue
		}
		if state == CHAR_STATE_LETTER {
			if isCharQuote(ch) {
				cur_string += string(ch)
				state = CHAR_STATE_END
				continue
			}
			state = ERROR
			continue
		}
		if state == CHAR_STATE_END {
			result = append(result, Token{CHAR, cur_string, line, cur})
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}
		if state == OPERATER_STATE {
			if isOperaterChar(ch) {
				cur_string += string(ch)
				state = OPERATER_STATE
				continue
			}
			if isOperaterString(cur_string) {
				result = append(result, Token{operaters[cur_string], cur_string, line, cur})
			} else {
				result = append(result, Token{UNKNOWN, cur_string, line, cur})
			}
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}
		if state == ERROR {
			result = append(result, Token{BADTOKEN, cur_string, line, cur})
			restart(&state, &cur_string, &i, ch, &cur)
		}
		if state == STOP {
			result = append(result, Token{operaters[cur_string], cur_string, line, cur})
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}
		if state == STRING_STATE_END {
			result = append(result, Token{STRING, cur_string, line, cur})
			restart(&state, &cur_string, &i, ch, &cur)
			continue
		}

	}
	for _, i := range result {
		fmt.Print(i.String())
	}
	fmt.Printf("Lexical finish get %d tokens\n", len(result))
	//print result to tokens.txt
	file, err := os.Create("tokens.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for _, i := range result {
		file.WriteString(i.String())
	}
}
