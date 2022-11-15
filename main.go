// Video de ilustración: https://www.youtube.com/watch?v=zlrnwGZMBbU
package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	// Serve files from multiple directories
	app.Static("/", "./client/dist")

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"data": "Usuario desde el backend como respuesta.",
		})
	})

	log.Fatal(app.Listen(":3000"))
	fmt.Println("Server on port 3000")
}
