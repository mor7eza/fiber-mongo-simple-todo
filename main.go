package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	//Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	//Connect to database
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MongoDB uri didn't set. Please check .env file")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	coll := client.Database("sample_todos").Collection("todos")

	//Init fiber
	app := fiber.New(fiber.Config{
		AppName: "Go Fiber Learning",
	})

	//Endpoints
	app.Get("/todo", func(c *fiber.Ctx) error {
		cursor, err := coll.Find(context.TODO(), bson.D{{}})
		if err != nil {
			if err == mongo.ErrNoDocuments {
				//Check for empty result
				return nil
			}
			panic(err)
		}

		var todos []Todo = make([]Todo, 0)

		if err = cursor.All(context.TODO(), &todos); err != nil {
			panic(err)
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "All data retrived ",
			"data":    todos,
		})
	})

	app.Get("/todo/:id", func(c *fiber.Ctx) error {
		coll := client.Database("sample_todos").Collection("todos")

		id, _ := primitive.ObjectIDFromHex(c.Params("id"))

		result := coll.FindOne(context.TODO(), bson.D{{"_id", id}})

		todo := new(Todo)

		result.Decode(&todo)

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Todo finded",
			"data":    todo,
		})
	})

	app.Post("/todo", func(c *fiber.Ctx) error {
		coll := client.Database("sample_todos").Collection("todos")

		todo := new(Todo)

		if err := c.BodyParser(todo); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "invalid data recieved",
			})
		}

		todo.ID = ""

		insertionResult, err := coll.InsertOne(c.Context(), todo)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "invalid data recieved",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"message": "Todo inserted successfully",
			"data":    insertionResult.InsertedID,
		})
	})

	app.Patch("/todo/:id", func(c *fiber.Ctx) error {

		coll := client.Database("sample_todos").Collection("todos")

		id, _ := primitive.ObjectIDFromHex(c.Params("id"))

		todo := new(Todo)

		c.BodyParser(todo)

		_, err := coll.UpdateOne(context.TODO(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"title", todo.Title}}}})

		if err != nil {
			panic(err)
		}
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Updated successfully",
			"data":    nil,
		})
	})

	app.Delete("/todo/:id", func(c *fiber.Ctx) error {
		coll := client.Database("sample_todos").Collection("todos")

		id, _ := primitive.ObjectIDFromHex(c.Params("id"))

		_, err := coll.DeleteOne(context.TODO(), bson.D{{"_id", id}})

		if err != nil {
			panic(err)
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Record deleted",
			"data":    nil,
		})
	})

	app.Listen(":3000")
}
