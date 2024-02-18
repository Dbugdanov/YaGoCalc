package models

type TokenType int

const (
	NumberToken TokenType = iota
	OperatorToken
	ParenthesisToken
)

type Token struct {
	Type   TokenType
	Value  string
	Number float64
}

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
