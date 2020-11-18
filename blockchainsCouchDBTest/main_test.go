package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

func TestSimple(t *testing.T){
	cc := new(SampleChaincode)
	stub := shimtest.NewMockStub("couchdb",cc)
	stub.MockInvoke("1",[][]byte{[]byte("simpleSave"),[]byte("1"),[]byte("xiaoxiao"),[]byte("aoi"),[]byte("20")})
	stub.MockInvoke("1",[][]byte{[]byte("simpleSave"),[]byte("2"),[]byte("dada"),[]byte("aka"),[]byte("18")})
	stub.MockInvoke("1",[][]byte{[]byte("simpleSave"),[]byte("3"),[]byte("xiaoqiang"),[]byte("midori"),[]byte("28")})
	stub.MockInvoke("1",[][]byte{[]byte("simpleSave"),[]byte("4"),[]byte("xiaoming"),[]byte("midori"),[]byte("22")})

	resp := stub.MockInvoke("1",[][]byte{[]byte("colorQuery"),[]byte("midori")})
	if resp.Status == 200{
		t.Log(resp.String())
	}else{
		t.Log(resp.Message)
	}

	resp = stub.MockInvoke("1",[][]byte{[]byte("simpleRichQuery"),[]byte("midori")})
	if resp.Status == 200{
		t.Log(resp.String())
	}else{
		t.Log(resp.Message)
	}

	resp = stub.MockInvoke("1",[][]byte{[]byte("simpleQueryHistory"),[]byte("1")})
	if resp.Status == 200{
		t.Log(resp.String())
	}else{
		t.Log(resp.Message)
	}
}