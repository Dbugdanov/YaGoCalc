package handler

import (
	expression "YaGoCalc/internal/parser"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

type request struct {
	Expression string `json:"expression"`
}

type response struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

var tasks sync.Map

// AddTask добавляет задачу на вычисление арифметического выражения
//func AddTask(c echo.Context) error {
//	// Логика добавления задачи
//}

// GetTasks возвращает список всех задач
//func GetTasks(c echo.Context) error {
//	// Логика получения списка задач
//}

// GetTaskByID возвращает задачу по её идентификатору
//func GetTaskByID(c echo.Context) error {
//	// Логика получения задачи по ID
//}

// ReceiveTaskResult принимает результат выполнения задачи
//func ReceiveTaskResult(c echo.Context) error {
//	// Логика приема результата задачи
//}

// CalculateHandler обрабатывает POST запросы на вычисление арифметического выражения
func CalculateHandler(c echo.Context) error {
	req := new(request)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response{Error: err.Error()})
	}

	result, err := expression.EvaluateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, response{Result: result})
}
