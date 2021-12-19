package rax

//@Title		rax
//@Description
//		前缀基数树的节点
//		可通过节点的分叉对string进行查找
//		增添string时候需要增删结点,同时将结点内置的map中增删删除对应的string即可
//		当string到终点时存储元素

//node树节点结构体
//该节点是rax的树节点
//结点存储到此时的string的前缀数量
//son存储其下属分叉的子结点指针
//该节点同时存储其元素
type node struct {
	name  string           //以当前结点的string内容
	num   int              //以当前结点为前缀的数量
	value interface{}      //当前结点存储的元素
	sons  map[string]*node //该结点下属结点的指针
}

//@title    newNode
//@description
//		新建一个前缀基数树节点并返回
//		将传入的元素e作为该节点的承载元素
//@receiver		nil
//@param    	name		string					该节点的名字,即其对应的string
//@param    	e			interface{}				承载元素e
//@return    	n        	*node					新建的单词查找树节点的指针
func newNode(name string, e interface{}) (n *node) {
	return &node{
		name:  name,
		num:   0,
		value: e,
		sons:  make(map[string]*node),
	}
}

//@title    inOrder
//@description
//		以node前缀基数树节点做接收者
//		遍历其分叉以找到其存储的所有string
//@receiver		n			*node					接受者node的指针
//@param    	s			string					到该结点时的前缀string
//@return    	es        	[]interface{}			以该前缀s为前缀的所有string的集合
func (n *node) inOrder(s string) (es []interface{}) {
	if n == nil {
		return es
	}
	if n.value != nil {
		es = append(es, s+n.name)
	}
	for _, son := range n.sons {
		es = append(es, son.inOrder(s+n.name+"/")...)
	}
	return es
}

//@title    insert
//@description
//		以node前缀基数树节点做接收者
//		从n节点中继续插入以s为索引的元素e,且当前抵达的string位置为p
//		若该层string为""时候视为失败
//		当到达s终点时进行插入,如果此时node承载了元素则插入失败,否则成功
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可
//		当插入失败且对应子结点为新建节点时则需要删除该子结点
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@param    	e			interface{}				待插入元素e
//@return    	b        	bool					是否插入成功?
func (n *node) insert(ss []string, p int, e interface{}) (b bool) {
	if p == len(ss) {
		if n.value != nil {
			return false
		}
		n.value = e
		n.num++
		return true
	}
	s := ss[p]
	if s == "" {
		return false
	}
	//从其子结点的map中找到对应的方向
	son, ok := n.sons[s]
	if !ok {
		//不存在,新建并放入map中
		son = newNode(s, nil)
		n.sons[s] = son
	}
	//从子结点对应方向继续插入
	b = son.insert(ss, p+1, e)
	if b {
		n.num++
	} else {
		if !ok {
			//插入失败且该子节点为新建结点则需要删除该子结点
			delete(n.sons, s)
		}
	}
	return b
}

//@title    erase
//@description
//		以node前缀基数树节点做接收者
//		从n节点中继续删除以s为索引的元素e,且当前抵达的string位置为p
//		若该层string为""时候视为失败
//		当到达s终点时进行删除,如果此时node未承载元素则删除失败,否则成功
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接失败
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	b        	bool					是否删除成功?
func (n *node) erase(ss []string, p int) (b bool) {
	if p == len(ss) {
		if n.value != nil {
			n.value = nil
			n.num--
			return true
		}
		return false
	}
	//从map中找到对应下子结点位置并递归进行删除
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		//未找到或son不存在,删除失败
		return false
	}
	b = son.erase(ss, p+1)
	if b {
		n.num--
		if son.num <= 0 {
			//删除后子结点的num<=0即该节点无后续存储元素,可以销毁
			delete(n.sons, s)
		}
	}
	return b
}

//@title    delete
//@description
//		以node前缀基数树节点做接收者
//		从n节点中继续删除以s为索引的元素e,且当前抵达的string位置为p
//		若该层string为""时候视为失败
//		当到达s终点时进行删除,删除所有后续元素,并返回其后续元素的数量
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接返回0
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	num        	int						被删除元素的数量
func (n *node) delete(ss []string, p int) (num int) {
	if p == len(ss) {
		return n.num
	}
	//从map中找到对应下子结点位置并递归进行删除
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	num = son.delete(ss, p+1)
	if num > 0 {
		n.num -= num
		if son.num <= 0 {
			//删除后子结点的num<=0即该节点无后续存储元素,可以销毁
			delete(n.sons, s)
		}
	}
	return num
}

//@title    count
//@description
//		以node前缀基数树节点做接收者
//		从n节点中继续查找以s为前缀索引的元素e,且当前抵达的string位置为p
//		若该层string为""时候视为查找失败
//		当到达s终点时返回其值即可
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,当其分叉为nil则直接返回0
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	num        	int						以该s为前缀的string的数量
func (n *node) count(ss []string, p int) (num int) {
	if p == len(ss) {
		return n.num
	}
	//从map中找到对应下子结点位置并递归进行查找
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	return son.count(ss, p+1)
}

//@title    find
//@description
//		以node前缀基数树节点做接收者
//		从n节点中继续查找以s为前缀索引的元素e,且当前抵达的string位置为p
//		若该层string为""时候视为查找失败
//		当到达s终点时返回其承载的元素即可
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,当其分叉为nil则直接返回nil
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	e			interface{}				该索引所指向的元素e
func (n *node) find(ss []string, p int) (e interface{}) {
	if p == len(ss) {
		return n.value
	}
	//从map中找到对应下子结点位置并递归进行查找
	s := ss[p]
	son, ok := n.sons[s]
	if !ok || son == nil {
		return 0
	}
	return son.find(ss, p+1)
}
