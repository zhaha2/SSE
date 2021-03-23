package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//l为kw最大长度
func GetRandomK(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	//！！！！！！！
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	lenK := rand.Intn(l)
	if lenK == 0 {
		lenK = 1
	}
	for i := 0; i < lenK; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//nummax为id数字最大值，nmax为最多每个kw有几个id
//func  GetRandomV(nummax int, nmax int) string {
//	result := ""
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	n:=rand.Intn(nmax)
//	if n == 0 {n = 1}
//	for i := 0; i < n; i++ {
//		result += strconv.Itoa(r.Intn(nummax))+" "
//	}
//	return result
//}

//nummax为id数字最大值，nmax为最多每个kw有几个id
func GetRandomV(nummax int, nmax int) string {
	var str []int
	var sortStr []int
	for i := 0; i <= nummax; i++ {
		str = append(str, i)
	}

	//！！！！！！！
	result := ""
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(str), func(i, j int) { str[i], str[j] = str[j], str[i] })

	n := rand.Intn(nmax)
	if n == 0 {
		n = 1
	}
	for i := 0; i < n; i++ {
		sortStr = append(sortStr, str[i])
	}
	sort.Ints(sortStr)
	for i := 0; i < n; i++ {
		result += strconv.Itoa(sortStr[i]) + " "
	}

	//去除句尾空格
	result = strings.TrimSpace(result)
	return result
}

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

func main() {
	m := make(map[string]string)
	var key string
	keys := ""
	//生成20个kw的数据
	for i := 0; i < 100; i++ {
		key = GetRandomK(15) //l为kw最大长度
		//如果该key已存在，重来
		for len(m[key]) != 0 {
			key = GetRandomK(10)
		}
		keys += key + " "

		m[key] = GetRandomV(100, 60) //nummax为id数字最大值，nmax为最多每个kw有几个id
		time.Sleep(100)
	}

	//打印字典m
	for k, v := range m {
		println(k + " : " + v)
	}

	//写文件
	adress := "db.txt"
	f, err := os.OpenFile(adress, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := json.Marshal(m)
	if err != nil {
		println("数据序列化失败...")
	}

	_, err = fmt.Fprint(f, string(data))
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

	//去除字符串首位空格
	keys = strings.TrimSpace(keys)
	//写Keywords
	writeFile("Keywords.txt", keys)
}
