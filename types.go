package middleware

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

// ReqHeader request header
type ReqHeader struct {
	Action string
	Time   string
	// other
}

// Request -
type Request struct {
	Header ReqHeader   // request header
	Body   interface{} // body
}

// Response -
type Response struct {
	Code    int         // custom status code
	Message string      // message
	Body    interface{} // body
}

// StructToString convert struct to string
func StructToString(data interface{}) string {
	byt, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(byt)
}

func (res *Response) String() string {
	str, err := json.Marshal(res)
	if err != nil {
		return ""
	}

	return string(str)
}

func (req *Request) String() string {
	str, err := json.Marshal(req)
	if err != nil {
		return ""
	}

	return string(str)
}

// RequestReader read request data
func RequestReader(input io.ReadCloser, body interface{}) (req *Request, err error) {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("body is null")
	}

	req = new(Request)

	req.Body = body

	err = json.Unmarshal(data, req)

	return
}

// ResponseWriter initial response data
// code: status code
// message: response messages
// body: response body
func ResponseWriter(code int, message string, body interface{}) (res *Response) {
	res = &Response{
		Code:    code,
		Message: message,
		Body:    body,
	}
	return
}
