package main

/*
PEG definition
	Expr    ← Sum
	Sum     ← Product (('+' / '-') Product)*
	Product ← Value (('*' / '/') Value)*
	Value   ← [0-9]+ / '(' Expr ')'

refer: https://en.wikipedia.org/wiki/Parsing_expression_grammar
 */

import (
	"strconv"
	"strings"
	"fmt"
)

const (
	TokenNone = iota
	TokenValue
	TokenPlus
	TokenMinus
	TokenMulti
	TokenDiv
	TokenLeftBracket
	TokenRightBracket
)

func IsValue(c uint8) bool {
	return c > uint8('0') && c < uint8('9')
}

func NextToken(s string) (int, string, string) {
	// skip space
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return TokenNone, "", s
	}

	switch string(s[0]) {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		i := 0
		for {
			if i < len(s) && IsValue(s[i]) {
				i++
			} else {
				return TokenValue, s[0:i], s[i:]
			}
		}

	case "+":
		return TokenPlus, s[0:1], s[1:]
	case "-":
		return TokenMinus, s[0:1], s[1:]
	case "*":
		return TokenMulti, s[0:1], s[1:]
	case "/":
		return TokenDiv, s[0:1], s[1:]
	case "(":
		return TokenLeftBracket, s[0:1], s[1:]
	case ")":
		return TokenRightBracket, s[0:1], s[1:]
	default:
		panic("unknown token: " + s)
	}
}

func EvalExpr(s string) (int, string) {
	return EvalSum(s)
}

func EvalSum(s string) (int, string) {
	var result int
	var token int
	var v string
	var temp int

	result, s = EvalProduct(s)
	for {
		token, v, s = NextToken(s)

		switch token {

		case TokenPlus:
			temp, s = EvalProduct(s)
			result = result + temp

		case TokenMinus:
			temp, s = EvalProduct(s)
			result = result - temp

		default:
			return result, v + s
		}
	}
}

func EvalProduct(s string) (int, string) {
	var result int
	var token int
	var v string
	var temp int

	result, s = EvalValue(s)
	for {
		token, v, s = NextToken(s)

		switch token {

		case TokenMulti:
			temp, s = EvalValue(s)
			result = result * temp

		case TokenDiv:
			temp, s = EvalValue(s)
			result = result / temp

		default:
			return result, v + s
		}
	}
}

func EvalValue(s string) (int, string) {
	var result int
	var err error
	var token int
	var v string

	for {
		token, v, s = NextToken(s)

		switch token {
		case TokenValue:
			result, err = strconv.Atoi(v)
			if err != nil {
				panic("wrong value " + v)
			}
			return result, s
		case TokenLeftBracket:
			result, s = EvalExpr(s)
			token, v, s = NextToken(s)
			if token != TokenRightBracket {
				panic("expect ')' but found " + v)
			}

			return result, s

		default:
			return result, v + s
		}
	}
}

func main() {
	result, s := EvalExpr("1+2*3/(4-5)+6")
	if len(s) > 0 {
		panic("unexpected " + s + " at end of expression")
	}
	fmt.Println(result)
}
