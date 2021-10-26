package cbTree

//@Title		cbTree
//@Description
//		完全二叉树的节点
//		可通过节点实现完全二叉树的节点上升与下沉
//		也可查找待插入的最后一个节点的父节点,即该待插入节点将放入该父节点的左右子节点中
//@author     	hlccd		2021-07-14
import "github.com/hlccd/goSTL/utils/comparator"

//node树节点结构体
//该节点是完全二叉树的树节点
//该节点除了保存承载元素外,还将保存父节点、左右子节点的指针
type node struct {
	value  interface{} //节点中存储的元素
	parent *node       //父节点指针
	left   *node       //左节点指针
	right  *node       //右节点指针
}

//@title    newNode
//@description
//		新建一个完全二叉树节点并返回
//		将传入的元素e作为该节点的承载元素
//		将传入的parent节点作为其父节点,左右节点设为nil
//@auth      	hlccd		2021-07-14
//@receiver		nil
//@param    	parent		*node					新建节点的父节点指针
//@param    	e			interface{}				承载元素e
//@return    	n        	*node					新建的完全二叉树节点的指针
func newNode(parent *node, e interface{}) (n *node) {
	return &node{
		value:  e,
		parent: parent,
		left:   nil,
		right:  nil,
	}
}

//@title    frontOrder
//@description
//		以node节点做接收者
//		以前缀序列返回节点集合
//@auth      	hlccd		2021-07-14
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	es        	[]interface{}			以该节点为起点的前缀序列
func (n *node) frontOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	es = append(es, n.value)
	if n.left != nil {
		es = append(es, n.left.frontOrder()...)
	}
	if n.right != nil {
		es = append(es, n.right.frontOrder()...)
	}
	return es
}

//@title    lastParent
//@description
//		以node节点做接收者
//		根据传入数值通过转化为二进制的方式模拟查找最后一个父节点
//		由于查找父节点的路径等同于转化为二进制后除开首位的中间值,故该方案是可行的
//@auth      	hlccd		2021-07-14
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	ans        	*node					查找到的最后一个父节点
func (n *node) lastParent(num int) (ans *node) {
	if num > 3 {
		//去掉末尾的二进制值
		arr := make([]int, 0, 32)
		ans = n
		for num > 0 {
			//转化为二进制
			arr = append(arr, num%2)
			num /= 2
		}
		//去掉首位的二进制值
		for i := len(arr) - 2; i >= 1; i-- {
			if arr[i] == 1 {
				ans = ans.right
			} else {
				ans = ans.left
			}
		}
		return ans
	}
	return n
}

//@title    insert
//@description
//		以node节点做接收者
//		从该节点插入元素e,并根据传入的num寻找最后一个父节点用于插入最后一位值
//		随后对插入值进行上升处理
//@auth      	hlccd		2021-07-14
//@receiver		n			*node					接受者node的指针
//@param    	num			int						插入后的元素数量,用于寻找最后一个父节点
//@param    	e			interface{}				待插入元素
//@param    	cmp			comparator.Comparator	比较器,在节点上升时使用
//@return    	nil
func (n *node) insert(num int, e interface{}, cmp comparator.Comparator) {
	if n == nil {
		return
	}
	//寻找最后一个父节点
	n = n.lastParent(num)
	//将元素插入最后一个节点
	if num%2 == 0 {
		n.left = newNode(n, e)
		n = n.left
	} else {
		n.right = newNode(n, e)
		n = n.right
	}
	//对插入的节点进行上升
	n.up(cmp)
}

//@title    up
//@description
//		以node节点做接收者
//		对该节点进行上升
//		当该节点存在且父节点存在时,若该节点小于夫节点
//		则在交换两个节点值后继续上升即可
//@auth      	hlccd		2021-07-14
//@receiver		n			*node					接受者node的指针
//@param    	cmp			comparator.Comparator	比较器,在节点上升时使用
//@return    	nil
func (n *node) up(cmp comparator.Comparator) {
	if n == nil {
		return
	}
	if n.parent == nil {
		return
	}
	//该节点和父节点都存在
	if cmp(n.parent.value, n.value) > 0 {
		//该节点值小于父节点值,交换两节点值,继续上升
		n.parent.value, n.value = n.value, n.parent.value
		n.parent.up(cmp)
	}
}

//@title    delete
//@description
//		以node节点做接收者
//		从删除该,并根据传入的num寻找最后一个父节点用于替换删除
//		随后对替换后的值进行下沉处理即可
//@auth      	hlccd		2021-07-14
//@receiver		n			*node					接受者node的指针
//@param    	num			int						删除前的元素数量,用于寻找最后一个父节点
//@param    	cmp			comparator.Comparator	比较器,在节点下沉时使用
//@return    	nil
func (n *node) delete(num int, cmp comparator.Comparator) {
	if n == nil {
		return
	}
	//寻找最后一个父节点
	ln := n.lastParent(num)
	if num%2 == 0 {
		n.value = ln.left.value
		ln.left = nil
	} else {
		n.value = ln.right.value
		ln.right = nil
	}
	//对交换后的节点进行下沉
	n.down(cmp)
}

//@title    down
//@description
//		以node节点做接收者
//		对该节点进行下沉
//		当该存在右节点且小于自身元素时,与右节点进行交换并继续下沉
//		否则当该存在左节点且小于自身元素时,与左节点进行交换并继续下沉
//		当左右节点都不存在或都大于自身时下沉停止
//@auth      	hlccd		2021-07-14
//@receiver		n			*node					接受者node的指针
//@param    	cmp			comparator.Comparator	比较器,在节点下沉时使用
//@return    	nil
func (n *node) down(cmp comparator.Comparator) {
	if n == nil {
		return
	}
	if n.right != nil && cmp(n.left.value, n.right.value) >= 0 {
		n.right.value, n.value = n.value, n.right.value
		n.right.down(cmp)
		return
	}
	if n.left != nil && cmp(n.value, n.left.value) >= 0 {
		n.left.value, n.value = n.value, n.left.value
		n.left.down(cmp)
		return
	}
}
