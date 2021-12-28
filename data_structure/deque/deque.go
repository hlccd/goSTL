package deque

//@Title		deque
//@Description
//		deque双队列容器包
//		区别于queue的动态数组实现方式,deque采取将数组和链表相结合的方案
//		该容器既可以在首部增删元素,也可以在尾部增删元素
//		deque在扩容和缩容时,都是固定增加2^10的空间,同时这一部分空间形成链表节点,并串成链表去保存
//		可接纳不同类型的元素
//		通过并发控制锁保证了在高并发过程中的数据一致性

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//deque双向队列结构体
//包含链表的头尾节点指针
//当删除节点时通过头尾节点指针进入链表进行删除
//当一个节点全部被删除后则释放该节点,同时首尾节点做相应调整
//当添加节点时若未占满节点空间时移动下标并做覆盖即可
//当添加节点时空间已使用完毕时,根据添加位置新建一个新节点补充上去
type Deque struct {
	first *node      //链表首节点指针
	last  *node      //链表尾节点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}

//deque双向队列容器接口
//存放了deque容器可使用的函数
//对应函数介绍见下方

type dequer interface {
	Iterator() (i *Iterator.Iterator) //返回包含双向队列中所有元素的迭代器
	Size() (size uint64)              //返回该双向队列中元素的使用空间大小
	Clear()                           //清空该双向队列
	Empty() (b bool)                  //判断该双向队列是否为空
	PushFront(e interface{})          //将元素e添加到该双向队列的首部
	PushBack(e interface{})           //将元素e添加到该双向队列的尾部
	PopFront() (e interface{})        //将该双向队列首元素弹出
	PopBack() (e interface{})         //将该双向队列首元素弹出
	Front() (e interface{})           //获取该双向队列首部元素
	Back() (e interface{})            //获取该双向队列尾部元素
}

//@title    New
//@description
//		新建一个deque双向队列容器并返回
//		初始deque双向队列的链表首尾节点为nil
//		初始size为0
//@receiver		nil
//@param    	nil
//@return    	d        	*Deque					新建的deque指针
func New() *Deque {
	return &Deque{
		first: nil,
		last:  nil,
		size:  0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以deque双向队列容器做接收者
//		将deque双向队列容器中所承载的元素放入迭代器中
//		节点的冗余空间不释放
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (d *Deque) Iterator() (i *Iterator.Iterator) {
	if d == nil {
		d = New()
	}
	tmp := make([]interface{}, 0, d.size)
	//遍历链表的所有节点,将其中承载的元素全部复制出来
	for m := d.first; m != nil; m = m.nextNode() {
		tmp = append(tmp, m.value()...)
	}
	return Iterator.New(&tmp)
}

//@title    Size
//@description
//		以deque双向队列容器做接收者
//		返回该容器当前含有元素的数量
//		该长度并非实际占用空间数量
//		若容器为空则返回0
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	size       	uint64					容器中实际使用元素所占空间大小
func (d *Deque) Size() (size uint64) {
	if d == nil {
		d = New()
	}
	return d.size
}

//@title    Clear
//@description
//		以deque双向队列容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的首尾指针均置nil,将size重置为0
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	nil
func (d *Deque) Clear() {
	if d == nil {
		d = New()
		return
	}
	d.mutex.Lock()
	d.first = nil
	d.last = nil
	d.size = 0
	d.mutex.Unlock()
}

//@title    Empty
//@description
//		以deque双向队列容器做接收者
//		判断该deque双向队列容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过size进行判断,为0则为true,否则为false
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (d *Deque) Empty() (b bool) {
	if d == nil {
		d = New()
	}
	return d.Size() == 0
}

//@title    PushFront
//@description
//		以deque双向队列向量容器做接收者
//		在容器首部插入元素
//		通过链表首节点进行添加
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (d *Deque) PushFront(e interface{}) {
	if d == nil {
		d = New()
	}
	d.mutex.Lock()
	d.size++
	//通过首节点进行添加
	if d.first == nil {
		d.first = createFirst()
		d.last = d.first
	}
	d.first = d.first.pushFront(e)
	d.mutex.Unlock()
}

//@title    PushBack
//@description
//		以deque双向队列向量容器做接收者
//		在容器尾部插入元素
//		通过链表尾节点进行添加
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (d *Deque) PushBack(e interface{}) {
	if d == nil {
		d = New()
	}
	d.mutex.Lock()
	d.size++
	//通过尾节点进行添加
	if d.last == nil {
		d.last = createLast()
		d.first = d.last
	}
	d.last = d.last.pushBack(e)
	d.mutex.Unlock()
}

//@title    PopFront
//@description
//		以deque双向队列容器做接收者
//		利用首节点进行弹出元素,可能存在首节点全部释放要进行首节点后移的情况
//		当元素全部删除后,释放全部空间,将首尾节点都设为nil
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	e			interface{}				首元素
func (d *Deque) PopFront() (e interface{}) {
	if d == nil {
		d = New()
	}
	if d.size == 0 {
		return nil
	}
	d.mutex.Lock()
	//利用首节点删除首元素
	//返回新的首节点
	e = d.first.front()
	d.first = d.first.popFront()
	d.size--
	if d.size == 0 {
		//全部删除完成,释放空间,并将首尾节点设为nil
		d.first = nil
		d.last = nil
	}
	d.mutex.Unlock()
	return e
}

//@title    PopBack
//@description
//		以deque双向队列容器做接收者
//		利用尾节点进行弹出元素,可能存在尾节点全部释放要进行尾节点前移的情况
//		当元素全部删除后,释放全部空间,将首尾节点都设为nil
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	e			interface{}				尾元素
func (d *Deque) PopBack() (e interface{}) {
	if d == nil {
		d = New()
	}
	if d.size == 0 {
		return nil
	}
	d.mutex.Lock()
	//利用尾节点删除首元素
	//返回新的尾节点
	d.last = d.last.popBack()
	e = d.last.back()
	d.size--
	if d.size == 0 {
		//全部删除完成,释放空间,并将首尾节点设为nil
		d.first = nil
		d.last = nil
	}
	d.mutex.Unlock()
	return e
}

//@title    Front
//@description
//		以deque双向队列容器做接收者
//		返回该容器的第一个元素,利用首节点进行寻找
//		若该容器当前为空,则返回nil
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	e			interface{}				容器的第一个元素
func (d *Deque) Front() (e interface{}) {
	if d == nil {
		d = New()
	}
	return d.first.front()
}

//@title    Back
//@description
//		以deque双向队列容器做接收者
//		返回该容器的最后一个元素,利用尾节点进行寻找
//		若该容器当前为空,则返回nil
//@receiver    	d        	*Deque					接收者的deque指针
//@param    	nil
//@return    	e			interface{}				容器的最后一个元素
func (d *Deque) Back() (e interface{}) {
	if d == nil {
		d = New()
	}
	return d.last.back()
}
