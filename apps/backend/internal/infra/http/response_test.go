package httpx

import "testing"

func TestSuccessResponse_CodeShouldBeZero(t *testing.T) {
	resp := Success(map[string]string{"ok": "1"}, "trace-1")
	if resp.Code != 0 {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}
}
