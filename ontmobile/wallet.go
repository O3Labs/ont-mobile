package ontmobile

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/o3labs/ont-mobile/ontmobile/ontrpc"
	"github.com/ontio/ontology-crypto/keypair"
	sig "github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/core/utils"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
)

const (
	ONGDECIMALS = 9
)

var ONTContractAddress = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
var ONGContractAddress = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}

//package sm2 need to be on commit
//ontology-crypto commit 39fe0d3acce904abff68045b94c7d45f8cc903c2

type ONTAccount struct {
	Address    string //base58
	WIF        string
	PrivateKey []byte
	PublicKey  []byte
	account    *account.Account
}

func accountToLocalAccount(account *account.Account) *ONTAccount {
	address := account.Address.ToBase58()
	wifBytes, err := keypair.Key2WIF(account.PrivateKey)
	if err != nil {
		log.Printf("err %v", err)
		return nil
	}
	publicKey := keypair.SerializePublicKey(account.PublicKey)
	privateKey := keypair.SerializePrivateKey(account.PrivateKey)
	return &ONTAccount{
		Address:    address,
		WIF:        string(wifBytes),
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		account:    account,
	}
}

func NewONTAccount() *ONTAccount {
	//default to SHA256withECDSA
	account := account.NewAccount("")
	return accountToLocalAccount(account)
}

func ONTAddressFromPublicKey(publicKeyBytes []byte) string {
	publicKey, err := keypair.DeserializePublicKey(publicKeyBytes)
	if err != nil || publicKey == nil {
		return ""
	}
	address := types.AddressFromPubKey(publicKey)
	return address.ToBase58()
}

func ONTAccountWithPrivateKey(privateKeyBytes []byte) *ONTAccount {
	pri, err := keypair.DeserializePrivateKey(privateKeyBytes)
	if err != nil || pri == nil {
		return nil
	}
	pub := pri.Public()
	address := types.AddressFromPubKey(pub)
	account := &account.Account{
		SigScheme:  sig.SHA256withECDSA,
		PrivateKey: pri,
		PublicKey:  pub,
		Address:    address,
	}

	return accountToLocalAccount(account)
}

func ONTAccountWithWIF(wif string) *ONTAccount {
	var err error

	pri, err := keypair.WIF2Key([]byte(wif))
	if err != nil || pri == nil {
		return nil
	}
	pub := pri.Public()
	address := types.AddressFromPubKey(pub)
	account := &account.Account{
		SigScheme:  sig.SHA256withECDSA,
		PrivateKey: pri,
		PublicKey:  pub,
		Address:    address,
	}

	return accountToLocalAccount(account)
}

type RawTransaction struct {
	TXID string
	Data []byte
}

func Transfer(gasPrice uint, gasLimit uint, senderWIF string, asset string, toAddress string, amount float64) (*RawTransaction, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("Amount must be greater than zero")
	}

	if asset == "" {
		return nil, fmt.Errorf("Invalid asset. ONT or ONG only")
	}

	if toAddress == "" {
		return nil, fmt.Errorf("To address cannot be empty")
	}

	if senderWIF == "" {
		return nil, fmt.Errorf("Sender WIF cannot be empty")
	}

	sender := ONTAccountWithWIF(senderWIF)

	fromAddress := sender.Address

	from, err := common.AddressFromBase58(fromAddress)

	if err != nil {
		return nil, err
	}

	to, err := common.AddressFromBase58(toAddress)
	if err != nil {
		return nil, err
	}

	var contractAddress common.Address

	value := uint64(0)
	switch strings.ToUpper(asset) {
	case "ONT":
		a, _ := common.AddressParseFromBytes(ONTContractAddress)
		contractAddress = a
		value = uint64(amount)
	case "ONG":
		a, _ := common.AddressParseFromBytes(ONGContractAddress)
		contractAddress = a
		value = uint64(RoundFixed(amount, ONGDECIMALS) * float64(math.Pow10(ONGDECIMALS)))
	default:
		return nil, fmt.Errorf("%s is neither ONT nor ONG", asset)
	}

	var sts []*ont.State
	sts = append(sts, &ont.State{
		From:  from,
		To:    to,
		Value: value,
	})

	cversion := byte(0)
	method := "transfer"
	params := []interface{}{sts}
	gasPriceUint64 := uint64(gasPrice)
	gasLimitUint64 := uint64(gasLimit)

	invokeCode, err := utils.BuildNativeInvokeCode(contractAddress, cversion, method, params)

	mutableTx := utils.NewInvokeTransaction(invokeCode)
	mutableTx.GasPrice = gasPriceUint64
	mutableTx.GasLimit = gasLimitUint64
	mutableTx.Nonce = uint32(time.Now().Unix())

	tx, err := mutableTx.IntoImmutable()
	if err != nil {
		return nil, fmt.Errorf("[Failed to convert tx to immutable: %s]", err)
	}

	signer := sender.account
	if signer == nil {
		return nil, fmt.Errorf("Account is null")
	}

	err = signToTransaction(tx, signer)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("Serialize error:%s", err)
	}

	hash := tx.Hash()
	raw := &RawTransaction{
		TXID: hash.ToHexString(),
		Data: buffer.Bytes(),
	}
	return raw, nil
}

func signToTransaction(tx *types.Transaction, signer *account.Account) error {
	tx.Payer = signer.Address
	txHash := tx.Hash()
	sigData, err := signToData(txHash.ToArray(), signer)
	if err != nil {
		return fmt.Errorf("signToData error:%s", err)
	}

	sig := types.RawSig{
		Invoke: keypair.SerializePublicKey(signer.PublicKey),
		Verify: sigData,
	}
	tx.Sigs = []types.RawSig{sig}

	return nil
}

func signToData(data []byte, signer *account.Account) ([]byte, error) {
	s, err := sig.Sign(signer.SigScheme, signer.PrivateKey, data, nil)
	if err != nil {
		return nil, err
	}
	sigData, err := sig.Serialize(s)
	if err != nil {
		return nil, fmt.Errorf("sig.Serialize error:%s", err)
	}
	return sigData, nil
}

func SendRawTransaction(endpoint string, rawTransactionHex string) (string, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.SendRawTransaction(rawTransactionHex)
	if err != nil {
		return "", err
	}
	return response.Result, nil
}

func SendPreExecRawTransaction(endpoint string, rawTransactionHex string) (string, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.SendPreExecRawTransaction(rawTransactionHex)
	if err != nil {
		return "", err
	}
	return response.Result.Result, nil
}

func WithdrawONG(gasPrice uint, gasLimit uint, endpoint string, wif string) (*RawTransaction, error) {

	sender := ONTAccountWithWIF(wif)

	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetUnboundONG(sender.Address)
	if err != nil {
		return nil, err
	}
	unboundONG := response.Result //this is string
	amount, err := strconv.ParseUint(unboundONG, 10, 64)
	if err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, fmt.Errorf("Don't have unbound ong\n")
	}

	signer := sender.account
	if signer == nil {
		return nil, fmt.Errorf("Account is null")
	}

	fromAddress, _ := common.AddressParseFromBytes(ONTContractAddress)

	transferFrom := &ont.TransferFrom{
		Sender: sender.account.Address,
		From:   fromAddress,
		To:     sender.account.Address,
		Value:  amount,
	}

	cversion := byte(0)
	method := "transferFrom"
	gasPriceUint64 := uint64(gasPrice)
	gasLimitUint64 := uint64(gasLimit)

	a, _ := common.AddressParseFromBytes(ONGContractAddress)
	contractAddress := a

	invokeCode, err := utils.BuildNativeInvokeCode(contractAddress, cversion, method, []interface{}{transferFrom})
	if err != nil {
		return nil, fmt.Errorf("build invoke code error:%s", err)
	}

	mutableTx := utils.NewInvokeTransaction(invokeCode)
	mutableTx.GasPrice = gasPriceUint64
	mutableTx.GasLimit = gasLimitUint64
	mutableTx.Nonce = uint32(time.Now().Unix())
	tx, err := mutableTx.IntoImmutable()
	if err != nil {
		return nil, fmt.Errorf("[Failed to convert tx to immutable: %s]", err)
	}

	err = signToTransaction(tx, signer)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("Serialize error:%s", err)
	}

	hash := tx.Hash()
	raw := &RawTransaction{
		TXID: hash.ToHexString(),
		Data: buffer.Bytes(),
	}
	return raw, nil
}
