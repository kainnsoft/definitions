package rmq

import "encoding/json"

type IMessage interface {
	GetJson() []byte
}

type TokenizedMessage struct {
	Token string
}

type EmptyResponse struct {
}

func (res *EmptyResponse) GetJson() (data []byte) {
	data, _ = json.Marshal(res)
	return
}
