package list

//@Title		list
//@Description
//		list链表容器包
//		该部分包含了链表的节点
//		链表的增删都通过节点的增删完成
//		结点间可插入其前后节点,并同时将两结点建立连接
//		增删之后会返回对应的首尾节点以辅助list容器仍持有首尾节点
//		为保证效率问题,设定了一定的冗余量,即每个节点设定2^10的空间以存放元素

//链表的node节点结构体
//pre和next是该节点的前后两个节点的指针
//用以保证链表整体是相连的
type node struct {
	data interface{} //结点所承载的元素
	pre  *node       //前结点指针
	next *node       //后结点指针
}

//node结点容器接口
//存放了node容器可使用的函数
//对应函数介绍见下方

type noder interface {
	preNode() (m *node)     //返回前结点指针
	nextNode() (m *node)    //返回后结点指针
	insertPre(pre *node)    //在该结点前插入结点并建立连接
	insertNext(next *node)  //在该结点后插入结点并建立连接
	erase()                 //删除该结点,并使该结点前后两结点建立连接
	value() (e interface{}) //返回该结点所承载的元素
	setValue(e interface{}) //修改该结点承载元素为e
}

//@title    newNode
//@description
//		新建一个结点并返回其指针
//		初始首结点的前后结点指针都为nil
//@receiver		nil
//@param    	nil
//@return    	n        	*node					新建的node指针
func newNode(e interface{}) (n *node) {
	return &node{
		data: e,
		pre:  nil,
		next: nil,
	}
}

//@title    preNode
//@description
//		以node结点做接收者
//		返回该结点的前结点
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	pre			*node					该结点的前结点指针
func (n *node) preNode() (pre *node) {
	if n == nil {
		return
	}
	return n.pre
}

//@title    nextNode
//@description
//		以node结点做接收者
//		返回该结点的后结点
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	next		*node					该结点的后结点指针
func (n *node) nextNode() (next *node) {
	if n == nil {
		return
	}
	return n.next
}

//@title    insertPre
//@description
//		以node结点做接收者
//		对该结点插入前结点
//		并建立前结点和该结点之间的连接
//@receiver		n        	*node					接收者的node指针
//@param    	pre			*node					该结点的前结点指针
//@return    	nil
func (n *node) insertPre(pre *node) {
	if n == nil || pre == nil {
		return
	}
	pre.next = n
	pre.pre = n.pre
	if n.pre != nil {
		n.pre.next = pre
	}
	n.pre = pre
}

//@title    insertNext
//@description
//		以node结点做接收者
//		对该结点插入后结点
//		并建立后结点和该结点之间的连接
//@receiver		n        	*node					接收者的node指针
//@param    	next		*node					该结点的后结点指针
//@return    	nil
func (n *node) insertNext(next *node) {
	if n == nil || next == nil {
		return
	}
	next.pre = n
	next.next = n.next
	if n.next != nil {
		n.next.pre = next
	}
	n.next = next
}

//@title    erase
//@description
//		以node结点做接收者
//		销毁该结点
//		同时建立该节点前后节点之间的连接
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	nil
func (n *node) erase() {
	if n == nil {
		return
	}
	if n.pre == nil && n.next == nil {
		return
	} else if n.pre == nil {
		n.next.pre = nil
	} else if n.next == nil {
		n.pre.next = nil
	} else {
		n.pre.next = n.next
		n.next.pre = n.pre
	}
	n = nil
}

//@title    value
//@description
//		以node结点做接收者
//		返回该结点所要承载的元素
//@receiver		n        	*node					接收者的node指针
//@param    	nil
//@return    	e			interface{}				该节点所承载的元素e
func (n *node) value() (e interface{}) {
	if n == nil {
		return nil
	}
	return n.data
}

//@title    setValue
//@description
//		以node结点做接收者
//		对该结点设置其承载的元素
//@receiver		n        	*node					接收者的node指针
//@param    	e			interface{}				该节点所要承载的元素e
//@return    	nil
func (n *node) setValue(e interface{}) {
	if n == nil {
		return
	}
	n.data = e
}
