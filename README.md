<h1>Тестовое задание в Авито на Golang</h1>
Тз: https://github.com/avito-tech/backend-trainee-assignment-2024

<h2>Запуск</h2>
Запуск приложения через контейнеры

1. docker-compose build

2. docker-compose up
   
Чтобы запустить приложение вне контейнера необходимо поменять соответствующие адреса в /internal/repositories/postgres/postgres.go и  /internal/repositories/redis/redis.go на localhost

Приложение работает на localhost:8080

<h2>Тестирование</h2>

1. Для тестов поднимается контейнер с бд на порту 5433 и контейнер с redis на порту 6380. Не забудьте освободить эти порты перед запуском.

2. Все тесты написаны в ./main_test.go 
   
3. Запуск интеграционных тестов: go test .
   
4. Запуск нагрузочных тестов: ./load-testing/run

Для нагрузочных тестов используется grafana/k6, для запуска нужно скачать k6, в свой docker-compose я его не помещал.

Результат одного из прогонов тестов: /load-testing/results

<h2>Мысли при выполнении</h2>

1. По условию задачи предположил, что токены админа и юзера - константы.
Если токен пустой, то пользователь не авторизован. У админа токен всегда будет равен "admin_token".

<h2>Замечания</h2>
1. База данных при запуске заполняется некоторыми данными, их можно посмотреть в resources/database
2. При отправке запроса через SwaggerUI по пути PATCH: /banner/{id} у меня возникает ошибка 418(I'm a teapot), при отправке через Postman все работает корректно, найти причину я так и не смог.
3. 




