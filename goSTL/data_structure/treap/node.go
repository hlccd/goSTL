package treap
//@Title		Treap
//@Description
//		树堆的节点
//		节点在创建是赋予一个随机的优先级,随后进行堆平衡,使得整个树堆依概率实现平衡
//		可通过节点实现树堆的添加删除
//		也可通过节点返回整个二叉搜索树的所有元素
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"math/rand"
)

//node树节点结构体
//该节点是树堆的树节点
//若该树堆允许重复则对节点num+1即可,否则对value进行覆盖
//树堆节点将针对堆的性质通过左右旋转的方式做平衡
type node struct {
	value    interface{} //节点中存储的元素
	priority uint32      //该节点的优先级,随机生成
	num      int         //该节点中存储的数量
	left     *node       //左节点指针
	right    *node       //右节点指针
}

//@title    newNode
//@description
//		新建一个树堆节点并返回
//		将传入的元素e作为该节点的承载元素
//		该节点的num默认为1,左右子节点设为nil
//		该节点优先级随机生成,范围在0~2^16内
//@receiver		nil
//@param    	e			interface{}				承载元素e
//@param    	rand		*rand.Rand				随机数生成器
//@return    	n        	*node					新建的树堆树节点的指针
func newNode(e interface{}, rand *rand.Rand) (n *node) {
	return &node{
		value:    e,
		priority: uint32(rand.Intn(4294967295)),
		num:      1,
		left:     nil,
		right:    nil,
	}
}

//@title    inOrder
//@description
//		以node树堆节点做接收者
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
	for i := 0; i < n.num; i++ {
		es = append(es, n.value)
	}
	if n.right != nil {
		es = append(es, n.right.inOrder()...)
	}
	return es
}


//@title    rightRotate
//@description
//		以node树堆节点做接收者
//		新建一个节点作为n节点的右节点,同时将n节点的数值放入新建节点中作为右转后的n节点
//		右转后的n节点的左节点是原n节点左节点的右节点,右转后的右节点保持不变
//		原n节点改为原n节点的左节点,同时右节点指向新建的节点即右转后的n节点
//		该右转方式可以保证n节点的双亲节点不用更换节点指向
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	nil
func (n *node) rightRotate() {
	if n == nil {
		return
	}
	if n.left == nil {
		return
	}
	//新建节点作为更换后的n节点
	tmp := &node{
		value:    n.value,
		priority: n.priority,
		num:      n.num,
		left:     n.left.right,
		right:    n.right,
	}
	//原n节点左节点上移到n节点位置
	n.right = tmp
	n.value = n.left.value
	n.priority = n.left.priority
	n.num = n.left.num
	n.left = n.left.left
}

//@title    leftRotate
//@description
//		以node树堆节点做接收者
//		新建一个节点作为n节点的左节点,同时将n节点的数值放入新建节点中作为左转后的n节点
//		左转后的n节点的右节点是原n节点右节点的左节点,左转后的左节点保持不变
//		原n节点改为原n节点的右节点,同时左节点指向新建的节点即左转后的n节点
//		该左转方式可以保证n节点的双亲节点不用更换节点指向
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	nil
func (n *node) leftRotate() {
	if n == nil {
		return
	}
	if n.right == nil {
		return
	}
	//新建节点作为更换后的n节点
	tmp := &node{
		value:    n.value,
		priority: n.priority,
		num:      n.num,
		left:     n.left,
		right:    n.right.left,
	}
	//原n节点右节点上移到n节点位置
	n.left = tmp
	n.value = n.right.value
	n.priority = n.right.priority
	n.num = n.right.num
	n.right = n.right.right
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
func (n *node) insert(e *node, isMulti bool, cmp comparator.Comparator) (b bool) {
	if cmp(n.value, e.value) > 0 {
		if n.left == nil {
			//将左节点直接设为e
			n.left = e
			b = true
		} else {
			//对左节点进行递归插入
			b = n.left.insert(e, isMulti, cmp)
		}
		if n.priority > e.priority {
			//对n节点进行右转
			n.rightRotate()
		}
		return b
	} else if cmp(n.value, e.value) < 0 {
		if n.right == nil {
			//将右节点直接设为e
			n.right = e
			b = true
		} else {
			//对右节点进行递归插入
			b = n.right.insert(e, isMulti, cmp)
		}
		if n.priority > e.priority {
			//对n节点进行左转
			n.leftRotate()
		}
		return b
	}
	if isMulti {
		//允许重复
		n.num++
		return true
	}
	//不允许重复,对值进行覆盖
	n.value = e.value
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
			//待删除节点无子节点,直接删除即可
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
			//待删除节点无子节点,直接删除即可
			if n.left.left == nil && n.left.right == nil {
				//左子树可直接删除
				n.left = nil
				return true
			}
		}
		//从左子树继续删除
		return n.left.delete(e, isMulti, cmp)
	}
	if isMulti && n.num > 1 {
		//允许重复且数量超过1
		n.num--
		return true
	}
	//删除该节点
	tmp := n
	//左右子节点都存在则选择优先级较小一个进行旋转
	for tmp.left != nil && tmp.right != nil {
		if tmp.left.priority < tmp.right.priority {
			tmp.rightRotate()
			if tmp.right.left == nil && tmp.right.right == nil {
				tmp.right = nil
				return false
			}
			tmp = tmp.right
		} else {
			tmp.leftRotate()
			if tmp.left.left == nil && tmp.left.right == nil {
				tmp.left = nil
				return false
			}
			tmp = tmp.left
		}
	}
	if tmp.left == nil && tmp.right != nil {
		//到左子树为nil时直接换为右子树即可
		tmp.value = tmp.right.value
		tmp.num = tmp.right.num
		tmp.priority = tmp.right.priority
		tmp.left = tmp.right.left
		tmp.right = tmp.right.right
	} else if tmp.right == nil && tmp.left != nil {
		//到右子树为nil时直接换为左子树即可
		tmp.value = tmp.left.value
		tmp.num = tmp.left.num
		tmp.priority = tmp.left.priority
		tmp.right = tmp.left.right
		tmp.left = tmp.left.left
	}
	//当左右子树都为nil时直接结束
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
//@return    	num        	int						待查找元素在二叉树中存储的数量
func (n *node) search(e interface{}, cmp comparator.Comparator) (num int) {
	if n == nil {
		return 0
	}
	if cmp(n.value, e) > 0 {
		return n.left.search(e, cmp)
	} else if cmp(n.value, e) < 0 {
		return n.right.search(e, cmp)
	}
	return n.num
}