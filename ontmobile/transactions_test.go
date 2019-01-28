package ontmobile_test

import (
  "log"
  "testing"
  "encoding/json"

  ont "github.com/o3labs/ont-mobile/ontmobile"
  "github.com/o3labs/ont-mobile/ontmobile/ontrpc"
)

func TestBuildTransaction(t *testing.T){
  wif := ""
  account := ont.ONTAccountWithWIF(wif)
  address := account.Address

  addr := ont.ParameterJSONForm{T: "Address", V: address}
  val := ont.ParameterJSONForm{T: "String", V: "Hi there"}

  jsondat := &ont.ParameterJSONArrayForm{A: []ont.ParameterJSONForm{addr, val}}
  argData, _ := json.Marshal(jsondat)
  argString := string(argData)

  gasPrice := uint(500)
  gasLimit := uint(20000)

  txData, err := ont.BuildInvocationTransaction("c168e0fb1a2bddcd385ad013c2c98358eca5d4dc", "put", argString, gasPrice, gasLimit, wif)
  if err != nil {
    log.Printf("Error creating invocation transaction: %s", err)
    t.Fail()
  } else {
    log.Printf("Raw transaction: %s", txData)
  }

  client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
  res, err := client.SendRawTransaction(txData)
  if err != nil {
    log.Printf("Error invoking: %s", err)
    t.Fail()
  } else {
    log.Printf("Response: %v", res)
  }
}
