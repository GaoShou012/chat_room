package im

import "encoding/json"

func Ack(message *Message, body interface{}) (j []byte, err error) {
	message.Head.Path = message.Head.Path + "/ack"
	message.Body = body
	return json.Marshal(message)
}
