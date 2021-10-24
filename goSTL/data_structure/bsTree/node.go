package bsTree

//@Title		bsTree
//@Description
//		二叉搜索树的节点
//		可通过节点实现二叉搜索树的添加删除
//		也可通过节点返回整个二叉搜索树的所有元素

import "github.com/hlccd/goSTL/utils/comparator"

//node树节点结构体
//该节点是二叉搜索树的树节点
//若该二叉搜索树允许重复则对节点num+1即可,否则对value进行覆盖
//二叉搜索树节点不做平衡
type node struct {
	value interface{} //节点中存储的元素
	num   uint64      //该元素数量
	left  *node       //左节点指针
	right *node       //右节点指针
}

//@title    newNode
//@description
//		新建一个二叉搜索树节点并返回
//		将传入的元素e作为该节点的承载元素
//		该节点的num默认为1,左右子节点设为nil
//@receiver		nil
//@param    	e			interface{}				承载元素e
//@return    	n        	*node					新建的二叉搜索树节点的指针
func newNode(e interface{}) (n *node) {
	return &node{
		value: e,
		num:   1,
		left:  nil,
		right: nil,
	}
}

//@title    inOrder
//@description
//		以node二叉搜索树节点做接收者
//		以中缀序列返回节点集合
//		若允许重复存储则对于重复元素进行多次放入
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	es        	[]interface{}			以该节点为起点的中缀序列
func (n *node) inOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	if n.left != nil {
		es = append(es, n.left.inOrder()...)
	}
	for i := uint64(0); i < n.num; i++ {
		es = append(es, n.value)
	}
	if n.right != nil {
		es = append(es, n.right.inOrder()...)
	}
	return es
}

//@title    insert
//@description
//		以node二叉搜索树节点做接收者
//		从n节点中插入元素e
//		如果n节点中承载元素与e不同则根据大小从左右子树插入该元素
//		如果n节点与该元素相等,且允许重复值,则将num+1否则对value进行覆盖
//		插入成功返回true,插入失败或不允许重复插入返回false
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待插入元素
//@param    	isMulti		bool					是否允许重复?
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	b        	bool					是否插入成功?
func (n *node) insert(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
	if n == nil {
		return false
	}
	//n中承载元素小于e,从右子树继续插入
	if cmp(n.value, e) < 0 {
		if n.right == nil {
			//右子树为nil,直接插入右子树即可
			n.right = newNode(e)
			return true
		} else {
			return n.right.insert(e, isMulti, cmp)
		}
	}
	//n中承载元素大于e,从左子树继续插入
	if cmp(n.value, e) > 0 {
		if n.left == nil {
			//左子树为nil,直接插入左子树即可
			n.left = newNode(e)
			return true
		} else {
			return n.left.insert(e, isMulti, cmp)
		}
	}
	//n中承载元素等于e
	if isMulti {
		//允许重复
		n.num++
		return true
	}
	//不允许重复,直接进行覆盖
	n.value = e
	return false
}

//@title    delete
//@description
//		以node二叉搜索树节点做接收者
//		从n节点中删除元素e
//		如果n节点中承载元素与e不同则根据大小从左右子树删除该元素
//		如果n节点与该元素相等,且允许重复值,则将num-1否则直接删除该元素
//		删除时先寻找该元素的前缀节点,若不存在则寻找其后继节点进行替换
//		替换后删除该节点
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待删除元素
//@param    	isMulti		bool					是否允许重复?
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	b        	bool					是否删除成功?
func (n *node) delete(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
	if n == nil {
		return false
	}
	//n中承载元素小于e,从右子树继续删除
	if cmp(n.value, e) < 0 {
		if n.right == nil {
			//右子树为nil,删除终止
			return false
		}
		if cmp(e, n.right.value) == 0 && (!isMulti || n.right.num == 1) {
			if n.right.left == nil && n.right.right == nil {
				//右子树可直接删除
				n.right = nil
				return true
			}
		}
		//从右子树继续删除
		return n.right.delete(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续删除
	if cmp(n.value, e) > 0 {
		if n.left == nil {
			//左子树为nil,删除终止
			return false
		}
		if cmp(e, n.left.value) == 0 && (!isMulti || n.left.num == 1) {
			if n.left.left == nil && n.left.right == nil {
				//左子树可直接删除
				n.left = nil
				return true
			}
		}
		//从左子树继续删除
		return n.left.delete(e, isMulti, cmp)
	}
	//n中承载元素等于e
	if (*n).num > 1 && isMulti {
		//允许重复且个数超过1,则减少num即可
		(*n).num--
		return true
	}
	if n.left == nil && n.right == nil {
		//该节点无前缀节点和后继节点,删除即可
		*(&n) = nil
		return true
	}
	if n.left != nil {
		//该节点有前缀节点,找到前缀节点进行交换并删除即可
		ln := n.left
		if ln.right == nil {
			n.value = ln.value
			n.num = ln.num
			n.left = ln.left
		} else {
			for ln.right.right != nil {
				ln = ln.right
			}
			n.value = ln.right.value
			n.num = ln.right.num
			ln.right = ln.right.left
		}
	} else if (*n).right != nil {
		//该节点有后继节点,找到后继节点进行交换并删除即可
		tn := n.right
		if tn.left == nil {
			n.value = tn.value
			n.num = tn.num
			n.right = tn.right
		} else {
			for tn.left.left != nil {
				tn = tn.left
			}
			n.value = tn.left.value
			n.num = tn.left.num
			tn.left = tn.left.right
		}
		return true
	}
	return true
}

//@title    search
//@description
//		以node二叉搜索树节点做接收者
//		从n节点中查找元素e并返回存储的个数
//		如果n节点中承载元素与e不同则根据大小从左右子树查找该元素
//		如果n节点与该元素相等,则直接返回其个数
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待查找元素
//@param    	isMulti		bool					是否允许重复?
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	num        	uint64						待查找元素在二叉树中存储的数量
func (n *node) search(e interface{}, isMulti bool, cmp comparator.Comparator) (num uint64) {
	if n == nil {
		return 0
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.search(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.search(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.num
}
