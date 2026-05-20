package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Widget", Price: 9.99})
	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Widget" || got.Price != 9.99 {
		t.Errorf("unexpected product: %+v", got)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestGetAllReturnsAll(t *testing.T) {
	s := NewMemoryStore()
	s.Create(model.Product{Name: "A", Price: 1})
	s.Create(model.Product{Name: "B", Price: 2})
	s.Create(model.Product{Name: "C", Price: 3})

	products := s.GetAll()
	if len(products) != 3 {
		t.Errorf("expected 3 products, got %d", len(products))
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateExisting(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Old", Price: 1})

	updated, err := s.Update(created.ID, model.Product{Name: "New", Price: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ID != created.ID {
		t.Errorf("expected ID %d preserved, got %d", created.ID, updated.ID)
	}
	if updated.Name != "New" || updated.Price != 2 {
		t.Errorf("unexpected updated product: %+v", updated)
	}

	got, _ := s.GetByID(created.ID)
	if got.Name != "New" {
		t.Errorf("update did not persist: %+v", got)
	}
}

func TestUpdateNonExistent(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Update(999, model.Product{Name: "X", Price: 1})
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteExisting(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "ToDelete", Price: 1})

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.GetByID(created.ID); err != ErrNotFound {
		t.Errorf("expected product gone, got err=%v", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}
