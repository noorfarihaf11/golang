package service

import (
    "context"
    "errors"
    "testing"

    "github.com/noorfarihaf11/clean-arc/app/model"
)

// ================= MOCK REPO =================

type MockAlumniRepo struct {
    data map[string]*model.Alumni
}

func NewMockAlumniRepo() *MockAlumniRepo {
    return &MockAlumniRepo{
        data: make(map[string]*model.Alumni),
    }
}

func (m *MockAlumniRepo) Create(ctx context.Context, a *model.Alumni) (*model.Alumni, error) {
    if a.Nama == "" {
        return nil, errors.New("nama tidak boleh kosong")
    }
    m.data[a.NIM] = a
    return a, nil
}

func (m *MockAlumniRepo) GetByID(ctx context.Context, id string) (*model.Alumni, error) {
    if val, ok := m.data[id]; ok {
        return val, nil
    }
    return nil, errors.New("not found")
}

func (m *MockAlumniRepo) GetAll(ctx context.Context) ([]model.Alumni, error) {
    arr := []model.Alumni{}
    for _, v := range m.data {
        arr = append(arr, *v)
    }
    return arr, nil
}

func (m *MockAlumniRepo) Update(ctx context.Context, id string, data *model.Alumni) (*model.Alumni, error) {
    if _, ok := m.data[id]; !ok {
        return nil, errors.New("not found")
    }
    m.data[id] = data
    return data, nil
}

func (m *MockAlumniRepo) Delete(ctx context.Context, id string) error {
    if _, ok := m.data[id]; !ok {
        return errors.New("not found")
    }
    delete(m.data, id)
    return nil
}

// ================= TEST CREATE =================

func TestCreateAlumni(t *testing.T) {
    mock := NewMockAlumniRepo()
    svc := NewAlumniService(mock) // ← buat instance service pakai mock repo
    ctx := context.Background()

    alumni := &model.Alumni{
        NIM:  "123",
        Nama: "Budi",
    }

    result, err := svc.Create(ctx, alumni) // ← panggil method pada instance
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if result.Nama != "Budi" {
        t.Errorf("expected Budi, got %v", result.Nama)
    }
}


func TestCreateAlumni_EmptyName(t *testing.T) {
    mock := NewMockAlumniRepo()
    service := NewAlumniService(mock)
    ctx := context.Background()

    alumni := &model.Alumni{
        NIM:  "123",
        Nama: "",
    }

    _, err := service.Create(ctx, alumni)
    if err == nil {
        t.Errorf("expected error, got nil")
    }
}

// ================= TEST GET BY ID =================

func TestGetAlumniByID(t *testing.T) {
    mock := NewMockAlumniRepo()
    service := NewAlumniService(mock)
    ctx := context.Background()

    mock.data["123"] = &model.Alumni{NIM: "123", Nama: "Budi"}

    res, err := service.GetByID(ctx, "123")
    if err != nil {
        t.Fatalf("unexpected: %v", err)
    }

    if res.Nama != "Budi" {
        t.Errorf("expected Budi, got %v", res.Nama)
    }
}

func TestGetAlumniByID_NotFound(t *testing.T) {
    mock := NewMockAlumniRepo()
    service := NewAlumniService(mock)
    ctx := context.Background()

    _, err := service.GetByID(ctx, "999")
    if err == nil {
        t.Error("expected error, got nil")
    }
}

// ================= TEST GET ALL =================

func TestGetAllAlumni(t *testing.T) {
    mock := NewMockAlumniRepo()
    service := NewAlumniService(mock)
    ctx := context.Background()

    mock.data["1"] = &model.Alumni{NIM: "1", Nama: "A"}
    mock.data["2"] = &model.Alumni{NIM: "2", Nama: "B"}

    res, err := service.GetAll(ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(res) != 2 {
        t.Errorf("expected 2 data, got %v", len(res))
    }
}

// ================= TEST UPDATE =================

func TestUpdateAlumni(t *testing.T) {
    mock := NewMockAlumniRepo()
    service := NewAlumniService(mock)
    ctx := context.Background()

    mock.data["123"] = &model.Alumni{NIM: "123", Nama: "Budi"}

    update := &model.Alumni{NIM: "123", Nama: "Budi Baru"}

    res, err := service.Update(ctx, "123", update)
    if err != nil {
        t.Fatalf("unexpected: %v", err)
    }

    if res.Nama != "Budi Baru" {
        t.Errorf("expected updated name, got %v", res.Nama)
    }
}

func TestUpdateAlumni_NotFound(t *testing.T) {
    mock := NewMockAlumniRepo()
    service := NewAlumniService(mock)
    ctx := context.Background()

    update := &model.Alumni{NIM: "999", Nama: "X"}

    _, err := service.Update(ctx, "999", update)
    if err == nil {
        t.Error("expected error, got nil")
    }
}

// ================= TEST DELETE =================

func TestDeleteAlumni(t *testing.T) {
    mock := NewMockAlumniRepo()
    service := NewAlumniService(mock)
    ctx := context.Background()

    mock.data["123"] = &model.Alumni{NIM: "123", Nama: "Budi"}

    err := service.Delete(ctx, "123")
    if err != nil {
        t.Fatalf("unexpected: %v", err)
    }

    if _, ok := mock.data["123"]; ok {
        t.Error("data still exists after delete")
    }
}
