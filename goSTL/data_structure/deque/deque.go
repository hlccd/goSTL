package deque

//@Title		deque
//@Description
//		deque双向队列容器包
//		以切片数组的形式实现
//		该容器可以在首部和尾部实现线性增减元素
//		通过interface实现泛型
//		可接纳不同类型的元素
//@author     	hlccd		2021-07-6
//@update		hlccd 		2021-08-01		增加互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//deque双向队列结构体
//包含泛型切片和该切片的首尾指针
//当删除节点时仅仅需要后移首指针一位即可
//当剩余长度小于实际占用空间长度的一半时会重新规划以释放掉多余占用的空间
//当添加节点时若未占满全部已分配空间则尾指针后移一位同时进行覆盖存放
//当添加节点时尾指针大于已分配空间长度,则新增空间
//首节点指针始终不能超过尾节点指针
type deque struct {
	data  []interface{} //泛型切片
	begin int           //首节点指针
	end   int           //尾节点指针
	mutex sync.Mutex    //并发控制锁
}

//deque双向队列容器接口
//存放了deque容器可使用的函数
//对应函数介绍见下方
type dequeer interface {
	Iterator() *Iterator.Iterator //返回一个包含deque中所有使用元素的迭代器
	Size() (num int)              //返回该双向队列中元素的使用空间大小
	Clear()                       //清空该双向队列
	Empty() (b bool)              //判断该双向队列是否为空
	PushFront(e interface{})      //将元素e添加到该队列首部
	PushBack(e interface{})       //将元素e添加到该队列末尾
	PopFront() (e interface{})    //将该队列首元素弹出并返回
	PopBack() (e interface{})     //将该队列尾元素弹出并返回
	Front() (e interface{})       //获取该队列首元素
	Back() (e interface{})        //获取该队列尾元素
}

//@title    New
//@description
//		新建一个deque双向队列容器并返回
//		初始deque的切片数组为空
//		初始deque的首尾指针均置零
//@author     	hlccd		2021-07-6
//@receiver		nil
//@param    	nil
//@return    	d        	*deque					新建的deque指针
func New() *deque {
	return &deque{
		data:  make([]interface{}, 0, 0),
		begin: 0,
		end:   0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以deque双向队列容器做接收者
//		将deque队列容器中不使用空间释放掉
//		返回一个包含容器中所有使用元素的迭代器
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (d *deque) Iterator() (i *Iterator.Iterator) {
	if d == nil {
		return Iterator.New(make([]interface{}, 0, 0))
	}
	d.mutex.Lock()
	d.data = d.data[d.begin:d.end]
	d.begin = 0
	d.end = len(d.data)
	i = iterator.New(d.data)
	d.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以deque双向队列容器做接收者
//		返回该容器当前含有元素的数量
//		该长度并非实际占用空间数量
//		当容器为nil时返回-1
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (d *deque) Size() (num int) {
	if d == nil {
		return -1
	}
	return d.end - d.begin
}

//@title    Clear
//@description
//		以deque双向队列容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的首尾指针均置0
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	nil
func (d *deque) Clear() {
	if d == nil {
		return
	}
	d.mutex.Lock()
	d.data = d.data[0:0]
	d.begin = 0
	d.end = 0
	d.mutex.Unlock()
}

//@title    Empty
//@description
//		以deque双向队列容器做接收者
//		判断该deque双向队列容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过首尾指针数值进行判断
//		当尾指针数值等于首指针时说明不含有元素
//		当尾指针数值大于首指针时说明含有元素
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (d *deque) Empty() bool {
	if d == nil {
		return true
	}
	return d.Size() <= 0
}

//@title    PushFront
//@description
//		以deque双向队列容器做接收者
//		在容器首部插入元素
//		若首指针大于0,则对当前指针位置进行覆盖,同时首指针前移一位
//		若首指针等于0,重新分配空间并将首指针置0尾指针设为总长度
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	e			interface{}				待插入首部的元素
//@return    	nil
func (d *deque) PushFront(e interface{}) {
	if d == nil {
		return
	}
	d.mutex.Lock()
	if d.begin > 0 {
		d.begin--
		d.data[d.begin] = e
	} else {
		d.data = append(append([]interface{}{}, e), d.data[:d.end]...)
		d.begin = 0
		d.end = len(d.data)
	}
	d.mutex.Unlock()
}

//@title    PushBack
//@description
//		以deque双向队列容器做接收者
//		在容器尾部插入元素
//		若尾指针小于切片实际使用长度,则对当前指针位置进行覆盖,同时尾指针后移一位
//		若尾指针等于切片实际使用长度,则新增切片长度同时使尾指针后移一位
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	e			interface{}				待插入尾部的元素
//@return    	nil
func (d *deque) PushBack(e interface{}) {
	if d == nil {
		return
	}
	d.mutex.Lock()
	if d.end < len(d.data) {
		d.data[d.end] = e
	} else {
		d.data = append(d.data, e)
	}
	d.end++
	d.mutex.Unlock()
}

//@title    PopFront
//@description
//		以deque双向队列容器做接收者
//		弹出容器第一个元素,同时首指针后移一位
//		当剩余元素数量小于容器切片实际使用空间的一半时,重新分配空间释放未使用部分
//		若容器为空,则不进行弹出
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	e 			interface{}				队尾元素
func (d *deque) PopFront() (e interface{}) {
	if d == nil {
		return nil
	}
	if d.Empty() {
		return nil
	}
	d.mutex.Lock()
	e = d.data[d.begin]
	d.begin++
	if d.begin*2 >= d.end {
		d.data = d.data[d.begin:d.end]
		d.begin = 0
		d.end = len(d.data)
	}
	d.mutex.Unlock()
	return e
}

//@title    PopBack
//@description
//		以deque双向队列容器做接收者
//		弹出容器最后一个元素,同时尾指针前移一位
//		当剩余元素数量小于容器切片实际使用空间的一半时,重新分配空间释放未使用部分
//		若容器为空,则不进行弹出
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	e 			interface{}				队首元素
func (d *deque) PopBack() (e interface{}) {
	if d == nil {
		return nil
	}
	if d.Empty() {
		return nil
	}
	d.mutex.Lock()
	d.end--
	e = d.data[d.end]
	if d.begin*2 >= d.end {
		d.data = d.data[d.begin:d.end]
		d.begin = 0
		d.end = len(d.data)
	}
	d.mutex.Unlock()
	return e
}

//@title    Front
//@description
//		以deque双向队列容器做接收者
//		返回该容器的第一个元素
//		若该容器当前为空,则返回nil
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	e			interface{}				容器的第一个元素
func (d *deque) Front() (e interface{}) {
	if d == nil {
		return nil
	}
	if d.Empty() {
		return nil
	}
	d.mutex.Lock()
	e = d.data[d.begin]
	d.mutex.Unlock()
	return e
}

//@title    Back
//@description
//		以deque双向队列容器做接收者
//		返回该容器的最后一个元素
//		若该容器当前为空,则返回nil
//@auth      	hlccd		2021-07-6
//@return    	d        	*deque					接收者的deque指针
//@param    	nil
//@return    	e			interface{}				容器的最后一个元素
func (d *deque) Back() (e interface{}) {
	if d == nil {
		return
	}
	if d.Empty() {
		return nil
	}
	d.mutex.Lock()
	e = d.data[d.end-1]
	d.mutex.Unlock()
	return e
}
