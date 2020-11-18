package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	shim "github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type SimpleAsset struct{

}

func (t *SimpleAsset)Init(stub shim.ChaincodeStubInterface) peer.Response{
	// 获取Init调用的参数, 并检查合法性
	_,args := stub.GetFunctionAndParameters()
	if len(args) != 2{
		return shim.Error("Incorrect argument")
	}
	// 将初始状态存入账本
	err := stub.PutState(args[0],[]byte(args[1]))
	if err != nil{
		return shim.Error(err.Error())
	}
	fmt.Printf("Save Data: (%s , %s)\n",args[0],args[1])
	// 返回一个peer.Response对象
	return shim.Success(nil)
}

func (t *SimpleAsset)Invoke(stub shim.ChaincodeStubInterface) peer.Response{
	// 为链码应用程序的方法解析方法名和参数
	fn,args := stub.GetFunctionAndParameters()
	// 验证函数名
	if fn == "set"{
		return t.set(stub,args)
	}else if fn == "get"{
		return t.get(stub,args)
	}else if fn == "set2"{
		return set2(stub,args)
	}else if fn == "get2"{
		return get2(stub,args)
	}else if fn == "getByIdAndName"{
		return getByIdAndName(stub,args)
	}else if fn == "getCreator" {
		return getCreator(stub)
	}
	return shim.Error("error function name")
}

func (t *SimpleAsset) set(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 2{
		return  shim.Error("incorrect arguments")
	}
	err := stub.PutState(args[0],[]byte(args[1]))
	if err != nil{
		return shim.Error(err.Error())
	}
	fmt.Printf("Save Data: (%s , %s)\n",args[0],args[1])
	return shim.Success(nil)
}

func (t *SimpleAsset)get(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1{
		return  shim.Error("incorrect arguments")
	}
	value,err := stub.GetState(args[0])
	if err != nil{
		return shim.Error(err.Error())
	}
	if value == nil{
		return shim.Error("result is nil")
	}
	return shim.Success(value)
}

type people struct{
	Id string `json:"id"`
	Name string `json:"name"`
}

func set2(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 2{
		return shim.Error("incorrect arguments")
	}
	indexName := "id~name"
	p := people{Id: args[0],Name: args[1]}
	idNameIndexKey,err := stub.CreateCompositeKey(indexName,[]string{p.Id,p.Name})
	if err != nil{
		return shim.Error(err.Error())
	}
	err = stub.PutState(idNameIndexKey,[]byte{0x00})
	if err != nil{
		return shim.Error(err.Error())
	}
	people ,err := json.Marshal(&people{Id: args[0],Name: args[1]})
	if err != nil{
		return shim.Error(err.Error())
	}
	err = stub.PutState(args[0],people)
	if err != nil{
		return shim.Error(err.Error())
	}
	fmt.Printf("Save Data: %s\n",people)
	return shim.Success(nil)
}

func get2(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 1{
		return shim.Error("incorrect arguments")
	}

	data,err := stub.GetState(args[0])
	if err != nil{
		return shim.Error(err.Error())
	}
	if err = json.Unmarshal(data,&people{}); err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}

func getByIdAndName(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	if len(args) != 2{
		return shim.Error("incorrect arguments")
	}
	indexName := "id~name"
	//p := people{Id: args[0],Name: args[1]}
	//idNameIndexKey,err := stub.CreateCompositeKey(indexName,[]string{p.Id,p.Name})
	//if err != nil{
	//	return shim.Error(err.Error())
	//}
	//err = stub.PutState(idNameIndexKey,[]byte{0x00})
	//if err != nil{
	//	return shim.Error(err.Error())
	//}
	idNameResultIterator,err := stub.GetStateByPartialCompositeKey(indexName,[]string{args[0]})
	if err != nil{
		return shim.Error(err.Error())
	}
	res := bytes.Buffer{}
	res.WriteString("[")
	defer idNameResultIterator.Close()
	for idNameResultIterator.HasNext(){
		idNameKey,err := idNameResultIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}
		objectType,compositeKey,err := stub.SplitCompositeKey(idNameKey.Key)
		returnId := compositeKey[0]
		returnName := compositeKey[1]
		fmt.Printf("objectType: %s,id: %s,name: %s\n",objectType,returnId,returnName)

		pState,err := stub.GetState(returnId)
		if err != nil{
			return shim.Error(err.Error())
		}
		fmt.Println(string(pState))
		res.Write(pState)
		if idNameResultIterator.HasNext() {
			res.WriteString(",")
		}
	}
	res.WriteString("]")
	return shim.Success(res.Bytes())
}

func getCreator(stub shim.ChaincodeStubInterface) peer.Response{
	creator,err :=  stub.GetCreator()
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success(creator)
}

func main(){
	if err := shim.Start(new(SimpleAsset));err != nil{
		fmt.Println(err)
	}
}