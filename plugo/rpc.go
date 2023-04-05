package plugo

import "reflect"

const (
	MSGPACK_RPC_TYPE_REQUEST  = 0
	MSGPACK_RPC_TYPE_RESPONSE = 1
	MSGPACK_RPC_TYPE_NOTIFY   = 2
)

type MessageHeader struct {
	Type int
	Id   int
}

// RequestMessage
type RequestMessage struct {
	Type   int
	Id     int
	Method string
	Params map[string]interface{}
}

// ResponseMessage
type ResponseMessage struct {
	Type   int
	Id     int
	Error  string
	Result map[string]interface{}
}

// NotifyMessage
type NotifyMessage struct {
	Type   int
	Method string
	Params map[string]interface{}
}

// CallMethodOfStruct calls a method of a struct by method name and params.
func CallMethodOfStruct(s interface{}, name string, params map[string]interface{}) error {
	f := reflect.ValueOf(s).MethodByName(name)
	f.Call([]reflect.Value{reflect.ValueOf(params)})
	return nil
}
