package main

import (
	"fmt"
)

func main() {
	err := fmt.Errorf("获取数据发生错误")
	println(err)
	println(err == fmt.Errorf("获取数据发生错误"))
}
