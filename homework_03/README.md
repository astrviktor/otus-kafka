## Выполнение домашнего задания

Реализовано на golang

Запуск через docker-compose

- Сборка приложений
```
make docker-build
```

- Запуск
```
make compose-up
```

- Проверка
```
docker logs producer
docker logs consumer
```

- Остановка
```
make compose-down
```

Скриншот выполнения:

![Alt text](./kafka-transactions.jpg?raw=true "")






