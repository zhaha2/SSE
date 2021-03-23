package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"io"
	"io/ioutil"
	"sort"
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
			panic(err)
			return shim.Error(err.Error())
		}

		t := strconv.FormatFloat(float64(timeCost)*0.000000001, 'f', 4, 64)
		//将搜索总时间写入文件
		writeFileSC("Time_q.txt", t)
		fmt.Println("Query time: " + t + " s\n")
	} else if fn == "tquery" {
		fmt.Println("tquery")
		timeCost := int64(0)
		result, err, timeCost = tquery(stub, args)
		if err != nil {
			return shim.Error(err.Error())
		}

		t := strconv.FormatFloat(float64(timeCost)*0.000000001, 'f', 4, 64)
		//将搜索总时间写入文件
		writeFileSC("Time_q.txt", t)
		fmt.Println("Tquery time: " + t + " s\n")
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

//查询之后判断结果正确
func tquery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error, int64) {
	start := time.Now()
	var result []string
	var dddl []string
	var A []string
	var isResult bool
	var isRoot bool
	var parentkk string
	var err error
	var newK2 string
	isR := false

	if len(args) != 3 {
		return nil, fmt.Errorf("给定的参数个数不符合要求"), 0
	}

	isResult, isRoot, parentkk, dddl, err, newK2 = fstPreQ(stub, args)
	if err != nil {
		panic(err)
		return nil, err, 0
	}
	//println("ok2"+args[1])
	//更新k2
	args[1] = newK2
	//println("nk2"+args[1])

	//直到找到根节点
	for isRoot != true {
		for isResult != true {
			//println("not root not result")

			//获取数组A
			if A == nil || len(A) == 0 {
				A1, err := stub.GetState("A")
				if err != nil {
					panic(err)
					return nil, fmt.Errorf("获取数据发生错误"), 0
				}
				if A1 == nil {
					panic(err)
					return nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), 0
				}
				//还原成字符串数组
				A = strings.Split(string(A1), " ")
			}

			var temp []string

			//要先预处理dddl
			//println("odddl")
			//fmt.Println(dddl)
			//println()

			for _, i := range dddl {
				ii, _ := strconv.Atoi(i)
				//fmt.Println("i: "+i+" Ai: "+A[ii])
				//fmt.Println(ii)
				dd := F(args[1], A[ii])
				temp = append(temp, dd)
			}
			dddl = temp

			//println("ndddl")
			//fmt.Println(dddl)
			//println("____ndddl")
			//println()

			//判断是结果还是地址
			if string(dddl[len(dddl)-1][len(dddl[len(dddl)-1])-1]) == "+" {
				isResult = true
				//先处理dddl最后一位
				dddl[len(dddl)-1], parentkk = preDddl(dddl[len(dddl)-1])
				if len(parentkk) == 0 {
					isRoot = true
				}
				//fmt.Print("res: ")
				//fmt.Println(result)
			} else {
				isResult = false

				var temp []string

				//fmt.Println(dddl)
				//fmt.Println()
				//fmt.Println(dddl[len(dddl)-1])
				//println()

				//先处理dddl最后一位
				dddl[len(dddl)-1], parentkk = preDddl(dddl[len(dddl)-1])

				//fmt.Println(dddl[len(dddl)-1])
				//
				//println("pre dddl")
				//fmt.Println(dddl)
				//
				//println("new dddl")
				//println()
				for _, i := range dddl {
					dd := strings.Split(i, " ")
					for _, ddi := range dd {
						//println(ddi)
						temp = append(temp, ddi)
					}
				}
				dddl = temp
			}
		}
		//println("not root is result")

		//找到result

		for _, i := range dddl {
			fmt.Printf("\033[34m%s\033[0m", i+" ")
			result = append(result, i)
		}

		//fmt.Print("res: ")
		//fmt.Println(result)

		if isRoot {
			//这种方式退出的说明已经找完根节点内容了，下面不用重复找了
			isR = true
			break
		}

		kk := parentkk
		k1 := kk[:len(kk)/2]
		k2 := kk[len(kk)/2:]
		args[0] = k1
		args[1] = k2

		isResult, isRoot, parentkk, dddl, err = preQ(stub, args)
		if err != nil {
			panic(err)
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
				panic(err)
				return nil, fmt.Errorf("获取数据发生错误"), 0
			}
			if A1 == nil {
				panic(err)
				return nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), 0
			}
			//还原成字符串数组
			A = strings.Split(string(A1), " ")
		}

		var temp []string

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

			var temp []string

			//fmt.Println(dddl)
			//fmt.Println()
			//fmt.Println(dddl[len(dddl)-1])

			//先处理dddl最后一位
			dddl[len(dddl)-1], parentkk = preDddl(dddl[len(dddl)-1])

			//fmt.Println(dddl[len(dddl)-1])
			//
			//println("pre dddl")
			//fmt.Println(dddl)
			//
			//println("new dddl")
			for _, i := range dddl {
				dd := strings.Split(i, " ")
				for _, ddi := range dd {
					//println(ddi)
					temp = append(temp, ddi)
				}
			}
			dddl = temp
		}
	}

	//非isR的要再找一次
	if !isR {
		//找到result
		for _, i := range dddl {
			fmt.Printf("\033[34m%s\033[0m", i+" ")
			result = append(result, i)
		}
	}

	fmt.Printf("\033[32m%s\033[0m", "\n Query complete!\n")

	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(result)
	byteSlice := []byte(buf2.Bytes())
	end := time.Now()
	//输出总时间
	b := end.Sub(start)
	fmt.Printf("Query time cost: %s \n", b)

	//gai 比较结果正确性
	//从文件读出db
	//读取指定文件内容，返回的data是[]byte类型数据
	var db map[string]string
	var trst []string
	var itrst []int
	var irst []int
	sum := 0
	flag := true

	fdb, err := ioutil.ReadFile("db.txt")
	if err != nil {
		panic(err)
		fmt.Print(err)
	}

	//Unmarshal将data数据转换成指定的结构体类型,经过json转换后data中的数据已写入map中
	err = json.Unmarshal(fdb, &db)
	if err != nil {
		panic(err)
	}

	trst = strings.Split(db[args[2]], " ")
	for _, i := range trst {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		itrst = append(itrst, j)
	}
	for _, i := range result {
		if i == " " {
			continue
		}
		ii := strings.Split(i, " ")
		for _, iii := range ii {
			if iii == " " || len(iii) == 0 {
				continue
			}
			j, err := strconv.Atoi(iii)
			if err != nil {
				panic(err)
			}
			irst = append(irst, j)
			sum += 1
		}
	}
	sort.Ints(irst)

	fmt.Println(irst)
	fmt.Println(itrst)

	for i, j := range itrst {
		if j != irst[i] {
			flag = false
		}
	}
	if flag {
		//fmt.Printf("\033[32m%s%d\033[0m", "Num of ids: ", sum)
		//fmt.Println()
		fmt.Printf("\033[32m%t\033[0m", flag)
	} else {
		fmt.Printf("\033[41m%t\033[0m", flag)
		panic("Result wrong!!!!!!!!!!!")
	}
	fmt.Println()

	return byteSlice, nil, b.Nanoseconds()
}

func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error, int64) {
	start := time.Now()
	var result []string
	var dddl []string
	var A []string
	var isResult bool
	var isRoot bool
	var parentkk string
	var err error
	var newK2 string
	isR := false

	if len(args) != 2 {
		return nil, fmt.Errorf("给定的参数个数不符合要求"), 0
	}

	isResult, isRoot, parentkk, dddl, err, newK2 = fstPreQ(stub, args)
	if err != nil {
		panic(err)
		return nil, err, 0
	}

	args[1] = newK2

	//直到找到根节点
	for isRoot != true {
		for isResult != true {
			//获取数组A
			if A == nil || len(A) == 0 {
				A1, err := stub.GetState("A")
				if err != nil {
					panic(err)
					return nil, fmt.Errorf("获取数据发生错误"), 0
				}
				if A1 == nil {
					panic(err)
					return nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), 0
				}
				//还原成字符串数组
				A = strings.Split(string(A1), " ")
			}

			var temp []string

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
				if len(parentkk) == 0 {
					isRoot = true
				}
			} else {
				isResult = false

				var temp []string

				//先处理dddl最后一位
				dddl[len(dddl)-1], parentkk = preDddl(dddl[len(dddl)-1])

				for _, i := range dddl {
					dd := strings.Split(i, " ")
					for _, ddi := range dd {
						temp = append(temp, ddi)
					}
				}
				dddl = temp
			}
		}
		//找到result

		for _, i := range dddl {
			fmt.Printf("\033[34m%s\033[0m", i+" ")
			result = append(result, i)
		}

		if isRoot {
			//这种方式退出的说明已经找完根节点内容了，下面不用重复找了
			isR = true
			break
		}

		kk := parentkk
		k1 := kk[:len(kk)/2]
		k2 := kk[len(kk)/2:]
		args[0] = k1
		args[1] = k2

		isResult, isRoot, parentkk, dddl, err = preQ(stub, args)
		if err != nil {
			panic(err)
			return nil, err, 0
		}
	}

	//找到根节点
	for isResult != true {
		//获取数组A
		if A == nil {
			A1, err := stub.GetState("A")
			if err != nil {
				panic(err)
				return nil, fmt.Errorf("获取数据发生错误"), 0
			}
			if A1 == nil {
				panic(err)
				return nil, fmt.Errorf("根据 %s 没有获取到相应的数据", args[0]), 0
			}
			//还原成字符串数组
			A = strings.Split(string(A1), " ")
		}

		var temp []string

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

			var temp []string

			//先处理dddl最后一位
			dddl[len(dddl)-1], parentkk = preDddl(dddl[len(dddl)-1])

			for _, i := range dddl {
				dd := strings.Split(i, " ")
				for _, ddi := range dd {
					temp = append(temp, ddi)
				}
			}
			dddl = temp
		}
	}

	//非isR的要再找一次
	if !isR {
		//找到result
		for _, i := range dddl {
			fmt.Printf("\033[34m%s\033[0m", i+" ")
			result = append(result, i)
		}
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
	if len(pak) == 1 || len(pak[1]) == 0 {
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
