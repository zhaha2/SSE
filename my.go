package main

import (
	"bytes"
	"cpabse"
	"encoding/gob"
	"fmt"
	//"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nik-U/pbc"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var curveParams = "type a\n" +
	"q 87807107996633125224377819847540498158068831994142082" +
	"1102865339926647563088022295707862517942266222142315585" +
	"8769582317459277713367317481324925129998224791\n" +
	"h 12016012264891146079388821366740534204802954401251311" +
	"822919615131047207289359704531102844802183906537786776\n" +
	"r 730750818665451621361119245571504901405976559617\n" +
	"exp2 159\n" + "exp1 107\n" + "sign1 1\n" + "sign0 1\n"

type MyChaincode struct {
}

//实现Init函数，初始化
func (t *MyChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init...")
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

func (t *MyChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke")
	var result []byte
	var err error
	if fn == "set" {
		fmt.Println("set")
		result, err = upload(stub, args)
	} else if fn == "query" {
		fmt.Println("query")
		timeCost := int64(0)
		result, err, timeCost = query(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}

		t := strconv.FormatFloat(float64(timeCost)*0.000000001, 'f', 4, 64)
		//将搜索总时间写入文件
		writeFile("Time_q.txt", t)
		fmt.Println("Half query time: " + t + " s\n")

	} else if fn == "shortquery" {
		fmt.Println("shortquery")
		timeCost := int64(0)
		result, err, timeCost = shortquery(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}

		t := strconv.FormatFloat(float64(timeCost)*0.000000001, 'f', 4, 64)
		//将搜索总时间写入文件
		writeFile("Time_q.txt", t)
		fmt.Println("Half shortquery time: " + t + " s\n")

	} else if fn == "chainquery" {
		fmt.Println("chainquery")
		timeCost := int64(0)
		result, err, timeCost = chainquery(stub, args[0]) //arg[0]为地址（Key1, Key2, ...）

		t := strconv.FormatFloat(float64(timeCost)*0.000000001, 'f', 4, 64)
		//将搜索总时间写入文件
		writeFile("Time_q.txt", t)
		fmt.Println("Chainquery time: " + t + " s\n")
		//} else if fn == "batchcq" {
		//	fmt.Println("batchcq")
		//	var totalTime = int64(0)
		//	result, err, totalTime = batchcq(stub, args[0])         //arg[0]为地址串，以空格分开（Key1, Key2, ...）
		//	tt := strconv.FormatFloat(float64(totalTime) * 0.000001, 'f', 2,64)
		//	fmt.Println("Total time: " + tt +" ms")
		//
		//	//将搜索总时间写入文件
		//	writeFile("totalTime_cq.txt", tt)
	} else {
		fmt.Println("Input err!")
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success(result)
}

//写文件函数
func writeFile(adress string, data string) {
	f, err := os.OpenFile(adress, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprint(f, data+" ")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Write file complete!")

}

//上传
func upload(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 4 {
		return []byte(""), fmt.Errorf("Incorrect arguments. Expecting ID ,a key and a value")
	}

	//用传来的key加密，可作为划分数据库用
	var key, _ = stub.CreateCompositeKey(args[3], []string{args[0], args[1]})

	err := stub.PutState(key, []byte(args[2]))
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to set asset: %s", args[0])
	}
	fmt.Printf("\033[32m%s\033[0m\n", "id: "+args[0]+" data: "+args[2]+" key: "+args[3])
	fmt.Println()
	return []byte(args[0]), nil
}

//func query(stub shim.ChaincodeStubInterface, args []string) peer.Response {
//	// 检查传递的参数个数是否为1
//	if len(args) != 1{
//		return shim.Error("指定的参数错误，必须且只能指定相应的Key")
//	}
//
//	// 根据指定的Key调用GetState方法查询数据
//	result, err := stub.GetState(args[0])
//	if err != nil {
//		return shim.Error("根据指定的 " + args[0] + " 查询数据时发生错误")
//	}
//	if result == nil {
//		return shim.Error("根据指定的 " + args[0] + " 没有查询到相应的数据")
//	}
//
//	// 返回查询结果
//	return shim.Success(result)
//}

//普通查询
func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error, int64) {
	start := time.Now()
	var p *pbc.Pairing
	params := new(pbc.Params)
	params, _ = pbc.NewParamsFromString(curveParams)
	p = pbc.NewPairing(params)

	//token解密
	tk := cpabse.TkDec(args[0], p)

	//找到所有数据库中用该key加密的所有数据信息（即该数据库所有数据信息，一个数据库用1个key加密）
	queryResultsIterator, err := stub.GetStateByPartialCompositeKey(args[1], []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect"), 0
	}
	defer queryResultsIterator.Close()

	var result []string

	//对所有数据遍历
	for queryResultsIterator.HasNext() {

		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect"), 0
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect"), 0
		}

		cph := cpabse.CphDec(compositeKeyParts[1], p)

		//判断cph是否符合，是则输出，不是则继续，最终遍历完整个数据库
		if cpabse.Check(tk, cph, p) {
			result = append(result, string(responseRange.Value))
			fmt.Println(string(responseRange.Value))
		} else {
			continue
		}
	}
	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())
	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Printf("Query time cost: %s \n", b)
	return byteSlice, nil, b.Nanoseconds()
}

//普通查询
func shortquery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error, int64) {
	start := time.Now()
	var p *pbc.Pairing
	params := new(pbc.Params)
	params, _ = pbc.NewParamsFromString(curveParams)
	p = pbc.NewPairing(params)

	//token解密
	tk := cpabse.TkDec(args[0], p)

	//找到所有数据库中用该key加密的所有数据信息（即该数据库所有数据信息，一个数据库用1个key加密）
	queryResultsIterator, err := stub.GetStateByPartialCompositeKey(args[1], []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect"), 0
	}
	defer queryResultsIterator.Close()

	var result []string

	//对所有数据遍历
	for queryResultsIterator.HasNext() {

		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect"), 0
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect"), 0
		}

		cph := cpabse.CphDec(compositeKeyParts[1], p)

		//判断cph是否符合，是则输出，不是则继续，最终遍历完整个数据库
		if cpabse.Check(tk, cph, p) {
			result = append(result, string(responseRange.Value))
			fmt.Println(string(responseRange.Value))
			break
		} else {
			continue
		}
	}
	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())
	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Printf("Query time cost: %s \n", b)
	return byteSlice, nil, b.Nanoseconds()
}

//文件链查询
func chainquery(stub shim.ChaincodeStubInterface, Key string) ([]byte, error, int64) {
	//开始时间，一会计时用
	start := time.Now()
	fmt.Printf("\033[33m%s\033[0m\n", "Key: "+Key)

	//用符合键中一个键进行查询
	queryResultsIterator, err := stub.GetStateByPartialCompositeKey(Key, []string{})
	if err != nil {
		return []byte(""), fmt.Errorf("Incorrect"), 0
	}
	defer queryResultsIterator.Close()

	var result []string

	//返回结果
	for queryResultsIterator.HasNext() {

		responseRange, err := queryResultsIterator.Next()
		if err != nil {
			return []byte(""), fmt.Errorf("Incorrect"), 0
		}
		fmt.Println(responseRange.Value)
		//分离数据和地址
		Value := strings.Split(string(responseRange.Value), "::")
		//输出结果（数据）
		result = append(result, Value[0])
		fmt.Printf("\033[34mResult: %s\033[0m\n", Value[0])
		fmt.Println()
		if Value[0] == "data1" {
			break
		}
		key_next := Value[1]

		//如果没有地址则为链尾，停止搜索，否则继续搜索文件链
		for {
			fmt.Printf("\033[33m%s\033[0m\n", "Key: "+key_next)
			//用复合键中一个键进行查询，这里用地址Key查询
			queryResultsIterator1, err := stub.GetStateByPartialCompositeKey(key_next, []string{})
			if err != nil {
				return []byte(""), fmt.Errorf("Incorrect"), 0
			}
			defer queryResultsIterator1.Close()

			responseRange, err = queryResultsIterator1.Next()
			if err != nil {
				return []byte(""), fmt.Errorf("Incorrect"), 0
			}
			fmt.Println(responseRange.Value)
			//分离数据和地址
			Value := strings.Split(string(responseRange.Value), "::")
			//输出结果（数据）
			result = append(result, Value[0])

			//找到链尾 搜索结束
			if Value[0] == "data1" {
				fmt.Printf("\033[34mResult: %s\033[0m\n", Value[0])
				fmt.Println()
				break
			}

			//否则，继续按地址搜索
			key_next = Value[1]
			fmt.Printf("\033[34mResult: %s\033[0m\n", Value[0])
			fmt.Println()
		}

		//if len(Value) != 2 {
		//      break
		//} else {
		//      var r []string
		//      res, _ := chainquery(stub, Value[1])
		//      r = append(r, strconv.Itoa(int(res[0])))
		//      for i := 1; i < len(res); i++ {
		//              r = append(r, strconv.Itoa(int(res[i])))
		//      }
		//      result = append(result, r...)
		//}
	}
	queryResultsIterator.Close()
	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())

	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Printf("Chainquery time cost: %s \n", b)
	return byteSlice, nil, b.Nanoseconds()
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("error start MyChaincode")
	}
}
