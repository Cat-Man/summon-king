package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTraceAndAuthAndIdempotencyContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(InjectTraceID())
	r.Use(InjectAuthToken())
	r.Use(InjectIdempotencyKey())
	r.GET("/context", func(c *gin.Context) {
		traceID := GetTraceID(c)
		if traceID == "" {
			t.Fatalf("trace ID missing")
		}
		if auth := GetAuthToken(c); auth != "token123" {
			t.Fatalf("unexpected auth token %s", auth)
		}
		if idemp := GetIdempotencyKey(c); idemp != "idem-1" {
			t.Fatalf("unexpected idempotency key %s", idemp)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/context", nil)
	req.Header.Set("Authorization", " Bearer token123 ")
	req.Header.Set("X-Idempotency-Key", "idem-1")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.Code)
	}
}
