package repository

import (
	"context"
	"errors"
	"time"

	"github.com/noorfarihaf11/clean-arc/app/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAlumniRepository struct {
	data map[primitive.ObjectID]*model.Alumni
}

func NewMockAlumniRepository() *MockAlumniRepository {
	return &MockAlumniRepository{
		data: make(map[primitive.ObjectID]*model.Alumni),
	}
}

func (m *MockAlumniRepository) GetAllAlumni(ctx context.Context) ([]model.Alumni, error) {
	var result []model.Alumni
	for _, v := range m.data {
		result = append(result, *v)
	}
	return result, nil
}

func (m *MockAlumniRepository) GetAlumniByID(ctx context.Context, id primitive.ObjectID) (*model.Alumni, error) {
	if v, ok := m.data[id]; ok {
		return v, nil
	}
	return nil, errors.New("alumni not found")
}

func (m *MockAlumniRepository) CreateAlumni(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error) {
	alumni.ID = primitive.NewObjectID()
	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()
	m.data[alumni.ID] = alumni
	return alumni, nil
}

func (m *MockAlumniRepository) UpdateAlumni(ctx context.Context, id primitive.ObjectID, data *model.Alumni) (*model.Alumni, error) {
	if _, exists := m.data[id]; !exists {
		return nil, errors.New("alumni not found")
	}
	data.ID = id
	data.UpdatedAt = time.Now()
	m.data[id] = data
	return data, nil
}

func (m *MockAlumniRepository) DeleteAlumni(ctx context.Context, id primitive.ObjectID) error {
	if _, exists := m.data[id]; !exists {
		return errors.New("alumni not found")
	}
	delete(m.data, id)
	return nil
}
