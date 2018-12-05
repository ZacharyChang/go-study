package codec

import (
	"encoding/xml"

	"golang.org/x/net/websocket"
)

var XMLCodec = websocket.Codec{
	Marshal:   xmlMarshal,
	Unmarshal: xmlUnmarshal,
}

type Person struct {
	Name   string
	Emails []string
}

func xmlMarshal(v interface{}) (msg []byte, payloadType byte, err error) {
	msg, err = xml.Marshal(v)
	return msg, websocket.TextFrame, nil
}

func xmlUnmarshal(msg []byte, playloadType byte, v interface{}) (err error) {
	err = xml.Unmarshal(msg, v)
	return err
}
