package hashMap

//@Title		hashMap
//@Description
//		哈希映射-hash map
//		分为两层,第一层以vector实现,当出现hash冲突时以avl树存储即第二层
//		若为基本数据类型可不用传入hash函数,否则需要传入自定义hash函数
//		所有key不可重复,但value可重复,key不应为nil
//		扩容因子为0.75,当存储数超过0.75倍总容量时应当扩容
//		使用互斥锁实现并发控制
import (
	"github.com/hlccd/goSTL/algorithm"
	"github.com/hlccd/goSTL/data_structure/avlTree"
	"github.com/hlccd/goSTL/data_structure/vector"
	"github.com/hlccd/goSTL/utils/comparator"
	"github.com/hlccd/goSTL/utils/iterator"
	"sync"
)

//hashMap哈希映射结构体
//该实例存储第一层vector的指针
//同时保存hash函数
//哈希映射中的hash函数在创建时传入,若不传入则在插入首个key-value时从默认hash函数中寻找
type hashMap struct {
	arr   *vector.Vector   //第一层的vector
	hash  algorithm.Hasher //hash函数
	size  uint64           //当前存储数量
	cap   uint64           //vector的容量
	mutex sync.Mutex       //并发控制锁
}

//索引结构体
//存储key-value结构
type indexes struct {
	key   interface{}
	value interface{}
}

//hashMap哈希映射容器接口
//存放了hashMap哈希映射可使用的函数
//对应函数介绍见下方
type hashMaper interface {
	Iterator() (i *Iterator.Iterator)        //返回一个包含hashMap容器中所有value的迭代器
	Size() (num uint64)                      //返回hashMap已存储的元素数量
	Cap() (num uint64)                       //返回hashMap中的存放空间的容量
	Clear()                                  //清空hashMap
	Empty() (b bool)                         //返回hashMap是否为空
	Insert(key, value interface{}) (b bool)  //向hashMap插入以key为索引的value,若存在会覆盖
	Erase(key interface{}) (b bool)          //删除hashMap中以key为索引的value
	GetKeys() (keys []interface{})           //返回hashMap中所有的keys
	Get(key interface{}) (value interface{}) //以key为索引寻找vlue
}

//@title    New
//@description
//		新建一个hashMap哈希映射容器并返回
//		初始vector长度为16
//		若有传入的hash函数,则将传入的第一个hash函数设为该hash映射的hash函数
//@receiver		nil
//@param    	Cmp			...algorithm.Hasher		hashMap的hash函数集
//@return    	hm			*hashMap				新建的hashMap指针
func New(hash ...algorithm.Hasher) (hm *hashMap) {
	var h algorithm.Hasher
	if len(hash) == 0 {
		h = nil
	} else {
		h = hash[0]
	}
	cmp := func(a, b interface{}) int {
		ka, kb := a.(*indexes), b.(*indexes)
		return comparator.GetCmp(ka.key)(ka.key, kb.key)
	}
	//新建vector并将其扩容到16
	v := vector.New()
	for i := 0; i < 16; i++ {
		//vector中嵌套avl树
		v.PushBack(avlTree.New(false, cmp))
	}
	return &hashMap{
		arr:   v,
		hash:  h,
		size:  0,
		cap:   16,
		mutex: sync.Mutex{},
	}
}

//@title    Iterator
//@description
//		以hashMap哈希映射做接收者
//		将该hashMap中所有保存的索引指针的value放入迭代器中
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	i        	*iterator.Iterator		新建的Iterator迭代器指针
func (hm *hashMap) Iterator() (i *Iterator.Iterator) {
	if hm == nil {
		return nil
	}
	if hm.arr == nil {
		return nil
	}
	hm.mutex.Lock()
	//取出hashMap中存放的所有value
	values := make([]interface{}, 0, 1)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avlTree.AvlTree)
		ite := avl.Iterator()
		es := make([]interface{}, 0, 1)
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idx := j.Value().(*indexes)
			es = append(es, idx.value)
		}
		values = append(values, es...)
	}
	//将所有value放入迭代器中
	i = Iterator.New(&values)
	hm.mutex.Unlock()
	return i
}

//@title    Size
//@description
//		以hashMap哈希映射做接收者
//		返回该容器当前含有元素的数量
//		如果容器为nil返回0
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	num        	uint64					当前存储的元素数量
func (hm *hashMap) Size() (num uint64) {
	if hm == nil {
		return 0
	}
	return hm.size
}

//@title    Cap
//@description
//		以hashMap哈希映射做接收者
//		返回该容器当前容量
//		如果容器为nil返回0
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	num        	int						容器中实际占用的容量大小
func (hm *hashMap) Cap() (num uint64) {
	if hm == nil {
		return 0
	}
	return hm.cap
}

//@title    Clear
//@description
//		以hashMap哈希映射做接收者
//		将该容器中所承载的元素清空
//		将该容器的size置0,容量置为16,vector重建并扩容到16
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	nil
func (hm *hashMap) Clear() {
	if hm == nil {
		return
	}
	hm.mutex.Lock()
	//重建vector并扩容到16
	v := vector.New()
	cmp := func(a, b interface{}) int {
		ka, kb := a.(*indexes), b.(*indexes)
		return comparator.GetCmp(ka.key)(ka.key, kb.key)
	}
	for i := 0; i < 16; i++ {
		v.PushBack(avlTree.New(false, cmp))
	}
	hm.arr = v
	hm.size = 0
	hm.cap = 16
	hm.mutex.Unlock()
}

//@title    Empty
//@description
//		以hashMap哈希映射做接收者
//		判断该哈希映射中是否含有元素
//		如果含有元素则不为空,返回false
//		如果不含有元素则说明为空,返回true
//		如果容器不存在,返回true
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	b			bool					该容器是空的吗?
func (hm *hashMap) Empty() (b bool) {
	if hm == nil {
		return false
	}
	return hm.size > 0
}

//@title    Insert
//@description
//		以hashMap哈希映射做接收者
//		向哈希映射入元素e,若已存在相同key则进行覆盖,覆盖仍然视为插入成功
//		若不存在相同key则直接插入即可
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	key			interface{}				待插入元素的key
//@param    	value		interface{}				待插入元素的value
//@return    	b			bool					添加成功?
func (hm *hashMap) Insert(key, value interface{}) (b bool) {
	if hm == nil {
		return false
	}
	if hm.arr == nil {
		return false
	}
	if hm.hash == nil {
		hm.hash = algorithm.GetHash(key)
	}
	if hm.hash == nil {
		return false
	}
	hm.mutex.Lock()
	//计算hash值并找到对应的avl树
	hash := hm.hash(key) % hm.cap
	avl := hm.arr.At(hash).(*avlTree.AvlTree)
	idx := &indexes{
		key:   key,
		value: value,
	}
	//判断是否存在该avl树中
	if avl.Count(idx) == 0 {
		//avl树中不存在相同key,插入即可
		avl.Insert(idx)
		hm.size++
		if hm.size >= hm.cap/4*3 {
			//当达到扩容条件时候进行扩容
			hm.expend()
		}
	} else {
		//覆盖
		avl.Insert(idx)
	}
	hm.mutex.Unlock()
	return true
}

//@title    expend
//@description
//		以hashMap哈希映射做接收者
//		对原vector进行扩容
//		将所有的key-value取出,让vector自行扩容并清空原有结点
//		扩容后将所有的key-value重新插入vector中
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	nil
func (hm *hashMap) expend() {
	//取出所有的key-value
	idxs := make([]*indexes, 0, hm.size)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avlTree.AvlTree)
		ite := avl.Iterator()
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idxs = append(idxs, j.Value().(*indexes))
		}
	}
	cmp := func(a, b interface{}) int {
		ka, kb := a.(*indexes), b.(*indexes)
		return comparator.GetCmp(ka.key)(ka.key, kb.key)
	}
	//对vector进行扩容,扩容到其容量上限即可
	hm.arr.PushBack(avlTree.New(false, cmp))
	for i := uint64(0); i < hm.arr.Size()-1; i++ {
		hm.arr.At(i).(*avlTree.AvlTree).Clear()
	}
	for i := hm.arr.Size(); i < hm.arr.Cap(); i++ {
		hm.arr.PushBack(avlTree.New(false, cmp))
	}
	//将vector容量设为hashMap容量
	hm.cap = hm.arr.Cap()
	//重新将所有的key-value插入到hashMap中去
	for i := 0; i < len(idxs); i++ {
		key, value := idxs[i].key, idxs[i].value
		hash := hm.hash(key) % hm.cap
		avl := hm.arr.At(hash).(*avlTree.AvlTree)
		idx := &indexes{
			key:   key,
			value: value,
		}
		avl.Insert(idx)
	}
}

//@title    Erase
//@description
//		以hashMap哈希映射做接收者
//		从hashMap中删除以key为索引的value
//		若存在则删除,否则直接结束,删除成功后size-1
//		删除后可能出现size<0.75/2*cap且大于16,此时要缩容
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	key			interface{}				待删除元素的key
//@return    	b			bool					删除成功?
func (hm *hashMap) Erase(key interface{}) (b bool) {
	if hm == nil {
		return false
	}
	if hm.arr == nil {
		return false
	}
	if hm.hash == nil {
		return false
	}
	hm.mutex.Lock()
	//计算该key的hash值
	hash := hm.hash(key) % hm.cap
	avl := hm.arr.At(hash).(*avlTree.AvlTree)
	idx := &indexes{
		key:   key,
		value: nil,
	}
	//从对应的avl树中删除该key-value
	b = avl.Erase(idx)
	if b {
		//删除成功,此时size-1,同时进行缩容判断
		hm.size--
		if hm.size < hm.cap/8*3 && hm.cap > 16 {
			hm.shrink()
		}
	}
	hm.mutex.Unlock()
	return b
}

//@title    shrink
//@description
//		以hashMap哈希映射做接收者
//		对原vector进行缩容
//		将所有的key-value取出,让vector自行缩容并清空所有结点
//		当vector容量与缩容开始时不同时则视为缩容结束
//		随容后将所有的key-value重新插入vector中
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	nil
func (hm *hashMap) shrink() {
	//取出所有key-value
	idxs := make([]*indexes, 0, hm.size)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avlTree.AvlTree)
		ite := avl.Iterator()
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idxs = append(idxs, j.Value().(*indexes))
		}
	}
	//进行缩容,当vector的cap与初始不同时,说明缩容结束
	cap := hm.arr.Cap()
	for ; cap == hm.arr.Cap(); {
		hm.arr.PopBack()
	}
	hm.cap = hm.arr.Cap()
	//将所有的key-value重新放入hashMap中
	for i := 0; i < len(idxs); i++ {
		key, value := idxs[i].key, idxs[i].value
		hash := hm.hash(key) % hm.cap
		avl := hm.arr.At(hash).(*avlTree.AvlTree)
		idx := &indexes{
			key:   key,
			value: value,
		}
		avl.Insert(idx)
	}
}

//@title    GetKeys
//@description
//		以hashMap哈希映射做接收者
//		返回该hashMap中所有的key
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	nil
//@return    	keys		[]interface{}			hashMap中的所有key
func (hm *hashMap) GetKeys() (keys []interface{}) {
	if hm == nil {
		return nil
	}
	if hm.arr == nil {
		return nil
	}
	hm.mutex.Lock()
	keys = make([]interface{}, 0, 1)
	for i := uint64(0); i < hm.arr.Size(); i++ {
		avl := hm.arr.At(i).(*avlTree.AvlTree)
		ite := avl.Iterator()
		es := make([]interface{}, 0, 1)
		for j := ite.Begin(); j.HasNext(); j.Next() {
			idx := j.Value().(*indexes)
			es = append(es, idx.key)
		}
		keys = append(keys, es...)
	}
	hm.mutex.Unlock()
	return keys
}

//@title    Get
//@description
//		以hashMap哈希映射做接收者
//		以key寻找到对应的value并返回
//@receiver		hm			*hashMap				接受者hashMap的指针
//@param    	keys		interface{}				待查找元素的key
//@return    	keys		interface{}				待查找元素的value
func (hm *hashMap) Get(key interface{}) (value interface{}) {
	if hm == nil {
		return
	}
	if hm.arr == nil {
		return
	}
	if hm.hash == nil {
		hm.hash = algorithm.GetHash(key)
	}
	if hm.hash == nil {
		return
	}
	hm.mutex.Lock()
	//计算hash值
	hash := hm.hash(key) % hm.cap
	//从avl树中找到对应该hash值的key-value
	info := hm.arr.At(hash).(*avlTree.AvlTree).Find(&indexes{key: key, value: nil})
	hm.mutex.Unlock()
	if info == nil {
		return nil
	}
	return info.(*indexes).value
}
