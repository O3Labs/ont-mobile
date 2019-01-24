package ontmobile_test

import (
  "log"
  "testing"

  ont "github.com/o3labs/ont-mobile/ontmobile"
  "github.com/o3labs/ont-mobile/ontmobile/ontrpc"
)

func TestTransfer(t *testing.T){
  account := ont.NewONTAccount()
  wif := account.WIF
  address := account.Address

  gasPrice := uint(500)
  gasLimit := uint(20000)

  tx, err := ont.Transfer(gasPrice, gasLimit, wif, "ONG", address, 1)
  if err != nil {
    log.Printf("Error creating transfer transaction: %s", err)
    t.Fail()
  } else {
    log.Printf("Transaction ID: %s\nData: %v", tx.TXID, tx.Data)
  }
}

func TestBuildTransaction(t *testing.T){
  wif := ""
  account := ont.ONTAccountWithWIF(wif)
  address := account.Address

  addr := ont.Parameter{ont.Address, address}
  val := ont.Parameter{ont.String, "Hi there"}

  args := []ont.Parameter{addr, val}

  gasPrice := uint(500)
  gasLimit := uint(20000)

  txData, err := ont.BuildInvocationTransaction("c168e0fb1a2bddcd385ad013c2c98358eca5d4dc", "put", args, gasPrice, gasLimit, wif)
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
