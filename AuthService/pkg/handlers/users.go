package handlers

import (
	"goServices/pkg/middleware"
	"goServices/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetMe obtiene la info del usuario autenticado
func GetMe(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var user models.User
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		return c.JSON(user)
	}
}

// UpdateMe actualiza la info del usuario autenticado
func UpdateMe(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var req models.UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		var user models.User
		if err := db.Model(&user).Where("id = ?", userID).Updates(req).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
		}

		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated user"})
		}

		return c.JSON(user)
	}
}

// GetUser obtiene info de un usuario específico (solo el propio usuario o admins)
func GetUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		userIDParam := c.Params("id_usuario")
		targetUserID, err := uuid.Parse(userIDParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
		}

		// Solo el propio usuario o admins pueden acceder
		if userID != targetUserID {
			var currentUser models.User
			if err := db.First(&currentUser, "id = ?", userID).Error; err != nil || currentUser.Rol != "admin" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
			}
		}

		var user models.User
		if err := db.First(&user, "id = ?", targetUserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		return c.JSON(user)
	}
}
