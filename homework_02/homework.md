## Домашнее задание

### Kraft и Kafka. Настройка безопасности.

Цель:
Научиться разворачивать kafka с помощью kraft и самостоятельно настраивать безопасность

Описание/Пошаговая инструкция выполнения домашнего задания:

Развернуть Kafka с KRaft и настроить безопасность:

1. Запустить Kafka с Kraft:
   - Сгенерировать UUID кластера
   - Отформатировать папки для журналов
   - Запустить брокер
2. Настроить аутентификацию SASL/PLAIN. Создать трёх пользователей с произвольными именами.
3. Настроить авторизацию. Создать топик. Первому пользователю выдать права на запись в этот топик. Второму пользователю выдать права на чтение этого топика. Третьему пользователю не выдавать никаких прав на этот топик.
4. От имени каждого пользователя выполнить команды:
   - Получить список топиков
   - Записать сообщения в топик
   - Прочитать сообщения из топика
   
В качестве результата ДЗ прислать документ с описанием настроек кластера, выполненных команд и снимков экрана с результатами команд.

Дополнительное задание: настроить SSL

Рекомендуем сдать до: 10.12.2023