// Referencia en https://github.com/gofiber/recipes/blob/master/mysql/main.go
// Api Rest con Fiber framework de Go.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// Database Connection
func conexionBD() (conexion *sql.DB) {

	godotenv.Load()

	Driver := os.Getenv("Driver")
	Usuario := os.Getenv("Usuario")
	Contrasena := os.Getenv("Contrasena")
	Nombre := os.Getenv("Nombre")

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasena+"@tcp(0.0.0.0:3306)/"+Nombre)
	if err != nil {
		log.Fatal(err.Error())
	}
	return conexion
}

// Servicio struct
type Servicio struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Correo      string `json:"correo"`
	Telefono    string `json:"telefono"`
	Equipo      string `json:"equipo"`
	Diagnostico string `json:"diagnostico"`
	Resultados  string `json:"resultados"`
	Decision    string `json:"decision"`
	Taller      string `json:"taller"`
	Servicio    string `json:"servicio"`
	Entrega     string `json:"entrega"`
	Fecha       string `json:"fecha"`
}

// Servicios struct
type Servicios struct {
	Servicios []Servicio `json:"servicios"`
}

func main() {
	// .env
	godotenv.Load()
	// Create a Fiber app
	app := fiber.New()

	// Cors Origin
	app.Use(cors.New())

	// Logger
	app.Use(logger.New())
	// Os .env
	port := os.Getenv("Port")

	if port == "" {
		port = "Port"
	}

	// Get all records from MySQL
	app.Get("/servicio", Get)

	// Get one record from MySQL
	app.Get("servicio/:id", GetServicio)

	// Add record into MySQL
	app.Post("/servicio", Post)

	// Update record into MySQL
	app.Put("/servicio/:id", Put)

	// Delete record from MySQL
	app.Delete("/servicio/:id", Delete)

	// Server running
	fmt.Println("Server running on Port" + port + " http://0.0.0.0:4002/servicio")
	log.Fatal(app.Listen(port))
}

// Code of successful
// Get all records from MySQL
func Get(c *fiber.Ctx) error {
	// Database Connection
	conexionEstablecida := conexionBD()
	// Get Employee list from database
	rows, err := conexionEstablecida.Query("SELECT id, nombre, correo, telefono, equipo, diagnostico, resultados, decision, taller, servicio, entrega, fecha FROM taller")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()
	result := Servicios{}

	for rows.Next() {
		servicio := Servicio{}
		if err := rows.Scan(&servicio.ID, &servicio.Nombre, &servicio.Correo, &servicio.Telefono, &servicio.Equipo, &servicio.Diagnostico, &servicio.Resultados, &servicio.Decision, &servicio.Taller, &servicio.Servicio, &servicio.Entrega, &servicio.Fecha); err != nil {
			return err // Exit if we get an error
		}

		// Append Employee to Employees
		result.Servicios = append(result.Servicios, servicio)
	}
	// Return Employees in JSON format
	return c.Status(fiber.StatusOK).JSON(result)
}

// Get a single record from MySQL
func GetServicio(c *fiber.Ctx) error {
	ID := c.Params("id")
	// Database Connection
	conexionEstablecida := conexionBD()
	// Get Employee list from database
	rows, err := conexionEstablecida.Query("SELECT id, nombre, correo, telefono, equipo, diagnostico, resultados, decision, taller, servicio, entrega, fecha FROM taller WHERE id=?", ID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()
	result := Servicios{}

	for rows.Next() {
		servicio := Servicio{}
		if err := rows.Scan(&servicio.ID, &servicio.Nombre, &servicio.Correo, &servicio.Telefono, &servicio.Equipo, &servicio.Diagnostico, &servicio.Resultados, &servicio.Decision, &servicio.Taller, &servicio.Servicio, &servicio.Entrega, &servicio.Fecha); err != nil {
			return err // Exit if we get an error
		}

		// Append Employee to Employees
		result.Servicios = append(result.Servicios, servicio)
	}
	// Return Employees in JSON format
	return c.Status(fiber.StatusOK).JSON(result)
}

// Successful Code
// Add record into MySQL
func Post(c *fiber.Ctx) error {
	// Database Connection
	conexionEstablecida := conexionBD()
	// New Employee struct
	u := new(Servicio)

	// Parse body into struct
	if err := c.BodyParser(u); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Insert Employee into database
	res, err := conexionEstablecida.Query("INSERT INTO taller (nombre, correo, telefono, equipo, diagnostico, resultados, decision, taller, servicio, entrega, fecha) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", u.Nombre, u.Correo, u.Telefono, u.Equipo, u.Diagnostico, u.Resultados, u.Decision, u.Taller, u.Servicio, u.Entrega, u.Fecha)
	if err != nil {
		return err
	}

	defer res.Close()
	// Print result
	log.Println(res)

	// Return Employee in JSON format
	return c.Status(fiber.StatusOK).JSON(u)
}

// Successful Code
// Update record into MySQL
func Put(c *fiber.Ctx) error {
	ID := c.Params("id")
	conexionEstablecida := conexionBD()
	// Obtener los datos del nuevo usuario desde el cuerpo de la solicitud
	var servicio Servicio
	if err := c.BodyParser(&servicio); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al parsear el cuerpo de la solicitud",
		})
	}

	// Consulta SQL para actualizar el registro
	result, err := conexionEstablecida.Exec("UPDATE taller SET nombre=?,correo=?,telefono=?,equipo=?,diagnostico=?,resultados=?,decision=?,taller=?,servicio=?,entrega=?,fecha=? WHERE id=?", servicio.Nombre, servicio.Correo, servicio.Telefono, servicio.Equipo, servicio.Diagnostico, servicio.Resultados, servicio.Decision, servicio.Taller, servicio.Servicio, servicio.Entrega, servicio.Fecha, ID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al actualizar el registro",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al obtener el número de filas afectadas",
		})
	}

	if rowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Registro no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Registro actualizado correctamente",
	})

}

// Successful Code
// Delete record from MySQL
func Delete(c *fiber.Ctx) error {
	ID := c.Params("id")

	conexionEstablecida := conexionBD()
	// Consulta SQL para eliminar el registro
	result, err := conexionEstablecida.Exec("DELETE FROM taller WHERE id=?", ID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al eliminar el registro",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al obtener el número de filas afectadas",
		})
	}

	if rowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Registro no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Registro eliminado correctamente",
	})
}
