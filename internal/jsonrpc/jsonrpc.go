package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Call struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	Id      interface{} `json:"id,omitempty"`
}

type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Id      interface{} `json:"id"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	ErrParseError           = -32700
	ErrInvalidRequest       = -32600
	ErrMethodNotFound       = -32601
	ErrInvalidParams        = -32602
	ErrInternalError        = -32603
	ErrServerErrorCodeStart = -32000
	ErrServerErrorCodeEnd   = -32099
)

func ParseCall(jsonStr string) (*Call, error) {
	var rpcCall Call
	err := json.Unmarshal([]byte(jsonStr), &rpcCall)
	if err != nil {
		return nil, err
	}

	if rpcCall.Jsonrpc != "2.0" {
		return nil, errors.New("Error: invalid JSON-RPC version")
	}

	if len(rpcCall.Method) == 0 {
		return nil, errors.New("Error: method is missing")
	}

	if rpcCall.Method[:4] == "rpc." {
		return nil, errors.New("Error: method cannot start with 'rpc.'")
	}

	return &rpcCall, nil
}

func GetIdFromCall(rpcCall *Call) (int, error) {
	if rpcCall.Id == nil {
		return 0, errors.New("Error: id is null")
	}

	switch id := rpcCall.Id.(type) {
	case float64:
		return int(id), nil
	case int:
		return id, nil
	case string:
		parsedInt, err := strconv.Atoi(id)
		if err != nil {
			return 0, errors.New("Error: id is a string but cannot be parsed as an integer")
		}
		return parsedInt, nil
	default:
		return 0, errors.New("Error: unsupported id type")
	}
}

func BuildResponse(id int, result interface{}) string {
	response := Response{
		Jsonrpc: "2.0",
		Result:  result,
		Id:      id,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"error":{"code":-32603,"message":"Internal error"}}`, id)
	}
	return string(jsonResponse)
}

func BuildTimeoutErrorResponse(id int) string {
	return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"error":{"code":-32000,"message":"Timeout (from byzantine-shield)"}}`, id)
}
