package ontmobile_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/o3labs/ont-mobile/ontmobile"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/program"
	"github.com/ontio/ontology/core/types"
)

func TestNewAccount(t *testing.T) {
	account := ontmobile.NewONTAccount()

	fmt.Printf("%+v", account)

	address := ontmobile.ONTAddressFromPublicKey(account.PublicKey)
	fmt.Printf("%+v", address)
}

func TestAccountFromWIF(t *testing.T) {
	wif := ""
	account := ontmobile.ONTAccountWithWIF(wif)
	fmt.Printf("address %v ", account.Address)
	fmt.Printf("account %+v", account)
	fmt.Printf("account public %x", account.PublicKey)
}

func TestMultiSigAddress(t *testing.T) {

	pri1, err := keypair.WIF2Key([]byte(""))
	if err != nil {
		return
	}
	pub1 := pri1.Public()

	pri2, err := keypair.WIF2Key([]byte(""))
	if err != nil {
		return
	}
	pub2 := pri2.Public()

	pri3, err := keypair.WIF2Key([]byte(""))
	if err != nil {
		return
	}
	pub3 := pri3.Public()
	log.Printf("%x", keypair.SerializePublicKey(pub3))
	pubKeys := []keypair.PublicKey{}
	pubKeys = append(pubKeys, pub1)
	pubKeys = append(pubKeys, pub2)
	pubKeys = append(pubKeys, pub3)

	program, err := program.ProgramFromMultiPubKey(pubKeys, 2)
	log.Printf("%x %v", program, err)

	address, err := types.AddressFromMultiPubKeys(pubKeys, 2)
	log.Printf("%v %v", address.ToBase58(), err)
}

func TestONGContractAddress(t *testing.T) {

	address, _ := common.AddressParseFromBytes(ontmobile.ONGContractAddress)
	log.Printf("%v", address.ToBase58())
}

func TestWithdrawONG(t *testing.T) {
	wif := ""

	gasPrice := uint(500)
	gasLimit := uint(20000)
	endpoint := "http://dappnode2.ont.io:20336"

	raw, err := ontmobile.WithdrawONG(gasPrice, gasLimit, endpoint, wif)
	if err != nil {
		log.Printf("error %v", err)
		t.Fail()
	}
	fmt.Printf("\n\ntx = %+v %x\n", raw.TXID, raw.Data)

}

func TestGetAddressUnboundOffsetKey(t *testing.T) {
	UNBOUND_TIME_OFFSET := "unboundTimeOffset"
	address, _ := common.AddressFromBase58("ASi48wqdF9avm91pWwdphcAmaDJQkPNdNt")

	temp := append([]byte(UNBOUND_TIME_OFFSET), address[:]...)
	log.Printf("%x", temp)
}
