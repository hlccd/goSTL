package list

//@Title		list
//@Description
//		list链表容器包
//		链表将所有结点通过指针的方式串联起来,从而使得其整体保持一个线性状态
//		对于链表的实现,其增删元素的过程只需要新建一个结点然后插入链表中即可
//		增删结点需要同步修改其相邻的元素的前后指针以保证其整体是联通的
//		可接纳不同类型的元素
//		通过并发控制锁保证了在高并发过程中的数据一致性

import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//list链表结构体
//包含链表的头尾节点指针
//当增删结点时只需要找到对应位置进行操作即可
//当一个节点进行增删时需要同步修改其临接结点的前后指针
//结构体中记录整个链表的首尾指针,同时记录其当前已承载的元素
//使用并发控制锁以保证数据一致性
type list struct {
	first *node      //链表首节点指针
	last  *node      //链表尾节点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}

//list链表容器接口
//存放了list容器可使用的函数
//对应函数介绍见下方

type lister interface {
	Iterator() (i *Iterator.Iterator)                              //创建一个包含链表中所有元素的迭代器并返回其指针
	Sort(Cmp ...comparator.Comparator)                             //将链表中所承载的所有元素进行排序
	Size() (size uint64)                                           //返回链表所承载的元素个数
	Clear()                                                        //清空该链表
	Empty() (b bool)                                               //判断该链表是否位空
	Insert(idx uint64, e interface{})                              //向链表的idx位(下标从0开始)插入元素组e
	Erase(idx uint64)                                              //删除第idx位的元素(下标从0开始)
	Get(idx uint64) (e interface{})                                //获得下标为idx的元素
	Set(idx uint64, e interface{})                                 //在下标为idx的位置上放置元素e
	IndexOf(e interface{}, Equ ...comparator.Equaler) (idx uint64) //返回和元素e相同的第一个下标
	SubList(begin, num uint64) (newList *list)                     //从begin开始复制最多num个元素以形成新的链表
}

//@title    New
//@description
//		新建一个list链表容器并返回
//		初始链表首尾节点为nil
//		初始size为0
//@receiver		nil
//@param    	nil
//@return    	l        	*list					新建的list指针
func New() (l *list) {
	return &list{
		first: nil,
		last:  nil,
		size:  0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以list链表容器做接收者
//		将list链表容器中所承载的元素放入迭代器中
//@receiver    	l        	*list					接收者的list指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (l *list) Iterator() (i *Iterator.Iterator) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//将所有元素复制出来放入迭代器中
	tmp := make([]interface{}, l.size, l.size)
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	i = Iterator.New(&tmp)
	l.mutex.Unlock()
	return i
}

//@title    Sort
//@description
//		以list链表容器做接收者
//		将list链表容器中所承载的元素利用比较器进行排序
//		可以自行传入比较函数,否则将调用默认比较函数
//@receiver    	l        	*list						接收者的list指针
//@param    	Cmp			...comparator.Comparator	比较函数
//@return		nil
func (l *list) Sort(Cmp ...comparator.Comparator) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//将所有元素复制出来用于排序
	tmp := make([]interface{}, l.size, l.size)
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	if len(Cmp) > 0 {
		comparator.Sort(&tmp, Cmp[0])
	} else {
		comparator.Sort(&tmp)
	}
	//将排序结果再放入链表中
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		n.setValue(tmp[idx])
	}
	l.mutex.Unlock()
}

//@title    Size
//@description
//		以list链表容器做接收者
//		返回该容器当前含有元素的数量
//@receiver    	l        	*list						接收者的list指针
//@param    	nil
//@return    	num        	int							容器中所承载的元素数量
func (l *list) Size() (size uint64) {
	if l == nil {
		l = New()
	}
	return l.size
}

//@title    Clear
//@description
//		以list链表容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的首尾指针均置nil,将size重置为0
//@receiver    	l        	*list						接收者的list指针
//@param    	nil
//@return    	nil
func (l *list) Clear() {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//销毁链表
	l.first = nil
	l.last = nil
	l.size = 0
	l.mutex.Unlock()
}

//@title    Empty
//@description
//		以list链表容器做接收者
//		判断该list链表容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过size进行判断,为0则为true,否则为false
//@receiver    	l        	*list						接收者的list指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (l *list) Empty() (b bool) {
	if l == nil {
		l = New()
	}
	return l.size == 0
}

//@title    Insert
//@description
//		以list链表容器做接收者
//		通过链表的首尾结点进行元素插入
//		插入的元素可以有很多个
//		通过判断idx
//@receiver    	l        	*list						接收者的list指针
//@param    	e			interface{}					待插入元素
//@return    	nil
func (l *list) Insert(idx uint64, e interface{}) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	n := newNode(e)
	if l.size == 0 {
		//链表中原本无元素,新建链表
		l.first = n
		l.last = n
	} else {
		//链表中存在元素
		if idx == 0 {
			//插入头节点
			n.insertNext(l.first)
			l.first = n
		} else if idx >= l.size {
			//插入尾节点
			l.last.insertNext(n)
			l.last = n
		} else {
			//插入中间节点
			//根据插入的位置选择从前或从后寻找
			if idx < l.size/2 {
				//从首节点开始遍历寻找
				m := l.first
				for i := uint64(0); i < idx-1; i++ {
					m = m.nextNode()
				}
				m.insertNext(n)
			} else {
				//从尾节点开始遍历寻找
				m := l.last
				for i := l.size - 1; i > idx; i-- {
					m = m.preNode()
				}
				m.insertPre(n)
			}
		}
	}
	l.size++
	l.mutex.Unlock()
}

//@title    Erase
//@description
//		以list链表容器做接收者
//		先判断是否为首尾结点,如果是首尾结点,在删除后将设置新的首尾结点
//		当链表所承载的元素全部删除后则销毁链表
//		删除时通过idx与总元素数量选择从前或从后进行遍历以找到对应位置
//		删除后,将该位置的前后结点连接起来,以保证链表不断裂
//@receiver    	l        	*list						接收者的list指针
//@param    	idx			uint64						被删除结点的下标(从0开始)
//@return    	nil
func (l *list) Erase(idx uint64) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	if l.size > 0 && idx < l.size {
		//链表中存在元素,且要删除的点在范围内
		if idx == 0 {
			//删除头节点
			l.first = l.first.next
		} else if idx == l.size-1 {
			//删除尾节点
			l.last = l.last.pre
		} else {
			//删除中间节点
			//根据删除的位置选择从前或从后寻找
			if idx < l.size/2 {
				//从首节点开始遍历寻找
				m := l.first
				for i := uint64(0); i < idx; i++ {
					m = m.nextNode()
				}
				m.erase()
			} else {
				//从尾节点开始遍历寻找
				m := l.last
				for i := l.size - 1; i > idx; i-- {
					m = m.preNode()
				}
				m.erase()
			}
		}
		l.size--
		if l.size == 0 {
			//所有节点都被删除,销毁链表
			l.first = nil
			l.last = nil
		}
	}
	l.mutex.Unlock()
}

//@title    Get
//@description
//		以list链表容器做接收者
//		获取第idx位结点所承载的元素,若不在链表范围内则返回nil
//@receiver    	l        	*list						接收者的list指针
//@param    	idx			uint64						被获取的结点位置(从0开始)
//@return    	e			interface{}					获取的元素
func (l *list) Get(idx uint64) (e interface{}) {
	if l == nil {
		l = New()
	}
	if idx >= l.size {
		return nil
	}
	l.mutex.Lock()
	if idx < l.size/2 {
		//从首节点开始遍历寻找
		m := l.first
		for i := uint64(0); i < idx; i++ {
			m = m.nextNode()
		}
		e = m.value()
	} else {
		//从尾节点开始遍历寻找
		m := l.last
		for i := l.size - 1; i > idx; i-- {
			m = m.preNode()
		}
		e = m.value()
	}
	l.mutex.Unlock()
	return e
}

//@title    Set
//@description
//		以list链表容器做接收者
//		修改第idx为结点所承载的元素,超出范围则不修改
//@receiver    	l        	*list						接收者的list指针
//@param    	idx			uint64						被修改的结点位置(从0开始)
//@param    	e			interface{}					修改后当元素
//@return		nil
func (l *list) Set(idx uint64, e interface{}) {
	if l == nil {
		l = New()
	}
	if idx >= l.size {
		return
	}
	l.mutex.Lock()
	if idx < l.size/2 {
		//从首节点开始遍历寻找
		m := l.first
		for i := uint64(0); i < idx; i++ {
			m = m.nextNode()
		}
		m.setValue(e)
	} else {
		//从尾节点开始遍历寻找
		m := l.last
		for i := l.size - 1; i > idx; i-- {
			m = m.preNode()
		}
		m.setValue(e)
	}
	l.mutex.Unlock()
}

//@title    IndexOf
//@description
//		以list链表容器做接收者
//		返回与e相同的元素的首个位置
//		可以自行传入用于判断相等的相等器进行处理
//		遍历从头至尾,如果不存在则返回l.size
//@receiver    	l        	*list						接收者的list指针
//@param    	e			interface{}					要查找的元素
//@param    	Equ			...comparator.Equaler		相等器
//@param    	idx			uint64						首下标
func (l *list) IndexOf(e interface{}, Equ ...comparator.Equaler) (idx uint64) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	var equ comparator.Equaler
	if len(Equ) > 0 {
		equ = Equ[0]
	} else {
		equ = comparator.GetEqual()
	}
	n := l.first
	//从头寻找直到找到相等的两个元素即可返回
	for idx = 0; idx < l.size && n != nil; idx++ {
		if equ(n.value(), e) {
			break
		}
		n = n.nextNode()
	}
	l.mutex.Unlock()
	return idx
}

//@title    SubList
//@description
//		以list链表容器做接收者
//		以begin为起点(包含),最多复制num个元素进入新链表
//		并返回新链表指针
//@receiver    	l        	*list						接收者的list指针
//@param    	begin		uint64						复制起点
//@param    	num			uint64						复制个数上限
//@param    	newList		*list						新链表指针
func (l *list) SubList(begin, num uint64) (newList *list) {
	if l == nil {
		l = New()
	}
	newList = New()
	l.mutex.Lock()
	if begin < l.size {
		//起点在范围内,可以复制
		n := l.first
		for i := uint64(0); i < begin; i++ {
			n = n.nextNode()
		}
		m := newNode(n.value())
		newList.first = m
		newList.size++
		for i := uint64(0); i < num-1 && i+begin < l.size-1; i++ {
			n = n.nextNode()
			m.insertNext(newNode(n.value()))
			m = m.nextNode()
			newList.size++
		}
		newList.last = m
	}
	l.mutex.Unlock()
	return newList
}
