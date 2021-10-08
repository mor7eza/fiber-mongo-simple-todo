package utils

import "github.com/gofiber/fiber/v2"

func Response(status bool, msg string, data interface{}) fiber.Map {
	return fiber.Map{
		"success": status,
		"message": msg,
		"data":    data,
	}
}
