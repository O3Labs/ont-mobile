package ontmobile

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"

	"github.com/ontio/ontology/core/payload"
	cutils "github.com/ontio/ontology/core/utils"
	"github.com/ontio/ontology/vm/neovm"
	"time"
)

type ParameterJSONArrayForm struct {
	A []ParameterJSONForm `json:"array"`
}

type ParameterJSONForm struct {
	T string      `json:"type"`
	V interface{} `json:"value"`
}

func buildParameters(argString string) []interface{} {
	var parameters []interface{}

	data := &ParameterJSONArrayForm{
		A: []ParameterJSONForm{},
	}

	err := json.Unmarshal([]byte(argString), data)

	for _, element := range data.A {
		var t = element.T
		var v = element.V
		var p interface{}
		if t == "Address" {
			p, err = common.AddressFromBase58(v.(string))
			if err != nil {
				log.Printf("Failed to convert string to address %s", err)
			}
		} else if t == "String" {
			p = v.(string)
		} else if t == "Integer" {
			p = uint(v.(float64))
		} else if t == "Fixed8" {
			p = uint(RoundFixed(v.(float64), 8) * float64(math.Pow10(8)))
		} else if t == "Fixed9" {
			p = uint(RoundFixed(v.(float64), ONGDECIMALS) * float64(math.Pow10(ONGDECIMALS)))
		} else if t == "Array" {
			p = buildParameters(v.(string))
		}
		parameters = append(parameters, p)
	}

	return parameters
}

// BuildInvocationTransaction : creates a raw transaction
func BuildInvocationTransaction(contractHex string, operation string, argString string, gasPrice uint, gasLimit uint, wif string, payer string) (string, error) {
	var contractAddress common.Address
	contractAddress, err := common.AddressFromHexString(contractHex)
	if err != nil {
		return "", fmt.Errorf("[Invalid contract hash error: %s]", err)
	}

	signer := ONTAccountWithWIF(wif).account
	parameters := buildParameters(argString)
	params := []interface{}{operation, parameters}

	tx, err := newNeovmInvokeTransaction(uint64(gasPrice), uint64(gasLimit), contractAddress, params)
	if err != nil {
		log.Printf("NewNeovmInvokeTransaction error:%s", err)
		return "", err
	}

	err = signTransaction(tx, signer, payer)
	if err != nil {
		log.Printf("SignTransaction error: %s", err)
		return "", err
	}

	immutTx, err := tx.IntoImmutable()

	if err != nil {
		log.Printf("IntoImmutable error: %s", err)
		return "", err
	}

	var buffer bytes.Buffer
	err = immutTx.Serialize(&buffer)
	if err != nil {
		log.Printf("Serialize error:%s", err)
		return "", err
	}

	txData := hex.EncodeToString(buffer.Bytes())
	return txData, nil
}

func signTransaction(tx *types.MutableTransaction, signer *account.Account, payer string) error {
	if tx.Payer == common.ADDRESS_EMPTY {
		addr, err := common.AddressFromBase58(payer)
		if err == nil {
			tx.Payer = addr
		}
	}
	txHash := tx.Hash()

	sigData, err := signToData(txHash.ToArray(), signer)
	if err != nil {
		return fmt.Errorf("signToData error:%s", err)
	}

	tx.Sigs = append(tx.Sigs, types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sigData},
	})

	return nil
}

func newNeovmInvokeTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, params []interface{}) (*types.MutableTransaction, error) {
	invokeCode, err := buildNeoVMInvokeCode(contractAddress, params)
	if err != nil {
		return nil, err
	}
	return newSmartContractTransaction(gasPrice, gasLimit, invokeCode)
}

func buildNeoVMInvokeCode(smartContractAddress common.Address, params []interface{}) ([]byte, error) {
	builder := neovm.NewParamsBuilder(new(bytes.Buffer))
	err := cutils.BuildNeoVMParam(builder, params)
	if err != nil {
		return nil, err
	}
	args := append(builder.ToArray(), 0x67)
	args = append(args, smartContractAddress[:]...)
	return args, nil
}

func newSmartContractTransaction(gasPrice, gasLimit uint64, invokeCode []byte) (*types.MutableTransaction, error) {
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     nil,
	}
	return tx, nil
}
