package rax

//@Title		rax
//@Description
//		前缀基数树-Rax
//		以多叉树的形式实现,根据'/'进行string分割,将分割后的string数组进行段存储
//		插入的string首字符必须为'/'
//		任意一段string不能为""
//		结点不允许覆盖,即插入值已经存在时会插入失败,需要先删除原值
//		使用互斥锁实现并发控制

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"strings"
	"sync"
)

//rax前缀基数树结构体
//该实例存储前缀基数树的根节点
//同时保存该树已经存储了多少个元素
//整个树不允许重复插入,若出现重复插入则直接失败
type rax struct {
	root  *node      //前缀基数树的根节点指针
	size  int        //当前已存放的元素数量
	mutex sync.Mutex //并发控制锁
}

//rax前缀基数树容器接口
//存放了rax前缀基数树可使用的函数
//对应函数介绍见下方
type raxer interface {
	Iterator() (i *Iterator.Iterator)        //返回包含该rax的所有string
	Size() (num int)                         //返回该rax中保存的元素个数
	Clear()                                  //清空该rax
	Empty() (b bool)                         //判断该rax是否为空
	Insert(s string, e interface{}) (b bool) //向rax中插入string并携带元素e
	Erase(s string) (b bool)                 //从rax中删除以s为索引的元素e
	Delete(s string) (num int)               //从rax中删除以s为前缀的所有元素
	Count(s string) (num int)                //从rax中寻找以s为前缀的string单词数
	Find(s string) (e interface{})           //从rax中寻找以s为索引的元素e
}

//@title    New
//@description
//		新建一个rax前缀基数树容器并返回
//		初始根节点为nil
//@receiver		nil
//@param    	nil
//@return    	r        	*rax						新建的rax指针
func New() (r *rax) {
	return &rax{
		root:  newNode("", nil),
		size:  0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以rax前缀基数树做接收者
//		将该rax中所有存放的string放入迭代器中并返回
//@receiver		r			*rax					接受者rax的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (r *rax) Iterator() (i *Iterator.Iterator) {
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
//		以rax前缀基数树做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回0
//@receiver		r			*rax					接受者rax的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (r *rax) Size() (num int) {
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
//		以rax前缀基数树做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@receiver		r			*rax					接受者rax的指针
//@param    	nil
//@return    	nil
func (r *rax) Clear() {
	if r == nil {
		return
	}
	r.mutex.Lock()
	r.root = newNode("", nil)
	r.size = 0
	r.mutex.Unlock()
}

//@title    Empty
//@description
//		以rax前缀基数树做接收者
//		判断该rax是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		r			*rax					接受者rax的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (r *rax) Empty() (b bool) {
	if r == nil {
		return true
	}
	return r.size == 0
}

//@title    Insert
//@description
//		以rax前缀基数树做接收者
//		向rax插入以string类型的s为索引的元素e
//		插入的string的首字符必须为'/',且中间按'/'分割的string不能为""
//		若存在重复的s则插入失败,不允许覆盖
//		否则插入成功
//@receiver		r			*rax					接受者rax的指针
//@param    	s			string					待插入元素的索引s
//@param    	e			interface{}				待插入元素e
//@return    	b			bool					添加成功?
func (r *rax) Insert(s string, e interface{}) (b bool) {
	if r == nil {
		return false
	}
	if len(s) == 0 {
		return false
	}
	if s[0] != '/' {
		return false
	}
	//将s按'/'进行分割,并去掉第一个即去掉"",随后一次按照分层结果进行插入
	ss := strings.Split(s, "/")[1:]
	r.mutex.Lock()
	if r.root == nil {
		//避免根节点为nil
		r.root = newNode("", nil)
	}
	//从根节点开始插入
	b = r.root.insert(ss, 0, e)
	if b {
		//插入成功,size+1
		r.size++
	}
	r.mutex.Unlock()
	return b
}

//@title    Erase
//@description
//		以rax前缀基数树做接收者
//		从rax树中删除元素以s为索引的元素e
//		用以删除的string索引的首字符必须为'/'
//@receiver		r			*rax					接受者rax的指针
//@param    	s			string					待删除元素的索引
//@return    	b			bool					删除成功?
func (r *rax) Erase(s string) (b bool) {
	if r.Empty() {
		return false
	}
	if len(s) == 0 {
		return false
	}
	if s[0] != '/' {
		return false
	}
	//将s按'/'进行分割,并去掉第一个即去掉"",随后一次按照分层结果进行删除
	ss := strings.Split(s, "/")[1:]
	if r.root == nil {
		//根节点为nil即无法删除
		return false
	}
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
//		以rax前缀基数树做接收者
//		从rax树中删除以s为前缀的所有元素
//		用以删除的string索引的首字符必须为'/'
//@receiver		r			*rax					接受者rax的指针
//@param    	s			string					待删除元素的前缀
//@return    	num			int						被删除的元素的数量
func (r *rax) Delete(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	if s[0] != '/' {
		return 0
	}
	if r.root == nil {
		return 0
	}
	//将s按'/'进行分割,并去掉第一个即去掉"",随后一次按照分层结果进行删除
	ss := strings.Split(s, "/")[1:]
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
//		以rax前缀基数树做接收者
//		从rax中查找以s为前缀的所有string的个数
//		用以查找的string索引的首字符必须为'/'
//		如果存在以s为前缀的则返回大于0的值即其数量
//		如果未找到则返回0
//@receiver		r			*rax					接受者rax的指针
//@param    	s			string					待查找的前缀s
//@return    	num			int						待查找前缀在rax树中存在的数量
func (r *rax) Count(s string) (num int) {
	if r.Empty() {
		return 0
	}
	if r.root == nil {
		return 0
	}
	if len(s) == 0 {
		return 0
	}
	if s[0] != '/' {
		return 0
	}
	//将s按'/'进行分割,并去掉第一个即去掉"",随后一次按照分层结果进行查找
	ss := strings.Split(s, "/")[1:]
	r.mutex.Lock()
	num = r.root.count(ss, 0)
	r.mutex.Unlock()
	return num
}

//@title    Find
//@description
//		以rax前缀基数树做接收者
//		从rax中查找以s为索引的元素e,找到则返回e
//		用以查找的string索引的首字符必须为'/'
//		如果未找到则返回nil
//@receiver		r			*rax					接受者rax的指针
//@param    	s			string					待查找索引s
//@return    	ans			interface{}				待查找索引所指向的元素
func (r *rax) Find(s string) (e interface{}) {
	if r.Empty() {
		return nil
	}
	if len(s) == 0 {
		return nil
	}
	if s[0] != '/' {
		return nil
	}
	if r.root == nil {
		return nil
	}
	//将s按'/'进行分割,并去掉第一个即去掉"",随后一次按照分层结果进行查找
	ss := strings.Split(s, "/")[1:]
	r.mutex.Lock()
	e = r.root.find(ss, 0)
	r.mutex.Unlock()
	return e
}
