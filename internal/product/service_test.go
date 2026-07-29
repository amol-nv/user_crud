package product

import (
	"testing"

	"amol-nv/user_crud/internal/repository"
)

type testRepo struct {
	*repository.InMemoryProductRepository
}

func TestService_Create_Get_Update_Delete(t *testing.T) {
	repo := repository.NewInMemoryProductRepository()
	svc := NewService(repo)

	id := "p1"
	createResp, err := svc.Create(id, CreateProductRequest{Name: "Phone", Price: 10.5, Description: "Nice"})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	if createResp.ID != id {
		t.Fatalf("expected id %s got %s", id, createResp.ID)
	}

	got, err := svc.GetByID(id)
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	if got.Name != "Phone" {
		t.Fatalf("expected name Phone got %s", got.Name)
	}

	_, err = svc.GetByID("missing")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound got %v", err)
	}

	upd, err := svc.Update(id, UpdateProductRequest{Name: "Phone2", Price: 11.0, Description: "Better"})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if upd.Name != "Phone2" {
		t.Fatalf("expected updated name Phone2 got %s", upd.Name)
	}

	_, err = svc.Update("missing", UpdateProductRequest{Name: "X", Price: 1, Description: ""})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound got %v", err)
	}

	if err := svc.Delete(id); err != nil {
		t.Fatalf("delete err: %v", err)
	}
	if err := svc.Delete(id); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound got %v", err)
	}
}
