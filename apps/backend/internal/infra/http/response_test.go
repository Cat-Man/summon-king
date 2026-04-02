package http

import "testing"

func TestSuccessResponse_CodeShouldBeZero(t *testing.T) {
    resp := Success(map[string]string{"ok": "1"}, "trace-1")
    if resp.Code != 0 {
        t.Fatalf("expected code 0, got %d", resp.Code)
    }
    if resp.TraceID != "trace-1" {
        t.Fatalf("expected trace id trace-1, got %s", resp.TraceID)
    }
}
