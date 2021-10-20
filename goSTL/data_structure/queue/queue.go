package queue

//@Title		queue
//@Description
//		queue队列容器包
//		以动态数组的形式实现
//		该容器可以在尾部实现线性增加元素,在首部实现线性减少元素
//		队列的扩容和缩容同vector一样,即数组小采用翻倍扩缩容/折半缩容,数组大时采用固定扩/缩容
//		该容器满足FIFO的先进先出模式
//		可接纳不同类型的元素
//		通过并发控制锁保证了在高并发过程中的数据一致性

import (
	"sync"
)

//queue队列结构体
//包含泛型切片和该切片的首尾位的下标
//当删除节点时仅仅需要后移首位一位即可
//当剩余长度较小时采用缩容策略进行缩容以释放空间
//当添加节点时若未占满全部已分配空间则尾指针后移一位同时进行覆盖存放
//当添加节点时尾指针大于已分配空间长度,则先去掉首部多出来的空间,如果还不足则进行扩容
//首节点指针始终不能超过尾节点指针
type queue struct {
	data  []interface{} //泛型切片
	begin uint64        //首节点下标
	end   uint64        //尾节点下标
	cap   uint64        //容量
	mutex sync.Mutex    //并发控制锁
}

//queue队列容器接口
//存放了queue容器可使用的函数
//对应函数介绍见下方

type queuer interface {
	Size() (num uint64)     //返回该队列中元素的使用空间大小
	Clear()                 //清空该队列
	Empty() (b bool)        //判断该队列是否为空
	Push(e interface{})     //将元素e添加到该队列末尾
	Pop() (e interface{})   //将该队列首元素弹出并返回
	Front() (e interface{}) //获取该队列首元素
	Back() (e interface{})  //获取该队列尾元素
}

//@title    New
//@description
//		新建一个queue队列容器并返回
//		初始queue的切片数组为空
//		初始queue的首尾指针均置零
//@receiver		nil
//@param    	nil
//@return    	q        	*queue					新建的queue指针
func New() (q *queue) {
	return &queue{
		data:  make([]interface{}, 1, 1),
		begin: 0,
		end:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}

//@title    Size
//@description
//		以queue队列容器做接收者
//		返回该容器当前含有元素的数量
//		该长度并非实际占用空间数量
//		若容器为空则返回0
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (q *queue) Size() (num uint64) {
	if q == nil {
		q = New()
	}
	return q.end - q.begin
}

//@title    Clear
//@description
//		以queue队列容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的首尾指针均置0,容量设为1
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	nil
func (q *queue) Clear() {
	if q == nil {
		q = New()
	}
	q.mutex.Lock()
	q.data = make([]interface{}, 1, 1)
	q.begin = 0
	q.end = 0
	q.cap = 1
	q.mutex.Unlock()
}

//@title    Empty
//@description
//		以queue队列容器做接收者
//		判断该queue队列容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过首尾指针数值进行判断
//		当尾指针数值等于首指针时说明不含有元素
//		当尾指针数值大于首指针时说明含有元素
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (q *queue) Empty() (b bool) {
	if q == nil {
		q = New()
	}
	return q.Size() <= 0
}

//@title    Push
//@description
//		以queue队列向量容器做接收者
//		在容器尾部插入元素
//		若尾指针小于切片实际使用长度,则对当前指针位置进行覆盖,同时尾下标后移一位
//		若尾指针等于切片实际使用长度,则对实际使用量和实际占用量进行判断
//		当首部还有冗余时则删将实际使用整体前移到首部,否则对尾部进行扩容即可
//@receiver		q			*queue					接受者queue的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (q *queue) Push(e interface{}) {
	if q == nil {
		q = New()
	}
	q.mutex.Lock()
	if q.end < q.cap {
		//不需要扩容
		q.data[q.end] = e
	} else {
		//需要扩容
		if q.begin > 0 {
			//首部有冗余,整体前移
			for i := uint64(0); i < q.end-q.begin; i++ {
				q.data[i] = q.data[i+q.begin]
			}
			q.end -= q.begin
			q.begin = 0
		} else {
			//冗余不足,需要扩容
			if q.cap <= 65536 {
				//容量翻倍
				if q.cap == 0 {
					q.cap = 1
				}
				q.cap *= 2
			} else {
				//容量增加2^16
				q.cap += 2 ^ 16
			}
			//复制扩容前的元素
			tmp := make([]interface{}, q.cap, q.cap)
			copy(tmp, q.data)
			q.data = tmp
		}
		q.data[q.end] = e
	}
	q.end++
	q.mutex.Unlock()
}

//@title    Pop
//@description
//		以queue队列容器做接收者
//		弹出容器第一个元素,同时首下标后移一位
//		若容器为空,则不进行弹出
//		弹出结束后,进行缩容判断,考虑到queue的冗余会存在于前后两个方向
//		所以需要对前后两方分别做判断, 但由于首部主要是减少,并不会增加,所以不需要太多冗余量,而尾部只做添加,所以需要更多的冗余
//		所以可以对首部预留2^10的冗余,当超过时直接对首部冗余清除即可,释放首部空间时尾部空间仍然保留不变
//		当首部冗余不足2^10时,但冗余超过实际使用空间,也会对首部进行缩容,尾部不变
//		同时返回队首元素
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	e 			interface{}				队首元素
func (q *queue) Pop() (e interface{}) {
	if q == nil {
		q = New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	q.mutex.Lock()
	e = q.data[q.begin]
	q.begin++
	if q.begin >= 1024 || q.begin*2>q.end {
		//首部冗余超过2^10或首部冗余超过实际使用
		q.cap -= q.begin
		q.end -= q.begin
		tmp := make([]interface{}, q.cap, q.cap)
		copy(tmp, q.data[q.begin:])
		q.data = tmp
		q.begin=0
	}
	q.mutex.Unlock()
	return e
}

//@title    Front
//@description
//		以queue队列容器做接收者
//		返回该容器的第一个元素
//		若该容器当前为空,则返回nil
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	e			interface{}				容器的第一个元素
func (q *queue) Front() (e interface{}) {
	if q == nil {
		q=New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	return q.data[q.begin]
}

//@title    Back
//@description
//		以queue队列容器做接收者
//		返回该容器的最后一个元素
//		若该容器当前为空,则返回nil
//@receiver		q			*queue					接受者queue的指针
//@param    	nil
//@return    	e			interface{}				容器的最后一个元素
func (q *queue) Back() (e interface{}) {
	if q == nil {
		q=New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	return q.data[q.end-1]
}
