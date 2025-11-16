package service

import (
    "context"
    "errors"

    "github.com/noorfarihaf11/clean-arc/app/model"
)

type IJobRepository interface {
    GetAll(ctx context.Context) ([]model.PekerjaanAlumni, error)
    GetByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error)
    GetByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error)
    Create(ctx context.Context, job *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error)
    Update(ctx context.Context, id string, job *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error)
    SoftDelete(ctx context.Context, id string) error
}

type JobService struct {
    repo IJobRepository
}

func NewJobService(repo IJobRepository) *JobService {
    return &JobService{repo: repo}
}

func (s *JobService) GetAll(ctx context.Context) ([]model.PekerjaanAlumni, error) {
    return s.repo.GetAll(ctx)
}

func (s *JobService) GetByID(ctx context.Context, id string) (*model.PekerjaanAlumni, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *JobService) GetByAlumniID(ctx context.Context, alumniID string) ([]model.PekerjaanAlumni, error) {
    return s.repo.GetByAlumniID(ctx, alumniID)
}

func (s *JobService) Create(ctx context.Context, job *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
    if job.NamaPerusahaan == "" {
        return nil, errors.New("nama perusahaan wajib diisi")
    }
    return s.repo.Create(ctx, job)
}

func (s *JobService) Update(ctx context.Context, id string, job *model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
    return s.repo.Update(ctx, id, job)
}

func (s *JobService) SoftDelete(ctx context.Context, id string) error {
    return s.repo.SoftDelete(ctx, id)
}
