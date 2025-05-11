package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/todos", getTodos)
	r.POST("/todos", createTodo)
	r.PATCH("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)
	return r
}

func TestTodoendpoints(t *testing.T) {
	router := SetupRouter()
	tests := []struct {
		name     string
		method   string
		url      string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "List todos",
			method:   http.MethodGet,
			url:      "/todos",
			body:     "",
			wantCode: http.StatusOK,
			wantBody: `{"id":1, "task":"Learn Go", "done":false}`,
		},
		{
			name:     "Create todo",
			method:   http.MethodPost,
			url:      "/todos",
			body:     `{"task":"Test"}`,
			wantCode: http.StatusCreated,
			wantBody: `{"id":2,"task":"Test","done":false}`,
		},
		{
			name:     "Update todo",
			method:   http.MethodPut,
			url:      "/todos/2",
			body:     `{"task":"X","done":true}`,
			wantCode: http.StatusOK,
			wantBody: `{"id":2,"task":"X","done":true}`,
		},
		{
			name:     "Delete todo",
			method:   http.MethodDelete,
			url:      "/todos/2",
			body:     ``,
			wantCode: http.StatusNoContent,
			wantBody: ``,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			if tc.body != "" {
				req.Header.Set("Content-type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.wantCode {
				t.Errorf("%s: expected status %d, got%d", tc.name, tc.wantCode, w.Code)
			}
			got := strings.TrimSpace(w.Body.String())
			if got != tc.wantBody {
				t.Errorf("%s: expected body %q, got %q", tc.name, tc.wantBody, got)
			}
		})
	}
}
