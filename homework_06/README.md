## Выполнение домашнего задания

1) Запускаем Kafka и Kafka Connect
```
docker-compose up -d
docker-compose ps -a
```

2) Проверям логи Kafka Connect
``` 
docker logs -f connect
```

3) Проверяем статус и плагины коннекторов
``` 
curl http://localhost:8083 | jq
curl http://localhost:8083/connector-plugins | jq
``` 

4) Проверяем список топиков
```
docker exec kafka1 kafka-topics --list --bootstrap-server kafka1:19092,kafka2:19093,kafka3:19094
```

5) Создаем таблицу в PostgreSQL и загружаем данные
```
docker exec -ti postgres psql -U postgres
CREATE TABLE customers (id int PRIMARY KEY, first_name text, last_name text, gender text, card_number text, bill numeric(7,2), created_date timestamp, modified_date timestamp);
COPY customers FROM '/data/Demo.csv' WITH (FORMAT csv, HEADER true);
SELECT * FROM customers LIMIT 5;
\q
```

6) Создаем коннектор customers-connector
```
curl -X POST --data-binary "@customers.json" -H "Content-Type: application/json" http://localhost:8083/connectors | jq
```

7) Проверяем коннектор customers-connector
```
curl http://localhost:8083/connectors | jq
curl http://localhost:8083/connectors/customers-connector/status | jq
```

8) Проверяем топики, читаем топик postgres.public.customers
``` 
docker exec kafka1 kafka-topics --list --bootstrap-server kafka1:19092,kafka2:19093,kafka3:19094
docker exec kafka1 kafka-console-consumer --topic postgres.public.customers --bootstrap-server kafka1:19092,kafka2:19093,kafka3:19094 --property print.offset=true --property print.key=true --from-beginning
```

9) Удаляем коннектор customers-connector
```
curl -X DELETE http://localhost:8083/connectors/customers-connector
```

10) Останавливаем Kafka и Kafka Connect
```
docker-compose stop
docker container prune -f
docker volume prune -f
```
