package ontmobile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/constants"
	"github.com/ontio/ontology/common/serialization"
)

//ASi48wqdF9avm91pWwdphcAmaDJQkPNdNt
func GetOffset(addressBase58 string, rpcEndpoint string) {
	rpc = rpcEndpoint
	address = addressBase58
	addr, err := common.AddressFromBase58(address)
	if err != nil {
		fmt.Printf("address %s invalid", addr)
		return
	}
	str := genAddressUnboundOffsetKey(addr)
	value, err := SendRpcRequest("getstorage", []interface{}{OntContractAddress.ToHexString(), common.ToHexString(str)})
	if err != nil {
		fmt.Println("rpc requset error:", err)
		return
	}
	v, err := serialization.ReadUint32(bytes.NewBuffer(value))
	if err != nil {
		fmt.Println("read timestamp error:", err)
	}

	fmt.Println("last timestamp:", v+constants.GENESIS_BLOCK_TIMESTAMP)
	fmt.Println("last transfer time: ", GetTimeFormat(int64(v+constants.GENESIS_BLOCK_TIMESTAMP), "2006-01-02 03:04:05 PM"))
}

const (
	UNBOUND_TIME_OFFSET = "unboundTimeOffset"
)

var (
	OntContractAddress, _ = common.AddressParseFromBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
	address               string
	rpc                   string
)

//JsonRpcRequest object in rpc
type JsonRpcRequest struct {
	Version string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

//JsonRpcResponse object response for JsonRpcRequest
type JsonRpcResponse struct {
	Error  int64  `json:"error"`
	Desc   string `json:"desc"`
	Result string `json:"result"`
}

func genAddressUnboundOffsetKey(address common.Address) []byte {
	return append([]byte(UNBOUND_TIME_OFFSET), address[:]...)
}

func GetTimeFormat(second int64, format string) string {
	return time.Unix(second, 0).Format(format)
}

func SendRpcRequest(method string, params []interface{}) ([]byte, error) {
	rpcReq := &JsonRpcRequest{
		Id:     "1",
		Method: method,
		Params: params,
	}
	data, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("JsonRpcRequest json.Marsha error:%s", err)
	}
	resp, err := http.Post(rpc, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("http post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rpc response body error:%s", err)
	}
	rpcRsp := &JsonRpcResponse{}
	err = json.Unmarshal(body, rpcRsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal JsonRpcResponse:%s error:%s", body, err)
	}
	if rpcRsp.Error != 0 {
		return nil, fmt.Errorf("error code:%d desc:%s", rpcRsp.Error, rpcRsp.Desc)
	}
	return common.HexToBytes(rpcRsp.Result)

}
