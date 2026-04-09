package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

const (
	productNameWidget   = "Widget"
	productNameOld      = "Old"
	productNameNew      = "New"
	productNameToDelete = "ToDelete"

	productPriceWidget   = 9.99
	productPriceOld      = 1.00
	productPriceNew      = 2.00
	productPriceToDelete = 5.00

	subtestZero     = "zero"
	subtestNegative = "negative"
	subtestLarge    = "large"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: productNameWidget, Price: productPriceWidget})

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != productNameWidget || got.Price != productPriceWidget {
		t.Errorf("got %+v, want Name=%s Price=%.2f", got, productNameWidget, productPriceWidget)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	cases := []struct {
		name string
		id   int
	}{
		{subtestZero, 0},
		{subtestNegative, -1},
		{subtestLarge, 99999},
	}

	s := NewMemoryStore()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.GetByID(tc.id)
			if err != ErrNotFound {
				t.Errorf("id=%d: expected ErrNotFound, got %v", tc.id, err)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: productNameOld, Price: productPriceOld})

	updated, err := s.Update(created.ID, model.Product{Name: productNameNew, Price: productPriceNew})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != productNameNew || updated.Price != productPriceNew {
		t.Errorf("got %+v, want Name=%s Price=%.2f", updated, productNameNew, productPriceNew)
	}

	got, _ := s.GetByID(created.ID)
	if got.Name != productNameNew {
		t.Errorf("store not updated: got %+v", got)
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: productNameToDelete, Price: productPriceToDelete})

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}
