package service

import (
	"context"
	api "test-zero-agency/internal/api/dto"
	"test-zero-agency/internal/entity"
	entityRepo "test-zero-agency/internal/repository/entity"
)

type PeopleService struct {
	PeopleRepo entity.IPeopleRepository
}

func NewPeopleService(PeopleRepo entity.IPeopleRepository) entity.IPeopleService {
	PeopleService := new(PeopleService)
	PeopleService.PeopleRepo = PeopleRepo
	return PeopleService
}

func (s PeopleService) ReadPeopleList(ctx context.Context, getPeopleListRequest *api.GetPeopleListRequest) ([]api.GetPeopleListResponse, error) {
	filter := entityRepo.NewFilter(getPeopleListRequest.Offset, getPeopleListRequest.Limit, getPeopleListRequest.Conditions)

	return s.PeopleRepo.SelectPeopleList(ctx, filter)
}
