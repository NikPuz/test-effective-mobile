package entity

import (
	"context"
	api "test-zero-agency/internal/api/dto"
	entityRepo "test-zero-agency/internal/repository/entity"
)

type IPeopleService interface {
	GatPeopleList(ctx context.Context, getPeopleListRequest *api.GetPeopleListRequest) ([]api.GetPeopleListResponse, error)
	DeletePeople(ctx context.Context, id int) error
	UpdatePeople(ctx context.Context, peopleRequest *api.PutPeopleRequest) (*api.PutPeopleResponse, error)
	CreatePeople(ctx context.Context, peopleRequest *api.PostPeopleRequest) (*api.PostPeopleResponse, error)
}

type IPeopleRepository interface {
	SelectPeopleList(ctx context.Context, filter *entityRepo.Filter) ([]api.GetPeopleListResponse, error)
	DeletePeople(ctx context.Context, id int) error
	UpdatePeople(ctx context.Context, peopleRequest *api.PutPeopleRequest) (*api.PutPeopleResponse, error)
	InsertPeople(ctx context.Context, peopleRequest *People) (*api.PostPeopleResponse, error)
}

type IPeopleGetaway interface {
	GetAgeByName(ctx context.Context, name string) (int, error)
	GetGenderByName(ctx context.Context, name string) (string, error)
	GetNationalityByName(ctx context.Context, name string) (string, error)
}

type People struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
