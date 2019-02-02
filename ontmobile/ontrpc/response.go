package ontrpc

type JSONRPCResponse struct {
	Desc   string `json:"desc"`
	Error  int    `json:"error"`
	ID     int    `json:"id"`
	Jsonpc string `json:"jsonpc"`
}

type GetBalanceResponse struct {
	JSONRPCResponse
	Result struct {
		Ont string `json:"ont"`
		Ong string `json:"ong"`
	} `json:"result"`
}

type GetBlockCountResponse struct {
	JSONRPCResponse
	Result int `json:"result"`
}

type GetUnboundONGResponse struct {
	JSONRPCResponse
	Result string `json:"result"`
}

type GetSmartCodeEventResponse struct {
	JSONRPCResponse
	Result struct {
		TxHash      string `json:"TxHash"`
		State       int    `json:"State"`
		GasConsumed int    `json:"GasConsumed"`
		Notify      []struct {
			ContractAddress string        `json:"ContractAddress"`
			States          []interface{} `json:"States"`
		} `json:"Notify"`
	} `json:"result"`
}

type SendRawTransactionResponse struct {
	JSONRPCResponse
	//this is when we send 1
	// Result struct {
	// 	State  int    `json:"State"`
	// 	Gas    int    `json:"Gas"`
	// 	Result string `json:"Result"`
	// } `json:"result"`
	Result string `json:"result"`
}

type GetStorageResponse struct {
	JSONRPCResponse
	Result string `json:"result"`
}

type GetRawTransactionResponse struct {
	JSONRPCResponse
	Result string `json:"result"`
}

type GetBlockResponse struct {
	JSONRPCResponse
	Result string `json:"result"`
}
