package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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
	str := ""

	db, err := ioutil.ReadFile("db.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(db, &m)
	if err != nil {
		println("数据读取失败...")
		return
	}

	for k := range m {
		str += k + " "
	}

	//去除字符串首位空格
	str = strings.TrimSpace(str)

	println(len(m))
	//打印字典m
	for k, v := range m {
		println(k + " : " + v)
	}
	writeFile("Keywordss.txt", str)
}
