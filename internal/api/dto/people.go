package api

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type GetPeopleListRequest struct {
	Offset     int
	Limit      int
	Conditions map[string]string
}

type GetPeopleListResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         int     `json:"age"`
	Gender      string  `json:"gender"`
	Nationality string  `json:"nationality"`
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
