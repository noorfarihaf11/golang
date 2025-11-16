package service

import (
	"context"
	"errors"
	"testing"

	"github.com/noorfarihaf11/clean-arc/app/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockJobRepo struct {
	data map[string]*model.PekerjaanAlumni
}

func NewMockJobRepo() *MockJobRepo {
	return &MockJobRepo{data: make(map[string]*model.PekerjaanAlumni)}
}

func (m *MockJobRepo) GetAll(ctx context.Context) ([]model.PekerjaanAlumni, error) {
	out := []model.PekerjaanAlumni{}
	for _, v := range m.data {
		out = append(out, *v)
	}
	return out, nil
}

func (m *MockJobRepo) GetByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error) {
	d, ok := m.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return d, nil
}

func (m *MockJobRepo) GetByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error) {
	out := []model.PekerjaanAlumni{}
	for _, v := range m.data {
		if v.AlumniIDStr == alumniID {
			out = append(out, *v)
		}
	}
	return out, nil
}

func (m *MockJobRepo) Create(ctx context.Context, job *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	m.data[job.ID.Hex()] = job
	return job, nil
}

func (m *MockJobRepo) Update(ctx context.Context, id string, job *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	m.data[id] = job
	return job, nil
}

func (m *MockJobRepo) SoftDelete(ctx context.Context, id string) error {
	if _, ok := m.data[id]; !ok {
		return errors.New("not found")
	}
	m.data[id].IsDeleted = true
	return nil
}

func TestCreateJob(t *testing.T) {
	mock := NewMockJobRepo()
	svc := NewJobService(mock)
	ctx := context.Background()

	job := &model.PekerjaanAlumni{
		ID:             primitive.NewObjectID(),
		NamaPerusahaan: "Google",
		AlumniIDStr:    "abc123",
	}

	res, err := svc.Create(ctx, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.NamaPerusahaan != "Google" {
		t.Errorf("expected Google, got %v", res.NamaPerusahaan)
	}
}

func TestCreateJob_EmptyCompany(t *testing.T) {
	mock := NewMockJobRepo()
	svc := NewJobService(mock)
	ctx := context.Background()

	job := &model.PekerjaanAlumni{NamaPerusahaan: ""}

	_, err := svc.Create(ctx, job)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetJobByID(t *testing.T) {
	mock := NewMockJobRepo()
	svc := NewJobService(mock)
	ctx := context.Background()

	id := primitive.NewObjectID()
	mock.data[id.Hex()] = &model.PekerjaanAlumni{
		ID:             id,
		NamaPerusahaan: "Tokopedia",
	}

	res, err := svc.GetByID(ctx, id.Hex())
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if res.NamaPerusahaan != "Tokopedia" {
		t.Errorf("expected Tokopedia, got %v", res.NamaPerusahaan)
	}
}

func TestGetJobByID_NotFound(t *testing.T) {
	mock := NewMockJobRepo()
	svc := NewJobService(mock)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, "random")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateJob(t *testing.T) {
	mock := NewMockJobRepo()
	svc := NewJobService(mock)
	ctx := context.Background()

	id := primitive.NewObjectID().Hex()

	mock.data[id] = &model.PekerjaanAlumni{
		NamaPerusahaan: "Lama",
	}

	upd := &model.PekerjaanAlumni{
		NamaPerusahaan: "Baru",
	}

	res, err := svc.Update(ctx, id, upd)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if res.NamaPerusahaan != "Baru" {
		t.Errorf("expected update, got %v", res.NamaPerusahaan)
	}
}

func TestSoftDelete(t *testing.T) {
	mock := NewMockJobRepo()
	svc := NewJobService(mock)
	ctx := context.Background()

	oid := primitive.NewObjectID()
	id := oid.Hex()

	mock.data[id] = &model.PekerjaanAlumni{
		ID:             oid, // ✅ benar
		NamaPerusahaan: "Apple",
		IsDeleted:      false,
	}

	err := svc.SoftDelete(ctx, id)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if mock.data[id].IsDeleted != true {
		t.Error("expected soft delete = true")
	}
}
