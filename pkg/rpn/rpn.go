package rpn

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	//fmt.Println(tokens)
	postfix, err := infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}
	return evaluatePostfix(postfix)
}

func tokenize(expression string) []string {
	var tokens []string
	var number strings.Builder

	for _, char := range expression {
		switch char {
		case '+', '-', '*', '/', '(', ')':
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
			}
			tokens = append(tokens, string(char))
		case ' ':
			continue
		default:
			number.WriteRune(char)
		}
	}

	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else {
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				stack = append(stack, a/b)
			default:
				return 0, errors.New("unknown operator")
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}

func isNumber(s string) bool {
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return false
	}
	return true
}

func main() {
	expression := "3 + (4.5 * 5) - (10 / 2)"
	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Result: %.2f\n", result) // Output: Result: 18.00
	}

	// Пример с ошибкой
	expressionWithError := "3 * (5-2)"
	resultWithError, err := Calc(expressionWithError)
	if err != nil {
		fmt.Println("Error:", err) // Output: Error: mismatched parentheses
	} else {
		fmt.Printf("Result: %.2f\n", resultWithError)
	}
}
