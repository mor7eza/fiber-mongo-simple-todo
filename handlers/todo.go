package handlers

import (
	"context"
	"go-fiber-learning/database"
	"go-fiber-learning/models"
	"go-fiber-learning/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllTodos(c *fiber.Ctx) error {
	cursor, err := database.DB.Collection("todos").Find(context.TODO(), bson.D{{}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//Check for empty result
			return c.Status(204).JSON(fiber.Map{
				"success": true,
				"message": "There is no Todo in databse",
				"data":    nil,
			})
		}
		panic(err)
	}

	var todos []models.Todo = make([]models.Todo, 0)

	if err = cursor.All(context.TODO(), &todos); err != nil {
		panic(err)
	}

	return c.Status(200).JSON(
		utils.Response(true, "All Data Retrieved.", todos),
	)
}

func GetTodo(c *fiber.Ctx) error {

	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	result := database.DB.Collection("todos").FindOne(context.TODO(), bson.D{{"_id", id}})

	todo := new(models.Todo)

	result.Decode(&todo)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Todo finded",
		"data":    todo,
	})
}

func InsertTodo(c *fiber.Ctx) error {

	todo := new(models.Todo)

	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "invalid data recieved",
			"data":    nil,
		})
	}

	todo.ID = ""

	insertionResult, err := database.DB.Collection("todos").InsertOne(c.Context(), todo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Server error. Todo didn't inserted",
			"data":    nil,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Todo inserted successfully",
		"data":    insertionResult.InsertedID,
	})
}

func UpdateTodo(c *fiber.Ctx) error {

	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	todo := new(models.Todo)

	c.BodyParser(todo)

	_, err := database.DB.Collection("todos").UpdateOne(context.TODO(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"title", todo.Title}}}})

	if err != nil {
		panic(err)
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Updated successfully",
		"data":    nil,
	})
}

func DeleteTodo(c *fiber.Ctx) error {

	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	_, err := database.DB.Collection("todos").DeleteOne(context.TODO(), bson.D{{"_id", id}})

	if err != nil {
		panic(err)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Record deleted",
		"data":    nil,
	})
}
