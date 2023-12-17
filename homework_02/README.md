## Выполнение домашнего задания

Запуск через docker-compose

- Запустить брокер с аутентификацией SASL/PLAIN
```
cd docker
make kafka-kraft-sasl-plain-up
docker ps
```

- Создать трёх пользователей с произвольными именами и настроить аутентификацию
```
Заданы в конфигурации docker-compose
```

- Создать топик
```
docker exec -ti kafka /usr/bin/kafka-topics --create --topic topic1 --bootstrap-server kafka:9092 --command-config /tmp/users/admin.properties

Created topic topic1.
```

- Первому пользователю выдать права на запись в этот топик
```
docker exec -ti kafka /usr/bin/kafka-acls --bootstrap-server kafka:9092 --add --allow-principal User:user1 --operation Write --topic topic1 --command-config /tmp/users/admin.properties

Adding ACLs for resource `ResourcePattern(resourceType=TOPIC, name=topic1, patternType=LITERAL)`: 
    (principal=User:user1, host=*, operation=WRITE, permissionType=ALLOW) 

Current ACLs for resource `ResourcePattern(resourceType=TOPIC, name=topic1, patternType=LITERAL)`: 
    (principal=User:user1, host=*, operation=WRITE, permissionType=ALLOW)
```

- Второму пользователю выдать права на чтение этого топика
```
docker exec -ti kafka /usr/bin/kafka-acls --bootstrap-server kafka:9092 --add --allow-principal User:user2 --operation Read --topic topic1 --command-config /tmp/users/admin.properties

Adding ACLs for resource `ResourcePattern(resourceType=TOPIC, name=topic1, patternType=LITERAL)`: 
    (principal=User:user2, host=*, operation=READ, permissionType=ALLOW) 

Current ACLs for resource `ResourcePattern(resourceType=TOPIC, name=topic1, patternType=LITERAL)`: 
    (principal=User:user1, host=*, operation=WRITE, permissionType=ALLOW)
    (principal=User:user2, host=*, operation=READ, permissionType=ALLOW) 
```

- Третьему пользователю не выдавать никаких прав на этот топик
```
-
```

От имени каждого пользователя выполнить команды:

Получить список топиков, Записать сообщения в топик, Прочитать сообщение из топика

#### user1:

- Получить список топиков
```
docker exec -ti kafka /usr/bin/kafka-topics --list --bootstrap-server kafka:9092 --command-config /tmp/users/user1.properties

topic1
```

- Записать сообщения в топик
```
docker exec -ti kafka /usr/bin/kafka-console-producer --topic topic1 --bootstrap-server kafka:9092 --producer.config /tmp/users/user1.properties
> message1
> message2
``` 

- Прочитать сообщения из топика
```
docker exec -ti kafka /usr/bin/kafka-console-consumer --from-beginning --topic topic1 -bootstrap-server kafka:9092 --consumer.config /tmp/users/user1.properties

WARN [Consumer clientId=console-consumer, groupId=console-consumer-83137] Not authorized to read from partition topic1-0. (org.apache.kafka.clients.consumer.internals.AbstractFetch)
ERROR Error processing message, terminating consumer process:  (kafka.tools.ConsoleConsumer$)
org.apache.kafka.common.errors.TopicAuthorizationException: Not authorized to access topics: [topic1]
Processed a total of 0 messages
```

#### user2:

- Получить список топиков
```
docker exec -ti kafka /usr/bin/kafka-topics --list --bootstrap-server kafka:9092 --command-config /tmp/users/user2.properties

topic1
```

- Записать сообщения в топик
```
docker exec -ti kafka /usr/bin/kafka-console-producer --topic topic1 --bootstrap-server kafka:9092 --producer.config /tmp/users/user2.properties

>message1
ERROR Error when sending message to topic topic1 with key: null, value: 8 bytes with error: (org.apache.kafka.clients.producer.internals.ErrorLoggingCallback)
org.apache.kafka.common.errors.TopicAuthorizationException: Not authorized to access topics: [topic1]
``` 

- Прочитать сообщения из топика
```
docker exec -ti kafka /usr/bin/kafka-console-consumer --from-beginning --topic topic1 -bootstrap-server kafka:9092 --consumer.config /tmp/users/user2.properties

message1
message2
```

#### user3:

- Получить список топиков
```
docker exec -ti kafka /usr/bin/kafka-topics --list --bootstrap-server kafka:9092 --command-config /tmp/users/user3.properties

пусто
```

- Записать сообщения в топик
```
docker exec -ti kafka /usr/bin/kafka-console-producer --topic topic1 --bootstrap-server kafka:9092 --producer.config /tmp/users/user3.properties

>message1
WARN [Producer clientId=console-producer] Error while fetching metadata with correlation id 4 : {topic1=TOPIC_AUTHORIZATION_FAILED} (org.apache.kafka.clients.NetworkClient)
ERROR [Producer clientId=console-producer] Topic authorization failed for topics [topic1] (org.apache.kafka.clients.Metadata)
ERROR Error when sending message to topic topic1 with key: null, value: 8 bytes with error: (org.apache.kafka.clients.producer.internals.ErrorLoggingCallback)
org.apache.kafka.common.errors.TopicAuthorizationException: Not authorized to access topics: [topic1]
``` 

- Прочитать сообщения из топика
```
docker exec -ti kafka /usr/bin/kafka-console-consumer --from-beginning --topic topic1 -bootstrap-server kafka:9092 --consumer.config /tmp/users/user3.properties

WARN [Consumer clientId=console-consumer, groupId=console-consumer-21483] Error while fetching metadata with correlation id 2 : {topic1=TOPIC_AUTHORIZATION_FAILED} (org.apache.kafka.clients.NetworkClient)
ERROR [Consumer clientId=console-consumer, groupId=console-consumer-21483] Topic authorization failed for topics [topic1] (org.apache.kafka.clients.Metadata)
ERROR Error processing message, terminating consumer process:  (kafka.tools.ConsoleConsumer$)
org.apache.kafka.common.errors.TopicAuthorizationException: Not authorized to access topics: [topic1]
```

- Остановить брокер
```
cd docker
make kafka-kraft-sasl-plain-down
docker ps
```

Скриншот выполнения команд:

![Alt text](./kafka-sasl.jpg?raw=true "")

Полезная информация: https://github.com/vdesabou/kafka-docker-playground




