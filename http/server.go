package http

import (
	//_ "YaGoCalc/cmd/Calculator/docs"
	"YaGoCalc/http/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func StartServer() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Регистрация обработчиков
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/calculate", handler.CalculateHandler)
	e.POST("/task", handler.AddTask)
	e.GET("/tasks", handler.GetTasks)
	e.GET("/task/:id", handler.GetTaskByID)
	e.POST("/task/:id/result", handler.ReceiveTaskResult)

	// Запуск сервера
	e.Logger.Fatal(e.Start(":8080"))
}
