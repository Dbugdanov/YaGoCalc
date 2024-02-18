package expression

import (
	"YaGoCalc/models"
	"fmt"
	"strconv"
	"unicode"
)

func Tokenize(input string) ([]models.Token, error) {
	var tokens []models.Token
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if unicode.IsDigit(rune(ch)) || ch == '.' {
			// Обработка числа
			numStr := string(ch)
			for i+1 < len(input) && (unicode.IsDigit(rune(input[i+1])) || input[i+1] == '.') {
				i++
				numStr += string(input[i])
			}
			value, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, models.Token{Type: models.NumberToken, Value: numStr, Number: value})
		} else if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			// Обработка оператора
			tokens = append(tokens, models.Token{Type: models.OperatorToken, Value: string(ch)})
		} else if ch == '(' || ch == ')' {
			// Обработка скобок
			tokens = append(tokens, models.Token{Type: models.ParenthesisToken, Value: string(ch)})
		}
		// Пропускаем пробелы и другие незначащие символы
	}
	return tokens, nil
}

func Parse(tokens []models.Token) (models.Expr, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}
	expr, _, err := parseAddSub(tokens)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func parseAddSub(tokens []models.Token) (models.Expr, []models.Token, error) {
	if len(tokens) == 0 {
		return nil, tokens, fmt.Errorf("expected an expression but got nothing")
	}

	expr, tokens, err := parseMulDiv(tokens)
	if err != nil {
		return nil, tokens, err
	}

	for len(tokens) > 0 && (tokens[0].Type == models.OperatorToken && (tokens[0].Value == "+" || tokens[0].Value == "-")) {
		if len(tokens) < 2 {
			return nil, tokens, fmt.Errorf("missing right operand after %s", tokens[0].Value)
		}
		op := tokens[0]
		var right models.Expr
		right, tokens, err = parseMulDiv(tokens[1:])
		if err != nil {
			return nil, tokens, err
		}
		expr = models.BinaryOp{Left: expr, Op: rune(op.Value[0]), Right: right}
	}

	return expr, tokens, nil
}

func parseMulDiv(tokens []models.Token) (models.Expr, []models.Token, error) {
	if len(tokens) == 0 {
		return nil, tokens, fmt.Errorf("expected an expression but got nothing")
	}

	expr, tokens, err := parseFactor(tokens)
	if err != nil {
		return nil, tokens, err
	}

	for len(tokens) > 0 && (tokens[0].Type == models.OperatorToken && (tokens[0].Value == "*" || tokens[0].Value == "/")) {
		if len(tokens) < 2 {
			return nil, tokens, fmt.Errorf("missing right operand after %s", tokens[0].Value)
		}
		op := tokens[0]
		var right models.Expr
		right, tokens, err = parseFactor(tokens[1:])
		if err != nil {
			return nil, tokens, err
		}
		expr = models.BinaryOp{Left: expr, Op: rune(op.Value[0]), Right: right}
	}

	return expr, tokens, nil
}

func parseFactor(tokens []models.Token) (models.Expr, []models.Token, error) {
	if len(tokens) == 0 {
		return nil, tokens, fmt.Errorf("expected a number or a parenthesis but got nothing")
	}

	if tokens[0].Type == models.NumberToken {
		return models.Number{Value: tokens[0].Number}, tokens[1:], nil
	} else if tokens[0].Type == models.ParenthesisToken && tokens[0].Value == "(" {
		if len(tokens) < 3 { // Нужно как минимум 3 токена для "(expr)"
			return nil, tokens, fmt.Errorf("incomplete expression within parentheses")
		}
		expr, remainingTokens, err := parseAddSub(tokens[1:]) // Пропускаем открывающую скобку
		if err != nil {
			return nil, remainingTokens, err
		}
		if len(remainingTokens) == 0 || remainingTokens[0].Type != models.ParenthesisToken || remainingTokens[0].Value != ")" {
			return nil, remainingTokens, fmt.Errorf("missing closing parenthesis")
		}
		return expr, remainingTokens[1:], nil // Пропускаем закрывающую скобку
	}

	return nil, tokens, fmt.Errorf("unexpected token: %v", tokens[0])
}

// EvaluateExpression принимает строку с математическим выражением,
// парсит её и возвращает результат вычисления или ошибку.
func EvaluateExpression(expression string) (float64, error) {
	tokens, err := Tokenize(expression)
	if err != nil {
		return 0, fmt.Errorf("tokenize error: %v", err)
	}

	ast, err := Parse(tokens)
	if err != nil {
		return 0, fmt.Errorf("parse error: %v", err)
	}

	return ast.Evaluate(), nil
}
