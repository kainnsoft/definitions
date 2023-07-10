//go:generate protoc -I=. --go_out=paths=source_relative:. queue.proto
package amqp_rpc_generator
