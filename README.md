# Запуск 
Нужно заполнить файл .env по примеру файла .env.example
## В консоли
```shell
go run main.go
```
## В Docker
Построение образа Docker
```shell
docker build -t meter-readings-bot .
```
Запуск контейнера с использованием файла .env
```shell
docker run -d --env-file .env --name meter-bot meter-readings-bot
```