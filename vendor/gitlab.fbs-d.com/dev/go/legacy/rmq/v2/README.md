RabbitMQ
========

Пакет для работы с RabbitMQ.

Зависимости
===========

github.com/streadway/amqp
github.com/golang/mock/gomock

Конфигурация
============

```yaml
rabbit:
    host: 127.0.0.1
    port: 5672
    username: guest
    password: guest
    virtualHost: vhost
    defaultQos: 1
    qos:
        queue1: 1
        queue2: 2
    defaultConsumersCount: -1
    consumersCount:
        queue1: 1
        queue2: 1
    defaultTimeout: 5 # seconds
    timeouts:
        queue1: 10
        queue2: 20
```

Описание параметров:

* host - адрес сервера RabbitMQ
* port - порт сервера
* username - логин пользователя
* password - пароль пользователя
* virtualHost - вирутальный хост сервера RabbitMQ
* connectionName - имя соединения с RabbitMQ, отображается в UI раббита, позволяет идентицифировать коннект
* defaultQos - количество сообщений, забираемых из очереди одним консьюмером (если в очереди 5 сообщений, а qos стоит 2, то первый консьюмер заберет сразу 2 сообщения, второй 2 и третий одно)
* qos - переопределение defaultQos для конкретных очередей
* defaultConsumersCount - количество консьюмеров на одну очередь (-1 - количество будет равно количеству CPU, которое отдает система, 0 - не создавать консьюмеров)
* consumersCount - переопределение количества консьюмеров для конкретных очередей
* defaultTimeout - таймаут rpc-запросов
* timeouts - переопределение таймаутов для конкретных очередей
