package entity

import (
	"context"
	api "test-zero-agency/internal/api/dto"
	entityRepo "test-zero-agency/internal/repository/entity"
)

type IPeopleService interface {
	ReadPeopleList(ctx context.Context, getPeopleListRequest *api.GetPeopleListRequest) ([]api.GetPeopleListResponse, error)
}

type IPeopleRepository interface {
	SelectPeopleList(ctx context.Context, filter *entityRepo.Filter) ([]api.GetPeopleListResponse, error)
}

type People struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Patronymic    string `json:"patronymic"`
	Age           int    `json:"age"`
	GenderId      int    `json:"genderId"`
	NationalityId int    `json:"nationalityId"`
}
