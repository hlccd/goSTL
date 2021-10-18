package vector

//@Title		vector
//@Description
//		vector向量容器包
//		以切片数组的形式实现
//		该容器可以在尾部实现线性增减元素
//		通过interface实现泛型
//		可接纳不同类型的元素
//		但建议在同一个vector中使用相同类型的元素
//		可通过配合比较器competitor和迭代器iterator对该vector容器进行排序或查找

import (
	"sync"
)

//vector向量结构体
//包含泛型切片和该切片的尾指针
//当删除节点时仅仅需要前移尾指针一位即可
//当剩余长度小于实际占用空间长度的一半时会重新规划以释放掉多余占用的空间
//当添加节点时若未占满全部已分配空间则尾指针后移一位同时进行覆盖存放
//当添加节点时尾指针大于已分配空间长度,则新增空间

type vector struct {
	data  []interface{} //泛型切片
	end   int           //尾指针
	mutex sync.Mutex    //并发控制锁
}

//vector向量容器接口
//存放了vector容器可使用的函数
//对应函数介绍见下方

type vectorer interface {
	Size() (num int)               //返回vector的长度
	Clear()                        //清空vector
	Empty() (b bool)               //返回vector是否为空,为空则返回true反之返回false
	PushBack(e interface{})        //向vector末尾插入一个元素
	PopBack()                      //弹出vector末尾元素
	Insert(idx int, e interface{}) //向vector第idx的位置插入元素e,同时idx后的其他元素向后退一位
	Erase(idx int)                 //删除vector的第idx个元素
	Reverse()                      //逆转vector中的数据顺序
	At(idx int) (e interface{})    //返回vector的第idx的元素
	Front() (e interface{})        //返回vector的第一个元素
	Back() (e interface{})         //返回vector的最后一个元素
}

//@title    New
//@description
//		新建一个vector向量容器并返回
//		初始vector的切片数组为空
//		初始vector的尾指针置0
//@receiver		nil
//@param    	nil
//@return    	v        	*vector					新建的vector指针
func New() (v *vector) {
	return &vector{
		data:  make([]interface{}, 1),
		end:   0,
		mutex: sync.Mutex{},
	}
}


//@title    Size
//@description
//		以vector向量容器做接收者
//		返回该容器当前含有元素的数量
//		该长度并非实际占用空间数量
//		如果容器为nil返回-1
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (v *vector) Size() (num int) {
	if v == nil {
		return -1
	}
	return v.end
}

//@title    Clear
//@description
//		以vector向量容器做接收者
//		将该容器中所承载的元素清空
//		将该容器的尾指针置0
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	nil
func (v *vector) Clear() {
	if v == nil {
		return
	}
	v.mutex.Lock()
	v.data = v.data[0:0]
	v.end = 0
	v.mutex.Unlock()
}

//@title    Empty
//@description
//		以vector向量容器做接收者
//		判断该vector向量容器是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//		该判断过程通过尾指针数值进行判断
//		当尾指针数值为0时说明不含有元素
//		当尾指针数值大于0时说明含有元素
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (v *vector) Empty() (b bool) {
	if v == nil {
		return true
	}
	return v.Size() <= 0
}

//@title    PushBack
//@description
//		以vector向量容器做接收者
//		在容器尾部插入元素
//		若尾指针小于切片实际使用长度,则对当前指针位置进行覆盖,同时尾指针后移一位
//		若尾指针等于切片实际使用长度,则新增切片长度同时使尾指针后移一位
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	e			interface{}				待插入元素
//@return    	nil
func (v *vector) PushBack(e interface{}) {
	if v == nil {
		return
	}
	v.mutex.Lock()
	if v.end < len(v.data) {
		v.data[v.end] = e
	} else {
		v.data = append(v.data, e)
	}
	v.end++
	v.mutex.Unlock()
}

//@title    PopBack
//@description
//		以vector向量容器做接收者
//		弹出容器最后一个元素,同时尾指针前移一位
//		当尾指针小于容器切片实际使用空间的一半时,重新分配空间释放未使用部分
//		若容器为空,则不进行弹出
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	nil
func (v *vector) PopBack() {
	if v == nil {
		return
	}
	if v.Empty() {
		return
	}
	v.mutex.Lock()
	v.end--
	if v.end*2 <= len(v.data) {
		v.data = v.data[0:v.end]
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
//		idx从0计算
//		尾指针同步后移一位
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	idx			int						待插入节点的位置
//@param		e			interface{}				待插入元素
//@return    	nil
func (v *vector) Insert(idx int, e interface{}) {
	if v == nil {
		return
	}
	v.mutex.Lock()
	if idx <= 0 {
		v.data = append(append([]interface{}{}, e), v.data[:v.end]...)
		v.end++
	} else if idx >= v.Size() {
		v.PushBack(e)
	} else {
		es := append([]interface{}{}, v.data[idx:v.end]...)
		v.data = append(append(v.data[:idx], e), es...)
		v.end++
	}
	v.mutex.Unlock()
}

//@title    Erase
//@description
//		以vector向量容器做接收者
//		向容器切片中删除一个元素
//		当idx不大于0时,则在容器切片的头部删除
//		当idx不小于切片使用长度时,在容器末尾删除
//		否则在切片中间第idx位删除元素,同时前移第idx位以后的元素
//		idx从0计算
//		尾指针同步后移以为
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	idx			int						待删除节点的位置
//@return    	nil
func (v *vector) Erase(idx int) {
	if v == nil {
		return
	}
	if v.Empty() {
		return
	}
	v.mutex.Lock()
	idx++
	if idx <= 1 {
		idx = 1
	} else if idx >= v.Size() {
		idx = v.Size()
	}
	es := append([]interface{}{}, v.data[:idx-1]...)
	v.data = append(es, v.data[idx:]...)
	v.end--
	v.mutex.Unlock()
}

//@title    Reverse
//@description
//		以vector向量容器做接收者
//		将vector容器中不使用空间释放掉
//		将该容器中的泛型切片中的所有元素顺序逆转
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param		nil
//@return    	nil
func (v *vector) Reverse() {
	if v == nil {
		return
	}
	v.mutex.Lock()
	if v.end > 0 {
		v.data = v.data[:v.end]
	} else {
		v.data = make([]interface{}, 0, 0)
	}
	for i := 0; i < v.end/2; i++ {
		v.data[i], v.data[v.end-i-1] = v.data[v.end-i-1], v.data[i]
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
//		idx从0计算
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	idx			int						待查找元素的位置
//@return    	e			interface{}				从容器中查找的第idx位元素
func (v *vector) At(idx int) (e interface{}) {
	if v == nil {
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
//@auth      	hlccd		2021-07-4
//@receiver		v			*vector					接受者vector的指针
//@param    	nil
//@return    	e			interface{}				容器的第一个元素
func (v *vector) Front() (e interface{}) {
	if v == nil {
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
		return nil
	}
	v.mutex.Lock()
	if v.Size() > 0 {
		e = v.data[v.end-1]
		v.mutex.Unlock()
		return e
	}
	v.mutex.Unlock()
	return nil
}
