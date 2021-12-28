package radix

//@Title		radix
//@Description
//		前缀基数树-radix
//		以多叉树的形式实现,根据'/'进行string分割,将分割后的string数组进行段存储
//		不存储其他元素,仅对string进行分段存储和模糊匹配
//		使用互斥锁实现并发控制

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//radix前缀基数树结构体
//该实例存储前缀基数树的根节点
//同时保存该树已经存储了多少个元素
type radix struct {
	root  *node      //前缀基数树的根节点指针
	size  int        //当前已存放的元素数量
	mutex sync.Mutex //并发控制锁
}

//radix前缀基数树容器接口
//存放了radix前缀基数树可使用的函数
//对应函数介绍见下方
type radixer interface {
	Iterator() (i *Iterator.Iterator)             //返回包含该radix的所有string
	Size() (num int)                              //返回该radix中保存的元素个数
	Clear()                                       //清空该radix
	Empty() (b bool)                              //判断该radix是否为空
	Insert(s string) (b bool)                     //向radix中插入string
	Erase(s string) (b bool)                      //从radix中删除string
	Delete(s string) (num int)                    //从radix中删除以s为前缀的所有string
	Count(s string) (num int)                     //从radix中寻找以s为前缀的string单词数
	Mate(s string) (m map[string]string, ok bool) //利用radix树中的string对s进行模糊匹配,':'可模糊匹配该层,'*'可模糊匹配后面所有
}

//@title    New
//@description
//		新建一个radix前缀基数树容器并返回
//		初始根节点为nil
//@receiver		nil
//@param    	nil
//@return    	r        	*radix						新建的radix指针
func New() (r *radix) {
	return &radix{
		root:  newNode(""),
		size:  0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以radix前缀基数树做接收者
//		将该radix中所有存放的string放入迭代器中并返回
//@receiver		r			*radix					接受者radix的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (r *radix) Iterator() (i *Iterator.Iterator) {
	if r == nil {
		return nil
	}
	r.mutex.Lock()
	es := r.root.inOrder("")
	i = Iterator.New(&es)
	r.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以radix前缀基数树做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回0
//@receiver		r			*radix					接受者radix的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (r *radix) Size() (num int) {
	if r == nil {
		return 0
	}
	if r.root == nil {
		return 0
	}
	return r.size
}

//@title    Clear
//@description
//		以radix前缀基数树做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@receiver		r			*radix					接受者radix的指针
//@param    	nil
//@return    	nil
func (r *radix) Clear() {
	if r == nil {
		return
	}
	r.mutex.Lock()
	r.root = newNode("")
	r.size = 0
	r.mutex.Unlock()
}

//@title    Empty
//@description
//		以radix前缀基数树做接收者
//		判断该radix是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		r			*radix					接受者radix的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (r *radix) Empty() (b bool) {
	if r == nil {
		return true
	}
	return r.size == 0
}

//@title    Insert
//@description
//		以radix前缀基数树做接收者
//		向radix插入string
//		将对string进行解析,按'/'进行分层,':'为首则为模糊匹配该层,'*'为首则为模糊匹配后面所有
//		已经存在则无法重复插入
//@receiver		r			*radix					接受者radix的指针
//@param    	s			string					待插入string
//@return    	b			bool					添加成功?
func (r *radix) Insert(s string) (b bool) {
	if r == nil {
		return false
	}
	//解析s并按规则重构s
	ss, s := analysis(s)
	r.mutex.Lock()
	if r.root == nil {
		//避免根节点为nil
		r.root = newNode("")
	}
	//从根节点开始插入
	b = r.root.insert(s, ss, 0)
	if b {
		//插入成功,size+1
		r.size++
	}
	r.mutex.Unlock()
	return b
}

//@title    Erase
//@description
//		以radix前缀基数树做接收者
//		从radix树中删除元素string
//@receiver		r			*radix					接受者radix的指针
//@param    	s			string					待删除的string
//@return    	b			bool					删除成功?
func (r *radix) Erase(s string) (b bool) {
	if r.Empty() {
		return false
	}
	if len(s) == 0 {
		return false
	}
	if r.root == nil {
		//根节点为nil即无法删除
		return false
	}
	//解析s并按规则重构s
	ss, _ := analysis(s)
	r.mutex.Lock()
	//从根节点开始删除
	b = r.root.erase(ss, 0)
	if b {
		//删除成功,size-1
		r.size--
		if r.size == 0 {
			//所有string都被删除,根节点置为nil
			r.root = nil
		}
	}
	r.mutex.Unlock()
	return b
}

//@title    Delete
//@description
//		以radix前缀基数树做接收者
//		从radix树中删除以s为前缀的所有string
//@receiver		r			*radix					接受者radix的指针
//@param    	s			string					待删除string的前缀
//@return    	num			int						被删除的元素的数量
func (r *radix) Delete(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	if r.root == nil {
		return 0
	}
	//解析s并按规则重构s
	ss, _ := analysis(s)
	r.mutex.Lock()
	//从根节点开始删除
	num = r.root.delete(ss, 0)
	if num > 0 {
		//删除成功
		r.size -= num
		if r.size <= 0 {
			//所有string都被删除,根节点置为nil
			r.root = nil
		}
	}
	r.mutex.Unlock()
	return num
}

//@title    Count
//@description
//		以radix前缀基数树做接收者
//		从radix中查找以s为前缀的所有string的个数
//		如果存在以s为前缀的则返回大于0的值即其数量
//		如果未找到则返回0
//@receiver		r			*radix					接受者radix的指针
//@param    	s			string					待查找的前缀s
//@return    	num			int						待查找前缀在radix树中存在的数量
func (r *radix) Count(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if r.root == nil {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	//解析s并按规则重构s
	ss, _ := analysis(s)
	r.mutex.Lock()
	num = r.root.count(ss, 0)
	r.mutex.Unlock()
	return num
}

//@title    Mate
//@description
//		以radix前缀基数树做接收者
//		从radix中查找以s为信息的第一个可以模糊匹配到的key和value的映射表
//		key是radix树中的段名,value是s中的段名
//		如果未找到则返回nil和false
//		否则返回一个映射表和true
//@receiver		r			*radix					接受者radix的指针
//@param    	s			string					待查找的信息s
//@return    	m			map[string]string		s从前缀基数树中利用模糊匹配到的所有key和value的映射
//@return    	ok			bool					匹配成功?
func (r *radix) Mate(s string) (m map[string]string, ok bool) {
	if r.Empty() {
		return nil, false
	}
	if len(s) == 0 {
		return nil, false
	}
	if r.root == nil {
		return nil, false
	}
	//将s按'/'进行分割,并去掉第一个即去掉"",随后一次按照分层结果进行查找
	m, ok = r.root.mate(s, 0)
	return m, ok
}
