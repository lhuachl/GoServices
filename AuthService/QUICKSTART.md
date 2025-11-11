# Microservicio de Transporte - Guía Rápida

## 🚀 Qué Se Creó

Un **microservicio profesional en Go** con:

- **Framework**: Fiber (rápido, ligero, similar a EXPRESS)
- **ORM**: GORM con PostgreSQL (Supabase)
- **Autenticación**: JWT de Supabase
- **Estructura escalable**: Handlers, Models, Middleware

## 📁 Estructura del Proyecto

```
goServices/
├── main.go                      # Punto de entrada, rutas
├── .env                         # Variables de entorno
├── go.mod / go.sum             # Dependencias
├── .vscode/
│   └── tasks.json              # Tareas de VS Code
├── .github/
│   └── copilot-instructions.md # Instrucciones para Copilot
└── pkg/
    ├── models/
    │   ├── user.go            # User, PerfilCliente, Direccion, DTOs
    │   └── transportista.go    # Transportista, DTOs
    ├── handlers/
    │   ├── users.go           # GetMe, UpdateMe, GetUser
    │   ├── addresses.go        # GetMyAddresses, CreateAddress, UpdateAddress, DeleteAddress
    │   └── transportistas.go   # GetTransportistas (con paginación), GetTransportista
    └── middleware/
        └── auth.go            # Middleware JWT de Supabase
```

## 🔌 Endpoints Implementados

### Públicos
- `GET /health` - Health check

### Autenticados (requieren JWT en header)

**Usuarios**:
- `GET /api/users/me` - Perfil del usuario autenticado
- `PUT /api/users/me` - Actualizar perfil
- `GET /api/users/:id_usuario` - Ver otro usuario (solo él mismo o admin)

**Direcciones**:
- `GET /api/users/me/addresses` - Listar direcciones
- `POST /api/users/me/addresses` - Crear dirección
- `PUT /api/users/me/addresses/:id_direccion` - Actualizar
- `DELETE /api/users/me/addresses/:id_direccion` - Eliminar

**Transportistas**:
- `GET /api/transportistas?page=1&page_size=10&estado=activo&ciudad=Quito&calificacion_min=3.5`
- `GET /api/transportistas/:id_transportista`

## 🔐 Autenticación

El middleware `AuthMiddleware` en `pkg/middleware/auth.go`:
- Extrae el JWT del header `Authorization: Bearer <token>`
- Decodifica el payload sin validar firma (en dev)
- Extrae el `sub` (user_id) y lo pasa al contexto
- En producción: validar con clave pública de Supabase

## 📝 Cómo Ejecutar

1. **Actualizar `.env`** con tu cadena de conexión Supabase:
   ```env
   DATABASE_URL=postgres://user:password@host:5432/database
   PORT=3000
   ```

2. **Instalar dependencias**:
   ```bash
   go mod tidy
   ```

3. **Ejecutar**:
   ```bash
   go run main.go
   ```

El servidor estará en `http://localhost:3000`

## 🎯 Características Principales

✅ **Modelos GORM** - Sincronizados con Supabase  
✅ **Autenticación JWT** - Supabase Auth integrado  
✅ **Paginación** - En endpoint de transportistas  
✅ **Validaciones** - Permisos (solo propio usuario o admin)  
✅ **DTOs** - Request/Response estructurados  
✅ **Manejo de errores** - Status HTTP apropiados  
✅ **CORS** - Habilitado para desarrollo  

## 🔗 Integración con Otros Servicios

Este microservicio está listo para:
- Comunicarse con **API Gateway** (envía JWT en headers)
- Conectarse con servicio de pagos/Rust
- Escalar con más endpoints

El JWT se pasa intacto en headers para que otros servicios lo validen localmente.

## ⚙️ Próximas Mejoras

- [ ] Endpoints de pedidos
- [ ] Validación JWT contra clave pública de Supabase (producción)
- [ ] gRPC para comunicación inter-servicios
- [ ] Documentación OpenAPI/Swagger
- [ ] Tests unitarios
- [ ] Rate limiting
- [ ] Logging estructurado
