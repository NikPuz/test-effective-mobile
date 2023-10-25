package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	api "test-zero-agency/internal/api/dto"
	"test-zero-agency/internal/entity"
	entityRepo "test-zero-agency/internal/repository/entity"
)

type PeopleRepository struct {
	db *pgxpool.Pool
}

func NewPeopleRepository(db *pgxpool.Pool) entity.IPeopleRepository {
	PeopleRepository := new(PeopleRepository)
	PeopleRepository.db = db
	return PeopleRepository
}

func (r PeopleRepository) SelectPeopleList(ctx context.Context, filter *entityRepo.Filter) ([]api.GetPeopleListResponse, error) {
	peopleList := make([]api.GetPeopleListResponse, 0)

	query := fmt.Sprintf("SELECT id, name, surname, patronymic, age, gender, nationality FROM people %s", filter.Build())

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, entity.NewError(err, 500)
	}

	for rows.Next() {
		var people api.GetPeopleListResponse
		err = rows.Scan(
			&people.Id,
			&people.Name,
			&people.Surname,
			&people.Patronymic,
			&people.Age,
			&people.Gender,
			&people.Nationality,
		)
		if err != nil {
			return nil, entity.NewError(err, 500)
		}

		peopleList = append(peopleList, people)
	}

	return peopleList, nil
}

//func (r PeopleRepository) buildLikeQueryPeople(people *api.GetPeopleListRequest) string {
//	queries := make([]string, 0, 6)
//
//	if len(people.Name) != 0 {
//		queries = append(queries, "p.name LIKE '"+people.Name+"'")
//	}
//	if len(people.Surname) != 0 {
//		queries = append(queries, "p.surname LIKE '"+people.Surname+"'")
//	}
//	if len(people.Patronymic) != 0 {
//		queries = append(queries, "p.patronymic LIKE '"+people.Patronymic+"'")
//	}
//	if people.Age != nil {
//		queries = append(queries, "p.age LIKE '"+strconv.Itoa(*people.Age)+"'")
//	}
//	if len(people.Gender) != 0 {
//		queries = append(queries, "g.name LIKE '"+people.Gender+"'")
//	}
//	if len(people.Nationality) != 0 {
//		queries = append(queries, "n.name LIKE '"+people.Nationality+"'")
//	}
//
//	return strings.Join(queries, " AND ")
//}
