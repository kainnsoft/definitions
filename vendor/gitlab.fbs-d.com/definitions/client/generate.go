//go:generate protoc --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative client.proto
//go:generate easyjson -all client.pb.go
package client
