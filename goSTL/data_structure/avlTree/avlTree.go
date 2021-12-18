package avlTree

//@Title		avlTree
//@Description
//		平衡二叉树-Balanced Binary Tree
//		以二叉树的形式实现
//		平衡二叉树实例保存根节点和比较器以及保存的数量
//		可以在创建时设置节点是否可重复
//		若节点可重复则增加节点中的数值,否则对节点存储元素进行覆盖
//		平衡二叉树在添加和删除时都将对节点进行平衡,以保证一个节点的左右子节点高度差不超过1
//		使用互斥锁实现并发控制

import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//avlTree平衡二叉树结构体
//该实例存储平衡二叉树的根节点
//同时保存该二叉树已经存储了多少个元素
//二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找
//创建时传入是否允许该二叉树出现重复值,如果不允许则进行覆盖,允许则对节点数目增加即可
type AvlTree struct {
	root    *node                 //根节点指针
	size    int                   //存储元素数量
	cmp     comparator.Comparator //比较器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}

//avlTree平衡二叉树容器接口
//存放了avlTree平衡二叉树可使用的函数
//对应函数介绍见下方
type avlTreer interface {
	Iterator() (i *Iterator.Iterator) //返回包含该二叉树的所有元素,重复则返回多个
	Size() (num int)                  //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Insert(e interface{}) (b bool)            //向二叉树中插入元素e
	Erase(e interface{}) (b bool)     //从二叉树中删除元素e
	Count(e interface{}) (num int)    //从二叉树中寻找元素e并返回其个数
}

//@title    New
//@description
//		新建一个avlTree平衡二叉树容器并返回
//		初始根节点为nil
//		传入该二叉树是否为可重复属性,如果为true则保存重复值,否则对原有相等元素进行覆盖
//		若有传入的比较器,则将传入的第一个比较器设为该二叉树的比较器
//@receiver		nil
//@param    	isMulti		bool						该二叉树是否保存重复值?
//@param    	Cmp			 ...comparator.Comparator	avlTree比较器集
//@return    	avl        	*avlTree						新建的avlTree指针
func New(isMulti bool, cmps ...comparator.Comparator) (avl *AvlTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	return &AvlTree{
		root:    nil,
		size:    0,
		cmp:     cmp,
		isMulti: isMulti,
	}
}

//@title    Iterator
//@description
//		以avlTree平衡二叉树做接收者
//		将该二叉树中所有保存的元素将从根节点开始以中缀序列的形式放入迭代器中
//		若允许重复存储则对于重复元素进行多次放入
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (avl *AvlTree) Iterator() (i *Iterator.Iterator) {
	if avl == nil {
		return nil
	}
	avl.mutex.Lock()
	es := avl.root.inOrder()
	i = Iterator.New(&es)
	avl.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以avlTree平衡二叉树做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回0
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (avl *AvlTree) Size() (num int) {
	if avl == nil {
		return 0
	}
	return avl.size
}

//@title    Clear
//@description
//		以avlTree平衡二叉树做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	nil
//@return    	nil
func (avl *AvlTree) Clear() {
	if avl == nil {
		return
	}
	avl.mutex.Lock()
	avl.root = nil
	avl.size = 0
	avl.mutex.Unlock()
}

//@title    Empty
//@description
//		以avlTree平衡二叉树做接收者
//		判断该二叉搜索树是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (avl *AvlTree) Empty() (b bool) {
	if avl == nil {
		return true
	}
	if avl.size > 0 {
		return false
	}
	return true
}

//@title    Insert
//@description
//		以avlTree平衡二叉树做接收者
//		向二叉树插入元素e,若不允许重复则对相等元素进行覆盖
//		如果二叉树为空则之间用根节点承载元素e,否则以根节点开始进行查找
//		当节点左右子树高度差超过1时将进行旋转以保持平衡
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	e			interface{}				待插入元素
//@return    	b			bool					添加成功?
func (avl *AvlTree) Insert(e interface{}) (b bool){
	if avl == nil {
		return false
	}
	avl.mutex.Lock()
	if avl.Empty() {
		if avl.cmp == nil {
			avl.cmp = comparator.GetCmp(e)
		}
		if avl.cmp == nil {
			return
		}
		//二叉树为空,用根节点承载元素e
		avl.root = newNode(e)
		avl.size = 1
		avl.mutex.Unlock()
		return true
	}
	//从根节点进行插入,并返回节点,同时返回是否插入成功
	avl.root, b = avl.root.insert(e, avl.isMulti, avl.cmp)
	if b {
		//插入成功,数量+1
		avl.size++
	}
	avl.mutex.Unlock()
	return b
}

//@title    Erase
//@description
//		以avlTree平衡二叉树做接收者
//		从平衡二叉树中删除元素e
//		若允许重复记录则对承载元素e的节点中数量记录减一即可
//		若不允许重复记录则删除该节点同时将前缀节点或后继节点更换过来以保证二叉树的不发送断裂
//		如果该二叉树仅持有一个元素且根节点等价于待删除元素,则将二叉树根节点置为nil
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	e			interface{}				待删除元素
//@return    	b			bool					删除成功
func (avl *AvlTree) Erase(e interface{}) (b bool) {
	if avl == nil {
		return false
	}
	if avl.Empty() {
		return false
	}
	avl.mutex.Lock()
	if avl.size == 1 && avl.cmp(avl.root.value, e) == 0 {
		//二叉树仅持有一个元素且根节点等价于待删除元素,将二叉树根节点置为nil
		avl.root = nil
		avl.size = 0
		avl.mutex.Unlock()
		return true
	}
	//从根节点进行插入,并返回节点,同时返回是否删除成功
	avl.root, b = avl.root.erase(e, avl.cmp)
	if b {
		avl.size--
	}
	avl.mutex.Unlock()
	return b
}

//@title    Count
//@description
//		以avlTree平衡二叉树做接收者
//		从搜素二叉树中查找元素e的个数
//		如果找到则返回该二叉树中和元素e相同元素的个数
//		如果不允许重复则最多返回1
//		如果未找到则返回0
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	e			interface{}				待查找元素
//@return    	num			int						待查找元素在二叉树中存储的个数
func (avl *AvlTree) Count(e interface{}) (num int) {
	if avl == nil {
		//二叉树为空,返回0
		return 0
	}
	if avl.Empty() {
		return 0
	}
	avl.mutex.Lock()
	num = avl.root.count(e, avl.isMulti, avl.cmp)
	avl.mutex.Unlock()
	return num
}

//@title    Find
//@description
//		以avlTree平衡二叉树做接收者
//		从搜素二叉树中查找以元素e为索引信息的全部信息
//		如果找到则返回该二叉树中和索引元素e相同的元素的全部信息
//		如果未找到则返回nil
//@receiver		avl			*avlTree				接受者avlTree的指针
//@param    	e			interface{}				待查找索引元素
//@return    	ans			interface{}				待查找索引元素所指向的元素
func (avl *AvlTree) Find(e interface{}) (ans interface{}) {
	if avl == nil {
		//二叉树为空,返回0
		return 0
	}
	if avl.Empty() {
		return 0
	}
	avl.mutex.Lock()
	ans = avl.root.find(e, avl.isMulti, avl.cmp)
	avl.mutex.Unlock()
	return ans
}
