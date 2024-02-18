package http

import (
	//_ "YaGoCalc/cmd/Calculator/docs"
	"YaGoCalc/http/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

func StartServer() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Регистрация обработчиков
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/calculate", handler.CalculateHandler)
	//e.POST("/task", handler.AddTask)
	//e.GET("/tasks", handler.GetTasks)
	//e.GET("/task/:id", handler.GetTaskByID)
	//e.POST("/task/:id/result", handler.ReceiveTaskResult)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
