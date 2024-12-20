package service

import (
	"context"

	"github.com/mhusainh/FastTix/internal/entity"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
)

type TicketService interface {
	GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error)
	GetById(ctx context.Context, id int) (entity.Product, error)
}

type ticketService struct {
	ticketRepository repository.TicketRepository
}

func NewTicketService(ticketRepository repository.TicketRepository) TicketService {
	return &ticketService{ticketRepository}
}

func (s *ticketService) GetAll(ctx context.Context, req dto.GetAllProductsRequest) ([]entity.Product, error) {
	return s.ticketRepository.GetAll(ctx, req)
}

func (s *ticketService) GetById(ctx context.Context, id int) (entity.Product, error) {
	return s.ticketRepository.GetById(ctx, id)
}
