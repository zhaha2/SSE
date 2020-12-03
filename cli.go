package main

import (
	"bufio"
	"cpabse"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//IndexGen
func set(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk, num int) {
	var policy, keyword, data string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter your policy:")
	scanner.Scan()
	policy = scanner.Text()
	fmt.Println("Please enter your keyword:")
	scanner.Scan()
	keyword = scanner.Text()
	fmt.Println("Please enter your data:")
	scanner.Scan()
	data = scanner.Text()

	//c为cph
	c, _ := cpabse.CP_Enc(pm, policy, msk, keyword)
	fmt.Println(c)
	//c1为字符串形式的c
	c1 := strconv.Itoa(int(c[0]))
	for i := 1; i < len(c); i++ {
		temp := strconv.Itoa(int(c[i]))
		c1 += " "
		c1 += temp
	}
	//ns为main()中的n，计数器（地址）
	ns := strconv.Itoa(num)

	//让chiancode执行操作
	comm := `peer chaincode invoke -n my -c '{"Args":["set","` + ns + `","` + c1 + `","` + string(data) + `"]}' -C myc`
	cmd := exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Printf("\033[32m%s\033[0m", "Upload Complete!\n")
}

//手动搜索， 即‘Q’操作
func queryManual(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk) {
	var attrs, keyword string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter you attrs:")
	scanner.Scan()
	attrs = scanner.Text()
	fmt.Println("Please enter your keyword:")
	scanner.Scan()
	keyword = scanner.Text()
	queryBase(pm, msk, attrs, keyword)

}

//基础的搜索功能
func queryBase(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk, attrs string, keyword string) {
	prv := cpabse.CP_Keygen(pm, msk, attrs)
	//tocken
	t, _ := cpabse.CP_TkEnc(prv, keyword, msk, pm)
	t1 := strconv.Itoa(int(t[0]))
	for i := 1; i < len(t); i++ {
		temp := strconv.Itoa(int(t[i]))
		t1 += " "
		t1 += temp
	}

	//分别搜索数据库的两部分
	k := "Key1"
	comm := `peer chaincode invoke -n my -c '{"Args":["query","` + t1 + `","` + k + `"]}' -C myc`
	cmd := exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	k = "Key2"
	comm = `peer chaincode invoke -n my -c '{"Args":["query","` + t1 + `","` + k + `"]}' -C myc`
	cmd = exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Printf("\033[32m%s\033[0m", "Query Complete!\n")

}

func shortQuery(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk, attrs string, keyword string) {
	prv := cpabse.CP_Keygen(pm, msk, attrs)
	//tocken
	t, _ := cpabse.CP_TkEnc(prv, keyword, msk, pm)
	t1 := strconv.Itoa(int(t[0]))
	for i := 1; i < len(t); i++ {
		temp := strconv.Itoa(int(t[i]))
		t1 += " "
		t1 += temp
	}

	k := "Key1"
	comm := `peer chaincode invoke -n my -c '{"Args":["shortquery","` + t1 + `","` + k + `"]}' -C myc`
	cmd := exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	fmt.Printf("\033[32m%s\033[0m", "Query Complete!\n")

}

func chainquery(Key string) {

	scanner := bufio.NewScanner(os.Stdin)
	if Key == "" {
		//Key为地址，即上传时的num，这里只能用这个搜
		fmt.Println("Please enter you Key:")
		scanner.Scan()
		Key = scanner.Text()
	}

	//Key字符转换为数字k
	k, _ := strconv.Atoi(Key)

	//这个循环为了模拟多次执行？，这里设置为不循环
	for i := k; i < k+1; i++ {
		ks := strconv.Itoa(i)
		//因为加密时地址为Key1，Key2...，这里也做成这个形式
		k := "Key" + ks

		comm := `peer chaincode invoke -n my -c '{"Args":["chainquery","` + k + `"]}' -C myc`
		cmd := exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		fmt.Printf("\033[32m%s\033[0m", "Chainquery Complete!\n")
	}
}

//一次搜索多个关键字
func batchq(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk) {

	//n为要返回多少结果
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the number of files you want to get:")
	scanner.Scan()
	n, err1 := strconv.Atoi(scanner.Text())
	if err1 != nil {
		fmt.Println("Input error!", err1)
		return
	}

	//搜索20次，为了实验取平均值
	for i := 0; i < 20; i++ {
		rand.Seed(int64(i))

		//随机搜索n/5个关键字，共返回n个文件
		for j := 0; j < (n / 5); j++ {
			kw := "kw" + strconv.Itoa(rand.Intn(1000))
			queryBase(pm, msk, "baf", kw)
		}

		fmt.Printf("\033[32mBatchQuery for %d times\n\033[0m", i+1)
	}
	fmt.Printf("\033[32m%s\033[0m", "BatchQuery Complete\n")
}

//一次搜索多个关键字
func batchcq(pm *cpabse.CpabePm, msk *cpabse.CpabeMsk) {

	//将kwStartIndex读入本地保存
	kwStartIndex, err := ioutil.ReadFile("kwStartIndex.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	kwStartIndexList := strings.Split(string(kwStartIndex), " ")

	//要返回多少结果
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the number of files you want to get:")
	scanner.Scan()
	n, err1 := strconv.Atoi(scanner.Text())
	if err1 != nil {
		fmt.Println("Input error!", err1)
		return
	}

	//搜索20次，为了实验取平均值
	for i := 0; i < 20; i++ {
		rand.Seed(int64(i))
		//随机搜索n/5个关键字，共返回n个文件
		for j := 0; j < (n / 10); j++ {
			a := rand.Intn(300)
			kw := "kw" + strconv.Itoa(a)

			shortQuery(pm, msk, "baf", kw)
			chainquery(kwStartIndexList[a])
		}
		fmt.Printf("\033[32mBatchChainquery for %d times\n\033[0m", i+1)
	}

	//搜索30次
	//for i := 0; i < 30; i++ {
	//	rand.Seed(int64(i))
	//	searchkey := ""
	//	//每次搜索n/5个关键字，共返回n个结果
	//	for j := 0; j < (n / 10); j++{
	//		searchkey = searchkey + kwStartIndexList[rand.Intn(1000)] + " "
	//	}
	//	comm := `peer chaincode invoke -n my -c '{"Args":["batchcq","` + string(searchkey) + `"]}' -C myc`
	//	cmd := exec.Command("/bin/sh", "-c", comm)
	//	cmd.Stdout = os.Stdout
	//	_ = cmd.Run()
	//
	fmt.Printf("\033[32m%s\033[0m", "BatchChainquery Complete\n")
	//}
}

func main() {
	//生成pm
	bpm := new(cpabse.BytePm)
	f, _ := os.Open("pm.txt")
	dec := gob.NewDecoder(f)
	_ = dec.Decode(&bpm)
	pm := new(cpabse.CpabePm)
	cpabse.Psetup(pm)
	cpabse.BpmToPm(pm, bpm)

	//生成msk
	bmsk := new(cpabse.ByteMsk)
	f1, _ := os.Open("msk.txt")
	dec1 := gob.NewDecoder(f1)
	_ = dec1.Decode(&bmsk)
	msk := new(cpabse.CpabeMsk)
	cpabse.BmskToMsk(msk, bmsk, pm)

	fmt.Println("Please choose the operation: upload/query/chainquery(U/Q/CQ/BQ/BCQ)")
	var c string
	fmt.Scanln(&c)
	n := 1 //n只有在 U 时才++，注意不要和plzidong的地址冲突
	if c == "U" || c == "u" {
		set(pm, msk, n)
		n++
	} else if c == "Q" || c == "q" {
		queryManual(pm, msk)
	} else if c == "CQ" || c == "cq" { //文件链查询
		chainquery("")
	} else if c == "BQ" || c == "bq" {
		batchq(pm, msk)
	} else if c == "BCQ" || c == "bcq" { //文件链查询
		batchcq(pm, msk)
	} else {
		fmt.Println("error input")
		return
	}
}
