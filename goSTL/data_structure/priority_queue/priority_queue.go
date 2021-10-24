package priority_queue

//@Title		priority_queue
//@Description
//		priority_queue优先队列集合容器包
//		以动态数组的形式实现
//		该容器可以增减元素使最值元素处于顶端
//		若使用默认比较器,顶端元素是最小元素
//		该集合只能对于相等元素可以存储多个
//		可接纳不同类型的元素,但为了便于比较,建议使用同一个类型
//		使用互斥锁实现并发控制,以保证在高并发情况下的数据一致性
import (
	"github.com/hlccd/goSTL/utils/comparator"
	"sync"
)

//priority_queue优先队列集合结构体
//包含动态数组和比较器,同时包含其实际使用长度和实际占用空间容量
//该数据结构可以存储多个相同的元素,并不会产生冲突
//增删节点后会使用比较器保持该动态数组的相对有序性
type priority_queue struct {
	data  []interface{}         //动态数组
	len   uint64                //实际使用长度
	cap   uint64                //实际占用的空间的容量
	cmp   comparator.Comparator //该优先队列的比较器
	mutex sync.Mutex            //并发控制锁
}

//priority_queue优先队列容器接口
//存放了priority_queue容器可使用的函数
//对应函数介绍见下方
type priority_queueer interface {
	Size() (num uint64)   //返回该容器存储的元素数量
	Clear()               //清空该容器
	Empty() (b bool)      //判断该容器是否为空
	Push(e interface{})   //将元素e插入该容器
	Pop()                 //弹出顶部元素
	Top() (e interface{}) //返回顶部元素
}

//@title    New
//@description
//		新建一个priority_queue优先队列容器并返回
//		初始priority_queue的切片数组为空
//		如果有传入比较器,则将传入的第一个比较器设为可重复集合默认比较器
//		如果不传入比较器,在后续的增删过程中将会去寻找默认比较器
//@receiver		nil
//@param    	Cmp			...comparator.Comparator	priority_queue的比较器
//@return    	pq        	*priority_queue				新建的priority_queue指针
func New(cmps ...comparator.Comparator) (pq *priority_queue) {
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	//比较器为nil时后续的增删将会去寻找默认比较器
	return &priority_queue{
		data:  make([]interface{}, 1, 1),
		len:   0,
		cap:   1,
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}

//@title    Size
//@description
//		以priority_queue容器做接收者
//		返回该容器当前含有元素的数量
//@receiver		pq			*priority_queue			接受者priority_queue的指针
//@param    	nil
//@return    	num        	uint64					容器中存储元素的个数
func (pq *priority_queue) Size() (num uint64) {
	if pq == nil {
		pq = New()
	}
	return pq.len
}

//@title    Clear
//@description
//		以priority_queue容器做接收者
//		将该容器中所承载的元素清空
//@receiver		pq			*priority_queue			接受者priority_queue的指针
//@param    	nil
//@return    	nil
func (pq *priority_queue) Clear() {
	if pq == nil {
		pq = New()
	}
	pq.mutex.Lock()
	//清空已分配的空间
	pq.data = make([]interface{}, 1, 1)
	pq.len = 0
	pq.cap = 1
	pq.mutex.Unlock()
}

//@title    Empty
//@description
//		以priority_queue容器做接收者
//		判断该priority_queue容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过含有元素个数进行判断
//@receiver		pq			*priority_queue			接受者priority_queue的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (pq *priority_queue) Empty() bool {
	if pq == nil {
		pq = New()
	}
	return pq.len == 0
}

//@title    Push
//@description
//		以priority_queue容器做接收者
//		在该优先队列中插入元素e,利用比较器和交换使得优先队列保持相对有序状态
//		插入时,首先将该元素放入末尾,然后通过比较其逻辑上的父结点选择是否上移
//		扩容策略同vector,先进行翻倍扩容,在进行固定扩容,界限为2^16
//@receiver		pq			*priority_queue			接受者priority_queue的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (pq *priority_queue) Push(e interface{}) {
	if pq == nil {
		pq = New()
	}
	pq.mutex.Lock()
	//判断是否存在比较器,不存在则寻找默认比较器,若仍不存在则直接结束
	if pq.cmp == nil {
		pq.cmp = comparator.GetCmp(e)
	}
	if pq.cmp == nil {
		pq.mutex.Unlock()
		return
	}
	//先判断是否需要扩容,同时使用和vector相同的扩容策略
	//即先翻倍扩容再固定扩容,随后在末尾插入元素e
	if pq.len < pq.cap {
		//还有冗余,直接添加
		pq.data[pq.len] = e
	} else {
		//冗余不足,需要扩容
		if pq.cap <= 65536 {
			//容量翻倍
			if pq.cap == 0 {
				pq.cap = 1
			}
			pq.cap *= 2
		} else {
			//容量增加2^16
			pq.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
		pq.data[pq.len] = e
	}
	pq.len++
	//到此时,元素以插入到末尾处,同时插入位的元素的下标为pq.len-1,随后将对该位置的元素进行上升
	//即通过比较它逻辑上的父结点进行上升
	pq.up(pq.len - 1)
	pq.mutex.Unlock()
}

//@title    up
//@description
//		以priority_queue容器做接收者
//		用于递归判断任意子结点和其父结点之间的关系，满足上升条件则递归上升
//		从而保证父节点必然都大于或都小于子节点
//@receiver		pq			*priority_queue					接受者priority_queue的指针
//@param    	p			uint64					待上升节点的位置
//@return    	nil
func (pq *priority_queue) up(p uint64) {
	if p == 0 {
		//以及上升到顶部,直接结束即可
		return
	}
	if pq.cmp(pq.data[(p-1)/2], pq.data[p]) > 0 {
		//判断该结点和其父结点的关系
		//满足给定的比较函数的关系则先交换该结点和父结点的数值,随后继续上升即可
		pq.data[p], pq.data[(p-1)/2] = pq.data[(p-1)/2], pq.data[p]
		pq.up((p - 1) / 2)
	}
}

//@title    Pop
//@description
//		以priority_queue容器做接收者
//		在该优先队列中删除顶部元素,利用比较器和交换使得优先队列保持相对有序状态
//		删除时首先将首结点移到最后一位进行交换,随后删除最后一位即可,然后对首节点进行下降即可
//		缩容时同vector一样,先进行固定缩容在进行折半缩容,界限为2^16
//@receiver		pq			*priority_queue					接受者priority_queue的指针
//@param    	nil
//@return    	nil
func (pq *priority_queue) Pop() {
	if pq == nil {
		pq = New()
	}
	if pq.Empty() {
		return
	}
	pq.mutex.Lock()
	//将最后一位移到首位,随后删除最后一位,即删除了首位,同时判断是否需要缩容
	pq.data[0] = pq.data[pq.len-1]
	pq.data[pq.len-1]=nil
	pq.len--
	//缩容判断,缩容策略同vector,即先固定缩容在折半缩容
	if pq.cap-pq.len >= 65536 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		pq.cap -= 65536
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
	} else if pq.len*2 < pq.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		pq.cap /= 2
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
	}
	//判断是否为空,为空则直接结束
	if pq.Empty() {
		pq.mutex.Unlock()
		return
	}
	//对首位进行下降操作,即对比其逻辑上的左右结点判断是否应该下降,再递归该过程即可
	pq.down(0)
	pq.mutex.Unlock()
}

//@title    down
//@description
//		以priority_queue容器做接收者
//		判断待下沉节点与其左右子节点的大小关系以确定是否进行递归上升
//		从而保证父节点必然都大于或都小于子节点
//@receiver		pq			*priority_queue					接受者priority_queue的指针
//@param    	p			uint64					待下沉节点的位置
//@return    	nil
func (pq *priority_queue) down(p uint64) {
	q := p
	//先判断其左结点是否在范围内,然后在判断左结点是否满足下降条件
	if 2*p+1 <= pq.len-1 && pq.cmp(pq.data[p], pq.data[2*p+1]) > 0 {
		q = 2*p + 1
	}
	//在判断右结点是否在范围内,同时若判断右节点是否满足下降条件
	if 2*p+2 <= pq.len-1 && pq.cmp(pq.data[q], pq.data[2*p+2]) > 0 {
		q = 2*p + 2
	}
	//根据上面两次判断,从最小一侧进行下降
	if p != q {
		//进行交互,递归下降
		pq.data[p], pq.data[q] = pq.data[q], pq.data[p]
		pq.down(q)
	}
}

//@title    Top
//@description
//		以priority_queue容器做接收者
//		返回该优先队列容器的顶部元素
//		如果容器不存在或容器为空,返回nil
//@receiver		pq			*priority_queue			接受者priority_queue的指针
//@param    	nil
//@return    	e			interface{}				优先队列顶元素
func (pq *priority_queue) Top() (e interface{}) {
	if pq == nil {
		pq = New()
	}
	if pq.Empty() {
		return nil
	}
	pq.mutex.Lock()
	e = pq.data[0]
	pq.mutex.Unlock()
	return e
}
