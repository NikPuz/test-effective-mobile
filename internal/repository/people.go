package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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
		return nil, err
	}
	defer rows.Close()

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
			return nil, err
		}

		peopleList = append(peopleList, people)
	}

	return peopleList, nil
}

func (r PeopleRepository) DeletePeople(ctx context.Context, id int) error {

	result, err := r.db.Exec(ctx, `DELETE FROM people WHERE id=$1`, id)
	if result.RowsAffected() == 0 {
		return entity.NotFoundErrorResponse
	} else if err != nil {
		return err
	}

	return nil
}

func (r PeopleRepository) UpdatePeople(ctx context.Context, peopleRequest *api.PutPeopleRequest) (*api.PutPeopleResponse, error) {
	var people api.PutPeopleResponse

	err := r.db.QueryRow(ctx,
		`UPDATE people SET name = $2, surname = $3, patronymic = $4, age = $5, gender = $6, nationality = $7 WHERE id = $1 
RETURNING id, name, surname, patronymic, age, gender, nationality`,
		peopleRequest.Id, peopleRequest.Name, peopleRequest.Surname, peopleRequest.Patronymic, peopleRequest.Age, peopleRequest.Gender, peopleRequest.Nationality).Scan(
		&people.Id,
		&people.Name,
		&people.Surname,
		&people.Patronymic,
		&people.Age,
		&people.Gender,
		&people.Nationality)

	if errors.As(err, &pgx.ErrNoRows) {
		return nil, entity.NotFoundErrorResponse
	} else if err != nil {
		return nil, err
	}

	return &people, nil
}

func (r PeopleRepository) InsertPeople(ctx context.Context, peopleRequest *entity.People) (*api.PostPeopleResponse, error) {
	var people api.PostPeopleResponse

	err := r.db.QueryRow(ctx,
		`INSERT INTO people(name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, name, surname, patronymic, age, gender, nationality`,
		peopleRequest.Name, peopleRequest.Surname, peopleRequest.Patronymic, peopleRequest.Age, peopleRequest.Gender, peopleRequest.Nationality).Scan(
		&people.Id,
		&people.Name,
		&people.Surname,
		&people.Patronymic,
		&people.Age,
		&people.Gender,
		&people.Nationality)

	if err != nil {
		return nil, err
	}

	return &people, nil
}
