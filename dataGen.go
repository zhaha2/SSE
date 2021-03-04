package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//l为kw最大长度
func GetRandomK(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	//！！！！！！！
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
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
		result += strconv.Itoa(str[i]) + " "
	}
	result = result[:len(result)-1]
	return result
}

func main() {
	m := make(map[string]string)
	//生成20个kw的数据
	for i := 0; i < 20; i++ {
		m[GetRandomK(10)] = GetRandomV(15, 12)
		time.Sleep(100)
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
}
