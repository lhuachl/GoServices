package models

import (
	"time"

	"github.com/google/uuid"
)

// EstadoTransportista define los estados del transportista
type EstadoTransportista string

const (
	EstadoVerificacionPendiente EstadoTransportista = "verificacion_pendiente"
	EstadoActivo                EstadoTransportista = "activo"
	EstadoInactivo              EstadoTransportista = "inactivo"
	EstadoSuspendido            EstadoTransportista = "suspendido"
)

// Transportista modelo de transportista
type Transportista struct {
	IDTransportista      uuid.UUID  `json:"id_transportista" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	IDUsuario            uuid.UUID  `json:"id_usuario" gorm:"type:uuid;uniqueIndex"`
	TipoVehiculo         string     `json:"tipo_vehiculo"`
	PlacaVehiculo        string     `json:"placa_vehiculo" gorm:"uniqueIndex"`
	CapacidadCarga       float64    `json:"capacidad_carga"`
	Estado               string     `json:"estado" gorm:"type:varchar(30);default:'verificacion_pendiente'"`
	IDZonaAsignada       *uuid.UUID `json:"id_zona_asignada"`
	CalificacionPromedio float64    `json:"calificacion_promedio" gorm:"default:0.0"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relaciones
	Usuario *User `json:"usuario,omitempty" gorm:"foreignKey:IDUsuario"`
}

// CreateTransportistaRequest DTO para crear transportista
type CreateTransportistaRequest struct {
	Nombre           string  `json:"nombre" binding:"required"`
	Apellido         string  `json:"apellido" binding:"required"`
	TipoVehiculo     string  `json:"tipo_vehiculo" binding:"required"`
	PlacaVehiculo    string  `json:"placa_vehiculo" binding:"required"`
	CapacidadCarga   float64 `json:"capacidad_carga" binding:"required,gt=0"`
}

// UpdateTransportistaRequest DTO para actualizar transportista
type UpdateTransportistaRequest struct {
	TipoVehiculo      string  `json:"tipo_vehiculo"`
	PlacaVehiculo     string  `json:"placa_vehiculo"`
	CapacidadCarga    float64 `json:"capacidad_carga"`
	Estado            string  `json:"estado"`
	IDZonaAsignada    *uuid.UUID `json:"id_zona_asignada"`
}

// TransportistaListResponse respuesta con paginación
type TransportistaListResponse struct {
	Data       []Transportista `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
