package deque

//@Title		deque
//@Description
//		deque双队列容器包
//		该部分包含了deque双向队列中链表节点
//		链表的增删都通过节点的增删完成
//		节点空间全部使用或全部废弃后将进行节点增删
//		增删之后会返回对应的首尾节点以辅助deque容器仍持有首尾节点
//		为保证效率问题,设定了一定的冗余量,即每个节点设定2^10的空间以存放元素

//deque双向队列中链表的node节点结构体
//包含一个2^10空间的固定数组用以承载元素
//使用begin和end两个下标用以表示新增的元素的下标,由于begin可能出现-1所以不选用uint16
//pre和next是该节点的前后两个节点
//用以保证链表整体是相连的
type node struct {
	data  [1024]interface{} //用于承载元素的股东数组
	begin int16             //该结点在前方添加结点的下标
	end   int16             //该结点在后方添加结点的下标
	pre   *node             //该结点的前一个结点
	next  *node             //该节点的后一个结点
}

//deque双向队列容器接口
//存放了deque容器可使用的函数
//对应函数介绍见下方
type noder interface {
	nextNode() (m *node)                   //返回下一个结点
	preNode() (m *node)                    //返回上一个结点
	value() (es []interface{})             //返回该结点所承载的所有元素
	pushFront(e interface{}) (first *node) //在该结点头部添加一个元素,并返回新首结点
	pushBack(e interface{}) (last *node)   //在该结点尾部添加一个元素,并返回新尾结点
	popFront() (first *node)               //弹出首元素并返回首结点
	popBack() (last *node)                 //弹出尾元素并返回尾结点
	front() (e interface{})                //返回首元素
	back() (e interface{})                 //返回尾元素
}

//@title    createFirst
//@description
//		新建一个冗余在前方的首结点并返回其指针
//		初始首结点的begin为1023,end为1024
//		该结点的前后结点均置为nil
//@receiver		nil
//@param    	nil
//@return    	n        	*node					新建的node指针
func createFirst() (n *node) {
	return &node{
		data:  [1024]interface{}{},
		begin: 1023,
		end:   1024,
		pre:   nil,
		next:  nil,
	}
}

//@title    createLast
//@description
//		新建一个冗余在后方的尾结点并返回其指针
//		初始首结点的begin为-1,end为0
//		该结点的前后结点均置为nil
//@receiver		nil
//@param    	nil
//@return    	n        	*node					新建的node指针
func createLast() (n *node) {
	return &node{
		data:  [1024]interface{}{},
		begin: -1,
		end:   0,
		pre:   nil,
		next:  nil,
	}
}

//@title    nextNode
//@description
//		以node结点做接收者
//		返回该结点的后一个结点
//		如果n为nil则返回nil
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	m        	*node					n结点的下一个结点m的指针
func (n *node) nextNode() (m *node) {
	if n == nil {
		return nil
	}
	return n.next
}

//@title    preNode
//@description
//		以node结点做接收者
//		返回该结点的前一个结点
//		如果n为nil则返回nil
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	m        	*node					n结点的上一个结点m的指针
func (n *node) preNode() (m *node) {
	if n == nil {
		return nil
	}
	return n.pre
}

//@title    value
//@description
//		以node结点做接收者
//		返回该结点所承载的所有元素
//		根据其begin和end来获取其元素
//		当该结点为nil时返回[]而非nil
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	es        	[]interface{}			该结点所承载的所有元素
func (n *node) value() (es []interface{}) {
	es = make([]interface{}, 0, 0)
	if n == nil {
		return es
	}
	if n.begin > n.end {
		return es
	}
	es = n.data[n.begin+1 : n.end]
	return es
}

//@title    pushFront
//@description
//		以node结点做接收者
//		向该节点前方添加元素e
//		当该结点空间已经使用完毕后,新建一个结点并将新结点设为首结点
//		将插入元素放入新结点并返回新结点作为新的首结点
//		否则插入当前结点并返回当前结点,首结点不变
//@receiver		n        	*node					接收者的node指针
//@param    	e			interface{}				待插入元素
//@return    	first     	*node					首节点指针
func (n *node) pushFront(e interface{}) (first *node) {
	if n == nil {
		return n
	}
	if n.begin >= 0 {
		//该结点仍有空间可用于承载元素
		n.data[n.begin] = e
		n.begin--
		return n
	}
	//该结点无空间承载,创建新的首结点用于存放
	m := createFirst()
	m.data[m.begin] = e
	m.next = n
	n.pre = m
	m.begin--
	return m
}

//@title    pushBack
//@description
//		以node结点做接收者
//		向该节点后方添加元素e
//		当该结点空间已经使用完毕后,新建一个结点并将新结点设为尾结点
//		将插入元素放入新结点并返回新结点作为新的尾结点
//		否则插入当前结点并返回当前结点,尾结点不变
//@receiver		n        	*node					接收者的node指针
//@param    	e			interface{}				待插入元素
//@return    	last     	*node					尾节点指针
func (n *node) pushBack(e interface{}) (last *node) {
	if n == nil {
		return n
	}
	if n.end < int16(len(n.data)) {
		//该结点仍有空间可用于承载元素
		n.data[n.end] = e
		n.end++
		return n
	}
	//该结点无空间承载,创建新的尾结点用于存放
	m := createLast()
	m.data[m.end] = e
	m.pre = n
	n.next = m
	m.end++
	return m
}

//@title    popFront
//@description
//		以node结点做接收者
//		利用首节点进行弹出元素,可能存在首节点全部释放要进行首节点后移的情况
//		当发生首结点后移后将会返回新首结点,否则返回当前结点
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	first     	*node					首节点指针
func (n *node) popFront() (first *node) {
	if n == nil {
		return nil
	}
	if n.begin < int16(len(n.data))-2 {
		//该结点仍有承载元素
		n.begin++
		n.data[n.begin] = nil
		return n
	}
	if n.next != nil {
		//清除该结点下一节点的前结点指针
		n.next.pre = nil
	}
	return n.next
}

//@title    popBack
//@description
//		以node结点做接收者
//		利用尾节点进行弹出元素,可能存在尾节点全部释放要进行尾节点前移的情况
//		当发生尾结点前移后将会返回新尾结点,否则返回当前结点
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	last     	*node					尾节点指针
func (n *node) popBack() (last *node) {
	if n == nil {
		return nil
	}
	if n.end > 1 {
		//该结点仍有承载元素
		n.end--
		n.data[n.end] = nil
		return n
	}
	if n.pre != nil {
		//清除该结点上一节点的后结点指针
		n.pre.next = nil
	}
	return n.pre
}

//@title    front
//@description
//		以node结点做接收者
//		返回该结点的第一个元素,利用首节点和begin进行查找
//		若该结点为nil,则返回nil
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	e			interface{}				该结点承载的的第一个元素
func (n *node) front() (e interface{}) {
	if n == nil {
		return nil
	}
	return n.data[n.begin+1]
}

//@title    back
//@description
//		以node结点做接收者
//		返回该结点的最后一个元素,利用尾节点和end进行查找
//		若该结点为nil,则返回nil
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	e			interface{}				该结点承载的的最后一个元素
func (n *node) back() (e interface{}) {
	if n == nil {
		return nil
	}
	return n.data[n.end-1]
}
