package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func TestHolaEndpoint(t *testing.T) {
	tmp := t.TempDir() + "/test.db"

	var err error
	db, err = sql.Open("sqlite", tmp)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS greetings (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	)`)
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/hola", holaHandler)

	req, _ := http.NewRequest("GET", "/hola", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	db.Close()

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	t.Logf("Response: %s", w.Body.String())
}
