package definitions

//go:generate mockgen -source client.go -destination ./mocks/client_mocks.go -package mocks
//go:generate protoc -I=. -I=./vendor --go_opt=paths=source_relative --go_out=. --amqprpc3_out=. proto/common_partner-transfers.rmq.client.proto
//go:generate gofmt -w ./proto/common_partner-transfers.rmq.client.pb.go ./proto/common_partner-transfers.rmq.client_amqprpc.gen.go
//go:generate find . -type f -name "*.pb.go" -exec sed -i.bak s/,omitempty// {} ;

import (
	_ "gitlab.fbs-d.com/definitions/client"
	_ "gitlab.fbs-d.com/dev/amqp-rpc-generator"
)
