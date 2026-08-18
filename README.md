# Matrix Services

Servicio de descomposición QR de matrices con arquitectura de microservicios.

## Demo desplegada
Servicio	URL
Frontend	https://testservices.vercel.app
API Go	https://resourceful-upliftment-production-5bec.up.railway.app
API Node	https://romantic-enchantment-production-8874.up.railway.app

El frontend hace login automático contra credenciales mock (candidate / challenge2024) al cargar la página — no requiere que el usuario ingrese nada para probarlo.

## Arquitectura

```
┌──────────┐      ┌──────────┐      ┌──────────┐
│ Frontend │─────▶│  API Go  │─────▶│ API Node │
│  (nginx) │      │  (8080)  │      │  (3000)  │
└──────────┘      └──────────┘      └──────────┘
   :5173              │
                      │ JWT
                      ▼
               POST /api/v1/matrix/qr
               (QR + stats combinados)
```

- **Frontend**: HTML/CSS/JS vanilla servido por nginx
- **API Go**: Gateway que recibe matrices, calcula QR (Gram-Schmidt modificado), llama a Node para estadísticas
- **API Node**: Calcula estadísticas (max, min, promedio, suma, diagonal) sobre matrices Q y R

## Decisiones de Diseño

1. **QR en lugar de rotación**: El enunciado original era ambiguo entre descomposición QR y rotación. Se eligió QR porque es más general y permite factorizar cualquier matriz rectangular.

2. **Orquestación desde Go**: La API Go actúa como gateway único. El frontend solo habla con Go, que internamente llama a Node. Esto simplifica el frontend y permite degradación controlada.

3. **Degradación a `stats: null`**: Cuando Node no responde (timeout/error), Go devuelve el QR calculado con `stats: null` y un `warning` en vez de fallar con 500. Esto prioriza disponibilidad sobre completitud.

4. **Gram-Schmidt modificado**: Más estable numéricamente que el clásico, ya que normaliza cada columna antes de proyectar sobre las siguientes, reduciendo la acumulación de errores de redondeo.

## Requisitos Previos

- Docker Desktop (o Docker Engine + Docker Compose v2)
- Comando `docker compose` (sin guion)

## Inicio Rápido

```bash
cp .env.example .env
docker compose up --build
```

## URLs

| Servicio   | URL                          |
|------------|------------------------------|
| Frontend   | http://localhost:5173         |
| API Go     | http://localhost:8080         |
| API Node   | http://localhost:3000         |

## Ejemplos cURL

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"candidate","password":"challenge2024"}'
```

### QR Decomposition

```bash
# Reemplazar <TOKEN> con el token obtenido del login
curl -X POST http://localhost:8080/api/v1/matrix/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"matrix":[[1,2],[0,1],[1,0]]}'
```

### Stats (Node directamente)

```bash
curl -X POST http://localhost:3000/api/v1/matrix/stats \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"matrices":{"q":[[1,0],[0,1]],"r":[[1,2],[0,1]]}}'
```

## Tests

### Go

```bash
cd api-go
go vet ./...
go test ./... -v
```

### Node

```bash
cd api-node
npm test
```

## Limitaciones Conocidas

- No soporta matrices con más columnas que filas (m < n)
- No persiste estado entre requests
- JWT de login es un mock con credenciales fijas (`candidate`/`challenge2024`)
- El frontend almacena el JWT en memoria (se pierde al recargar)
