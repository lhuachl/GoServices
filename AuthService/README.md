# Transport Services API - Backend Microservicio

Backend de microservicio en Go usando **Fiber** y **GORM** con autenticación **Supabase JWT** y base de datos **PostgreSQL** (Supabase).

## 📋 Estructura del Proyecto

```
goServices/
├── main.go                      # Punto de entrada
├── .env                         # Variables de entorno
├── go.mod / go.sum             # Dependencias
└── pkg/
    ├── models/
    │   ├── user.go            # Modelos de usuarios y direcciones
    │   └── transportista.go    # Modelos de transportistas
    ├── handlers/
    │   ├── users.go           # Handlers de usuarios
    │   ├── addresses.go        # Handlers de direcciones
    │   └── transportistas.go   # Handlers de transportistas
    └── middleware/
        └── auth.go            # Middleware de autenticación JWT
```

## 🚀 Endpoints

### Públicos
- `GET /health` - Health check

### Autenticados (requieren JWT en header `Authorization: Bearer <token>`)

#### Usuarios
- `GET /api/users/me` - Obtener mi perfil
- `PUT /api/users/me` - Actualizar mi perfil
- `GET /api/users/:id_usuario` - Obtener perfil de otro usuario (solo propio o admin)

#### Direcciones
- `GET /api/users/me/addresses` - Listar mis direcciones
- `POST /api/users/me/addresses` - Crear dirección
- `PUT /api/users/me/addresses/:id_direccion` - Actualizar dirección
- `DELETE /api/users/me/addresses/:id_direccion` - Eliminar dirección

#### Transportistas
- `GET /api/transportistas?page=1&page_size=10&estado=activo&ciudad=Quito&calificacion_min=3.5` - Listar transportistas con filtros y paginación
- `GET /api/transportistas/:id_transportista` - Obtener detalles de transportista

## 📦 Dependencias

- [Fiber v2](https://docs.gofiber.io/) - Framework web
- [GORM](https://gorm.io/) - ORM para Go
- [PostgreSQL Driver](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL) - Driver de PostgreSQL
- [godotenv](https://github.com/joho/godotenv) - Carga de .env
- [uuid](https://github.com/google/uuid) - Generación de UUIDs

## 🔧 Configuración

### 1. Crear archivo `.env`
```env
DATABASE_URL=postgres://user:password@host:port/database
PORT=3000
```

### 2. Instalar dependencias
```bash
go mod tidy
```

### 3. Ejecutar
```bash
go run main.go
```

## 🔐 Autenticación

El proyecto valida JWT de Supabase:
- Token enviado en header: `Authorization: Bearer <token>`
- El middleware extrae el `sub` (user_id) del token
- Se asocia automáticamente a cada petición

## 📝 Modelos Principales

### User
- `id` (UUID) - PK sincronizado con auth.users
- `nombre`, `apellido`
- `rol` (cliente, transportista, admin)
- `foto_perfil`

### PerfilCliente
- `id_perfil` (UUID) - PK
- `id_usuario` (UUID) - FK a User
- `documento_identidad`, `telefono`

### Direccion
- `id_direccion` (UUID) - PK
- `id_perfil` (UUID) - FK a PerfilCliente
- `calle`, `ciudad`, `pais`
- `latitud`, `longitud`
- `es_predeterminada`

### Transportista
- `id_transportista` (UUID) - PK
- `id_usuario` (UUID) - FK a User
- `tipo_vehiculo`, `placa_vehiculo`
- `capacidad_carga`
- `estado` (verificacion_pendiente, activo, inactivo, suspendido)
- `calificacion_promedio`

## 🔗 Integración con Gateway

Este microservicio es el primer servicio en la arquitectura. Los headers JWT se extraen y validan localmente, preparados para comunicación con otros servicios vía HTTP o RPC.

## 📚 Próximos Pasos

- [ ] Agregar endpoints de pedidos
- [ ] Integrar con servicio de pagos (Rust)
- [ ] Configurar gRPC para comunicación inter-servicios
- [ ] Documentación OpenAPI/Swagger
- [ ] Tests unitarios

