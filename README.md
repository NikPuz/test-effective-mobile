В сервисе дополнительно реализованы Graceful shutdown, Migrations, Middlewares, Configs, Logging

При разработке использовалось:
- Среда разработки GoLand IDEA
- Библиотеки pgx, chi, zap, viper, goose
- Система управления базами данных PostgreSQL
- Система контроля версий GitHub
- Система контейнеризации Docker

Сервис поддерживает запросы:
- GET /people
- DELETE /people/{id}
- PUT /people
- POST /people

Примеры запросов:
- POST
  {
  "name": "Misha",
  "surname": "Olegovich"
  }
- GET /people?limit=20&offset=0&age=50
- PUT
    {
    "id": 1,
    "name": "Nikita",
    "surname": "Puzirey",
    "patronymic": "Sergeevich",
    "age": 21,
    "gender": "male",
    "nationality": "RU"
    }
- DELETE /people/1



Сборка прокета:
1. ``docker compose up``
2. дождаться запуска сервиса
3. ``docker exec em_server goose -dir migrations postgres "user=root password=rpass dbname=em_db host=em_pg port=5432 sslmode=disable" up``