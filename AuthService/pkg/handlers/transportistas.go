package handlers

import (
	"goServices/pkg/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetTransportistas obtiene lista de transportistas con paginación y filtros
func GetTransportistas(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("page_size", 10)
		estado := c.Query("estado")
		ciudad := c.Query("ciudad")
		calificacionMin := c.Query("calificacion_min")

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize

		// Base query
		query := db.Joins("LEFT JOIN users ON transportistas.id_usuario = users.id")

		// Filtro por estado
		if estado != "" {
			query = query.Where("transportistas.estado = ?", estado)
		}

		// Filtro por ciudad
		if ciudad != "" {
			query = query.Where("zonas.ciudad = ?", ciudad)
		}

		// Filtro por calificación mínima
		if calificacionMin != "" {
			minCalif, err := strconv.ParseFloat(calificacionMin, 64)
			if err == nil {
				query = query.Where("transportistas.calificacion_promedio >= ?", minCalif)
			}
		}

		// Solo transportistas activos o con verificación pending si no está filtrado
		if estado == "" {
			query = query.Where("transportistas.estado IN ?", []models.EstadoTransportista{models.EstadoActivo, models.EstadoVerificacionPendiente})
		}

		var total int64
		if err := query.Model(&models.Transportista{}).Count(&total).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count transportistas"})
		}

		var transportistas []models.Transportista
		if err := query.
			Preload("Usuario").
			Offset(offset).
			Limit(pageSize).
			Find(&transportistas).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch transportistas"})
		}

		totalPages := (int(total) + pageSize - 1) / pageSize

		response := models.TransportistaListResponse{
			Data:       transportistas,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		}

		return c.JSON(response)
	}
}

// GetTransportista obtiene detalles de un transportista específico
func GetTransportista(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		transportistaID := c.Params("id_transportista")
		_, err := uuid.Parse(transportistaID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transportista ID"})
		}

		var transportista models.Transportista
		if err := db.Preload("Usuario").First(&transportista, "id_transportista = ?", transportistaID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transportista not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		return c.JSON(transportista)
	}
}
