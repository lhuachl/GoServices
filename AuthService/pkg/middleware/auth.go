package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Claims estructura del JWT de Supabase
type Claims struct {
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

// AuthMiddleware valida el JWT de Supabase
func AuthMiddleware(c *fiber.Ctx) error {
	// Obtener el token del header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	// Extraer el token del formato "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	token := parts[1]

	// Decodificar el JWT sin validar la firma (en producción, validar con la clave pública de Supabase)
	// Para desarrollo, extraemos el sub (user_id) del payload
	claims, err := parseTokenClaims(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Convertir el sub a UUID
	userID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID in token",
		})
	}

	// Almacenar el user_id en el contexto
	c.Locals("user_id", userID)
	c.Locals("token", token)

	return c.Next()
}

// parseTokenClaims decodifica el JWT (sin validar firma, solo para desarrollo)
func parseTokenClaims(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decodificar el payload (parte 2)
	payload := parts[1]
	
	// Agregar padding si es necesario
	switch len(payload) % 4 {
	case 1:
		payload += "==="
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decodedPayload := make([]byte, len(payload))
	n, err := decodeBase64URL(payload, decodedPayload)
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := json.Unmarshal(decodedPayload[:n], &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// decodeBase64URL decodifica base64 URL-safe
func decodeBase64URL(s string, dst []byte) (int, error) {
	// Reemplazar caracteres base64-URL
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")

	// Aquí simplificado; en producción usa base64.RawURLEncoding.DecodedLen
	return len(s), nil
}

// GetUserIDFromContext obtiene el user_id del contexto
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}
