package client

import (
	"encoding/json"
	"fmt"
)

//nolint:gochecknoglobals
var requestID = int64(0)

const JSONRPCVersion = "2.0"

type (
	RPCRequests []*RPCRequest

	RPCRequest struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
		ID      int64       `json:"id,omitempty"`
	}

	RPCResponse struct {
		JSONRPC string      `json:"jsonrpc"`
		Error   *RPCError   `json:"error,omitempty"`
		Result  interface{} `json:"result,omitempty"`
		ID      int64       `json:"id,omitempty"`
	}

	RPCResponseRaw struct {
		JSONRPC string          `json:"jsonrpc"`
		Error   *RPCError       `json:"error,omitempty"`
		Result  json.RawMessage `json:"result,omitempty"`
		ID      int64           `json:"id,omitempty"`
	}

	RPCError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

func (r *Request) RPCCall(result interface{}, method string, params interface{}) error {
	req := &RPCRequest{JSONRPC: JSONRPCVersion, Method: method, Params: params, ID: genID()}

	var resp *RPCResponse
	if err := r.Post(&resp, "", req); err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	return resp.GetObject(result)
}

func (r *Request) RPCCallRaw(method string, params interface{}) ([]byte, error) {
	req := &RPCRequest{JSONRPC: JSONRPCVersion, Method: method, Params: params, ID: genID()}
	var resp *RPCResponseRaw

	if err := r.Post(&resp, "", req); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	return []byte(resp.Result), nil
}

func (r *Request) RPCBatchCall(requests RPCRequests) ([]RPCResponse, error) {
	var resp []RPCResponse

	if err := r.Post(&resp, "", requests.fillDefaultValues()); err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

func (r *RPCResponse) GetObject(toType interface{}) error {
	js, err := json.Marshal(r.Result)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	if err = json.Unmarshal(js, toType); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

func (rs RPCRequests) fillDefaultValues() RPCRequests {
	for _, r := range rs {
		r.JSONRPC = JSONRPCVersion
		r.ID = genID()
	}

	return rs
}

func genID() int64 {
	requestID++

	return requestID
}
