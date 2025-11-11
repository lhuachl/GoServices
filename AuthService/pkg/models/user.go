package models

import (
	"time"

	"github.com/google/uuid"
)

// RolUsuario define los roles disponibles
type RolUsuario string

const (
	RolCliente       RolUsuario = "cliente"
	RolTransportista RolUsuario = "transportista"
	RolAdmin         RolUsuario = "admin"
)

// User modelo de usuario (sincronizado con auth.users de Supabase)
type User struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Nombre     string     `json:"nombre" gorm:"type:text"`
	Apellido   string     `json:"apellido" gorm:"type:text"`
	Rol        string     `json:"rol" gorm:"type:varchar(20);default:'cliente'"`
	FotoPerfil string     `json:"foto_perfil"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// PerfilCliente perfil adicional del cliente
type PerfilCliente struct {
	IDPerfil            uuid.UUID `json:"id_perfil" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IDUsuario           uuid.UUID `json:"id_usuario" gorm:"type:uuid;uniqueIndex"`
	DocumentoIdentidad  string    `json:"documento_identidad" gorm:"uniqueIndex"`
	Telefono            string    `json:"telefono"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	Usuario     *User        `json:"usuario,omitempty" gorm:"foreignKey:IDUsuario"`
	Direcciones []Direccion  `json:"direcciones,omitempty"`
}

// Direccion dirección del cliente
type Direccion struct {
	IDDireccion            uuid.UUID `json:"id_direccion" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IDPerfil               uuid.UUID `json:"id_perfil" gorm:"type:uuid"`
	Calle                  string    `json:"calle"`
	Ciudad                 string    `json:"ciudad"`
	ReferenciasAdicionales string    `json:"referencias_adicionales"`
	Pais                   string    `json:"pais" gorm:"default:'Ecuador'"`
	Latitud                float64   `json:"latitud"`
	Longitud               float64   `json:"longitud"`
	EsPredeterminada       bool      `json:"es_predeterminada" gorm:"default:false"`
	CreatedAt              time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relación
	PerfilCliente *PerfilCliente `json:"-" gorm:"foreignKey:IDPerfil"`
}

// CreateUserRequest DTO para crear usuario
type CreateUserRequest struct {
	Nombre  string `json:"nombre" binding:"required"`
	Apellido string `json:"apellido" binding:"required"`
}

// UpdateUserRequest DTO para actualizar usuario
type UpdateUserRequest struct {
	Nombre     string `json:"nombre"`
	Apellido   string `json:"apellido"`
	FotoPerfil string `json:"foto_perfil"`
}

// CreateDireccionRequest DTO para crear dirección
type CreateDireccionRequest struct {
	Calle                  string  `json:"calle" binding:"required"`
	Ciudad                 string  `json:"ciudad" binding:"required"`
	ReferenciasAdicionales string  `json:"referencias_adicionales"`
	Pais                   string  `json:"pais" binding:"required"`
	Latitud                float64 `json:"latitud" binding:"required"`
	Longitud               float64 `json:"longitud" binding:"required"`
	EsPredeterminada       bool    `json:"es_predeterminada"`
}

// UpdateDireccionRequest DTO para actualizar dirección
type UpdateDireccionRequest struct {
	Calle                  string  `json:"calle"`
	Ciudad                 string  `json:"ciudad"`
	ReferenciasAdicionales string  `json:"referencias_adicionales"`
	Pais                   string  `json:"pais"`
	Latitud                float64 `json:"latitud"`
	Longitud               float64 `json:"longitud"`
	EsPredeterminada       bool    `json:"es_predeterminada"`
}
