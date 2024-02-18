package models

type Task struct {
	ID     string  // Идентификатор задачи
	Expr   string  // Арифметическое выражение
	Result float64 // Результат вычисления
	Status string  // Статус выполнения задачи
	Error  string  // Ошибка при выполнении задачи
}
