package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	childs []*Node
	fs     string
	ss     string
	val    string
}

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
	var delete []*Node
	var dbw string
	var p *Node

	for w := range db {
		p = nil
		max := 0
		flag := 0
		s = []*Node{}

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
			t.val = w
			t.childs = nil

			roots = append(roots, t)
		} else if flag == 1 {
			//没想好
			p.val += w
		} else if flag == 2 {
			t := new(Node)
			t.fs = dbw
			//取差集
			a := strings.Split(dbw, " ")
			b := strings.Split(p.fs, " ")
			t.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			t.val = w
			t.childs = nil

			p.childs = append(p.childs, t)
		} else if flag == 3 {
			t := new(Node)
			t.fs = dbw
			t.ss = t.fs
			t.val = w
			t.childs = append(t.childs, p)

			roots = append(roots, t)
			delete = append(delete, p)
		} else if flag == 4 {
			//没想好
			//父母节点的指向还是没变
			t := new(Node)

			//取交集
			a := strings.Split(p.fs, " ")
			b := strings.Split(dbw, " ")

			//数组转字符串
			//strings.Replace(strings.Trim(fmt.Sprint(Intersect(a, b)), "[]"), " ", ",", -1)
			t.fs = strings.Trim(fmt.Sprint(intersect(a, b)), "[]")
			t.ss = t.fs
			t.val = ""
			t.childs = nil

			u := new(Node)
			u.fs = dbw
			//取差集
			a = strings.Split(dbw, " ")
			b = strings.Split(t.fs, " ")
			u.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")
			u.val = w
			u.childs = nil

			a = strings.Split(p.fs, " ")
			b = strings.Split(t.fs, " ")
			p.ss = strings.Trim(fmt.Sprint(difference(a, b)), "[]")

			t.childs = append(t.childs, u)
			t.childs = append(t.childs, p)
		}
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

func indexGen(roots []*Node) {
	//初始化
	K := 1
	var L map[string]string
	top := len(roots) - 1
	//实验确定
	b := 20
	B := 10

	for top >= 0 {
		n := roots[top]
		top--
		//预定义A大小
		var A []string

		k1, k2 := F(K, n.val)
		t := len(n.ss) / B
		db := strings.Split(n.fs, " ")

		if len(n.ss) < b {
			//如果是空值的怎么办
			str := partition(len(n.ss), b, db)

			d := H(k2, str)
			l := H(k1, 1)
			L[l] = d
		} else if len(n.ss) > b && len(n.ss) < b*B {
			buf := strings.Split(partition2(len(n.ss), B, db), " ")

			//数组随机位置存储
			var arr []int
			for i := 0; i < 100; i++ {
				arr = append(arr, i)
			}
			ii, _ := Random(arr, 100)

			//这里还是只能连续存储
			for _, v := range ii {
				A = append(A, H(k2, buf[v]))
			}

			//把ii打包分块
			var iistr []string
			for i := 0; i < 100; i++ {
				iistr = append(iistr, strconv.Itoa(ii[i]))
			}
			rst := partition(len(iistr), b, iistr)

			d := H(k2, rst)
			l := H(k1, 1)
			L[l] = d
		} else {
			buf := strings.Split(partition2(len(n.ss), B, db), " ")

			//数组随机位置存储
			var arr []int
			for i := 0; i < 100; i++ {
				arr = append(arr, i)
			}
			ii, _ := Random(arr, 100)

			//这里还是只能连续存储
			for _, v := range ii {
				A = append(A, H(k2, buf[v]))
			}

			//把ii打包分块
			var iistr []string
			for i := 0; i < 100; i++ {
				iistr = append(iistr, strconv.Itoa(ii[i]))
			}
			rst := partition(len(iistr), b, iistr)

			d := H(k2, rst)
			l := H(k1, 1)
			L[l] = d
		}
	}
}

func partition(len int, b int, db []string) []string {
	var rst []string
	p := 0
	for p < len {
		q := p + b

		//分块
		if q < len {
			for ; p < q; p++ {
				rst = append(rst, db[p])
			}
		} else {
			for ; p < len; p++ {
				rst = append(rst, db[p])
			}
			//Pad
			//填充什么
			for ; p < q; p++ {
				rst = append(rst, strconv.Itoa(rand.Int()))
			}
		}
	}
	return rst
}

func partition2(len int, b int, db []string) string {
	str := ""
	p := 0
	for p < len {
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
			str += ","
			//填充什么
			for ; p < q; p++ {
				str += " " + strconv.Itoa(rand.Int())
			}
		}
	}
	return str
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

func main() {
	//treeBuild()
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	b, _ := Random(arr, 7)

	println(len(arr))
	for i := 0; i < 7; i++ {
		print(b[i])
	}
}
