# Тестовое задание

## Проект
### Сборка проекта 
``` 
make build
```
### Запуск тестов
```
make test
```

## Docker
### Создание образа
```
make build-docker name=<your-image-name>
```
### Запуск контейнера
```
docker -rm -v /absolute/local/path/file.txt:/app/input.txt <your-image-name> input.txt
```
