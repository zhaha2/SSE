package main

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"io"
	"strings"

	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"
)

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
	if fn == "uploadA" {
		fmt.Println("uploadA")
		result, err = uploadA(stub, args)
	} else if fn == "uploadL" {
		fmt.Println("uploadL")
		result, err = uploadL(stub, args)
		fmt.Printf("\033[32m%s\033[0m", "L upload complete!\n")
	} else if fn == "uploadFst" {
		fmt.Println("uploadFst")
		result, err = uploadFst(stub, args)
		fmt.Printf("\033[32m%s\033[0m", "Fst upload complete!\n")
	} else if fn == "query" {
		fmt.Println("query")
		timeCost := int64(0)
		result, err, timeCost = query(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}

		t := strconv.FormatFloat(float64(timeCost)*0.000000001, 'f', 4, 64)
		//将搜索总时间写入文件
		writeFileSC("Time_q.txt", t)
		fmt.Println("Query time: " + t + " s\n")
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
func writeFileSC(adress string, data string) {
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

//AES解密用
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

//解密
func F(key, enData string) (rst string) {
	bytEn, err := base64.StdEncoding.DecodeString(enData)
	if err != nil {
		panic(err)
	}

	AesKey := key //秘钥长度为16的倍数
	origin, err := AesDecrypt(bytEn, []byte(AesKey))
	if err != nil {
		panic(err)
	}
	return string(origin)
}

func H(k string, val string) (rst string) {
	t := sha256.New()
	io.WriteString(t, k+val)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//上传数组,以字符串形式存储
func uploadA(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return []byte(""), fmt.Errorf("Incorrect arguments. ")
	}

	err := stub.PutState("A", []byte(args[0]))
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to set asset: %s", args[0])
	}
	fmt.Printf("\033[32m%s\033[0m", "Array upload complete!\n")
	fmt.Println()
	return []byte(args[0]), nil
}

//上传字典，依次发送、存储每队键值对
func uploadL(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return []byte(""), fmt.Errorf("Incorrect arguments. ")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to set asset: %s", args[0])
	}
	fmt.Printf("\033[32m%s\033[0m\n", "key: "+args[0]+" data: "+args[1])
	fmt.Println()
	return []byte(args[0]), nil
}

//上传顶层索引，依次发送、存储每队键值对
func uploadFst(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return []byte(""), fmt.Errorf("Incorrect arguments. ")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return []byte(""), fmt.Errorf("Failed to set asset: %s", args[0])
	}
	fmt.Printf("\033[32m%s\033[0m\n", "key: "+args[0]+" data: "+args[1])
	fmt.Println()
	return []byte(args[0]), nil
}

func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error, int64) {
	start := time.Now()
	var result []string
	var dddl []string
	var A []string
	var temp []string
	var isResult bool
	var isRoot bool
	var parentkk string
	var err error
	var newK2 string

	if len(args) != 2 {
		return nil, fmt.Errorf("给定的参数个数不符合要求"), 0
	}

	isResult, isRoot, parentkk, dddl, err, newK2 = fstPreQ(stub, args)
	if err != nil {
		return nil, err, 0
	}
	//更新k2
	args[1] = newK2

	//直到找到根节点
	for isRoot != true {
		for isResult != true {
			//println("not root not result")

			//获取数组A
			if A == nil || len(A) == 0 {
				A1, err := stub.GetState("A")
				if err != nil {
					return nil, fmt.Errorf("获取数据发生错误"), 0
				}
				if A1 == nil {
					return nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), 0
				}
				//还原成字符串数组
				A = strings.Split(string(A1), " ")
			}

			for _, i := range dddl {
				ii, _ := strconv.Atoi(i)
				dd := F(args[1], A[ii])
				temp = append(temp, dd)
			}
			dddl = temp

			//判断是结果还是地址
			if string(dddl[len(dddl)-1][len(dddl[len(dddl)-1])-1]) == "+" {
				isResult = true
				//先处理dddl最后一位
				dddl[len(dddl)-1], parentkk = preDddl(dddl[len(dddl)-1])
			} else {
				isResult = false
			}
		}
		//println("not root is result")

		//找到result

		for _, i := range dddl {
			fmt.Printf("\033[34m%s\033[0m", i+" ")
			result = append(result, i)
		}

		kk := parentkk
		k1 := kk[:len(kk)/2]
		k2 := kk[len(kk)/2:]
		args[0] = k1
		args[1] = k2

		isResult, isRoot, parentkk, dddl, err = preQ(stub, args)
		if err != nil {
			return nil, err, 0
		}
	}

	//找到根节点
	for isResult != true {
		//println("is root not result")
		//获取数组A
		if A == nil {
			A1, err := stub.GetState("A")
			if err != nil {
				return nil, fmt.Errorf("获取数据发生错误"), 0
			}
			if A1 == nil {
				return nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), 0
			}
			//还原成字符串数组
			A = strings.Split(string(A1), " ")
		}

		for _, i := range dddl {
			ii, _ := strconv.Atoi(i)
			dd := F(args[1], A[ii])
			temp = append(temp, dd)
		}
		dddl = temp

		//判断是结果还是地址
		if string(dddl[len(dddl)-1][len(dddl[len(dddl)-1])-1]) == "+" {
			isResult = true
			//先处理dddl最后一位
			dddl[len(dddl)-1], _ = preDddl(dddl[len(dddl)-1])
		} else {
			isResult = false
		}
	}

	//找到result

	for _, i := range dddl {
		fmt.Printf("\033[34m%s\033[0m", i+" ")
		result = append(result, i)
	}
	fmt.Printf("\033[32m%s\033[0m", "\n Query complete!\n")

	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())
	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Printf("Query time cost: %s \n", b)
	return byteSlice, nil, b.Nanoseconds()
}

//处理dddl最后一位
func preDddl(dddlast string) (string, string) {
	var parentkk string
	dddlast = dddlast[:len(dddlast)-1]

	pak := strings.Split(dddlast, "||")
	if len(pak[1]) == 0 {
		dddlast = pak[0]
		parentkk = ""
	} else {
		dddlast = pak[0]
		parentkk = pak[1]
	}

	//处理pad
	dddlast = strings.Split(dddlast, "*")[0]

	return dddlast, parentkk
}

//第一次搜索，先搜Fst
func fstPreQ(stub shim.ChaincodeStubInterface, args []string) (bool, bool, string, []string, error, string) {
	var isResult bool
	var isRoot bool
	var parentkk string
	var fstStr string

	k1 := args[0]
	k2 := args[1]

	//先查Fst
	rst, err := stub.GetState(H(k1, "0"))
	if err != nil {
		return true, true, "", nil, fmt.Errorf("获取数据发生错误"), ""
	}
	if rst == nil {
		fmt.Printf("\nFst中根据 %s 没有获取到相应的数据\n", args[0])
		return true, true, "", nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), ""
	}

	//更新k2
	rstPak := strings.Split(F(k2, string(rst)), "|||")
	if len(rstPak) != 2 {
		return true, true, "", nil, fmt.Errorf("第一次获取数据发生错误"), ""
	} else {
		fstStr = rstPak[0]
		k2 = rstPak[1]
	}

	//再查L
	rst, err = stub.GetState(fstStr)
	if err != nil {
		return true, true, "", nil, fmt.Errorf("获取数据发生错误"), k2
	}
	if rst == nil {
		fmt.Printf("\nL中根据 %s 没有获取到相应的数据\n", args[0])
		return true, true, "", nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), k2
	}

	ddd := F(k2, string(rst))

	//判断是结果还是地址
	if string(ddd[len(ddd)-1]) == "+" {
		isResult = true
	} else {
		isResult = false
	}
	ddd = ddd[:len(ddd)-1]

	//获得父节点地址 gai 如果是叶子节点才要判断
	if isResult {
		pak := strings.Split(ddd, "||")
		if len(pak[1]) == 0 {
			ddd = pak[0]
			isRoot = true
		} else {
			ddd = pak[0]
			parentkk = pak[1]
			isRoot = false
		}
	}

	//处理pad
	ddd = strings.Split(ddd, "*")[0]

	//结果集
	dddl := strings.Split(ddd, " ")

	return isResult, isRoot, parentkk, dddl, nil, k2
}

func preQ(stub shim.ChaincodeStubInterface, args []string) (bool, bool, string, []string, error) {
	var isResult bool
	var isRoot bool
	var parentkk string

	k1 := args[0]
	k2 := args[1]

	rst, err := stub.GetState(H(k1, "1"))
	if err != nil {
		return true, true, "", nil, fmt.Errorf("获取数据发生错误")
	}
	if rst == nil {
		println("没有获取到相应的数据")
		return true, true, "", nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0])
	}

	ddd := F(k2, string(rst))

	//判断是结果还是地址
	if string(ddd[len(ddd)-1]) == "+" {
		isResult = true
	} else {
		isResult = false
	}
	ddd = ddd[:len(ddd)-1]

	//获得父节点地址 gai 如果是叶子节点才要判断
	if isResult {
		pak := strings.Split(ddd, "||")
		if len(pak[1]) == 0 {
			ddd = pak[0]
			isRoot = true
		} else {
			ddd = pak[0]
			parentkk = pak[1]
			isRoot = false
		}
	}

	//处理pad
	ddd = strings.Split(ddd, "*")[0]

	//结果集
	dddl := strings.Split(ddd, " ")

	return isResult, isRoot, parentkk, dddl, nil
}

func main() {
	err := shim.Start(new(MyChaincode))
	if err != nil {
		fmt.Printf("启动 PaymentChaincode 时发生错误: %s", err)
	}
}
