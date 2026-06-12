package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"trama/internal/api/database"
	"trama/internal/api/handlers"
	"trama/internal/core"
	coregen "trama/internal/gen/core"
)

func TestHolaEndpoint(t *testing.T) {
	tmp := t.TempDir() + "/test.db"

	db, err := database.Open(tmp)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}

	q := coregen.New(db)
	h := handlers.New(db,
		core.NewGameSystemRepository(q),
		core.NewEditionRepository(q),
		core.NewFactionRepository(q),
	)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/hola", h.Hola)

	req, _ := http.NewRequest("GET", "/hola", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	t.Logf("Response: %s", w.Body.String())
}
