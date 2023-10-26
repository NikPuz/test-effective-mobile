package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/url"
	"strconv"
)

type GetPeopleListRequest struct {
	Offset     int
	Limit      int
	Conditions map[string]string
}

type GetPeopleListResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type PutPeopleRequest struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type PutPeopleResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type PostPeopleRequest struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string `json:"patronymic"`
}

type PostPeopleResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

func NewGetPeopleListRequest(url url.Values) (*GetPeopleListRequest, error) {
	getPeopleListRequest := GetPeopleListRequest{Conditions: make(map[string]string)}

	for key, value := range url {
		switch key {
		case "offset":
			offset, err := strconv.Atoi(value[0])
			if err != nil {
				return nil, err
			}

			getPeopleListRequest.Offset = offset
		case "limit":
			limit, err := strconv.Atoi(value[0])
			if err != nil {
				return nil, err
			}

			getPeopleListRequest.Limit = limit
		case "name", "surname", "patronymic", "age", "gender", "nationality":
			getPeopleListRequest.Conditions[key] = value[0]
		default:
			return nil, errors.New(fmt.Sprintf("invalid parameter: %s", key))
		}
	}

	return &getPeopleListRequest, nil
}

func NewPostPeopleRequest(decoder *json.Decoder) (*PostPeopleRequest, error) {
	var peopleRequest PostPeopleRequest

	err := decoder.Decode(&peopleRequest)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(peopleRequest)
	if err != nil {
		return nil, err
	}

	return &peopleRequest, nil
}
