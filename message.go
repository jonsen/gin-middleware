package middleware

import (
	"encoding/json"
)

type (
	//
	Msg struct {
		Code    int    // status code
		Message string // status message
		Url     string // redirect url
	}
)

func DecodeMsg(body []byte) (msg Msg, err error) {
	err = json.Unmarshal(body, &msg)
	if err != nil {
		return
	}

	return
}

func EncodeMsg(msg Msg) (body []byte, err error) {
	body, err = json.Marshal(msg)
	return
}
