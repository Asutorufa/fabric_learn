package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"log"
	"strconv"
)

type SampleChaincode struct{

}

type CDBT struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
	Size int `json:"size"`
}

func main(){
	err := shim.Start(new(SampleChaincode))
	if err != nil{
		log.Println(err)
	}
}

func (s *SampleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn,args := stub.GetFunctionAndParameters()
	if fn == "simpleSave"{
		return  s.simpleSave(stub,args)
	} else if fn == "simpleQuery"{
		return s.simpleQuery(stub,args)
	}else if fn == "simpleDelete" {
		return s.simpleDelete(stub,args)
	}else if fn == "simpleRichQuery"{
		return s.simpleRichQuery(stub,args)
	}else if fn == "colorQuery"{
		return s.colorQuery(stub,args)
	}else if fn == "queryCreator"{
		return s.queryCreator(stub)
	}else if fn == "simpleQueryHistory" {
		return s.simpleQueryHistory(stub,args)
	}

	return shim.Error("SampleChaincode:Invoke() -> not support Invoke -> "+fn)
}

func (s *SampleChaincode) simpleSave(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 4 {
		return shim.Error("simpleSave() -> incorrect arguments")
	}
	var err error
	data := &CDBT{}
	data.Id = args[0]
	data.Name = args[1]
	data.Color = args[2]
	data.Size,err = strconv.Atoi(args[3])
	if err != nil{
		return shim.Error("simpleSave() -> "+err.Error())
	}

	dataJson,err := json.Marshal(data)
	if err != nil{
		return shim.Error("simpleSave() -> "+err.Error())
	}

	fmt.Printf("save data: %s\n",dataJson)

	err = stub.PutState(data.Id,dataJson)
	if err != nil{
		return shim.Error("simpleSave() -> "+err.Error())
	}

	compositeName := "color~id"
	colorIdIndexKey ,err := stub.CreateCompositeKey(compositeName,[]string{data.Color,data.Id})
	if err != nil{
		return shim.Error(err.Error())
	}
	err = stub.PutState(colorIdIndexKey,[]byte{0x00})
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SampleChaincode) simpleDelete(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1{
		return shim.Error("incorrect arguments")
	}

	dataJson,err := stub.GetState(args[0])
	if err != nil{
		return shim.Error(err.Error())
	}
	if dataJson == nil{
		return shim.Error("result is nil")
	}
	data := &CDBT{}
	err = json.Unmarshal(dataJson,data)
	if err != nil{
		return shim.Error(err.Error())
	}

	err = stub.DelState(args[0])
	if err != nil{
		return shim.Error(err.Error())
	}

	compositeName := "color~id"
	compositeKey,err := stub.CreateCompositeKey(compositeName,[]string{data.Color,data.Id})
	if err != nil{
		return shim.Error(err.Error())
	}
	err = stub.DelState(compositeKey)
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

func (s *SampleChaincode)simpleQuery(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1{
		return shim.Error("simpleQuery() -> incorrect arguments")
	}

	res,err := stub.GetState(args[0])
	if err != nil{
		return shim.Error(err.Error())
	}
	if res == nil{
		return shim.Error("result is nil")
	}
	return shim.Success(res)
}

func (s *SampleChaincode) simpleRichQuery(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1 {
		return shim.Error("incorrect arguments")
	}
	data,err := stub.GetQueryResult(fmt.Sprintf("{\"selector\":{\"color\":\"%s\"}}",args[0]))
	if err != nil{
		return shim.Error(err.Error())
	}

	res := bytes.NewBufferString("[")

	for data.HasNext(){
		x,err := data.Next()
		if err != nil{
			return shim.Error(err.Error())
		}
		res.WriteString("{\"key\":")
		res.WriteString(x.Key)
		res.WriteString(",\"record\":")
		res.Write(x.Value)
		res.WriteString("}")
		if data.HasNext() {
			res.WriteString(",")
		}
	}
	res.WriteString("]")
	return shim.Success(res.Bytes())
}

func (s *SampleChaincode)colorQuery(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1{
		return shim.Error("incorrect arguments")
	}
	compositeName := "color~id"
	resultIterator,err := stub.GetStateByPartialCompositeKey(compositeName,args)
	if err != nil{
		return shim.Error(err.Error())
	}

	res := bytes.NewBuffer(nil)
	res.WriteString("[")
	for resultIterator.HasNext(){
		x,err := resultIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}
		_,y,err := stub.SplitCompositeKey(x.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		z,err := stub.GetState(y[1])
		if err != nil{
			return shim.Error(err.Error())
		}

		res.Write(z)
		if resultIterator.HasNext(){
			res.WriteString(",")
		}
	}
	res.WriteString("]")
	return shim.Success(res.Bytes())
}

func (s *SampleChaincode)queryCreator(stub shim.ChaincodeStubInterface) peer.Response{
	channelID := stub.GetChannelID()
	creator,err := stub.GetCreator()
	if err != nil{
		log.Println(err)
	}
	decoration := stub.GetDecorations()

	id,err := cid.GetID(stub)
	if err != nil{
		log.Println(err)
	}
	mspid,err := cid.GetMSPID(stub)
	if err != nil{
		log.Println(err)
	}
	cert, err := cid.GetX509Certificate(stub)
	if err != nil{
		log.Println(err)
	}
	return shim.Success([]byte(fmt.Sprintf("channel_id: %s, creator: %s, decoration: %v, id: %s, mspid: %s, cert: %v",channelID,creator,decoration,id,mspid,cert)))
}

func (s *SampleChaincode)simpleQueryHistory(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1{
		return shim.Error("incorrect arguments")
	}

	historyItrator,err := stub.GetHistoryForKey(args[0])
	if err != nil{
		return shim.Error(err.Error())
	}

	res := bytes.NewBuffer(nil)

	for historyItrator.HasNext() {
		x,err := historyItrator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		xx,err := json.Marshal(x)
		if err != nil{
			return shim.Error(err.Error())
		}

		res.Write(xx)
		if historyItrator.HasNext(){
			res.Write([]byte(","))
		}
	}
	res.WriteString("]")
	return shim.Success(res.Bytes())
}
//
//func (s *SampleChaincode)ClientIdentityChaincodeExample(stub shim.ChaincodeStubInterface) peer.Response{
//	id,err := cid.GetID(stub)
//	if err != nil{
//		return shim.Error(err.Error())
//	}
//	fmt.Println(id)
//
//	mspid,err := cid.GetMSPID(stub)
//	if err != nil{
//		return shim.Error(err.Error())
//	}
//	switch mspid {
//	case "org1MSP":
//	case "org2MSP":
//
//	}

	//val, ok, err := cid.GetAttributeValue(stub, "attr1")
	//if err != nil {
		// There was an error trying to retrieve the attribute
	//}
	//if !ok {
		// The client identity does not possess the attribute
	//}
	// Do something with the value of 'val'
//}