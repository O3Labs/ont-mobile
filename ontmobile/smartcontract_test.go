package ontmobile

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"

	"github.com/ontio/ontology/cmd/utils"
	"github.com/ontio/ontology/common"
	httpcom "github.com/ontio/ontology/http/base/common"
)

var contractAddress = ""

func TestSerializeInvoke(t *testing.T) {
	gasPrice := uint64(500)
	gasLimit := uint64(20000)
	wif := ""
	signer := ONTAccountWithWIF(wif)
	log.Printf(signer.Address)
	log.Printf("pb %x", signer.PublicKey)

	smartcodeAddress, err := common.AddressFromHexString("c168e0fb1a2bddcd385ad013c2c98358eca5d4dc")

	if err != nil {
		log.Printf("AddressFromHexString error:%s", err)
		return
	}

	operation := "put"
	args := []interface{}{}
	args = append(args, common.Address(signer.account.Address))
	args = append(args, string("hey hey hey"))
	params := []interface{}{operation, args}
	tx, err := httpcom.NewNeovmInvokeTransaction(gasPrice, gasLimit, smartcodeAddress, params)
	tx.Payer = signer.account.Address
	tx.Nonce = uint32(0)
	if err != nil {
		log.Printf("NewNeovmInvokeTransaction error:%s", err)
		return
	}

	err = utils.SignTransaction(signer.account, tx)
	if err != nil {
		log.Printf("signToTransaction error:%s", err)
		return
	}

	immutTx, err := tx.IntoImmutable()

	if err != nil {
		log.Printf("IntoImmutable error:%s", err)
		return
	}
	log.Printf("%+v", tx.Sigs)
	var buffer bytes.Buffer
	err = immutTx.Serialize(&buffer)
	if err != nil {
		log.Printf("serialize error:%s", err)
		return
	}
	txData := hex.EncodeToString(buffer.Bytes())

	log.Printf("%v", txData)
}

func TestGetRyuCoin(t *testing.T) {
	gasPrice := uint64(500)
	gasLimit := uint64(20000)
	// wif := ""
	// signer := ONTAccountWithWIF(wif)
	// log.Printf(signer.Address)

	smartcodeAddress, err := common.AddressFromHexString("c168e0fb1a2bddcd385ad013c2c98358eca5d4dc")

	if err != nil {
		log.Printf("AddressFromHexString error:%s", err)
		return
	}
	a, _ := common.AddressFromBase58("ANkfzLQyEYH9GwW3DJVNQgYRwG58rfXXzM")

	args := []interface{}{}
	args = append(args, a)

	params := []interface{}{string("get"), args}
	tx, err := httpcom.NewNeovmInvokeTransaction(gasPrice, gasLimit, smartcodeAddress, params)
	// tx.Payer = signer.account.Address

	if err != nil {
		log.Printf("NewNeovmInvokeTransaction error:%s", err)
		return
	}

	// err = utils.SignTransaction(signer.account, tx)
	// if err != nil {
	// 	log.Printf("signToTransaction error:%s", err)
	// 	return
	// }

	immutTx, err := tx.IntoImmutable()

	if err != nil {
		log.Printf("IntoImmutable error:%s", err)
		return
	}

	var buffer bytes.Buffer
	err = immutTx.Serialize(&buffer)
	if err != nil {
		log.Printf("serialize error:%s", err)
		return
	}
	txData := hex.EncodeToString(buffer.Bytes())

	log.Printf("%v", txData)
}
