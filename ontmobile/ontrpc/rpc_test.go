package ontrpc_test

import (
	"log"
	"testing"

	"github.com/o3labs/ont-mobile/ontmobile/ontrpc"
)

func TestGetONTBalance(t *testing.T) {
	log.Printf("TestGetONTBalance")
	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
	response, err := client.GetBalance("AULUHtmLSATxrJPMGoiWrgWyCKfzTkwEV2")
	if err != nil {
		t.Fail()
	}
	log.Printf("%v", response)
}

func TestGetBlockCount(t *testing.T) {
	log.Printf("TestGetBlockCount")
	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
	response, err := client.GetBlockCount()
	if err != nil {
		t.Fail()
	}
	log.Printf("%v", response)
}

func TestGetSmartCodeEvent(t *testing.T) {
	log.Printf("TestGetSmartCodeEvent")
	client := ontrpc.NewRPCClient("http://polaris2.ont.io:20336")
	response, err := client.GetSmartCodeEvent("41ee265bf50952cd0445d0f612bf2574af523b741c9cc82617bd27c0f7404b14")
	if err != nil {
		t.Fail()
	}
	log.Printf("%v", response)
}

func TestGetUnboundONG(t *testing.T) {
	log.Printf("TestGetUnboundONG")
	client := ontrpc.NewRPCClient("http://dappnode2.ont.io:20336")
	response, err := client.GetUnboundONG("AeNkbJdiMx49kBStQdDih7BzfDwyTNVRfb")
	if err != nil {
		t.Fail()
	}
	log.Printf("%v", response)
}
