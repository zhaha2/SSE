package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Node struct {
	childs []*Node
	fs     string
	ss     string
	val    string
}

func forestGen() {
	var db map[string]string
	db = map[string]string{"one": "1", "two": "2"}
	fmt.Print(db["w"])
	sort(db)

	//定义变量
	m := 10
	var s []*Node
	top := -1
	var roots []*Node
	var dbw string
	var p *Node

	for w := range db {
		p = nil
		max := 0
		flag := 0
		s = []*Node{}
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
			is := Intersect(a, b)

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
			} else if len(is) > max && len(is) > m {
				p = n
				max = len(is)
				flag = 3
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
			t.ss = dbw + p.ss
			t.val = w
			t.childs = nil

			p.childs = append(p.childs, t)
		} else if flag == 3 {
			//父母节点的指向还是没变
			t := new(Node)

			//取交集
			a := strings.Split(p.fs, " ")
			b := strings.Split(dbw, " ")

			//数组转字符串
			//strings.Replace(strings.Trim(fmt.Sprint(Intersect(a, b)), "[]"), " ", ",", -1)
			t.fs = strings.Trim(fmt.Sprint(Intersect(a, b)), "[]")
			t.ss = t.fs
			t.val = ""
			t.childs = nil

			u := new(Node)
			u.fs = dbw
			//取差集
			u.ss = dbw + t.ss
			u.val = w

			p.ss = p.fs + t.ss

			t.childs = append(t.childs, u)
			t.childs = append(t.childs, p)
		}
	}

	//print(len(node.childs))
}

//判断交集
func Intersect(a []string, b []string) []string {
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

func indexGen(r []*Node) {
	//初始化
	K := 1
	var L map[string]string
	L = map[string]string{"one": "1", "two": "2"}
	top := len(r) - 1
	b := 20
	B := 10

	for top >= 0 {
		n := r[top]
		top--

		k1, k2 := F(K, n.val)
		t := len(n.ss) / B
		db := strings.Split(n.fs, " ")

		if len(n.ss) < b {
			//如果是空值的怎么办
			for p <= len(n.ss) {
				q := p + b

				//分块
				if q < len(n.ss) {
					for ; p < q; p++ {
						str += " " + db[p]
					}
				} else {
					for ; p < len(n.ss); p++ {
						str += " " + db[p]
					}
					//Pad
					str += ","
					for ; p < q; p++ {
						str += " " + strconv.Itoa(rand.Int())
					}
				}

				d := H(k2, str)

				l := H(k1, 1)
				L[l] = d
			}
		} else if len(n.ss) > b && len(n.ss) < b*B {
			str := partition(len(n.ss), b, db)
			print(str)
		}
	}
}

func partition(len int, b int, db []string) string {
	str := ""
	p := 0
	for p <= len {
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
			//Pad
			str += ","
			for ; p < q; p++ {
				str += " " + strconv.Itoa(rand.Int())
			}
		}
	}
	return str
}

func main() {
	//treeBuild()
}
