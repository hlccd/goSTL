package vector

//@Title		vector
//@Description
//		vector向量容器包
//		以动态数组的形式实现
//		该容器可以在尾部实现线性增减元素
//		通过interface实现泛型
//		可接纳不同类型的元素
//		但建议在同一个vector中使用相同类型的元素
//		可通过配合比较器competitor和迭代器iterator对该vector容器进行排序查找或遍历

import (
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//vector向量结构体
//包含泛型切片和该切片的尾指针
//当删除节点时仅仅需要长度-1即可
//当剩余长度较小时会采取缩容策略释放空间
//当添加节点时若未占满全部已分配空间则长度+1同时进行覆盖存放
//当添加节点时尾指针大于已分配空间长度,则按照扩容策略进行扩容
//并发控制锁用以保证在高并发过程中不会出现错误
//使用比较器重载了Sort

type vector struct {
	data  []interface{} //动态数组
	len   uint64        //当前已用数量
	cap   uint64        //可容纳元素数量
	mutex sync.Mutex    //并发控制锁
}

//vector向量容器接口
//存放了vector容器可使用的函数
//对应函数介绍见下方

type vectorer interface {
	Iterator() (i *Iterator.Iterator)  //返回一个包含vector所有元素的迭代器
	Sort(Cmp ...comparator.Comparator) //利用比较器对其进行排序
	Size() (num uint64)                //返回vector的长度
	Clear()                            //清空vector
	Empty() (b bool)                   //返回vector是否为空,为空则返回true反之返回false
	PushBack(e interface{})            //向vector末尾插入一个元素
	PopBack()                          //弹出vector末尾元素
	Insert(idx uint64, e interface{})  //向vector第idx的位置插入元素e,同时idx后的其他元素向后退一位
	Erase(idx uint64)                  //删除vector的第idx个元素
	Reverse()                          //逆转vector中的数据顺序
	At(idx uint64) (e interface{})     //返回vector的第idx的元素
	Front() (e interface{})            //返回vector的第一个元素
	Back() (e interface{})             //返回vector的最后一个元素
}

//@title    New
//@description
//		新建一个vector向量容器并返回
//		初始vector的切片数组为空
//		初始vector的长度为0,容量为1
//@receiver		nil
//@param    	nil
//@return    	v        	*vector					新建的vector指针
func New() (v *vector) {
	return &vector{
		data:  make([]interface{}, 1, 1),
		len:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以vector向量容器做接收者
//		释放未使用的空间,并将已使用的部分用于创建迭代器
//		返回一个包含容器中所有使用元素的迭代器
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (v *vector) Iterator() (i *Iterator.Iterator) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.data == nil {
		//data不存在,新建一个
		v.data = make([]interface{}, 1, 1)
		v.len = 0
		v.cap = 1
	} else if v.len < v.cap {
		//释放未使用的空间
		tmp := make([]interface{}, v.len, v.len)
		copy(tmp, v.data)
		v.data = tmp
	}
	//创建迭代器
	i = Iterator.New(&v.data)
	v.mutex.Unlock()
	return i
}

//@title    Sort
//@description
//		以vector向量容器做接收者
//		将vector向量容器中不使用空间释放掉
//		对元素中剩余的部分进行排序
//@receiver		v			*vector						接受者vector的指针
//@param    	Cmp			...comparator.Comparator	比较函数
//@return    	i        	*iterator.Iterator			新建的Iterator迭代器指针
func (v *vector) Sort(Cmp ...comparator.Comparator) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.data == nil {
		//data不存在,新建一个
		v.data = make([]interface{}, 1, 1)
		v.len = 0
		v.cap = 1
	} else if v.len < v.cap {
		//释放未使用空间
		tmp := make([]interface{}, v.len, v.len)
		copy(tmp, v.data)
		v.data = tmp
		v.cap = v.len
	}
	//调用比较器的Sort进行排序
	if len(Cmp) == 0 {
		comparator.Sort(&v.data)
	} else {
		comparator.Sort(&v.data, Cmp[0])
	}
	v.mutex.Unlock()
}

//@title    Size
//@description
//		以vector向量容器做接收者
//		返回该容器当前含有元素的数量
//		该长度并非实际占用空间数量,而是实际使用空间
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (v *vector) Size() (num uint64) {
	if v == nil {
		v = New()
	}
	return v.len
}

//@title    Clear
//@description
//		以vector向量容器做接收者
//		将该容器中的动态数组重置,所承载的元素清空
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	nil
func (v *vector) Clear() {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	//清空data
	v.data = make([]interface{}, 1, 1)
	v.len = 0
	v.cap = 1
	v.mutex.Unlock()
}

//@title    Empty
//@description
//		以vector向量容器做接收者
//		判断该vector向量容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		该判断过程通过长度进行判断
//		当长度为0时说明不含有元素
//		当长度大于0时说明含有元素
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (v *vector) Empty() (b bool) {
	if v == nil {
		v = New()
	}
	return v.Size() <= 0
}

//@title    PushBack
//@description
//		以vector向量容器做接收者
//		在容器尾部插入元素
//		若长度小于容量时,则对以长度为下标的位置进行覆盖,同时len++
//		若长度等于容量时,需要进行扩容
//		对于扩容而言,当容量小于2^16时,直接将容量翻倍,否则将容量增加2^16
//@receiver		v			*vector					接受者vector的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (v *vector) PushBack(e interface{}) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.len < v.cap {
		//还有冗余,直接添加
		v.data[v.len] = e
	} else {
		//冗余不足,需要扩容
		if v.cap < 2^16 {
			//容量翻倍
			if v.cap == 0 {
				v.cap = 1
			}
			v.cap *= 2
		} else {
			//容量增加2^16
			v.cap += 2 ^ 16
		}
		//复制扩容前的元素
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
		v.data[v.len] = e
	}
	v.len++
	v.mutex.Unlock()
}

//@title    PopBack
//@description
//		以vector向量容器做接收者
//		弹出容器最后一个元素,同时长度--即可
//		若容器为空,则不进行弹出
//		当弹出元素后,可能进行缩容
//		当容量和实际使用差值超过2^16时,容量直接减去2^16
//		否则,当实际使用长度是容量的一半时,进行折半缩容
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	nil
func (v *vector) PopBack() {
	if v == nil {
		v = New()
	}
	if v.Empty() {
		return
	}
	v.mutex.Lock()
	v.len--
	if v.cap-v.len >= 2^16 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		v.cap -= 2 ^ 16
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	} else if v.len*2 < v.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		v.cap /= 2
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	}
	v.mutex.Unlock()
}

//@title    Insert
//@description
//		以vector向量容器做接收者
//		向容器切片中插入一个元素
//		当idx不大于0时,则在容器切片的头部插入元素
//		当idx不小于切片使用长度时,在容器末尾插入元素
//		否则在切片中间第idx位插入元素,同时后移第idx位以后的元素
//		根据冗余量选择是否扩容,扩容策略同上
//		插入后len++
//@receiver		v			*vector					接受者vector的指针
//@param    	idx			uint64					待插入节点的位置(下标从0开始)
//@param		e			interface{}				待插入元素
//@return    	nil
func (v *vector) Insert(idx uint64, e interface{}) {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.len >= v.cap {
		//冗余不足,进行扩容
		if v.cap < 2^16 {
			//容量翻倍
			if v.cap == 0 {
				v.cap = 1
			}
			v.cap *= 2
		} else {
			//容量增加2^16
			v.cap += 2 ^ 16
		}
		//复制扩容前的元素
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	}
	//从后往前复制,即将idx后的全部后移一位即可
	var p uint64
	for p = v.len; p > 0 && p > uint64(idx); p-- {
		v.data[p] = v.data[p-1]
	}
	v.data[p] = e
	v.len++
	v.mutex.Unlock()
}

//@title    Erase
//@description
//		以vector向量容器做接收者
//		向容器切片中删除一个元素
//		当idx不大于0时,则删除头部
//		当idx不小于切片使用长度时,则删除尾部
//		否则在切片中间第idx位删除元素,同时前移第idx位以后的元素
//		长度同步--
//		进行缩容判断,缩容策略同上
//@receiver		v			*vector					接受者vector的指针
//@param    	idx			uint64					待删除节点的位置(下标从0开始)
//@return    	nil
func (v *vector) Erase(idx uint64) {
	if v == nil {
		v = New()
	}
	if v.Empty() {
		return
	}
	v.mutex.Lock()
	for p := idx; p < v.len-1; p++ {
		v.data[p] = v.data[p+1]
	}
	v.len--
	if v.cap-v.len >= 2^16 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		v.cap -= 2 ^ 16
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	} else if v.len*2 < v.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		v.cap /= 2
		tmp := make([]interface{}, v.cap, v.cap)
		copy(tmp, v.data)
		v.data = tmp
	}
	v.mutex.Unlock()
}

//@title    Reverse
//@description
//		以vector向量容器做接收者
//		将vector容器中不使用空间释放掉
//		将该容器中的泛型切片中的所有元素顺序逆转
//@receiver		v			*vector					接受者vector的指针
//@param		nil
//@return    	nil
func (v *vector) Reverse() {
	if v == nil {
		v = New()
	}
	v.mutex.Lock()
	if v.data == nil {
		//data不存在,新建一个
		v.data = make([]interface{}, 1, 1)
		v.len = 0
		v.cap = 1
	} else if v.len < v.cap {
		//释放未使用的空间
		tmp := make([]interface{}, v.len, v.len)
		copy(tmp, v.data)
		v.data = tmp
	}
	for i := uint64(0); i < v.len/2; i++ {
		v.data[i], v.data[v.len-i-1] = v.data[v.len-i-1], v.data[i]
	}
	v.mutex.Unlock()
}

//@title    At
//@description
//		以vector向量容器做接收者
//		根据传入的idx寻找位于第idx位的元素
//		当idx不在容器中泛型切片的使用范围内
//		即当idx小于0或者idx大于容器所含有的元素个数时返回nil
//		反之返回对应位置的元素
//@receiver		v			*vector					接受者vector的指针
//@param    	idx			uint64					待查找元素的位置(下标从0开始)
//@return    	e			interface{}				从容器中查找的第idx位元素
func (v *vector) At(idx uint64) (e interface{}) {
	if v == nil {
		v=New()
		return nil
	}
	v.mutex.Lock()
	if idx < 0 && idx >= v.Size() {
		v.mutex.Unlock()
		return nil
	}
	if v.Size() > 0 {
		e = v.data[idx]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}

//@title    Front
//@description
//		以vector向量容器做接收者
//		返回该容器的第一个元素
//		若该容器当前为空,则返回nil
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	e			interface{}				容器的第一个元素
func (v *vector) Front() (e interface{}) {
	if v == nil {
		v=New()
		return nil
	}
	v.mutex.Lock()
	if v.Size() > 0 {
		e = v.data[0]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}

//@title    Back
//@description
//		以vector向量容器做接收者
//		返回该容器的最后一个元素
//		若该容器当前为空,则返回nil
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	e			interface{}				容器的最后一个元素
func (v *vector) Back() (e interface{}) {
	if v == nil {
		v=New()
		return nil
	}
	v.mutex.Lock()
	if v.Size() > 0 {
		e = v.data[v.len-1]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}