package ontrpc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type RPCInterface interface {
	makeRequest(method string, params []interface{}, out interface{}) error
	GetBlockCount() (GetBlockCountResponse, error)
	GetBalance(ontAddress string) (GetBalanceResponse, error)
	GetSmartCodeEvent(txHash string) (GetSmartCodeEventResponse, error)
	SendRawTransaction(rawTransactionHex string) (SendRawTransactionResponse, error)
	GetUnboundONG(ontAddress string) (GetUnboundONGResponse, error)
	GetStorage(scriptHash string, key string) (GetStorageResponse, error)
	GetRawTransaction(txID string) (GetRawTransactionResponse, error)
	GetBlockWithHash(blockHash string) (GetBlockResponse, error)
	GetBlockWithHeight(blockHeight int) (GetBlockResponse, error)
}

type RPCClient struct {
	Endpoint   url.URL
	httpClient *http.Client
}

//make sure all method interface is implemented
var _ RPCInterface = (*RPCClient)(nil)

func NewRPCClient(endpoint string) *RPCClient {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil
	}

	var netClient = &http.Client{
		Timeout: time.Second * 60,
		// Transport: netTransport,
	}

	return &RPCClient{Endpoint: *u, httpClient: netClient}
}

func (n *RPCClient) makeRequest(method string, params []interface{}, out interface{}) error {
	request := NewRequest(method, params)

	jsonValue, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", n.Endpoint.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	res, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return err
	}

	return nil
}

func (n *RPCClient) GetBalance(ontAddress string) (GetBalanceResponse, error) {
	response := GetBalanceResponse{}
	params := []interface{}{ontAddress, 1}
	err := n.makeRequest("getbalance", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) GetBlockCount() (GetBlockCountResponse, error) {
	response := GetBlockCountResponse{}
	params := []interface{}{}
	err := n.makeRequest("getblockcount", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) GetSmartCodeEvent(txHash string) (GetSmartCodeEventResponse, error) {
	response := GetSmartCodeEventResponse{}
	params := []interface{}{txHash}
	err := n.makeRequest("getsmartcodeevent", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) SendRawTransaction(rawTransactionHex string) (SendRawTransactionResponse, error) {
	response := SendRawTransactionResponse{}
	params := []interface{}{rawTransactionHex}
	err := n.makeRequest("sendrawtransaction", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) GetUnboundONG(ontAddress string) (GetUnboundONGResponse, error) {
	response := GetUnboundONGResponse{}
	params := []interface{}{ontAddress}
	err := n.makeRequest("getunboundong", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) GetStorage(scriptHash string, key string) (GetStorageResponse, error) {
	response := GetStorageResponse{}
	params := []interface{}{scriptHash, key}
	err := n.makeRequest("getstorage", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) GetRawTransaction(txID string) (GetRawTransactionResponse, error) {
	response := GetRawTransactionResponse{}
	params := []interface{}{txID}
	err := n.makeRequest("getrawtransaction", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) GetBlockWithHash(blockHash string) (GetBlockResponse, error) {
	response := GetBlockResponse{}
	params := []interface{}{blockHash}
	err := n.makeRequest("getblock", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (n *RPCClient) GetBlockWithHeight(blockHeight int) (GetBlockResponse, error) {
	response := GetBlockResponse{}
	params := []interface{}{blockHeight}
	err := n.makeRequest("getblock", params, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
