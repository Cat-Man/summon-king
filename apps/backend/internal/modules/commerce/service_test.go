package commerce

import (
	"context"
	"testing"
)

func TestCreatePayOrder_SameIdempotencyKeyReturnsSameOrderNo(t *testing.T) {
	svc := newTestCommerceService(t)
	o1, _ := svc.CreateOrder(context.Background(), 1001, "idem-1", 1001)
	o2, _ := svc.CreateOrder(context.Background(), 1001, "idem-1", 1001)
	if o1.OrderNo != o2.OrderNo {
		t.Fatalf("expected same order no, got %s and %s", o1.OrderNo, o2.OrderNo)
	}
}

func newTestCommerceService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryRepository())
}
