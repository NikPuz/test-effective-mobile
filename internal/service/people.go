package service

import (
	"context"
	api "test-zero-agency/internal/api/dto"
	"test-zero-agency/internal/app/config"
	"test-zero-agency/internal/entity"
	entityRepo "test-zero-agency/internal/repository/entity"
)

type PeopleService struct {
	PeopleRepo    entity.IPeopleRepository
	PeopleGetaway entity.IPeopleGetaway
	Cfg           *config.Config
}

func NewPeopleService(cfg *config.Config, peopleRepo entity.IPeopleRepository, peopleGetaway entity.IPeopleGetaway) entity.IPeopleService {
	PeopleService := new(PeopleService)
	PeopleService.PeopleRepo = peopleRepo
	PeopleService.PeopleGetaway = peopleGetaway
	PeopleService.Cfg = cfg
	return PeopleService
}

func (s PeopleService) GatPeopleList(ctx context.Context, getPeopleListRequest *api.GetPeopleListRequest) ([]api.GetPeopleListResponse, error) {
	filter := entityRepo.NewFilter(getPeopleListRequest.Offset, getPeopleListRequest.Limit, getPeopleListRequest.Conditions)

	return s.PeopleRepo.SelectPeopleList(ctx, filter)
}

func (s PeopleService) DeletePeople(ctx context.Context, id int) error {
	return s.PeopleRepo.DeletePeople(ctx, id)
}

func (s PeopleService) UpdatePeople(ctx context.Context, peopleRequest *api.PutPeopleRequest) (*api.PutPeopleResponse, error) {
	return s.PeopleRepo.UpdatePeople(ctx, peopleRequest)
}

func (s PeopleService) CreatePeople(ctx context.Context, peopleRequest *api.PostPeopleRequest) (*api.PostPeopleResponse, error) {

	age, err := s.PeopleGetaway.GetAgeByName(ctx, peopleRequest.Name)
	if err != nil {
		return nil, err
	}

	gender, err := s.PeopleGetaway.GetGenderByName(ctx, peopleRequest.Name)
	if err != nil {
		return nil, err
	}

	nationality, err := s.PeopleGetaway.GetNationalityByName(ctx, peopleRequest.Name)
	if err != nil {
		return nil, err
	}

	people := &entity.People{
		Name:        peopleRequest.Name,
		Surname:     peopleRequest.Surname,
		Patronymic:  peopleRequest.Patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	return s.PeopleRepo.InsertPeople(ctx, people)
}
