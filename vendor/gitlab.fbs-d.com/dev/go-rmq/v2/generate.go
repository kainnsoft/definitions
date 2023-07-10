//go:generate protoc -I=. -I=$GOPATH/src --go_out=. middleware/transport/message.proto
//go:generate easyjson -all middleware/transport/message.pb.go
//go:generate easyjson middleware/client/client.go
//go:generate mockgen -source rmq.go -destination rmq_mocks.go --package rmq

package rmq
