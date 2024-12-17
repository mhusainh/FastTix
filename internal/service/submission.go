package service

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/repository"
)

type SubmissionService interface {
	GetAll(ctx context.Context) ([]entity.Product, error)
	GetById(ctx context.Context, id int64) (*entity.Product, error)
	// Approve(ctx context.Context, submission *entity.Product) error
	// Reject(ctx context.Context, submission *entity.Product) error
}

type submissionService struct {
	submissionRepository repository.SubmissionRepository
}

func NewSubmissionService(submissionRepository repository.SubmissionRepository) SubmissionService {
	return &submissionService{submissionRepository}
}

func (s submissionService) GetAll(ctx context.Context) ([]entity.Product, error) {
	return s.submissionRepository.GetAll(ctx)
}

func (s submissionService) GetById(ctx context.Context, id int64) (*entity.Product, error) {
	return s.submissionRepository.GetById(ctx, id)
}
