package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Gulner-GI/BookList/db"
	"github.com/Gulner-GI/BookList/models"
	"github.com/Gulner-GI/BookList/routes"
	"github.com/gin-gonic/gin"
)

func setupRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db.InitDB(":memory:")
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter()
}

func TestFullCRUDIntegration(t *testing.T) {
	router := setupRouter(t)
	{
		req := httptest.NewRequest(http.MethodOptions, "/books", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("OPTIONS /books returned %d, want %d", w.Code, http.StatusOK)
		}
		allow := w.Header().Get("Allow")
		if allow == "" {
			t.Error("OPTIONS: missing Allow header")
		}
	}

	{
		req := httptest.NewRequest(http.MethodHead, "/books", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("HEAD /books returned %d, want %d", w.Code, http.StatusOK)
		}
		if w.Body.Len() != 0 {
			t.Errorf("HEAD /books: expected empty body, got %q", w.Body.String())
		}
	}

	{
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("GET /books returned %d, want %d", w.Code, http.StatusOK)
		}
		var list []models.Book
		if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
			t.Fatalf("GET /books: cannot unmarshal body: %v", err)
		}
		if len(list) != 0 {
			t.Errorf("GET /books: expected empty list, got %v", list)
		}
	}

	var created models.Book
	{
		in := models.Book{
			Title:  "IntegrationTest",
			Year:   2025,
			Genre:  "Science",
			Status: false,
			Link:   nil,
		}
		body, _ := json.Marshal(in)
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("POST /books returned %d, want %d; body=%s",
				w.Code, http.StatusCreated, w.Body)
		}
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("POST /books: cannot unmarshal body: %v", err)
		}
		if created.ID == 0 {
			t.Fatalf("POST /books: expected non‑zero ID, got %d", created.ID)
		}
		if created.Title != "IntegrationTest" {
			t.Errorf("POST /books: expected Title IntegrationTest, got %q", created.Title)
		}
	}

	{
		url := "/books?id=" + strconv.Itoa(created.ID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GET %s returned %d, want %d", url, w.Code, http.StatusOK)
		}
		var fetched models.Book
		if err := json.Unmarshal(w.Body.Bytes(), &fetched); err != nil {
			t.Fatalf("GET by ID: cannot unmarshal body: %v", err)
		}
		if fetched.ID != created.ID {
			t.Errorf("GET by ID: expected ID %d, got %d", created.ID, fetched.ID)
		}
	}

	{
		patchBody := `{"genre":"Fantasy","status":true}`
		url := "/books/" + strconv.Itoa(created.ID)
		req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader([]byte(patchBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("PATCH %s returned %d, want %d; body=%s",
				url, w.Code, http.StatusOK, w.Body)
		}
		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("PATCH: cannot unmarshal body: %v", err)
		}
		if resp["genre"] != "Fantasy" || resp["status"] != true {
			t.Errorf("PATCH: expected genre=Fantasy,status=true; got %v", resp)
		}
	}

	{
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("GET /books after PATCH returned %d, want %d", w.Code, http.StatusOK)
		}
		var list []models.Book
		json.Unmarshal(w.Body.Bytes(), &list)
		if len(list) != 1 {
			t.Errorf("GET /books after PATCH: expected 1 book, got %d", len(list))
		} else if list[0].Genre != "Fantasy" {
			t.Errorf("GET /books after PATCH: expected genre Fantasy, got %q", list[0].Genre)
		}
	}

	{
		url := "/books/" + strconv.Itoa(created.ID)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Fatalf("DELETE %s returned %d, want %d", url, w.Code, http.StatusNoContent)
		}
	}

	{
		url := "/books?id=" + strconv.Itoa(created.ID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("GET after DELETE returned %d, want %d", w.Code, http.StatusNotFound)
		}
	}
}
