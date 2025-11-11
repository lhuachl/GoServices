package main

import (
	"goServices/pkg/handlers"
	"goServices/pkg/middleware"
	"goServices/pkg/models"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func setupRoutes(app *fiber.App, db *gorm.DB) {
	// Rutas públicas
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Rutas autenticadas
	api := app.Group("/api", middleware.AuthMiddleware)

	// Users endpoints
	api.Get("/users/me", handlers.GetMe(db))
	api.Put("/users/me", handlers.UpdateMe(db))
	api.Get("/users/:id_usuario", handlers.GetUser(db))

	// Addresses endpoints
	api.Get("/users/me/addresses", handlers.GetMyAddresses(db))
	api.Post("/users/me/addresses", handlers.CreateAddress(db))
	api.Put("/users/me/addresses/:id_direccion", handlers.UpdateAddress(db))
	api.Delete("/users/me/addresses/:id_direccion", handlers.DeleteAddress(db))

	// Transportistas endpoints
	api.Get("/transportistas", handlers.GetTransportistas(db))
	api.Get("/transportistas/:id_transportista", handlers.GetTransportista(db))
}

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Conectar a la base de datos
	db, err := GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Realizar migraciones solo para tablas que no existen
	// Si las tablas ya existen con RLS, GORM no puede modificarlas
	if err := db.AutoMigrate(
		&models.PerfilCliente{},
		&models.Direccion{},
		&models.Transportista{},
	); err != nil {
		log.Printf("Warning during migrations: %v (puede ser por RLS en Supabase)", err)
		// No paniquear, continuamos
	}

	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName: "Transport Services API",
	})

	// Middleware global
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Configurar rutas
	setupRoutes(app, db)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
