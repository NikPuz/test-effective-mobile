-- +goose Up

--
-- таблица `people`
--

CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE people (
                           id SERIAL PRIMARY KEY,
                           name VARCHAR(64) NOT NULL,
                           surname VARCHAR(64) NOT NULL,
                           patronymic VARCHAR(64) NOT NULL,
                           age INT NOT NULL,
                           gender gender NOT NULL,
                           nationality VARCHAR(64) NOT NULL
                       );