package ontmobile_test

import (
	"encoding/json"
	"log"
	"math"
	"testing"

	ont "github.com/o3labs/ont-mobile/ontmobile"
	"github.com/o3labs/ont-mobile/ontmobile/ontrpc"
)

func TestTransferOEP4(t *testing.T) {
	//oep4 tetnet 35666bb22c59d20925d7c761a4d1088be52f000d
	contract := "35666bb22c59d20925d7c761a4d1088be52f000d"
	wif := ""
	account := ont.ONTAccountWithWIF(wif)
	address := "AMgVktoAhY8wX8byyNjUX3Jhiq94T7hSak" //another wallet on cyano
	payer := account.Address
	decimals := 8
	amount := 900.00099

	transferringAmount := uint(ont.RoundFixed(float64(amount), decimals) * float64(math.Pow10(decimals)))
	fromAddress := ont.ParameterJSONForm{T: "Address", V: account.Address}
	addr := ont.ParameterJSONForm{T: "Address", V: address}
	val := ont.ParameterJSONForm{T: "Integer", V: transferringAmount}

	jsondat := &ont.ParameterJSONArrayForm{A: []ont.ParameterJSONForm{fromAddress, addr, val}}
	argData, _ := json.Marshal(jsondat)
	argString := string(argData)

	gasPrice := uint(500)
	gasLimit := uint(20000)

	txData, err := ont.BuildInvocationTransaction(contract, "transfer", argString, gasPrice, gasLimit, wif, payer)
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

// func TestBuildTransaction(t *testing.T) {
// 	wif := ""
// 	account := ont.ONTAccountWithWIF(wif)
// 	address := account.Address

// 	addr := ont.ParameterJSONForm{T: "Address", V: address}
// 	val := ont.ParameterJSONForm{T: "String", V: "Hi there"}

// 	jsondat := &ont.ParameterJSONArrayForm{A: []ont.ParameterJSONForm{addr, val}}
// 	argData, _ := json.Marshal(jsondat)
// 	argString := string(argData)

// 	gasPrice := uint(500)
// 	gasLimit := uint(20000)

// 	txData, err := ont.BuildInvocationTransaction("c168e0fb1a2bddcd385ad013c2c98358eca5d4dc", "put", argString, gasPrice, gasLimit, wif)
// 	if err != nil {
// 		log.Printf("Error creating invocation transaction: %s", err)
// 		t.Fail()
// 	} else {
// 		log.Printf("Raw transaction: %s", txData)
// 	}

// 	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
// 	res, err := client.SendRawTransaction(txData)
// 	if err != nil {
// 		log.Printf("Error invoking: %s", err)
// 		t.Fail()
// 	} else {
// 		log.Printf("Response: %v", res)
// 	}
// }
