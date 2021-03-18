package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	childs []*Node
	fs     string
	ss     string
	val    []string
	parent *Node
}

//森林生成之后生成索引
func forest_IndexGen() {
	//森林生成
	var db map[string]string

	//从文件读出db
	//读取指定文件内容，返回的data是[]byte类型数据
	b, err := ioutil.ReadFile("db.txt")
	if err != nil {
		fmt.Print(err)
	}

	//Unmarshal将data数据转换成指定的结构体类型,经过json转换后data中的数据已写入map中
	err = json.Unmarshal(b, &db)
	if err != nil {
		panic(err)
	}

	////对每个文件对应的ids也按大小排序,否则每次输出结果不一样,!!gai应该在生成阶段做
	//sortIds(db)

	//先排序，大的在前先合并
	var newDbKeys = sortDb(db)

	//定义变量
	//gai自定
	m := 1
	var s []*Node
	var roots []*Node
	//var delete []*Node
	var dbw string
	var p *Node

	for _, w := range newDbKeys {
		println("_______________________________")
		p = nil
		max := 0
		flag := 0
		s = []*Node{}

		//若是第一次
		if len(roots) == 0 {
			t := new(Node)
			dbw = db[w]
			t.fs = dbw
			t.ss = dbw
			t.val = append(t.val, w)
			t.childs = nil

			delete(db, w)

			roots = append(roots, t)
			continue
		}

		//所有根节点加入栈中
		for _, r := range roots {
			s = append(s, r)
		}

		dbw = db[w]

		for len(s) != 0 {
			if len(s) <= 0 {
				break
			}
			n := s[len(s)-1]
			s = s[:len(s)-1]

			//取交集
			a := strings.Split(n.fs, " ")
			b := strings.Split(dbw, " ")
			is := intersect(a, b)

			//若当前合并节点的父节点的全集不是新结果集子集，不进行合并
			if n.parent != nil {
				c := strings.Split(n.parent.fs, " ")
				d := strings.Split(dbw, " ")
				//他的孩子节点也不判断了, flag为0
				if len(intersect(c, d)) != len(c) {
					continue
				}
			}

			//判断是哪种合并情况
			if len(is) == len(a) {
				if len(is) == len(b) {
					p = n
					flag = 1
					break
				} else if len(is) > max {
					p = n
					max = len(is)
					flag = 2
				}
			} else if len(is) == len(b) && len(is) > max {
				p = n
				max = len(is)
				flag = 3
			} else if len(is) > max && len(is) > m {
				p = n
				max = len(is)
				flag = 4
			}

			print("n.v: ")
			for _, v := range n.val {
				print(v + " ")
			}
			print("w: " + w)
			print("max: " + strconv.Itoa(max))

			println(" flag: " + strconv.Itoa(flag))

			//之后，将n所有的孩子节点也放入栈中
			//gai如果某种情况就不用加了
			for _, ch := range n.childs {
				s = append(s, ch)
			}
		}

		//从db中删除这个关键字和合并的关键字
		delete(db, w)

		//根据情况构造树
		if flag == 0 {
			t := new(Node)
			t.fs = dbw
			t.ss = dbw
			t.val = append(t.val, w)
			t.childs = nil
			t.parent = nil
			//t.val = append(t.val, "f0")

			roots = append(roots, t)
		} else if flag == 1 {
			//
			p.val = append(p.val, w)
			//p.val = append(p.val, "f1")
		} else if flag == 2 {
			t := new(Node)
			t.fs = dbw
			//取差集
			a := strings.Split(dbw, " ")
			b := strings.Split(p.fs, " ")
			t.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			t.val = append(t.val, w)
			t.childs = nil
			t.parent = p

			//t.val = append(t.val, "f2")
			//p.val = append(p.val, "f2")

			p.childs = append(p.childs, t)
		} else if flag == 3 {
			t := new(Node)
			t.fs = p.fs
			b := strings.Split(dbw, " ")
			a := strings.Split(p.fs, " ")
			t.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			t.val = p.val
			//t.val = append(t.val, "f3")
			t.parent = p
			t.childs = p.childs

			p.fs = dbw
			p.ss = p.fs
			p.val = nil
			p.val = append(p.val, w)
			p.childs = nil
			p.childs = append(p.childs, t)

			//roots = append(roots, t)
			//delete = append(delete, p)
		} else if flag == 4 {
			//t := new(Node)
			//t.fs = p.fs
			//t.ss = p.ss
			//t.val = p.val
			//t.childs = p.childs
			//t.parent = p
			//
			////取交集
			//a := strings.Split(p.fs, " ")
			//b := strings.Split(dbw, " ")
			//
			////数组转字符串
			////strings.Replace(strings.Trim(fmt.Sprint(Intersect(a, b)), "[]"), " ", ",", -1)
			//p.fs = strings.Trim(fmt.Sprint(intersect(a, b)), "[]")
			//p.ss = p.fs
			//p.val = nil
			//p.childs = nil
			//
			//u := new(Node)
			//u.fs = dbw
			////取差集
			//a = strings.Split(dbw, " ")
			//b = strings.Split(p.fs, " ")
			//u.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			//u.val = append(u.val, w)
			//u.childs = nil
			//u.parent = p
			//
			//a = strings.Split(t.fs, " ")
			//b = strings.Split(p.fs, " ")
			//t.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			//
			//p.childs = append(t.childs, u)
			//p.childs = append(p.childs, t)

			print("最终 n.v: ")
			for _, v := range p.val {
				print(v + " ")
			}
			print("w: " + w)
			print("max: " + strconv.Itoa(max))

			println(" flag: " + strconv.Itoa(flag))

			t := new(Node)
			t.parent = p.parent
			if p.parent == nil {
				roots = append(roots, t)
				//从roots删除p
				for i, n := range roots {
					if n == p {
						roots = append(roots[:i], roots[i+1:]...)
						break
					}
				}
			} else {
				//从p.parent.child删除p加入t
				for i, n := range p.parent.childs {
					if n == p {
						p.parent.childs = append(p.parent.childs[:i], p.parent.childs[i+1:]...)
						break
					}
				}
				p.parent.childs = append(p.parent.childs, t)
			}

			//取交集
			a := strings.Split(p.fs, " ")
			b := strings.Split(dbw, " ")

			//数组转字符串
			//strings.Replace(strings.Trim(fmt.Sprint(Intersect(a, b)), "[]"), " ", ",", -1)
			t.fs = strings.Trim(fmt.Sprint(intersect(a, b)), "[]")

			//取差集
			if t.parent == nil {
				t.ss = t.fs
			} else {
				a = strings.Split(t.fs, " ")
				b = strings.Split(t.parent.fs, " ")
				t.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			}

			t.val = nil
			t.childs = nil

			u := new(Node)
			u.fs = dbw
			//取差集
			a = strings.Split(dbw, " ")
			b = strings.Split(t.fs, " ")
			u.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			u.val = append(u.val, w)
			u.childs = nil
			u.parent = t

			a = strings.Split(p.fs, " ")
			b = strings.Split(t.fs, " ")
			p.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			p.parent = t

			//t.val = append(t.val, "f4")
			//u.val = append(u.val, "f4")
			//p.val = append(p.val, "f4")

			t.childs = append(t.childs, u)
			t.childs = append(t.childs, p)
		}

		println("\n_______________________________")
	}

	//输出
	for _, t := range roots {
		bfs(t)
		writeFile("bfsTree.txt", "\n")
		writeFile("bfsTree.txt", "\n")
	}

	//索引生成
	indexGen(roots)
}

//func sortIds(db map[string]string) {
//	for k, v := range db {
//		intIds := []int{}
//		ids := strings.Split(v, " ")
//		for _, i := range ids {
//			num, _ := strconv.Atoi(i)
//			intIds = append(intIds, num)
//		}
//		sort.Ints(intIds)
//
//		sortIds := ""
//		for _, i := range intIds {
//			sortIds += strconv.Itoa(i) + " "
//		}
//		db[k] = sortIds
//	}
//}

//先排序，大的在前先合并
func sortDb(db map[string]string) []string {
	var numDb = make([]int, 0)
	var newDbKey = make([]string, 0)

	for oldk, v := range db {
		newDbKey = append(newDbKey, oldk)
		//numdb为db中各k对应文件id个数
		numDb = append(numDb, len(strings.Split(v, " ")))
	}

	//冒泡排序
	var isDone = false
	for !isDone {
		isDone = true
		var i = 0
		for i < len(numDb)-1 {
			if numDb[i] < numDb[i+1] {
				var temp = numDb[i]
				numDb[i] = numDb[i+1]
				numDb[i+1] = temp

				//key也对应交换
				var tem = newDbKey[i]
				newDbKey[i] = newDbKey[i+1]
				newDbKey[i+1] = tem

				isDone = false
			}
			i++
		}
	}

	return newDbKey
}

//取交集
func intersect(a []string, b []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range a {
		m[v]++
	}

	for _, v := range b {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

//求差集 slice1-slice2并集
func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

var K = "randstr123456"

//预定义A大小
var lenA = 1000

//设为全局变量来递归
var A [1000]string

//索引生成
func indexGen(roots []*Node) {
	//初始化
	L := make(map[string]string)
	top := len(roots) - 1
	//实验确定
	b := 20
	B := 10

	//gai和L一起存就行了 不用多存一次第一个关键字, 有没有必要加,会不会泄露信息
	//顶层索引
	Fst := make(map[string]string)

	//出于性能考虑，可以预定义map大小
	//记录数组A的空位
	emp := make(map[int]int)
	for i := 0; i < lenA; i++ {
		emp[i] = i
	}

	for top >= 0 {
		n := roots[top]
		top--

		//处理每棵树
		nodeProcess(n, b, B, L, emp, "", Fst)

	}

	//把A中剩余空位填满
	empad := ""
	num := 0
	for num < B {
		empad += "1"
		num++
	}
	kk := H(K, "empad")
	k2 := kk[len(kk)/2:]
	d := F(k2, empad)

	for _, em := range emp {
		A[em] = d
	}

	//结果写入文件
	Astr := ""
	for _, v := range A {
		Astr += v + " "
	}
	writeFile("A.txt", Astr)
	//gai用json?
	writeFileMap("Fst.txt", Fst)
	writeFileMap("L.txt", L)
}

//递归处理一棵树
func nodeProcess(n *Node, b int, B int, L map[string]string,
	emp map[int]int, parentkk string, Fst map[string]string) {

	var nkw string
	var kk string
	var k1 string
	var k2 string
	var l string
	var fl string

	db := strings.Split(n.ss, " ")

	//判断当前节点有几个kw,计算kk
	if len(n.val) == 0 {

		//kw为空的话随机分配一个值
		rand.Seed(time.Now().Unix())
		now := time.Now()            //获取当前时间
		timestamp2 := now.UnixNano() //纳秒时间戳
		tsstr := strconv.Itoa(int(timestamp2))
		nkw = "emp/" + tsstr + strconv.Itoa(rand.Intn(99999999))

		kk = H(K, nkw)
		k1 = kk[:len(kk)/2]
		k2 = kk[len(kk)/2:]

		l = H(k1, "1")
	} else if len(n.val) == 1 {
		nkw = n.val[0]

		kk = H(K, nkw)
		k1 = kk[:len(kk)/2]
		k2 = kk[len(kk)/2:]

		l = H(k1, "1")

		fl = H(k1, "0")
		Fst[fl] = F(k2, l+"|||"+k2)
	} else {
		nkw = n.val[0]

		kk = H(K, nkw)
		k1 = kk[:len(kk)/2]
		k2 = kk[len(kk)/2:]

		l = H(k1, "1")

		fl = H(k1, "0")
		Fst[fl] = F(k2, l+"|||"+k2)

		//所有相同结果集的关键字都和第一个关键字地址一样
		for _, val := range n.val[1:] {
			bkk := H(K, val)
			bk1 := bkk[:len(bkk)/2]
			bk2 := bkk[len(bkk)/2:]

			bfl := H(bk1, "0")
			Fst[bfl] = F(bk2, l+"|||"+k2)
		}
	}

	//不sleep会乱序
	time.Sleep(100)

	if len(db) < b {

		str := partition2(len(db), b, db)

		//父节点地址
		strr := ""
		if parentkk != "" {
			strr = parentkk
		}
		str += "||" + strr

		//标明是结果
		str += "+"

		d := F(k2, str)
		L[l] = d
	} else if len(db) > b && len(db) < b*B {

		buf := partition(len(db), B, db)

		//父节点地址
		strr := ""
		if parentkk != "" {
			strr = parentkk
		}
		buf[len(buf)-1] += "||" + strr

		//标明是结果
		buf[len(buf)-1] += "+"

		//数组随机位置存储
		var arr []int
		for k := range emp {
			arr = append(arr, k)
		}
		ii, _ := Random(arr, len(arr))

		//将buf中1，2，3，4...位置的ids存到A中随机空位
		for j, v := range ii {
			A[v] = F(k2, buf[j])
		}

		//将存过数据的位置从emp中删除
		for _, v := range ii {
			delete(emp, v)
		}

		//把ii打包分块
		var iistr []string
		for i := 0; i < len(ii); i++ {
			iistr = append(iistr, strconv.Itoa(ii[i])+" ")
		}
		rst := partition2(len(iistr), b, iistr)

		//标明是地址
		rst += "-"

		d := F(k2, rst)
		L[l] = d
	} else {

		buf := partition(len(db), B, db)

		//父节点地址
		strr := ""
		if parentkk != "" {
			strr = parentkk
		}
		buf[len(buf)-1] += "||" + strr

		//标明是结果
		buf[len(buf)-1] += "+"

		//数组随机位置存储
		var arr []int
		for k := range emp {
			arr = append(arr, k)
		}
		ii, _ := Random(arr, len(arr))

		//将buf中1，2，3，4...位置的ids存到A中随机空位
		for j, v := range ii {
			A[v] = F(k2, buf[j])
		}

		//将存过数据的位置从emp中删除
		for _, v := range ii {
			delete(emp, v)
		}

		//把ii打包分块
		var iistr []string
		for i := 0; i < len(ii); i++ {
			iistr = append(iistr, strconv.Itoa(ii[i]))
		}

		//再次打包
		buf2 := partition(len(iistr), b, iistr)

		//标明是地址
		buf[len(buf)-1] += "-"

		arr = []int{}
		for k := range emp {
			arr = append(arr, k)
		}
		ii, _ = Random(arr, len(arr))

		//将buf2中1，2，3，4...位置的ids存到A中随机空位
		for j, v := range ii {
			A[v] = F(k2, buf2[j])
		}

		//将存过数据的位置从emp中删除
		for _, v := range ii {
			delete(emp, v)
		}

		//把ii打包分块
		iistr = []string{}
		for i := 0; i < len(ii); i++ {
			iistr = append(iistr, strconv.Itoa(ii[i]))
		}
		rst := partition2(len(iistr), b, iistr)

		//标明是地址
		rst += "-"

		d := F(k2, rst)
		L[l] = d
	}

	if len(n.childs) != 0 {
		//不sleep会乱序
		time.Sleep(100)
		//处理所有孩子节点
		for _, c := range n.childs {
			nodeProcess(c, b, B, L, emp, kk, Fst)
		}
		//不sleep会乱序
		time.Sleep(100)
	} else {
		//不sleep会乱序
		time.Sleep(100)
		return
	}
}

//分块 返回字符串数组
func partition(len int, b int, db []string) []string {
	var rst []string
	p := 0
	for p < len {
		str := ""
		q := p + b

		//分块
		if q < len {
			for ; p < q; p++ {
				//加分隔符
				str += " " + db[p]
			}
		} else {
			for ; p < len; p++ {
				str += " " + db[p]
			}
			//Pad
			//填充什么
			str += "*"
			for ; p < q; p++ {
				str += " " + strconv.Itoa(rand.Intn(999))
			}
		}
		rst = append(rst, str)
	}
	return rst
}

//分块 返回字符串
func partition2(len int, b int, db []string) string {
	rst := ""
	p := 0
	for p < len {
		str := ""
		q := p + b

		//分块
		if q < len {
			for ; p < q; p++ {
				str += " " + db[p]
			}
		} else {
			for ; p < len; p++ {
				str += " " + db[p]
			}
			//Pad(怎么区分)
			//填充什么
			str += "*"
			for ; p < q; p++ {
				str += " " + strconv.Itoa(rand.Intn(999))
			}
		}

		//不同块之间用"，"隔开
		rst += str + ","
	}
	return rst
}

//随机打乱数组
func Random(ints []int, length int) ([]int, error) {
	rand.Seed(time.Now().Unix())
	if len(ints) <= 0 {
		return nil, errors.New("the length of the parameter strings should not be less than 0")
	}

	if length <= 0 || len(ints) < length {
		return nil, errors.New("the size of the parameter length illegal")
	}

	for i := len(ints) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		ints[i], ints[num] = ints[num], ints[i]
	}

	var rst []int
	for i := 0; i < length; i++ {
		rst = append(rst, ints[i])
	}
	return ints, nil
}

//aes加密用
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//加密
func F(key string, origData string) (rst string) {
	text := origData
	AesKey := key //秘钥长度为16的倍数
	encrypted, err := AesEncrypt([]byte(text), []byte(AesKey))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

func H(k string, val string) (rst string) {
	t := sha256.New()
	io.WriteString(t, k+val)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//写文件函数,写入字符串
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

//写文件函数,写入字典
func writeFileMap(adress string, m map[string]string) {
	f, err := os.OpenFile(adress, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range m {
		_, err = fmt.Fprint(f, k+","+v+"\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

	fmt.Println("Write file complete!")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

//读文件函数,逐行读取返回字典
func readFile(adress string) map[string]string {

	rst := make(map[string]string)

	file, err := os.Open(adress)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		if len(lineText) == 0 {
			continue
		}

		data := strings.Split(lineText, ",")
		rst[data[0]] = data[1]
	}

	fmt.Println("Read file complete!")
	return rst
}

//dfs
func dfs(t *Node) {
	if t == nil {
		return
	}
	writeFile("dfsTree.txt", "fs:"+t.fs+"    ss:"+t.ss+"\n")

	for _, n := range t.childs {
		dfs(n)
	}
	writeFile("dfsTree.txt", "\n")
}

////dfs
//func bfs(t *Node) {
//
//	if t == nil {
//		return
//	}
//
//	que := []*Node{t}
//
//	hei := 1
//	m := make(map[int][]*Node)
//	m[hei] = []*Node{t}
//
//	for len(que) != 0 {
//		if len(que) == 0 {break}
//		t = que[0]
//		que = que[1:]
//
//		//if t == nil {
//		//	writeFile("dfsTree.txt", " ||***"+"\n")
//		//	continue
//		//}
//
//
//		hei += 1
//		for _, c := range t.childs {
//			que = append(que, c)
//
//			m[hei] = append(m[hei], c)
//		}
//	}
//
//	for i := 1; i <= len(m); i++ {
//		for _, t = range m[i]{
//
//			writeFile("bfsTree.txt", "    fs:"+t.fs+"    ss:"+t.ss+" || ")
//		}
//		writeFile("bfsTree.txt", "\n")
//	}
//}

//bfs对新结构体排序
type sortNode struct {
	n   *Node
	hei int
}

type snlist []*sortNode

//为结构体排序重写方法
func (sl snlist) Len() int {
	return len(snlist{})
}
func (sl snlist) Less(i, j int) bool {
	return sl[i].hei < sl[j].hei
}
func (sl snlist) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

func initSn(n *Node, hei int) *sortNode {
	sn := new(sortNode)
	sn.n = n
	sn.hei = hei

	return sn
}

//bfs
func bfs(t *Node) {

	if t == nil {
		return
	}

	hei := 1
	sn := initSn(t, hei)
	sortlist := []*sortNode{sn}
	slist := []*sortNode{sn}

	for len(sortlist) != 0 {
		if len(sortlist) == 0 {
			break
		}
		sn = sortlist[0]
		sortlist = sortlist[1:]

		hei = sn.hei + 1
		for _, c := range sn.n.childs {
			sn = initSn(c, hei)
			sortlist = append(sortlist, sn)
			slist = append(slist, sn)
		}
	}

	sort.Sort(snlist(slist))

	bf := 1
	af := bf
	for _, sn = range slist {
		nkww := ""
		bf = af
		af = sn.hei
		if af > bf {
			writeFile("bfsTree.txt", "\n")
		}

		for _, w := range sn.n.val {
			nkww += w + " "
		}
		writeFile("bfsTree.txt", "    fs:"+sn.n.fs+"    ss:"+sn.n.ss+"    kw:"+nkww+" || ")

	}
}

func main() {
	forest_IndexGen()
}
