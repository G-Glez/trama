# Testing Alternatives

## 1. Go Unit Tests with `testing` package + `httptest`

Escribe tests en Go usando el paquete estándar `testing` y `net/http/httptest` para probar handlers de Gin sin levantar el servidor.

**Ejemplo** (`main_test.go`):

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHolaEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/hola", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hola mundo"})
	})

	req, _ := http.NewRequest("GET", "/hola", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Hola mundo")
}
```

Ejecutar: `go test ./... -v`

## 2. HTTP Client tests (integration)

Pruebas que levantan el servidor real y hacen peticiones HTTP.

**Ejemplo** (`integration_test.go`):

```go
package main

import (
	"net/http"
	"testing"
	"time"
)

func TestServerLive(t *testing.T) {
	go main()
	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8080/hola")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
```

## 3. Test containers (dockertest)

Usa [dockertest](https://github.com/ory/dockertest) para levantar contenedores Docker reales (por ejemplo, la propia API o una DB de testing) y testear contra ellos.

## 4. OpenAPI / Swagger

Puedes generar documentación OpenAPI desde Go usando:

- **[swaggo/swag](https://github.com/swaggo/swag)** – anotaciones en los handlers que generan `swagger.json`, luego puedes servir la UI con `gin-swagger`.
- **[openapi-tools](https://github.com/chanced/openapi)** – librerías para validar requests/responses contra un schema OpenAPI.

**Ejemplo con swaggo**:

```go
// @Summary Say hello
// @Description Returns a greeting
// @Produce json
// @Success 200 {object} map[string]string
// @Router /hola [get]
func holaHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hola mundo"})
}
```

Luego generas la spec con `swag init` y sirves Swagger UI con `gin-swagger`.

Una vez generada la spec OpenAPI, puedes:

- Validar requests/responses automáticamente con middlewares que comparen contra el schema.
- Compartir la spec con frontends u otros servicios.
- Usar herramientas como [Postman](https://postman.com), [Insomnia](https://insomnia.rest), o [Hoppscotch](https://hoppscotch.io) importando el JSON para testing manual/automatizado.

**¿Se puede ejecutar OpenAPI desde Go en función del entorno?**

Sí. Puedes condicionar la carga de la ruta `/swagger/*` según `GIN_MODE`:

```go
if cfg.GinMode == "debug" {
	// registrar swagger solo en desarrollo
	r.GET("/swagger/*any", swaggerHandler)
}
```

## 5. Herramientas externas

- **[curl](https://curl.se)**: `curl http://localhost:8080/hola`
- **[httpie](https://httpie.io)**: `http :8080/hola`
- **[Postman](https://postman.com)** / **[Insomnia](https://insomnia.rest)**: GUI para colecciones de requests
- **[Hoppscotch](https://hoppscotch.io)**: alternativa online/open-source a Postman
- **[Newman](https://github.com/postmanlabs/newman)**: ejecuta colecciones de Postman desde CLI (ideal para CI/CD)
- **[ginkgo](https://github.com/onsi/ginkgo) + [gomega](https://github.com/onsi/gomega)**: BDD-style testing para Go

## 6. CI/CD

Integra tests en GitHub Actions, GitLab CI, etc:

```yaml
# .github/workflows/test.yml
name: Test
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: go test ./... -v -race -count=1
```

## Resumen

| Enfoque | Tipo | Velocidad | Dependencias |
|---|---|---|---|
| `httptest` | Unitario | Rápido | Solo Go stdlib |
| HTTP client | Integración | Medio | Ninguna |
| dockertest | E2E | Lento | Docker |
| OpenAPI/Swagger | Documentación + validación | Medio | swaggo |
| curl/httpie/Postman | Manual | Instantáneo | Herramienta externa |
