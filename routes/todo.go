package routes

import (
	"go-fiber-learning/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(api fiber.Router) {
	api.Get("/todo", handlers.GetAllTodos)
	api.Get("/todo/:id", handlers.GetTodo)
	api.Post("/todo", handlers.InsertTodo)
	api.Patch("/todo:id", handlers.UpdateTodo)
	api.Delete("todo/:id", handlers.DeleteTodo)
}
