package ring

//@Title		ring
//@Description
//		ring环容器包
//		环将所有结点通过指针的方式串联起来,从而使得其整体保持一个线性状态
//		不同于链表首尾不相连的情况,环将首尾结点连接起来,从而摒弃孤立的首尾结点
//		可以利用其中的任何一个结点遍历整个环,也可以在任何位置进行插入
//		增删结点需要同步修改其相邻的元素的前后指针以保证其整体是联通的
//		可接纳不同类型的元素
//		通过并发控制锁保证了在高并发过程中的数据一致性

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//ring环结构体
//包含环的头尾节点指针
//当增删结点时只需要移动到对应位置进行操作即可
//当一个节点进行增删时需要同步修改其临接结点的前后指针
//结构体中记录该环中当前所持有的结点的指针即可
//同时记录该环中存在多少元素即size
//使用并发控制锁以保证数据一致性
type ring struct {
	now   *node      //环当前持有的结点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}

//ring环容器接口
//存放了ring容器可使用的函数
//对应函数介绍见下方

type ringer interface {
	Iterator() (i *Iterator.Iterator) //创建一个包含环中所有元素的迭代器并返回其指针
	Size() (size uint64)              //返回环所承载的元素个数
	Clear()                           //清空该环
	Empty() (b bool)                  //判断该环是否位空
	Insert(e interface{})             //向环当前位置后方插入元素e
	Erase()                           //删除当前结点并持有下一结点
	Value() (e interface{})           //返回当前持有结点的元素
	Set(e interface{})                //在当前结点设置其承载的元素为e
	Next()                            //持有下一节点
	Pre()                             //持有上一结点
}

//@title    New
//@description
//		新建一个ring环容器并返回
//		初始持有的结点不存在,即为nil
//		初始size为0
//@receiver		nil
//@param    	nil
//@return    	r        	*ring					新建的ring指针
func New() (r *ring) {
	return &ring{
		now:   nil,
		size:  0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以ring环容器做接收者
//		将ring环容器中所承载的元素放入迭代器中
//		从该结点开始向后遍历获取全部承载的元素
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (r *ring) Iterator() (i *Iterator.Iterator) {
	if r == nil {
		r = New()
	}
	r.mutex.Lock()
	//将所有元素复制出来放入迭代器中
	tmp := make([]interface{}, r.size, r.size)
	//从当前结点开始向后遍历
	for n, idx := r.now, uint64(0); n != nil && idx < r.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	i = Iterator.New(&tmp)
	r.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以ring环容器做接收者
//		返回该容器当前含有元素的数量
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return    	num        	int						容器中所承载的元素数量
func (r *ring) Size() (size uint64) {
	if r == nil {
		r = New()
	}
	return r.size
}

//@title    Clear
//@description
//		以ring环容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的当前持有的结点置为nil,长度初始为0
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return    	nil
func (r *ring) Clear() {
	if r == nil {
		r = New()
	}
	r.mutex.Lock()
	//销毁环
	r.now = nil
	r.size = 0
	r.mutex.Unlock()
}

//@title    Empty
//@description
//		以ring环容器做接收者
//		判断该ring环容器是否含有元素
//		该判断过程通过size进行判断,size为0则为true,否则为false
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (r *ring) Empty() (b bool) {
	if r == nil {
		r = New()
	}
	return r.size == 0
}

//@title    Insert
//@description
//		以ring环容器做接收者
//		通过环中当前持有的结点进行添加
//		如果环为建立,则新建一个自环结点设为环
//		存在持有的结点,则在其后方添加即可
//@receiver    	r        	*ring					接收者的ring指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (r *ring) Insert(e interface{}) {
	if r == nil {
		r = New()
	}
	r.mutex.Lock()
	//新建自环结点
	n := newNode(e)
	if r.size == 0 {
		//原本无环,设为新环
		r.now = n
	} else {
		//持有结点,在后方插入
		r.now.insertNext(n)
	}
	r.size++
	r.mutex.Unlock()
}

//@title    Erase
//@description
//		以ring环容器做接收者
//		先判断是否仅持有一个结点
//		若仅有一个结点,则直接销毁环
//		否则将当前持有结点设为下一节点,并前插原持有结点的前结点即可
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return    	nil
func (r *ring) Erase() {
	if r == nil {
		r = New()
	}
	if r.size == 0 {
		return
	}
	r.mutex.Lock()
	//删除开始
	if r.size == 1 {
		//环内仅有一个结点,销毁环即可
		r.now = nil
	} else {
		//环内还有其他结点,将持有结点后移一位
		//后移后将当前结点前插原持有结点的前结点
		r.now = r.now.nextNode()
		r.now.insertPre(r.now.preNode().preNode())
	}
	r.size--
	r.mutex.Unlock()
}

//@title    Value
//@description
//		以ring环容器做接收者
//		获取环中当前持有节点所承载的元素
//		若环中持有的结点不存在,直接返回nil
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return    	e			interface{}				获取的元素
func (r *ring) Value() (e interface{}) {
	if r == nil {
		r = New()
	}
	if r.now == nil {
		//无持有结点,直接返回nil
		return nil
	}
	return r.now.value()
}

//@title    Set
//@description
//		以ring环容器做接收者
//		修改当前持有结点所承载的元素
//		若未持有结点,直接结束即可
//@receiver    	r        	*ring					接收者的ring指针
//@param    	e			interface{}					修改后当元素
//@return		nil
func (r *ring) Set(e interface{}) {
	if r == nil {
		r = New()
	}
	if r.now == nil {
		return
	}
	r.mutex.Lock()
	r.now.setValue(e)
	r.mutex.Unlock()
}

//@title    Next
//@description
//		以ring环容器做接收者
//		将当前持有的结点后移一位
//		若当前无持有结点,则直接结束
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return		nil
func (r *ring) Next() {
	if r == nil {
		r = New()
	}
	if r.now == nil {
		return
	}
	r.mutex.Lock()
	r.now = r.now.nextNode()
	r.mutex.Unlock()
}

//@title    Pre
//@description
//		以ring环容器做接收者
//		将当前持有的结点前移一位
//		若当前无持有结点,则直接结束
//@receiver    	r        	*ring					接收者的ring指针
//@param    	nil
//@return		nil
func (r *ring) Pre() {
	if r == nil {
		r = New()
	}
	if r.size == 0 {
		return
	}
	r.mutex.Lock()
	r.now = r.now.preNode()
	r.mutex.Unlock()
}
