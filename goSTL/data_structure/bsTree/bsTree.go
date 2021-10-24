package bsTree

//@Title		bsTree
//@Description
//		二叉搜索树-Binary Search Tree
//		以二叉树的形式实现
//		二叉树实例保存根节点和比较器以及保存的数量
//		可以在创建时设置节点是否可重复
//		若节点可重复则增加节点中的数值,否则对节点存储元素进行覆盖
//		二叉搜索树不进行平衡
//		增加互斥锁实现并发控制

import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//bsTree二叉搜索树结构体
//该实例存储二叉树的根节点
//同时保存该二叉树已经存储了多少个元素
//二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找
//创建时传入是否允许该二叉树出现重复值,如果不允许则进行覆盖,允许则对节点数目增加即可
type bsTree struct {
	root    *node                 //根节点指针
	size    uint64                //存储元素数量
	cmp     comparator.Comparator //比较器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}

//bsTree二叉搜索树容器接口
//存放了bsTree二叉搜索树可使用的函数
//对应函数介绍见下方
type bsTreeer interface {
	Iterator() (i *Iterator.Iterator) //返回包含该二叉树的所有元素,重复则返回多个
	Size() (num uint64)               //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Insert(e interface{})             //向二叉树中插入元素e
	Erase(e interface{})              //从二叉树中删除元素e
	Count(e interface{}) (num uint64) //从二叉树中寻找元素e并返回其个数
}

//@title    New
//@description
//		新建一个bsTree二叉搜索树容器并返回
//		初始根节点为nil
//		传入该二叉树是否为可重复属性,如果为true则保存重复值,否则对原有相等元素进行覆盖
//		若有传入的比较器,则将传入的第一个比较器设为该二叉树的比较器
//@receiver		nil
//@param    	isMulti		bool						该二叉树是否保存重复值?
//@param    	Cmp			 ...comparator.Comparator	bsTree比较器集
//@return    	bs        	*bsTree						新建的bsTree指针
func New(isMulti bool, Cmp ...comparator.Comparator) (bs *bsTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(Cmp) == 0 {
		cmp = nil
	} else {
		cmp = Cmp[0]
	}
	return &bsTree{
		root:    nil,
		size:    0,
		cmp:     cmp,
		isMulti: isMulti,
		mutex:   sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以bsTree二叉搜索树做接收者
//		将该二叉树中所有保存的元素将从根节点开始以中缀序列的形式放入迭代器中
//		若允许重复存储则对于重复元素进行多次放入
//@receiver		bt			*bsTree					接受者bsTree的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (bs *bsTree) Iterator() (i *Iterator.Iterator) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	es := bs.root.inOrder()
	i = Iterator.New(&es)
	bs.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以bsTree二叉搜索树做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil则创建一个并返回其承载的元素个数
//@receiver		bt			*bsTree					接受者bsTree的指针
//@param    	nil
//@return    	num        	uint64					容器中实际使用元素所占空间大小
func (bs *bsTree) Size() (num uint64) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	return bs.size
}

//@title    Clear
//@description
//		以bsTree二叉搜索树做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@receiver		bt			*bsTree					接受者bsTree的指针
//@param    	nil
//@return    	nil
func (bs *bsTree) Clear() {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	bs.root = nil
	bs.size = 0
	bs.mutex.Unlock()
}

//@title    Empty
//@description
//		以bsTree二叉搜索树做接收者
//		判断该二叉搜索树是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		bt			*bsTree					接受者bsTree的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (bs *bsTree) Empty() (b bool) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	return bs.size == 0
}

//@title    Insert
//@description
//		以bsTree二叉搜索树做接收者
//		向二叉树插入元素e,若不允许重复则对相等元素进行覆盖
//		如果二叉树为空则之间用根节点承载元素e,否则以根节点开始进行查找
//		不做平衡
//@receiver		bt			*bsTree					接受者bsTree的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (bs *bsTree) Insert(e interface{}) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	if bs.Empty() {
		//二叉树为空,用根节点承载元素e
		if bs.cmp == nil {
			bs.cmp = comparator.GetCmp(e)
		}
		if bs.cmp == nil {
			bs.mutex.Unlock()
			return
		}
		bs.root = newNode(e)
		bs.size++
		bs.mutex.Unlock()
		return
	}
	//二叉树不为空,从根节点开始查找添加元素e
	if bs.root.insert(e, bs.isMulti, bs.cmp) {
		bs.size++
	}
	bs.mutex.Unlock()
}

//@title    Erase
//@description
//		以bsTree二叉搜索树做接收者
//		从搜素二叉树中删除元素e
//		若允许重复记录则对承载元素e的节点中数量记录减一即可
//		若不允许重复记录则删除该节点同时将前缀节点或后继节点更换过来以保证二叉树的不发送断裂
//		如果该二叉树仅持有一个元素且根节点等价于待删除元素,则将二叉树根节点置为nil
//@receiver		bt			*bsTree					接受者bsTree的指针
//@param    	e			interface{}				待删除元素
//@return    	nil
func (bs *bsTree) Erase(e interface{}) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	if bs.size == 0 {
		return
	}
	bs.mutex.Lock()
	if bs.size == 1 && bs.cmp(bs.root.value, e) == 0 {
		//二叉树仅持有一个元素且根节点等价于待删除元素,将二叉树根节点置为nil
		bs.root = nil
		bs.size = 0
		bs.mutex.Unlock()
		return
	}
	//从根节点开始删除元素e
	//如果删除成功则将size-1
	if bs.root.delete(e, bs.isMulti, bs.cmp) {
		bs.size--
	}
	bs.mutex.Unlock()
}

//@title    Count
//@description
//		以bsTree二叉搜索树做接收者
//		从搜素二叉树中查找元素e的个数
//		如果找到则返回该二叉树中和元素e相同元素的个数
//		如果不允许重复则最多返回1
//		如果未找到则返回0
//@receiver		bt			*bsTree					接受者bsTree的指针
//@param    	e			interface{}				待查找元素
//@return    	num			uint64					待查找元素在二叉树中存储的个数
func (bs *bsTree) Count(e interface{}) (num uint64) {
	if bs == nil {
		//二叉树不存在,返回0
		return 0
	}
	if bs.Empty() {
		//二叉树为空,返回0
		return 0
	}
	bs.mutex.Lock()
	//从根节点开始查找并返回查找结果
	num = bs.root.search(e, bs.isMulti, bs.cmp)
	bs.mutex.Unlock()
	return num
}
