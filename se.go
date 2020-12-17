package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
)

//func main() {
//	h := md5.New()
//	io.WriteString(h, "qwwasas")
//	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil))
//	fmt.Println(seed)
//	rand.Seed(int64(seed))
//	fmt.Println(rand.Int())
//}

//对字符串进行MD5哈希
func a(data string) string {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//对字符串进行SHA1哈希
func b(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
func main() {
	var data string = "abc"
	fmt.Printf("MD5 : %s\n", a(data))
	fmt.Printf("SHA1 : %s\n", b(data))
}
