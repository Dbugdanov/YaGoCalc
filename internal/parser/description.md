Expression Parser
## задачи

* количество открывающих скобок равно количеству закрывающих.
* Целая часть числа отделена от дробной с помощью точки.
* В строке присутствуют только допустимые символы: цифры 0...9, операторы +-*/, скобки, точка.

## Описание
Парсер должен построить дерево лексем, пригодное для вычисления численного значения входного выражения. Значения параметров 
будем передавать методу, реализующему расчет численного значения.


Каждая цельная лексема представляет собой скобку, оператор либо операнд.

## Лексема как объект

Каждая цельная лексема в составе древовидной структуры описывается объектом. Любой объект «дерева» обладает набором свойств (полей) и определенным поведением. Перечислим предполагаемый набор полей данных для каждого объекта, моделирующего лексему:


```Go
package main

type TokenType int

const (
	NumberToken TokenType = iota
	OperatorToken
	ParenthesisToken
) 

type Token struct {
	Type TokenType
	Value string
	
}
```

## Структура дерева выражений

Для представления дерева выражений создадим интерфейс Expr и конкретные типы, реализующие этот интерфейс:

```go
package main

type Expr interface {
    Evaluate() float64 // Метод для вычисления численного значения выражения
}

type Number struct {
    Value float64
}

func (n Number) Evaluate() float64 {
    return n.Value
}

type BinaryOp struct {
    Left  Expr
    Op    rune // Используем rune для представления оператора
    Right Expr
}

func (b BinaryOp) Evaluate() float64 {
    switch b.Op {
    case '+':
        return b.Left.Evaluate() + b.Right.Evaluate()
    case '-':
        return b.Left.Evaluate() - b.Right.Evaluate()
    case '*':
        return b.Left.Evaluate() * b.Right.Evaluate()
    case '/':
        return b.Left.Evaluate() / b.Right.Evaluate()
    default:
        panic("unsupported operator")
    }
}

```

## Парсер

Парсер будет анализировать строку и строить на её основе дерево выражений. 
Для этого сначала реализуем лексический анализатор, который преобразует входную строку в 
список токенов, а затем синтаксический анализатор, который построит дерево выражений на основе списка токенов.

## Лексический анализ

Реализация функции Tokenize будет включать в себя итерацию по символам строки и классификацию каждого символа как числа, оператора или скобки. Важно также обрабатывать числа с плавающей точкой.
В этом коде числа считываются полностью, включая дробную часть, и преобразуются в float64. Операторы и скобки добавляются в список токенов как есть. Предполагается, что в Token добавлено поле Number float64 для хранения числового значения токена-числа.
```go
import (
    "strconv"
    "unicode"
)

func Tokenize(input string) ([]Token, error) {
    var tokens []Token
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
            tokens = append(tokens, Token{Type: NumberToken, Value: numStr, Number: value})
        } else if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
            // Обработка оператора
            tokens = append(tokens, Token{Type: OperatorToken, Value: string(ch)})
        } else if ch == '(' || ch == ')' {
            // Обработка скобок
            tokens = append(tokens, Token{Type: ParenthesisToken, Value: string(ch)})
        }
        // Пропускаем пробелы и другие незначащие символы
    }
    return tokens, nil
}
```

## Синтаксический анализ
Для синтаксического анализа воспользуемся методом рекурсивного спуска, который позволяет легко учитывать приоритет операций за счет структуры вызова функций.
Чтобы избежать ошибок index out of range, мы должны внимательно проверять наличие токенов перед их использованием в каждой функции анализа.
```go
func Parse(tokens []Token) (Expr, error) {
    if len(tokens) == 0 {
        return nil, fmt.Errorf("no tokens to parse")
    }
    expr, _, err := parseAddSub(tokens)
    if err != nil {
        return nil, err
    }
    return expr, nil
}

// parseAddSub обрабатывает сложение и вычитание
func parseAddSub(tokens []Token) (Expr, []Token, error) {
    if len(tokens) == 0 {
        return nil, tokens, fmt.Errorf("expected an expression but got nothing")
    }

    expr, tokens, err := parseMulDiv(tokens)
    if err != nil {
        return nil, tokens, err
    }

    for len(tokens) > 0 && (tokens[0].Type == OperatorToken && (tokens[0].Value == "+" || tokens[0].Value == "-")) {
        if len(tokens) < 2 {
            return nil, tokens, fmt.Errorf("missing right operand after %s", tokens[0].Value)
        }
        op := tokens[0]
        var right Expr
        right, tokens, err = parseMulDiv(tokens[1:])
        if err != nil {
            return nil, tokens, err
        }
        expr = BinaryOp{Left: expr, Op: rune(op.Value[0]), Right: right}
    }

    return expr, tokens, nil
}

// parseMulDiv обрабатывает умножение и деление
func parseMulDiv(tokens []Token) (Expr, []Token, error) {
    expr, tokens, err := parseFactor(tokens)
    if err != nil {
        return nil, nil, err
    }

    for len(tokens) > 0 && (tokens[0].Type == OperatorToken && (tokens[0].Value == "*" || tokens[0].Value == "/")) {
        op := tokens[0]
        var right Expr
        right, tokens, err = parseFactor(tokens[1:])
        if err != nil {
            return nil, nil, err
        }
        expr = BinaryOp{Left: expr, Op: rune(op.Value[0]), Right: right}
    }

    return expr, tokens, nil
}

// parseFactor обрабатывает числа и выражения в скобках
func parseFactor(tokens []Token) (Expr, []Token, error) {
    if len(tokens) == 0 {
        return nil, tokens, fmt.Errorf("expected a number or a parenthesis but got nothing")
    }

    if tokens[0].Type == NumberToken {
        return Number{Value: tokens[0].Number}, tokens[1:], nil
    } else if tokens[0].Type == ParenthesisToken && tokens[0].Value == "(" {
        if len(tokens) < 3 { // Нужно как минимум 3 токена для "(expr)"
            return nil, tokens, fmt.Errorf("incomplete expression within parentheses")
        }
        expr, remainingTokens, err := parseAddSub(tokens[1:]) // Пропускаем открывающую скобку
        if err != nil {
            return nil, remainingTokens, err
        }
        if len(remainingTokens) == 0 || remainingTokens[0].Type != ParenthesisToken || remainingTokens[0].Value != ")" {
            return nil, remainingTokens, fmt.Errorf("missing closing parenthesis")
        }
        return expr, remainingTokens[1:], nil // Пропускаем закрывающую скобку
    }

    return nil, tokens, fmt.Errorf("unexpected token: %v", tokens[0])
}

```