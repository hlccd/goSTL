package trie

//@Title		trie
//@Description
//		单词查找树的节点
//		可通过节点的分叉对string进行查找
//		增添string时候只需要增删结点即可
//		当string到终点时存储元素

//node树节点结构体
//该节点是trie的树节点
//结点存储到此时的string的前缀数量
//以son为分叉存储下属的string
//该节点同时存储其元素
type node struct {
	num   int
	son   [64]*node
	value interface{}
}

//@title    newNode
//@description
//		新建一个单词查找树节点并返回
//		将传入的元素e作为该节点的承载元素
//@receiver		nil
//@param    	e			interface{}				承载元素e
//@return    	n        	*node					新建的单词查找树节点的指针
func newNode(e interface{}) (n *node) {
	return &node{
		num:   0,
		value: e,
	}
}

//@title    inOrder
//@description
//		以node单词查找树节点做接收者
//		遍历其分叉以找到其存储的所有string
//@receiver		n			*node					接受者node的指针
//@param    	s			string					到该结点时的前缀string
//@return    	es        	[]interface{}			以该前缀s为前缀的所有string的集合
func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.value != nil {
		es = append(es, s)
	}
	for i, p := 0, 0; i < 62 && p < n.num; i++ {
		if n.son[i] != nil {
			if i < 26 {
				es = append(es, n.son[i].inOrder(s+string(i+'a'))...)
			} else if i < 52 {
				es = append(es, n.son[i].inOrder(s+string(i-26+'A'))...)
			} else {
				es = append(es, n.son[i].inOrder(s+string(i-52+'0'))...)
			}
			p++
		}
	}
	return es
}

//@title    getIdx
//@description
//		传入一个byte并根据其值返回其映射到分叉的值
//		当不属于'a'~'z','A'~'Z','0'~'9','+','/'时返回-1
//@receiver		nil
//@param    	c			byte					待映射的ASCII码
//@return    	idx        	int						以c映射出的分叉下标
func getIdx(c byte) (idx int) {
	if c >= 'a' && c <= 'z' {
		idx = int(c - 'a')
	} else if c >= 'A' && c <= 'Z' {
		idx = int(c-'A') + 26
	} else if c >= '0' && c <= '9' {
		idx = int(c-'0') + 52
	} else if c == '+' {
		idx = 62
	} else if c == '/' {
		idx = 63
	} else {
		idx = -1
	}
	return idx
}

//@title    insert
//@description
//		以node单词查找树节点做接收者
//		从n节点中继续插入以s为索引的元素e,且当前抵达的string位置为p
//		当到达s终点时进行插入,如果此时node承载了元素则插入失败,否则成功
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可
//@receiver		n			*node					接受者node的指针
//@param    	s			string					待插入元素的索引s
//@param    	p			int						索引当前抵达的位置
//@param    	e			interface{}				待插入元素e
//@return    	b        	bool					是否插入成功?
func (n *node) insert(s string, p int, e interface{}) (b bool) {
	if p == len(s) {
		if n.value != nil {
			return false
		}
		n.value = e
		n.num++
		return true
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return false
	}
	if n.son[idx] == nil {
		n.son[idx] = newNode(nil)
	}
	b = n.son[idx].insert(s, p+1, e)
	if b {
		n.num++
	}
	return b
}

//@title    erase
//@description
//		以node单词查找树节点做接收者
//		从n节点中继续删除以s为索引的元素e,且当前抵达的string位置为p
//		当到达s终点时进行删除,如果此时node未承载元素则删除失败,否则成功
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接失败
//@receiver		n			*node					接受者node的指针
//@param    	s			string					待删除元素的索引s
//@param    	p			int						索引当前抵达的位置
//@return    	b        	bool					是否删除成功?
func (n *node) erase(s string, p int) (b bool) {
	if p == len(s) {
		if n.value != nil {
			n.value = nil
			n.num--
			return true
		}
		return false
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return false
	}
	if n.son[idx] == nil {
		return false
	}
	b = n.son[idx].erase(s, p+1)
	if b {
		n.num--
		if n.son[idx].num == 0 {
			n.son[idx] = nil
		}
	}
	return b
}

//@title    count
//@description
//		以node单词查找树节点做接收者
//		从n节点中继续查找以s为前缀索引的元素e,且当前抵达的string位置为p
//		当到达s终点时返回其值即可
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,当其分叉为nil则直接返回0
//@receiver		n			*node					接受者node的指针
//@param    	s			string					待查找元素的前缀索引
//@param    	p			int						索引当前抵达的位置
//@return    	num        	int						以该s为前缀的string的数量
func (n *node) count(s string, p int) (num int) {
	if p == len(s) {
		return n.num
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return 0
	}
	if n.son[idx] == nil {
		return 0
	}
	return n.son[idx].count(s, p+1)
}

//@title    find
//@description
//		以node单词查找树节点做接收者
//		从n节点中继续查找以s为前缀索引的元素e,且当前抵达的string位置为p
//		当到达s终点时返回其承载的元素即可
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,当其分叉为nil则直接返回nil
//@receiver		n			*node					接受者node的指针
//@param    	s			string					待查找元素的前缀索引
//@param    	p			int						索引当前抵达的位置
//@return    	e			interface{}				该索引所指向的元素e
func (n *node) find(s string, p int) (e interface{}) {
	if p == len(s) {
		return n.value
	}
	idx := getIdx(s[p])
	if idx == -1 {
		return nil
	}
	if n.son[idx] == nil {
		return nil
	}
	return n.son[idx].find(s, p+1)
}
