package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

func TestFunc(t *testing.T) {
	cc := new(SimpleAsset)
	stub := shimtest.NewMockStub("sacc",cc)
	stub.MockInit("1",[][]byte{[]byte("a"),[]byte("bb")})
	res := stub.MockInvoke("1",[][]byte{[]byte("get"),[]byte("a")})
	t.Log(string(res.Payload))
	stub.MockInvoke("1",[][]byte{[]byte("set"),[]byte("a"),[]byte("aa")})

	res = stub.MockInvoke("1",[][]byte{[]byte("get"),[]byte("a")})
	t.Log(string(res.Payload))
}

func TestPeople(t *testing.T){
	cc := new(SimpleAsset)
	stub := shimtest.NewMockStub("sacc",cc)
	stub.MockInvoke("1",[][]byte{[]byte("set2"),[]byte("1"),[]byte("xiaoxiao")})
	res := stub.MockInvoke("1",[][]byte{[]byte("get"),[]byte("1")})
	if res.Status == 200 {
		t.Log(string(res.Payload))
	}else{
		t.Error(res.Message)
	}
	stub.MockInvoke("1",[][]byte{[]byte("set"),[]byte("a"),[]byte("b")})
	res = stub.MockInvoke("1",[][]byte{[]byte("get2"),[]byte("a")})
	if res.Status == 200 {
		t.Log(string(res.Payload))
	}else{
		t.Error(res.Message)
	}

	res = stub.MockInvoke("1",[][]byte{[]byte("getByIdAndName"),[]byte("1"),[]byte("xiaoxiao")})
	if res.Status == 200{
		t.Log(string(res.Payload))
	}else{
		t.Error(res.Message)
	}

	res = stub.MockInvoke("1",[][]byte{[]byte("getCreator")})
	if res.Status == 200{
		t.Logf("creator: %s,Status: %d,Message: %s",res.Payload,res.Status,res.Message)
	}else{
		t.Error(res.Message)
	}
}