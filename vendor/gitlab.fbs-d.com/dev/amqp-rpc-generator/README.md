# AMQP RPC Generator

Генератор сервера и клиента для работы с очередями RabbitMQ. 
Проект изначально был создан копированием функционала из библиотеки github.com/grpc-ecosystem/grpc-gateway  
Версии protoc-gen-amqprpc3 и старше используют первую версию grpc-gateway как зависимость. В данный момент grpc-gateway v1 помечена deprecated.  
В grpc-gateway v2 старые методы стали недоступны для использования извне, поэтому для версии protoc-gen-amqprpc4 были скопированы нужные файлы из grpc-gateway v2.15.2 в protoc-gen-amqprpc4/descriptor 
Эта зависимость нужна, т.к. содержит довольно сложную логику парсинга proto файлов для составления зависимостей. Без этой логики не получится добавлять все нужные импорты при генерации.  


## ConsulRoute 
Генерирует список http методов для регистрации в consul.

### Install

```
go get -u gitlab.fbs-d.com/dev/amqp-rpc-generator
cd $GOPATH/src/gitlab.fbs-d.com/dev/amqp-rpc-generator
go install gitlab.fbs-d.com/dev/amqp-rpc-generator/protoc-gen-amqprpc
go install gitlab.fbs-d.com/dev/amqp-rpc-generator/protoc-gen-amqprpc2
go install gitlab.fbs-d.com/dev/amqp-rpc-generator/protoc-gen-amqprpc3
go install gitlab.fbs-d.com/dev/amqp-rpc-generator/protoc-gen-amqprpc4
go install gitlab.fbs-d.com/dev/amqp-rpc-generator/protoc-gen-consulroute
```

### Usage

`protoc -I=. -I=$GOPATH/src -I=./vendor --amqprpc4_out=. example/example.proto`

### Debug

Для того, чтобы проверить ваш фикс или, если просто хотите посмотреть, что получится
можно использовать пакет [protoc-gen-debug](https://github.com/lyft/protoc-gen-star/tree/master/protoc-gen-debug)
Устанавливаем его себе, генерируем бин файл `code_generator_request.pb.bin` и используем его как ввод данных

```sh
protoc \
  --plugin=protoc-gen-debug=path/to/protoc-gen-debug \
  --debug_out=".:." \
  *.proto
```
