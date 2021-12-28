package radix

import (
	"strings"
)

//@Title		radix
//@Description
//		前缀基数树的节点
//		可通过节点的分叉对string进行查找
//		增添string时候需要增删结点,同时将结点内置的map中增删对应的string即可
//		当string到终点时存储元素

//node树节点结构体
//该节点是radix的树节点
//结点存储到此时的string的前缀数量
//son存储其下属分叉的子结点指针
//该节点同时存储其元素
type node struct {
	pattern string           //到终点时不为"",其他都为""
	part    string           //以当前结点的string内容
	num     int              //以当前结点为前缀的数量
	sons    map[string]*node //该结点下属结点的指针
	fuzzy   bool             //模糊匹配?该结点首字符为':'或'*'为模糊匹配
}

//@title    newNode
//@description
//		新建一个前缀基数树节点并返回
//		将传入的元素e作为该节点的承载元素
//@receiver		nil
//@param    	name		string					该节点的名字,即其对应的string
//@return    	n        	*node					新建的单词查找树节点的指针
func newNode(part string) (n *node) {
	fuzzy := false
	if len(part) > 0 {
		fuzzy = part[0] == ':' || part[0] == '*'
	}
	return &node{
		pattern: "",
		part:    part,
		num:     0,
		sons:    make(map[string]*node),
		fuzzy:   fuzzy,
	}
}

//@title    analysis
//@description
//		将string按'/'进行分段解析
//		为""部分直接舍弃,返回解析结果
//		同时按规则重组用以解析是string并返回
//@receiver		nil
//@param    	s			string					待解析的string
//@return    	ss			[]string				按'/'进行分层解析后的结果
//@return    	newS		string					按符合规则解析结果重组后的s
func analysis(s string) (ss []string, newS string) {
	vs := strings.Split(s, "/")
	ss = make([]string, 0)
	newS = "/"
	for _, item := range vs {
		if item != "" {
			ss = append(ss, item)
			newS = newS + "/" + item
			if item[0] == '*' {
				break
			}
		}
	}
	return ss, newS
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
	if n.pattern != "" {
		es = append(es, s+n.part)
	}
	for _, son := range n.sons {
		es = append(es, son.inOrder(s+n.part+"/")...)
	}
	return es
}

//@title    insert
//@description
//		以node前缀基数树节点做接收者
//		从n节点中继续插入以s为索引的元素e,且当前抵达的string位置为p
//		当到达s终点时进行插入,如果此时node承载了string则插入失败,否则成功
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可
//		当插入失败且对应子结点为新建节点时则需要删除该子结点
//@receiver		n			*node					接受者node的指针
//@param    	pattern		string					待插入的string整体
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	b        	bool					是否插入成功?
func (n *node) insert(pattern string, ss []string, p int) (b bool) {
	if p == len(ss) {
		if n.pattern != "" {
			//该节点承载了string
			return false
		}
		//成功插入
		n.pattern = pattern
		n.num++
		return true
	}
	//找到该层的string
	s := ss[p]
	//从其子结点的map中找到对应的方向
	son, ok := n.sons[s]
	if !ok {
		//不存在,新建并放入map中
		son = newNode(s)
		n.sons[s] = son
	}
	//从子结点对应方向继续插入
	b = son.insert(pattern, ss, p+1)
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
//		当到达s终点时进行删除,如果此时node未承载元素则删除失败,否则成功
//		当未到达终点时,根据当前抵达的位置去寻找其子结点继续遍历即可,若其分叉为nil则直接失败
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	b        	bool					是否删除成功?
func (n *node) erase(ss []string, p int) (b bool) {
	if p == len(ss) {
		if n.pattern != "" {
			//该结点承载是string是,删除成功
			n.pattern = ""
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
		son.num -= num
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

//@title    mate
//@description
//		以node前缀基数树节点做接收者
//		先从radix树的根节点开始找到第一个可以满足该模糊匹配方案的string结点
//		随后将s和结点的pattern进行模糊映射,将模糊查找的值和匹配值进行映射并返回即可
//		若该结点未找到则直接返回nil和false即可
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	m			map[string]string		s从结点中利用模糊匹配到的所有key和value的映射
//@return    	ok			bool					匹配成功?
func (n *node) mate(s string, p int) (m map[string]string, ok bool) {
	//解析url
	searchParts, _ := analysis(s)
	//从该请求类型中寻找对应的路由结点
	q := n.find(searchParts, 0)
	if q != nil {
		//解析该结点的pattern
		parts, _ := analysis(q.pattern)
		//动态参数映射表
		params := make(map[string]string)
		for index, part := range parts {
			if part[0] == ':' {
				//动态匹配,将参数名和参数内容的映射放入映射表内
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				//通配符,将后续所有内容全部添加到映射表内同时结束遍历
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return params, true
	}
	return nil, false
}

//@title    find
//@description
//		以node前缀基数树节点做接收者
//		从radix树的根节点开始找到第一个可以满足该模糊匹配方案的string结点
//		若该结点未找到则直接返回nil
//@receiver		n			*node					接受者node的指针
//@param    	ss			[]string				待删除元素的索引s的按'/'进行分层的索引集合
//@param    	p			int						索引当前抵达的位置
//@return    	m			map[string]string		s从结点中利用模糊匹配到的所有key和value的映射
//@return    	ok			bool					匹配成功?
func (n *node) find(parts []string, height int) (q *node) {
	//根据长度和局部string的首字符进行判断
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			//匹配失败,该结点处无匹配的信息
			return nil
		}
		//匹配成功,返回该结点
		return n
	}
	//从该结点的所有子结点中查找可用于递归查找的结点
	//当局部string信息和当前层string相同时可用于递归查找
	//当该子结点是动态匹配时也可以用于递归查找
	part := parts[height]
	//从所有子结点中找到可用于递归查找的结点
	children := make([]*node, 0, 0)
	for _, child := range n.sons {
		if child.part == part || child.fuzzy {
			//局部string相同或动态匹配
			children = append(children, child)
		}
	}
	for _, child := range children {
		//递归查询,并根据结果进行判断
		result := child.find(parts, height+1)
		if result != nil {
			//存在一个满足时就可以返回
			return result
		}
	}
	return nil
}
