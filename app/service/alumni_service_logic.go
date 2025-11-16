package service

import (
    "context"
    "errors"

    "github.com/noorfarihaf11/clean-arc/app/model"
)

// ================= INTERFACE REPOSITORY ==============

type IAlumniRepository interface {
    Create(ctx context.Context, a *model.Alumni) (*model.Alumni, error)
    GetByID(ctx context.Context, id string) (*model.Alumni, error)
    GetAll(ctx context.Context) ([]model.Alumni, error)
    Update(ctx context.Context, id string, data *model.Alumni) (*model.Alumni, error)
    Delete(ctx context.Context, id string) error
}

// ================= SERVICE STRUCT ====================

type AlumniService struct {
    repo IAlumniRepository
}

func NewAlumniService(repo IAlumniRepository) *AlumniService {
    return &AlumniService{repo: repo}
}

// ================= SERVICE METHODS ===================

func (s *AlumniService) Create(ctx context.Context, a *model.Alumni) (*model.Alumni, error) {
    if a.Nama == "" {
        return nil, errors.New("nama tidak boleh kosong")
    }
    return s.repo.Create(ctx, a)
}

func (s *AlumniService) GetByID(ctx context.Context, id string) (*model.Alumni, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *AlumniService) GetAll(ctx context.Context) ([]model.Alumni, error) {
    return s.repo.GetAll(ctx)
}

func (s *AlumniService) Update(ctx context.Context, id string, data *model.Alumni) (*model.Alumni, error) {
    return s.repo.Update(ctx, id, data)
}

func (s *AlumniService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}
