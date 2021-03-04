package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//IndexGen
func upload() {

	//读入A
	data, err := ioutil.ReadFile("A.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of A:", string(data))

	//将字符串形式的A上传至区块链
	comm := `peer chaincode invoke -n my -c '{"Args":["uploadA","` + string(data) + `"]}' -C myc`
	cmd := exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	fmt.Printf("\033[32m%s\033[0m", "Array upload Complete!\n")

	//逐行读L
	file, err := os.Open("L.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		lineText := s.Text()
		if len(lineText) == 0 {
			continue
		}

		pak := strings.Split(lineText, ",")
		key := pak[0]
		val := pak[1]

		//将L中的每对kv分别上传至区块链
		comm = `peer chaincode invoke -n my -c '{"Args":["uploadL","` + key + `","` + val + `"]}' -C myc`
		cmd = exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		fmt.Printf("\033[32m%s key:%s val:%s\033[0m", "UploadL:", key, val+"\n")
	}

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	//逐行读Fst
	file, err = os.Open("Fst.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	s = bufio.NewScanner(file)
	for s.Scan() {
		lineText := s.Text()
		if len(lineText) == 0 {
			continue
		}

		pak := strings.Split(lineText, ",")
		key := pak[0]
		val := pak[1]

		//将L中的每对kv分别上传至区块链
		comm = `peer chaincode invoke -n my -c '{"Args":["uploadFst","` + key + `","` + val + `"]}' -C myc`
		cmd = exec.Command("/bin/sh", "-c", comm)
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
		fmt.Printf("\033[32m%s key:%s val:%s\033[0m", "UploadFst:", key, val+"\n")
	}

	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
}

//手动搜索， 一次搜索一个关键字
func queryManual() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter your keyword:")
	scanner.Scan()
	keyword := scanner.Text()

	queryBase(keyword)
}

//基础的搜索功能
func queryBase(keyword string) {
	K := "randstr123456"

	kk := H(K, keyword)
	k1 := kk[:len(kk)/2]
	k2 := kk[len(kk)/2:]
	comm := `peer chaincode invoke -n my -c '{"Args":["query","` + k1 + `","` + k2 + `"]}' -C myc`
	cmd := exec.Command("/bin/sh", "-c", comm)
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	fmt.Printf("\033[32m%s\033[0m", "Query Complete!\n")

}

//一次搜索多个关键字
func batchq() {

	//n为要返回多少结果
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the number of files you want to get:")
	scanner.Scan()
	n, err1 := strconv.Atoi(scanner.Text())
	if err1 != nil {
		fmt.Println("Input error!", err1)
		return
	}

	//读入所有关键字
	kw, err := ioutil.ReadFile("kws.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	kws := strings.Split(string(kw), " ")
	lenkws := len(kws)

	//搜索20次，为了实验取平均值
	for i := 0; i < 20; i++ {
		//gai成时间更好？
		rand.Seed(int64(i))

		//随机搜索n/5个关键字，共返回n个文件
		for j := 0; j < (n / 5); j++ {
			kw := kws[rand.Intn(lenkws)]
			queryBase(kw)
		}

		fmt.Printf("\033[32mBatchQuery for %d times\n\033[0m", i+1)
	}
	fmt.Printf("\033[32m%s\033[0m", "BatchQuery Complete\n")
}

func main() {

	fmt.Println("Please choose the operation: upload/query/batchquery(U/Q/BQ)")
	var op string
	fmt.Scanln(&op)

	if op == "U" || op == "u" {
		upload()
	} else if op == "Q" || op == "q" {
		queryManual()
	} else if op == "BQ" || op == "bq" {
		batchq()
	} else {
		fmt.Println("error input")
		return
	}
}
