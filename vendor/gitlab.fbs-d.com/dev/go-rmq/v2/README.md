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
    defaultQos: -1
    qos:
        queue1: 1
        queue2: 2
    defaultConsumersCount: -1
    consumersCount:
        queue1: 1
        queue2: 1
    countChannelsPublisher: 10
    countChannelsConsumeRPC: 10
    defaultTimeout: 5 # seconds
    timeouts:
        queue1: 10
        queue2: 20
```
