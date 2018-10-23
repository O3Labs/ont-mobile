package ontrpc_test

import (
	"log"
	"testing"

	"github.com/o3labs/ont-mobile/ontmobile/ontrpc"
)

func TestGetONTBalance(t *testing.T) {
	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
	response, err := client.GetBalance("AULUHtmLSATxrJPMGoiWrgWyCKfzTkwEV2")
	if err != nil {
		t.Fail()
	}
	log.Printf("%v", response)
}

func TestGetBlockCount(t *testing.T) {
	log.Printf("test hello")
	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
	response, err := client.GetBlockCount()
	if err != nil {
		t.Fail()
	}

	t.Logf("%v", response)
}

func TestGetSmartCodeEvent(t *testing.T) {
	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
	response, err := client.GetSmartCodeEvent("8bd07b3f78c1b57f7a839d964e5faccdaf7eb1e4c24b164bb65184e73345c028")
	if err != nil {
		t.Fail()
	}
	log.Printf("%v", response)
}

func TestGetUnboundONG(t *testing.T) {
	client := ontrpc.NewRPCClient("http://dappnode2.ont.io:20336")
	response, err := client.GetUnboundONG("AeNkbJdiMx49kBStQdDih7BzfDwyTNVRfb")
	if err != nil {
		t.Fail()
	}
	log.Printf("%v", response)
}
