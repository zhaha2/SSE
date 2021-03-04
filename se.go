package main

import (
	"crypto/sha1"
	"crypto/sha256"
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
	t := sha256.New()
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
	var data string = "abcddas12"
	fmt.Printf("MD5 : %s\n", a(data))
	println("a")

	println(a(data))
	k1 := a(data)[:len(a(data))/2]
	println(k1)
	k2 := a(data)[len(a(data))/2:]
	println(k2)
	fmt.Printf("SHA1 : %s\n", b(data))
}
