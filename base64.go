package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"strings"
)

const (
	//BASE64字符表,不要有重复
	base64Table = "IJjkKLMNO567PQX12RVW3YZaDEFGbcdefghiABCHlSTUmnopqrxyz04stuvw89+/"
	//base64Table        = "<>:;',./?~!@#$CDVWX%^&*ABYZabcghijklmnopqrstuvwxyz01EFGHIJKLMNOPQRSTU2345678(def)_+|{}[]9/"
	hashFunctionHeader = "zh.ife.iya"
	hashFunctionFooter = "09.O25.O20.78"
)

var coder = base64.NewEncoding(base64Table)

/**
 * base64加密
 */
func Base64Encode(str string) string {
	var src []byte = []byte(hashFunctionHeader + str + hashFunctionFooter)
	return string([]byte(coder.EncodeToString(src)))
}

/**
 * base64解密
 */
func Base64Decode(str string) (string, error) {
	var src []byte = []byte(str)
	by, err := coder.DecodeString(string(src))
	return strings.Replace(strings.Replace(string(by), hashFunctionHeader, "", -1), hashFunctionFooter, "", -1), err
}

func main() {

	a := "asss"
	b := "22"

	c := Base64Encode(a)
	println(c)
	d, _ := Base64Decode(c)
	println(d)

	c = Base64Encode(b)
	println(c)
	d, _ = Base64Decode(c)
	println(d)

	buf2 := &bytes.Buffer{}
	gob.NewEncoder(buf2).Encode(a)
	println()
	println(buf2.Bytes())
}
