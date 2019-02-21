package ontmobile_test

import (
	"log"
	"testing"

	"github.com/o3labs/ont-mobile/ontmobile"
	"github.com/o3labs/ont-mobile/ontmobile/ontrpc"
)

func TestMakeRegister(t *testing.T) {
	ontidWif := ontmobile.NewONTAccount().WIF
	payerWif := ontmobile.NewONTAccount().WIF

  gasPrice := uint(500)
  gasLimit := uint(20000)

	raw, err := ontmobile.MakeRegister(gasPrice, gasLimit, ontidWif, payerWif)
	if err != nil {
		log.Printf("Error calling make register: %s", err)
		t.Fail()
	} else {
		log.Printf("Raw transaction: %s", raw)
	}

	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
	res, err := client.SendRawTransaction(raw)
	if err != nil {
		log.Printf("Error invoking: %s", err)
		t.Fail()
	} else {
		log.Printf("Response: %v", res)
	}
}
