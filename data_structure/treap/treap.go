package treap

//@Title		treap
//@Description
//		Treap树堆容器包
//		树堆本身是一个二叉树,同时赋予随机的节点优先级
//		通过旋转使树堆中节点既满足存储元素的组成符合二叉搜索树的性质,同时也使得优先级满足堆的的性质
//		同时由于每个节点的优先级随机生成,使得整个二叉树得以实现随机平衡
//		该树堆依概率实现平衡
//		可接纳不同类型的元素,但建议在同一个树堆中使用相同类型的元素
//		配合比较器实现元素之间的大小比较

import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"math/rand"
	"sync"
	"time"
)

//treap树堆结构体
//该实例存储树堆的根节点
//同时保存该树堆中已经存储了多少个元素
//二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找
//该树堆实例中存储随机数生成器,用于后续新建节点时生成随机数
//创建时传入是否允许该树堆出现重复值,如果不允许则进行覆盖,允许则对节点数目增加即可
type treap struct {
	root    *node                 //根节点指针
	size    int                   //存储元素数量
	cmp     comparator.Comparator //比较器
	rand    *rand.Rand            //随机数生成器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}

//treap树堆容器接口
//存放了treap树堆可使用的函数
//对应函数介绍见下方
type treaper interface {
	Iterator() (i *Iterator.Iterator) //返回包含该树堆的所有元素,重复则返回多个
	Size() (num int)                  //返回该树堆中保存的元素个数
	Clear()                           //清空该树堆
	Empty() (b bool)                  //判断该树堆是否为空
	Insert(e interface{})             //向树堆中插入元素e
	Erase(e interface{})              //从树堆中删除元素e
	Count(e interface{}) (num int)    //从树堆中寻找元素e并返回其个数
}

//@title    New
//@description
//		新建一个treap树堆容器并返回
//		初始根节点为nil
//		传入该树堆是否为可重复属性,如果为true则保存重复值,否则对原有相等元素进行覆盖
//		若有传入的比较器,则将传入的第一个比较器设为该树堆的比较器
//@receiver		nil
//@param    	isMulti		bool						该树堆是否保存重复值?
//@param    	Cmp			 ...comparator.Comparator	treap比较器集
//@return    	t        	*treap						新建的treap指针
func New(isMulti bool, Cmp ...comparator.Comparator) (t *treap) {
	//设置默认比较器
	var cmp comparator.Comparator
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	}
	//创建随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &treap{
		root:    nil,
		size:    0,
		cmp:     cmp,
		rand:    r,
		isMulti: isMulti,
		mutex:   sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以treap树堆做接收者
//		将该树堆中所有保存的元素将从根节点开始以中缀序列的形式放入迭代器中
//		若允许重复存储则对于重复元素进行多次放入
//@receiver		t			*treap					接受者treap的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (t *treap) Iterator() (i *Iterator.Iterator) {
	if t == nil {
		return nil
	}
	t.mutex.Lock()
	es := t.root.inOrder()
	i = Iterator.New(&es)
	t.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以treap树堆做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回0
//@receiver		t			*treap					接受者treap的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (t *treap) Size() (num int) {
	if t == nil {
		return 0
	}
	return t.size
}
//@title    Clear
//@description
//		以treap树堆做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@receiver		t			*treap					接受者treap的指针
//@param    	nil
//@return    	nil
func (t *treap) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root = nil
	t.size = 0
	t.mutex.Unlock()
}
//@title    Empty
//@description
//		以treap树堆做接收者
//		判断该二叉搜索树是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		t			*treap					接受者treap的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (t *treap) Empty() (b bool) {
	if t == nil {
		return true
	}
	if t.size > 0 {
		return false
	}
	return true
}

//@title    Insert
//@description
//		以treap树堆做接收者
//		向二叉树插入元素e,若不允许重复则对相等元素进行覆盖
//		如果二叉树为空则之间用根节点承载元素e,否则以根节点开始进行查找
//		对于该树堆来说,通过赋予随机的优先级根据堆的性质来实现平衡
//@receiver		t			*treap					接受者treap的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (t *treap) Insert(e interface{}) {
	//判断容器是否存在
	if t == nil {
		return
	}
	t.mutex.Lock()
	if t.Empty() {
		//判断比较器是否存在
		if t.cmp == nil {
			t.cmp = comparator.GetCmp(e)
		}
		if t.cmp == nil {
			t.mutex.Unlock()
			return
		}
		//插入到根节点
		t.root = newNode(e, t.rand)
		t.size = 1
		t.mutex.Unlock()
		return
	}
	//从根节点向下插入
	if t.root.insert(newNode(e, t.rand), t.isMulti, t.cmp) {
		t.size++
	}
	t.mutex.Unlock()
}


//@title    Erase
//@description
//		以treap树堆做接收者
//		从树堆中删除元素e
//		若允许重复记录则对承载元素e的节点中数量记录减一即可
//		若不允许重复记录则删除该节点同时将前缀节点或后继节点更换过来以保证树堆的不发送断裂
//		交换后根据优先级进行左右旋转以保证符合堆的性质
//		如果该树堆仅持有一个元素且根节点等价于待删除元素,则将根节点置为nil
//@receiver		t			*treap					接受者treap的指针
//@param    	e			interface{}				待删除元素
//@return    	nil
func (t *treap) Erase(e interface{}) {
	if t == nil {
		return
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	if t.size == 1 && t.cmp(t.root.value, e) == 0 {
		//该树堆仅持有一个元素且根节点等价于待删除元素,则将根节点置为nil
		t.root = nil
		t.size = 0
		t.mutex.Unlock()
		return
	}
	//从根节点开始删除元素
	if t.root.delete(e, t.isMulti, t.cmp) {
		//删除成功
		t.size--
	}
	t.mutex.Unlock()
}

//@title    Count
//@description
//		以treap树堆做接收者
//		从树堆中查找元素e的个数
//		如果找到则返回该树堆中和元素e相同元素的个数
//		如果不允许重复则最多返回1
//		如果未找到则返回0
//@receiver		t			*treap					接受者treap的指针
//@param    	e			interface{}				待查找元素
//@return    	num			int						待查找元素在树堆中存储的个数
func (t *treap) Count(e interface{}) (num int) {
	if t == nil {
		//树堆不存在,直接返回0
		return 0
	}
	if t.Empty() {
		return
	}
	t.mutex.Lock()
	num = t.root.search(e, t.cmp)
	t.mutex.Unlock()
	//树堆存在,从根节点开始查找该元素
	return num
}