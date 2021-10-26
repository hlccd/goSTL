package cbTree

//@Title		cbTree
//@Description
//		完全二叉树-Complete Binary Tree
//		以完全二叉树的形式实现的堆
//		通过比较器判断传入元素的大小
//		将最小的元素放在堆顶
//		该结构只保留整个树的根节点,其他节点通过根节点进行查找获得
//@author     	hlccd		2021-07-14
//@update		hlccd 		2021-08-01		增加互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//cbTree完全二叉树树结构体
//该实例存储二叉树的根节点
//同时保存该二叉树已经存储了多少个元素
//二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找
type cbTree struct {
	root  *node                 //根节点指针
	size  int                   //存储元素数量
	cmp   comparator.Comparator //比较器
	mutex sync.Mutex            //并发控制锁
}

//cbTree二叉搜索树容器接口
//存放了cbTree二叉搜索树可使用的函数
//对应函数介绍见下方
type cbTreer interface {
	Iterator() (i *iterator.Iterator) //返回包含该二叉树的所有元素
	Size() (num int)                  //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Push(e interface{})               //向二叉树中插入元素e
	Pop()                             //从二叉树中弹出顶部元素
	Top() (e interface{})             //返回该二叉树的顶部元素
}

//@title    New
//@description
//		新建一个cbTree完全二叉树容器并返回
//		初始根节点为nil
//		若有传入的比较器,则将传入的第一个比较器设为该二叉树的比较器
//@author     	hlccd		2021-07-14
//@receiver		nil
//@param    	Cmp			 ...comparator.Comparator	cbTree比较器集
//@return    	cb        	*cbTree						新建的cbTree指针
func New(Cmp ...comparator.Comparator) (cb *cbTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	}
	return &cbTree{
		root:  nil,
		size:  0,
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以cbTree完全二叉树做接收者
//		将该二叉树中所有保存的元素将从根节点开始以前缀序列的形式放入迭代器中
//@auth      	hlccd		2021-07-14
//@receiver		cb			*cbTree					接受者cbTree的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (cb *cbTree) Iterator() (i *iterator.Iterator) {
	if cb == nil {
		return iterator.New(make([]interface{}, 0, 0))
	}
	cb.mutex.Lock()
	i = iterator.New(cb.root.frontOrder())
	cb.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以cbTree完全二叉树做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回-1
//@auth      	hlccd		2021-07-14
//@receiver		cb			*cbTree					接受者cbTree的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (cb *cbTree) Size() (num int) {
	if cb == nil {
		return -1
	}
	return cb.size
}

//@title    Clear
//@description
//		以cbTree完全二叉树做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@auth      	hlccd		2021-07-14
//@receiver		cb			*cbTree					接受者cbTree的指针
//@param    	nil
//@return    	nil
func (cb *cbTree) Clear() {
	if cb == nil {
		return
	}
	cb.mutex.Lock()
	cb.root = nil
	cb.size = 0
	cb.mutex.Unlock()
}

//@title    Empty
//@description
//		以cbTree完全二叉树做接收者
//		判断该完全二叉树树是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@auth      	hlccd		2021-07-14
//@receiver		cb			*cbTree					接受者cbTree的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (cb *cbTree) Empty() (b bool) {
	if cb == nil {
		return true
	}
	if cb.size > 0 {
		return false
	}
	return true
}

//@title    Push
//@description
//		以cbTree完全二叉树做接收者
//		向二叉树插入元素e,将其放入完全二叉树的最后一个位置,随后进行元素上升
//		如果二叉树本身为空,则直接将根节点设为插入节点元素即可
//@auth      	hlccd		2021-07-14
//@receiver		cb			*cbTree					接受者cbTree的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (cb *cbTree) Push(e interface{}) {
	if cb == nil {
		return
	}
	cb.mutex.Lock()
	if cb.Empty() {
		if cb.cmp == nil {
			cb.cmp = comparator.GetCmp(e)
		}
		if cb.cmp == nil {
			cb.mutex.Unlock()
			return
		}
		cb.root = newNode(nil, e)
		cb.size++
	} else {
		cb.size++
		cb.root.insert(cb.size, e, cb.cmp)
	}
	cb.mutex.Unlock()
}

//@title    Pop
//@description
//		以cbTree完全二叉树做接收者
//		从完全二叉树中删除顶部元素e
//		将该顶部元素于最后一个元素进行交换
//		随后删除最后一个元素
//		再将顶部元素进行下沉处理即可
//@auth      	hlccd		2021-07-14
//@receiver		cb			*cbTree					接受者cbTree的指针
//@param    	nil
//@return    	nil
func (cb *cbTree) Pop() {
	if cb == nil {
		return
	}
	if cb.Empty() {
		return
	}
	cb.mutex.Lock()
	if cb.size == 1 {
		//该二叉树仅剩根节点,直接删除即可
		cb.root = nil
	} else {
		//该二叉树删除根节点后还有其他节点可生为跟节点
		cb.root.delete(cb.size, cb.cmp)
	}
	cb.size--
	cb.mutex.Unlock()
}

//@title    Top
//@description
//		以cbTree完全二叉树做接收者
//		返回该完全二叉树的顶部元素
//		当该完全二叉树不存在或根节点不存在时返回nil
//@auth      	hlccd		2021-07-14
//@receiver		cb			*cbTree					接受者cbTree的指针
//@param    	nil
//@return    	e 			interface{}				该完全二叉树的顶部元素
func (cb *cbTree) Top() (e interface{}) {
	if cb == nil {
		return nil
	}
	cb.mutex.Lock()
	e = cb.root.value
	cb.mutex.Unlock()
	return e
}
