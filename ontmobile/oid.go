package ontmobile

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	cutils "github.com/ontio/ontology/core/utils"
)

var OIDContractAddress = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}

func MakeRegister(gasPrice uint, gasLimit uint, ontidWif string, payerWif string) (string, error) {
	ontidAccount := ONTAccountWithWIF(ontidWif)
	if ontidAccount == nil {
		return "", fmt.Errorf("Invalid ontidWif")
	}

	payerAccount := ONTAccountWithWIF(payerWif)
	if payerAccount == nil {
		return "", fmt.Errorf("Invalid payerWif")
	}
	payer := payerAccount.account.Address

	contractAddress, _ := common.AddressParseFromBytes(OIDContractAddress)

	ontid := fmt.Sprintf("did:ont:%s", ontidAccount.Address)
	log.Printf("ontid: %s", ontid)

	cversion := byte(0)
	method := "regIDWithPublicKey"
	structs := []interface{}{ontid, ontidAccount.PublicKey}
	params := []interface{}{structs}

	mutableTx, err := newNativeInvokeTransaction(uint64(gasPrice), uint64(gasLimit), contractAddress, cversion, method, params)
	if err != nil {
		log.Printf("NewNativeInvokeTransaction error: %s", err)
		return "", err
	}

	signers := []*account.Account{ontidAccount.account, payerAccount.account}

	err = signMultiTransaction(mutableTx, signers, payer)
	if err != nil {
		log.Printf("SignTransaction error: %s", err)
		return "", err
	}

	tx, err := mutableTx.IntoImmutable()
	if err != nil {
		return "", fmt.Errorf("[Failed to convert tx to immutable: %s]", err)
	}

	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		log.Printf("Serialize error:%s", err)
		return "", err
	}

	txData := hex.EncodeToString(buffer.Bytes())
	return txData, nil
}

func BuildGetDDO(ontid string) (string, error) {
	contractAddress, _ := common.AddressParseFromBytes(OIDContractAddress)

	cversion := byte(0)
	method := "getDDO"
	params := []interface{}{ontid}

	mutableTx, err := newNativeInvokeTransaction(0, 0, contractAddress, cversion, method, params)
	if err != nil {
		log.Printf("NewNativeInvokeTransaction error: %s", err)
		return "", err
	}

	tx, err := mutableTx.IntoImmutable()
	if err != nil {
		return "", fmt.Errorf("[Failed to convert tx to immutable: %s]", err)
	}

	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		log.Printf("Serialize error:%s", err)
		return "", err
	}

	txData := hex.EncodeToString(buffer.Bytes())
	return txData, nil
}

func newNativeInvokeTransaction(gasPrice uint64, gasLimit uint64, contractAddress common.Address, version byte,
	method string, params []interface{}) (*types.MutableTransaction, error) {
	invokeCode, err := cutils.BuildNativeInvokeCode(contractAddress, version, method, params)
	if err != nil {
		return nil, err
	}
	return NewSmartContractTransaction(gasPrice, gasLimit, invokeCode)
}

func signMultiTransaction(tx *types.MutableTransaction, signers []*account.Account, payer common.Address) error {
	if tx.Payer == common.ADDRESS_EMPTY {
		tx.Payer = payer
	}
	txHash := tx.Hash()
	for _, signer := range signers {
		sigData, err := signToData(txHash.ToArray(), signer)
		if err != nil {
			return fmt.Errorf("signToData error:%s", err)
		}

		tx.Sigs = append(tx.Sigs, types.Sig{
			PubKeys: []keypair.PublicKey{signer.PublicKey},
			M:       1,
			SigData: [][]byte{sigData},
		})
	}
	return nil
}
