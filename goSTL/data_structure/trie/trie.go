package trie

//@Title		trie
//@Description
//		单词查找树-Trie
//		以多叉树的形式实现,本次实现中有64叉,即a~z,A~Z,0~9,'+','/'共64个,对应base64的字符
//		当存储的string中出现其他字符则无法存储
//		存储的string可以携带一个元素
//		结点不允许覆盖,即插入值已经存在时会插入失败,需要先删除原值
//		使用互斥锁实现并发控制

import (
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//trie单词查找树结构体
//该实例存储单词查找树的根节点
//同时保存该树已经存储了多少个元素
//整个树不允许重复插入,若出现重复插入则直接失败
type trie struct {
	root  *node      //根节点指针
	size  int        //存放的元素数量
	mutex sync.Mutex //并发控制锁
}

//trie单词查找树容器接口
//存放了trie单词查找树可使用的函数
//对应函数介绍见下方
type trieer interface {
	Iterator() (i *Iterator.Iterator)        //返回包含该trie的所有string
	Size() (num int)                         //返回该trie中保存的元素个数
	Clear()                                  //清空该trie
	Empty() (b bool)                         //判断该trie是否为空
	Insert(s string, e interface{}) (b bool) //向trie中插入string并携带元素e
	Erase(s string) (b bool)                 //从trie中删除以s为索引的元素e
	Count(s string) (num int)                //从trie中寻找以s为前缀的string单词数
	Find(s string) (e interface{})           //从trie中寻找以s为索引的元素e
}

//@title    New
//@description
//		新建一个trie单词查找树容器并返回
//		初始根节点为nil
//@receiver		nil
//@param    	nil
//@return    	t        	*trie						新建的trie指针
func New() (t *trie) {
	return &trie{
		root:  newNode(nil),
		size:  0,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以trie单词查找树做接收者
//		将该trie中所有存放的string放入迭代器中并返回
//@receiver		t			*trie					接受者trie的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (t *trie) Iterator() (i *Iterator.Iterator) {
	if t == nil {
		return nil
	}
	t.mutex.Lock()
	//找到trie中存在的所有string
	es := t.root.inOrder("")
	i = Iterator.New(&es)
	t.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以trie单词查找树做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回0
//@receiver		t			*trie					接受者trie的指针
//@param    	nil
//@return    	num        	int						容器中实际使用元素所占空间大小
func (t *trie) Size() (num int) {
	if t == nil {
		return 0
	}
	return t.size
}

//@title    Clear
//@description
//		以trie单词查找树做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0
//@receiver		t			*trie					接受者trie的指针
//@param    	nil
//@return    	nil
func (t *trie) Clear() {
	if t == nil {
		return
	}
	t.mutex.Lock()
	t.root = newNode(nil)
	t.size = 0
	t.mutex.Unlock()
}

//@title    Empty
//@description
//		以trie单词查找树做接收者
//		判断该trie是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		t			*trie					接受者trie的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (t *trie) Empty() (b bool) {
	if t.Size() > 0 {
		return false
	}
	return true
}

//@title    Insert
//@description
//		以trie单词查找树做接收者
//		向trie插入以string类型的s为索引的元素e
//		若存在重复的s则插入失败,不允许覆盖
//		否则插入成功
//@receiver		t			*trie					接受者trie的指针
//@param    	s			string					待插入元素的索引s
//@param    	e			interface{}				待插入元素e
//@return    	b			bool					添加成功?
func (t *trie) Insert(s string, e interface{}) (b bool) {
	if t == nil {
		return
	}
	if len(s) == 0 {
		return false
	}
	t.mutex.Lock()
	if t.root == nil {
		//避免根节点为nil
		t.root = newNode(nil)
	}
	//从根节点开始插入
	b = t.root.insert(s, 0, e)
	if b {
		//插入成功,size+1
		t.size++
	}
	t.mutex.Unlock()
	return b
}

//@title    Erase
//@description
//		以trie单词查找树做接收者
//		从trie树中删除元素以s为索引的元素e
//@receiver		t			*trie					接受者trie的指针
//@param    	s			string					待删除元素的索引
//@return    	b			bool					删除成功?
func (t *trie) Erase(s string) (b bool) {
	if t == nil {
		return false
	}
	if t.Empty() {
		return false
	}
	if len(s) == 0 {
		//长度为0无法删除
		return false
	}
	if t.root == nil {
		//根节点为nil即无法删除
		return false
	}
	t.mutex.Lock()
	//从根节点开始删除
	b = t.root.erase(s, 0)
	if b {
		//删除成功,size-1
		t.size--
		if t.size == 0 {
			//所有string都被删除,根节点置为nil
			t.root = nil
		}
	}
	t.mutex.Unlock()
	return b
}

//@title    Count
//@description
//		以trie单词查找树做接收者
//		从trie中查找以s为前缀的所有string的个数
//		如果存在以s为前缀的则返回大于0的值即其数量
//		如果未找到则返回0
//@receiver		t			*trie					接受者trie的指针
//@param    	s			string					待查找的前缀s
//@return    	num			int						待查找前缀在trie中存在的数量
func (t *trie) Count(s string) (num int) {
	if t == nil {
		return 0
	}
	if t.Empty() {
		return 0
	}
	if t.root == nil {
		return 0
	}
	t.mutex.Lock()
	//统计所有以s为前缀的string的数量并返回
	num = int(t.root.count(s, 0))
	t.mutex.Unlock()
	return num
}

//@title    Find
//@description
//		以trie单词查找树做接收者
//		从trie中查找以s为索引的元素e,找到则返回e
//		如果未找到则返回nil
//@receiver		t			*trie					接受者trie的指针
//@param    	s			string					待查找索引s
//@return    	ans			interface{}				待查找索引所指向的元素
func (t *trie) Find(s string) (e interface{}) {
	if t == nil {
		return nil
	}
	if t.Empty() {
		return nil
	}
	if t.root == nil {
		return nil
	}
	t.mutex.Lock()
	//从根节点开始查找以s为索引的元素e
	e = t.root.find(s, 0)
	t.mutex.Unlock()
	return e
}
