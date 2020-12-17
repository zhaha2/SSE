package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	childs []*Node
	fs     string
	ss     string
	val    []string
}

//森林生成
func forestGen() {
	var db map[string]string

	//从文件读出db
	db = map[string]string{"one": "1", "two": "2"}

	//之后实现
	//sort(db)

	//定义变量
	m := 10
	var s []*Node
	top := -1
	var roots []*Node
	//var delete []*Node
	var dbw string
	var p *Node

	for w := range db {
		p = nil
		max := 0
		flag := 0
		s = []*Node{}

		//若是第一次
		if len(roots) == 0 {
			t := new(Node)
			t.fs = dbw
			t.ss = dbw
			t.val = append(t.val, w)
			t.childs = nil

			roots = append(roots, t)
			break
		}

		//所有根节点加入栈中
		for _, r := range roots {
			s = append(s, r)
			top++
		}

		dbw = db[w]

		for len(s) != 0 {
			top--
			n := s[top]

			//取交集
			a := strings.Split(n.fs, " ")
			b := strings.Split(dbw, " ")
			is := intersect(a, b)

			//判断是哪种合并情况
			if len(is) == len(a) {
				if len(is) == len(b) {
					p = n
					flag = 1
				} else if len(is) > max {
					p = n
					max = len(is)
					flag = 2
				}
			} else if len(is) == len(b) {
				p = n
				max = len(is)
				flag = 3
			} else if len(is) > max && len(is) > m {
				p = n
				max = len(is)
				flag = 4
			}

			//之后，将n所有的孩子节点也放入栈中
			for _, c := range n.childs {
				s = append(s, c)
				top++
			}
		}
		print(p.fs, flag)

		//根据情况构造树
		if flag == 0 {
			t := new(Node)
			t.fs = dbw
			t.ss = dbw
			t.val = append(t.val, w)
			t.childs = nil

			roots = append(roots, t)
		} else if flag == 1 {
			//
			p.val = append(p.val, w)
		} else if flag == 2 {
			t := new(Node)
			t.fs = dbw
			//取差集
			a := strings.Split(dbw, " ")
			b := strings.Split(p.fs, " ")
			t.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			t.val = append(t.val, w)
			t.childs = nil

			p.childs = append(p.childs, t)
		} else if flag == 3 {
			t := new(Node)
			t.fs = p.fs
			t.ss = p.ss
			t.val = p.val
			t.childs = p.childs

			p.fs = dbw
			//
			p.ss = p.fs
			p.val = nil
			p.val = append(p.val, w)
			p.childs = nil
			p.childs = append(p.childs, t)

			//roots = append(roots, t)
			//delete = append(delete, p)
		} else if flag == 4 {
			t := new(Node)
			t.fs = p.fs
			t.ss = p.ss
			t.val = p.val
			t.childs = p.childs

			//取交集
			a := strings.Split(p.fs, " ")
			b := strings.Split(dbw, " ")

			//数组转字符串
			//strings.Replace(strings.Trim(fmt.Sprint(Intersect(a, b)), "[]"), " ", ",", -1)
			p.fs = strings.Trim(fmt.Sprint(intersect(a, b)), "[]")
			p.ss = p.fs
			p.val = nil
			p.childs = nil

			u := new(Node)
			u.fs = dbw
			//取差集
			a = strings.Split(dbw, " ")
			b = strings.Split(p.fs, " ")
			u.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			u.val = append(u.val, w)
			u.childs = nil

			a = strings.Split(t.fs, " ")
			b = strings.Split(p.fs, " ")
			t.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")

			p.childs = append(t.childs, u)
			p.childs = append(t.childs, t)
		}
	}

	//输出
	for _, t := range roots {
		dfs(t)
	}
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

//求差集 slice1-并集
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

//索引生成
func indexGen(roots []*Node) {
	//初始化
	K := 1
	var L map[string]string
	top := len(roots) - 1
	//实验确定
	b := 20
	B := 10

	//预定义A大小
	lenA := 1000
	var A [1000]string

	//出于性能考虑，可以预定义map大小
	emp := make(map[int]int)
	for i := 0; i < lenA; i++ {
		emp[i] = i
	}

	for top >= 0 {
		n := roots[top]
		top--

		//字符串怎么转int, 要多长
		k1, k2 := F(K, n.val)
		db := strings.Split(n.ss, " ")

		if len(n.ss) < b {
			//如果是空值的怎么办
			str := partition2(len(n.ss), b, db)

			d := F(k2, str)
			l := H(k1, 1)
			L[l] = d
		} else if len(n.ss) > b && len(n.ss) < b*B {
			buf := partition(len(n.ss), B, db)

			//数组随机位置存储
			var arr []int
			for k := range emp {
				arr = append(arr, k)
			}
			ii, _ := Random(arr, len(arr))

			//将buf中1，2，3，4...位置的ids存到A中随机空位
			for j, v := range ii {
				A[v] = H(k2, buf[j])
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
			rst := partition2(len(iistr), b, iistr)

			d := F(k2, rst)
			l := H(k1, 1)
			L[l] = d
		} else {
			buf := partition(len(n.ss), B, db)

			//数组随机位置存储
			var arr []int
			for k := range emp {
				arr = append(arr, k)
			}
			ii, _ := Random(arr, len(arr))

			//将buf中1，2，3，4...位置的ids存到A中随机空位
			for j, v := range ii {
				A[v] = H(k2, buf[j])
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

			arr = []int{}
			for k := range emp {
				arr = append(arr, k)
			}
			ii, _ = Random(arr, len(arr))

			//将buf2中1，2，3，4...位置的ids存到A中随机空位
			for j, v := range ii {
				A[v] = H(k2, buf2[j])
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

			d := F(k2, rst)
			l := H(k1, 1)
			L[l] = d
		}
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
			for ; p < q; p++ {
				str += " " + strconv.Itoa(rand.Int())
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
			for ; p < q; p++ {
				str += " " + strconv.Itoa(rand.Int())
			}
		}
		rst += "," + str
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

func F(k int64) (val int64) {
	rand.Seed(k)
	return rand.Int63()
}

func H(k int64) (val int64) {
	rand.Seed(k)
	return rand.Int63()
}

//写文件函数
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

//dfs
func dfs(t *Node) {
	if t == nil {
		return
	}
	writeFile("abc.txt", "fs:"+t.fs+" ss:"+t.ss+"\n")

	for _, n := range t.childs {
		dfs(n)
	}
}

func main() {
	//treeBuild()
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	b, _ := Random(arr, 7)

	println(len(arr))
	for i := 0; i < 7; i++ {
		print(b[i])
	}
}
