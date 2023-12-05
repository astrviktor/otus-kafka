## Выполнение домашнего задания

- Установить Java JDK
```
java -version
sudo apt install default-jre
java -version
```

- Скачать Kafka с сайта kafka.apache.org и развернуть на локальном диске
```
cd /tmp
wget https://downloads.apache.org/kafka/3.6.0/kafka_2.13-3.6.0.tgz

cd /opt
sudo tar xzvf /tmp/kafka_2.13-3.6.0.tgz

sudo ln -s /opt/kafka_2.13-3.6.0 kafka

cd /opt/kafka
```

- Запустить Zookeeper
```
sudo bin/zookeeper-server-start.sh -daemon config/zookeeper.properties
```

- Запустить Kafka Broker
```
sudo bin/kafka-server-start.sh -daemon config/server.properties
```

- Создать топик test
```
sudo bin/kafka-topics.sh --list --bootstrap-server localhost:9092
sudo bin/kafka-topics.sh --create --topic test --bootstrap-server localhost:9092
sudo bin/kafka-topics.sh --describe --topic test --bootstrap-server localhost:9092

```

- Записать несколько сообщений в топик
```
sudo bin/kafka-console-producer.sh --topic test --bootstrap-server localhost:9092
message 1
message 2
message 3
``` 

- Прочитать сообщения из топика
```
sudo bin/kafka-console-consumer.sh --topic test --from-beginning ---bootstrap-server localhost:9092
message 1
message 2
message 3
```

Скриншот выполнения команд:

![Alt text](./kafka-cmd.jpg?raw=true "")