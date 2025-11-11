package handlers

import (
	"goServices/pkg/middleware"
	"goServices/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetMyAddresses obtiene las direcciones del usuario autenticado
func GetMyAddresses(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var perfil models.PerfilCliente
		if err := db.First(&perfil, "id_usuario = ?", userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client profile not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		var direcciones []models.Direccion
		if err := db.Where("id_perfil = ?", perfil.IDPerfil).Find(&direcciones).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch addresses"})
		}

		return c.JSON(direcciones)
	}
}

// CreateAddress crea una nueva dirección para el usuario autenticado
func CreateAddress(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var req models.CreateDireccionRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Obtener el perfil del cliente
		var perfil models.PerfilCliente
		if err := db.First(&perfil, "id_usuario = ?", userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Client profile not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		// Si esta es la primera dirección, marcarla como predeterminada
		var count int64
		db.Model(&models.Direccion{}).Where("id_perfil = ?", perfil.IDPerfil).Count(&count)
		if count == 0 {
			req.EsPredeterminada = true
		}

		// Si se marca como predeterminada, desmarcar las otras
		if req.EsPredeterminada {
			db.Model(&models.Direccion{}).Where("id_perfil = ?", perfil.IDPerfil).Update("es_predeterminada", false)
		}

		direccion := models.Direccion{
			IDDireccion:            uuid.New(),
			IDPerfil:               perfil.IDPerfil,
			Calle:                  req.Calle,
			Ciudad:                 req.Ciudad,
			ReferenciasAdicionales: req.ReferenciasAdicionales,
			Pais:                   req.Pais,
			Latitud:                req.Latitud,
			Longitud:               req.Longitud,
			EsPredeterminada:       req.EsPredeterminada,
		}

		if err := db.Create(&direccion).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create address"})
		}

		return c.Status(fiber.StatusCreated).JSON(direccion)
	}
}

// UpdateAddress actualiza una dirección del usuario autenticado
func UpdateAddress(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		addressID := c.Params("id_direccion")
		_, err = uuid.Parse(addressID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid address ID"})
		}

		var req models.UpdateDireccionRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Verificar que la dirección pertenece al usuario
		var direccion models.Direccion
		if err := db.First(&direccion, "id_direccion = ?", addressID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Address not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		// Verificar que el perfil pertenece al usuario
		var perfil models.PerfilCliente
		if err := db.First(&perfil, "id_perfil = ?", direccion.IDPerfil).Error; err != nil || perfil.IDUsuario != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		// Si se marca como predeterminada, desmarcar las otras
		if req.EsPredeterminada {
			db.Model(&models.Direccion{}).Where("id_perfil = ? AND id_direccion != ?", perfil.IDPerfil, addressID).Update("es_predeterminada", false)
		}

		if err := db.Model(&direccion).Updates(req).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update address"})
		}

		return c.JSON(direccion)
	}
}

// DeleteAddress elimina una dirección del usuario autenticado
func DeleteAddress(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := middleware.GetUserIDFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		addressID := c.Params("id_direccion")
		_, err = uuid.Parse(addressID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid address ID"})
		}

		// Verificar que la dirección pertenece al usuario
		var direccion models.Direccion
		if err := db.First(&direccion, "id_direccion = ?", addressID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Address not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		// Verificar que el perfil pertenece al usuario
		var perfil models.PerfilCliente
		if err := db.First(&perfil, "id_perfil = ?", direccion.IDPerfil).Error; err != nil || perfil.IDUsuario != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		if err := db.Delete(&direccion).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete address"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
