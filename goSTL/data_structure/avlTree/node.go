package avlTree

//@Title		avlTree
//@Description
//		平衡二叉树的节点
//		可通过节点实现平衡二叉树的添加删除
//		也可通过节点返回整个平衡二叉树的所有元素
//		增减节点后通过左右旋转的方式保持平衡二叉树的平衡
import (
	"github.com/hlccd/goSTL/utils/comparator"
)

//node树节点结构体
//该节点是平衡二叉树的树节点
//若该平衡二叉树允许重复则对节点num+1即可,否则对value进行覆盖
//平衡二叉树节点当左右子节点深度差超过1时进行左右旋转以实现平衡
type node struct {
	value interface{} //节点中存储的元素
	num   int         //该元素数量
	depth int         //该节点的深度
	left  *node       //左节点指针
	right *node       //右节点指针
}

//@title    newNode
//@description
//		新建一个平衡二叉树节点并返回
//		将传入的元素e作为该节点的承载元素
//		该节点的num和depth默认为1,左右子节点设为nil
//@receiver		nil
//@param    	e			interface{}				承载元素e
//@return    	n        	*node					新建的二叉搜索树节点的指针
func newNode(e interface{}) (n *node) {
	return &node{
		value: e,
		num:   1,
		depth: 1,
		left:  nil,
		right: nil,
	}
}

//@title    inOrder
//@description
//		以node平衡二叉树节点做接收者
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

//@title    getDepth
//@description
//		以node平衡二叉树节点做接收者
//		返回该节点的深度,节点不存在返回0
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	depth       int						该节点的深度
func (n *node) getDepth() (depth int) {
	if n == nil {
		return 0
	}
	return n.depth
}

//@title    max
//@description
//		返回a和b中较大的值
//@receiver		nil
//@param    	a			int 					带比较的值
//@param    	b			int 					带比较的值
//@return    	m        	int						两个值中较大的值
func max(a, b int) (m int) {
	if a > b {
		return a
	} else {
		return b
	}
}

//@title    leftRotate
//@description
//		以node平衡二叉树节点做接收者
//		将该节点向左节点方向转动,使右节点作为原来节点,并返回右节点
//		同时将右节点的左节点设为原节点的右节点
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m       	*node					旋转后的原节点
func (n *node) leftRotate() (m *node) {
	//左旋转
	headNode := n.right
	n.right = headNode.left
	headNode.left = n
	//更新结点高度
	n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
	headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
	return headNode
}

//@title    rightRotate
//@description
//		以node平衡二叉树节点做接收者
//		将该节点向右节点方向转动,使左节点作为原来节点,并返回左节点
//		同时将左节点的右节点设为原节点的左节点
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m       	*node					旋转后的原节点
func (n *node) rightRotate() (m *node) {
	//右旋转
	headNode := n.left
	n.left = headNode.right
	headNode.right = n
	//更新结点高度
	n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
	headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
	return headNode
}

//@title    rightLeftRotate
//@description
//		以node平衡二叉树节点做接收者
//		将原节点的右节点先进行右旋再将原节点进行左旋
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m       	*node					旋转后的原节点
func (n *node) rightLeftRotate() (m *node) {
	//右旋转,之后左旋转
	//以失衡点右结点先右旋转
	sonHeadNode := n.right.rightRotate()
	n.right = sonHeadNode
	//再以失衡点左旋转
	return n.leftRotate()
}

//@title    leftRightRotate
//@description
//		以node平衡二叉树节点做接收者
//		将原节点的左节点先进行左旋再将原节点进行右旋
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m       	*node					旋转后的原节点
func (n *node) leftRightRotate() (m *node) {
	//左旋转,之后右旋转
	//以失衡点左结点先左旋转
	sonHeadNode := n.left.leftRotate()
	n.left = sonHeadNode
	//再以失衡点左旋转
	return n.rightRotate()
}

//@title    getMin
//@description
//		以node平衡二叉树节点做接收者
//		返回n节点的最小元素的承载的元素和数量
//		即该节点父节点的后缀节点的元素和数量
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	e       	interface{}				前缀节点元素
func (n *node) getMin() (e interface{}, num int) {
	if n == nil {
		return nil, 0
	}
	if n.left == nil {
		return n.value, n.num
	} else {
		return n.left.getMin()
	}
}

//@title    adjust
//@description
//		以node平衡二叉树节点做接收者
//		对n节点进行旋转以保持节点左右子树平衡
//@receiver		n			*node					接受者node的指针
//@param    	nil
//@return    	m       	*node					调整后的n节点
func (n *node) adjust() (m *node) {
	if n.right.getDepth()-n.left.getDepth() >= 2 {
		//右子树高于左子树且高度差超过2,此时应当对n进行左旋
		if n.right.right.getDepth() > n.right.left.getDepth() {
			//由于右右子树高度超过右左子树,故可以直接左旋
			n = n.leftRotate()
		} else {
			//由于右右子树不高度超过右左子树
			//所以应该先右旋右子树使得右子树高度不超过左子树
			//随后n节点左旋
			n = n.rightLeftRotate()
		}
	} else if n.left.getDepth()-n.right.getDepth() >= 2 {
		//左子树高于右子树且高度差超过2,此时应当对n进行右旋
		if n.left.left.getDepth() > n.left.right.getDepth() {
			//由于左左子树高度超过左右子树,故可以直接右旋
			n = n.rightRotate()
		} else {
			//由于左左子树高度不超过左右子树
			//所以应该先左旋左子树使得左子树高度不超过右子树
			//随后n节点右旋
			n = n.leftRightRotate()
		}
	}
	return n
}

//@title    insert
//@description
//		以node平衡二叉树节点做接收者
//		从n节点中插入元素e
//		如果n节点中承载元素与e不同则根据大小从左右子树插入该元素
//		如果n节点与该元素相等,且允许重复值,则将num+1否则对value进行覆盖
//		插入成功返回true,插入失败或不允许重复插入返回false
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待插入元素
//@param    	isMulti		bool					是否允许重复?
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	b        	bool					是否插入成功?
func (n *node) insert(e interface{}, isMulti bool, cmp comparator.Comparator) (m *node, b bool) {
	//节点不存在,应该创建并插入二叉树中
	if n == nil {
		return newNode(e), true
	}
	if cmp(e, n.value) < 0 {
		//从左子树继续插入
		n.left, b = n.left.insert(e, isMulti, cmp)
		if b {
			//插入成功,对该节点进行平衡
			n = n.adjust()
		}
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		return n, b
	}
	if cmp(e, n.value) > 0 {
		//从右子树继续插入
		n.right, b = n.right.insert(e, isMulti, cmp)
		if b {
			//插入成功,对该节点进行平衡
			n = n.adjust()
		}
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		return n, b
	}
	//该节点元素与待插入元素相同
	if isMulti {
		//允许重复,数目+1
		n.num++
		return n, true
	}
	//不允许重复,对值进行覆盖
	n.value = e
	return n, false
}

//@title    erase
//@description
//		以node平衡二叉树节点做接收者
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
func (n *node) erase(e interface{}, cmp comparator.Comparator) (m *node, b bool) {
	if n == nil {
		//待删除值不存在,删除失败
		return n, false
	}
	if cmp(e, n.value) < 0 {
		//从左子树继续删除
		n.left, b = n.left.erase(e, cmp)
	} else if cmp(e, n.value) > 0 {
		//从右子树继续删除
		n.right, b = n.right.erase(e, cmp)
	} else if cmp(e, n.value) == 0 {
		//存在相同值,从该节点删除
		b = true
		if n.num > 1 {
			//有重复值,节点无需删除,直接-1即可
			n.num--
		} else {
			//该节点需要被删除
			if n.left != nil && n.right != nil {
				//找到该节点后继节点进行交换删除
				n.value, n.num = n.right.getMin()
				//从右节点继续删除,同时可以保证删除的节点必然无左节点
				n.right, b = n.right.erase(n.value, cmp)
			} else if n.left != nil {
				n = n.left
			} else {
				n = n.right
			}
		}
	}
	//当n节点仍然存在时,对其进行调整
	if n != nil {
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		n = n.adjust()
	}
	return n, b
}

//@title    count
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
func (n *node) count(e interface{}, isMulti bool, cmp comparator.Comparator) (num int) {
	if n == nil {
		return 0
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.count(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.count(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.num
}

//@title    find
//@description
//		以node二叉搜索树节点做接收者
//		从n节点中查找元素e并返回其存储的全部信息
//		如果n节点中承载元素与e不同则根据大小从左右子树查找该元素
//		如果n节点与该元素相等,则直接返回其信息即可
//		若未找到该索引信息则返回nil
//@receiver		n			*node					接受者node的指针
//@param    	e			interface{}				待查找元素
//@param    	isMulti		bool					是否允许重复?
//@param    	cmp			comparator.Comparator	判断大小的比较器
//@return    	ans			interface{}				待查找索引元素所指向的元素
func (n *node) find(e interface{}, isMulti bool, cmp comparator.Comparator) (ans interface{}) {
	if n == nil {
		return nil
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.count(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.count(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.value
}
