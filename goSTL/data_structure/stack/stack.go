package stack

//@Title		stack
//@Description
//		stack栈容器包
//		以动态数组的形式实现,扩容方式同vector
//		该容器可以在顶部实现线性增减元素
//		通过interface实现泛型，可接纳不同类型的元素
//		互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//stack栈结构体
//包含动态数组和该数组的顶部指针
//顶部指针指向实际顶部元素的下一位置
//当删除节点时仅仅需要下移顶部指针一位即可
//当新增结点时优先利用冗余空间
//当冗余空间不足时先倍增空间至2^16，超过后每次增加2^16的空间
//删除结点后如果冗余超过2^16,则释放掉
//删除后若冗余量超过使用量，也释放掉冗余空间
type stack struct {
	data  []interface{} //用于存储元素的动态数组
	top   uint64        //顶部指针
	cap   uint64        //动态数组的实际空间
	mutex sync.Mutex    //并发控制锁
}

//stack栈容器接口
//存放了stack容器可使用的函数
//对应函数介绍见下方
type stacker interface {
	Iterator() (i *Iterator.Iterator) //返回一个包含栈中所有元素的迭代器
	Size() (num uint64)               //返回该栈中元素的使用空间大小
	Clear()                           //清空该栈容器
	Empty() (b bool)                  //判断该栈容器是否为空
	Push(e interface{})               //将元素e添加到栈顶
	Pop()                             //弹出栈顶元素
	Top() (e interface{})             //返回栈顶元素
}

//@title    New
//@description
//		新建一个stack栈容器并返回
//		初始stack的动态数组容量为1
//		初始stack的顶部指针置0，容量置1
//@receiver		nil
//@param    	nil
//@return    	s        	*stack					新建的stack指针
func New() (s *stack) {
	return &stack{
		data:  make([]interface{}, 1, 1),
		top:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以stack栈容器做接收者
//		将stack栈容器中不使用空间释放掉
//		返回一个包含容器中所有使用元素的迭代器
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (s *stack) Iterator() (i *Iterator.Iterator) {
	if s == nil {
		s = New()
	}
	s.mutex.Lock()
	if s.data == nil {
		//data不存在,新建一个
		s.data = make([]interface{}, 1, 1)
		s.top = 0
		s.cap = 1
	} else if s.top < s.cap {
		//释放未使用的空间
		tmp := make([]interface{}, s.top, s.top)
		copy(tmp, s.data)
		s.data = tmp
	}
	//创建迭代器
	i = Iterator.New(&s.data)
	s.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以stack栈容器做接收者
//		返回该容器当前含有元素的数量
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (s *stack) Size() (num uint64) {
	if s == nil {
		s = New()
	}
	return s.top
}

//@title    Clear
//@description
//		以stack栈容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的尾指针置0
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	nil
func (s *stack) Clear() {
	if s == nil {
		s = New()
	}
	s.mutex.Lock()
	s.data = make([]interface{}, 0, 0)
	s.top = 0
	s.cap = 1
	s.mutex.Unlock()
}

//@title    Empty
//@description
//		以stack栈容器做接收者
//		判断该stack栈容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过顶部指针数值进行判断
//		当顶部指针数值为0时说明不含有元素
//		当顶部指针数值大于0时说明含有元素
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (s *stack) Empty() (b bool) {
	if s == nil {
		return true
	}
	return s.Size() == 0
}

//@title    Push
//@description
//		以stack栈容器做接收者
//		在容器顶部插入元素
//		若存储冗余空间，则在顶部指针位插入元素，随后上移顶部指针
//		否则进行扩容，扩容后获得冗余空间重复上一步即可。
//@receiver		s			*stack					接受者stack的指针
//@param    	e			interface{}				待插入顶部的元素
//@return    	nil
func (s *stack) Push(e interface{}) {
	if s == nil {
		s = New()
	}
	s.mutex.Lock()
	if s.top < s.cap {
		//还有冗余,直接添加
		s.data[s.top] = e
	} else {
		//冗余不足,需要扩容
		if s.cap <= 65536 {
			//容量翻倍
			if s.cap == 0 {
				s.cap = 1
			}
			s.cap *= 2
		} else {
			//容量增加2^16
			s.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
		s.data[s.top] = e
	}
	s.top++
	s.mutex.Unlock()
}

//@title    Pop
//@description
//		以stack栈容器做接收者
//		弹出容器顶部元素,同时顶部指针下移一位
//		当顶部指针小于容器切片实际使用空间的一半时,重新分配空间释放未使用部分
//		若容器为空,则不进行弹出
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	nil
func (s *stack) Pop() {
	if s == nil {
		s = New()
		return
	}
	if s.Empty() {
		return
	}
	s.mutex.Lock()
	s.top--
	if s.cap-s.top >= 65536 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		s.cap -= 65536
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	} else if s.top*2 < s.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		s.cap /= 2
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	}
	s.mutex.Unlock()
}

//@title    Top
//@description
//		以stack栈容器做接收者
//		返回该容器的顶部元素
//		若该容器当前为空,则返回nil
//@receiver		s			*stack					接受者stack的指针
//@param    	nil
//@return    	e			interface{}				容器的顶部元素
func (s *stack) Top() (e interface{}) {
	if s == nil {
		return nil
	}
	if s.Empty() {
		return nil
	}
	s.mutex.Lock()
	e = s.data[s.top-1]
	s.mutex.Unlock()
	return e
}
