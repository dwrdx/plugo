package plugo

const (
	MSGPACK_RPC_TYPE_REQUEST  = 0
	MSGPACK_RPC_TYPE_RESPONSE = 1
	MSGPACK_RPC_TYPE_NOTIFY   = 2
)

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
