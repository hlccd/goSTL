github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		哈希映射（hash map），它是一个两层结构，即第一层以动态数组作为桶去存储元素，第二层存储hash值冲突的元素。

​		对于插入其中的任意一个元素来说，都可以计算其key的hash值然后将其映射到桶内对应位置，随后再插入即可。

​		hash映射最大的特点在于其查找、插入和删除都是O（1）的，但可能存在扩容和缩容的问题，此时其时间开销会增加。

### 原理

​		对于哈希映射来说，它需要做的主要是对key进行hash计算，将计算值映射到hash桶的对应位置，插入其中。

​		本次实现中，动态数组即hash桶复用之前实现过的vector，而当出现hash冲突时候以avl树存储存在冲突的值。

​		对于其桶的扩容和缩容来说，都可以依靠vector来完成，扩容缩容后重新将所有key-value插入即可。

​		当然，对于hash映射来说，最优的情况就是所有hash都不会冲突，即avl树仅有根节点，这样对于插入和删除都会更快一些，所以就需要尽可能的避免出现hash冲突的情况，即设计一个良好的hash函数。

​		一般来说，hash函数可以采用**平方折半法、逐位累加法**，考虑到数值类的遍历在诸位累计会消耗大量时间，故对于数值类型采用平方折半法，对于string类型采用逐位累加法。

​		vector的初始容量和最小容量为**16**，该值为经验值，设为16的时候hash冲突情况最小。

​		而对于扩容和缩容来说，需要一个**扩容因子**，即当前存储容量达到总容量的扩容因子的倍数时候需要扩容，反之当当前存储容量小于缩容后的扩容因子的倍数时则需要缩容。对于扩容因子来说，若设置的过小，虽然冲突的概率降低，但空间利用率也急剧下降，而当设置过大时，空间利用率提高了，但hash冲突的情况增多，需要的时间成本提高。本次实现中将扩容因子设为**0.75**。该值是根据**泊松分布和指数分布**计算出的出现冲突的可能性最小的最优解。

### 实现

​		hashMap哈希映射结构体，该实例存储第一层vector的指针，同时保存hash函数，哈希映射中的hash函数在创建时传入,若不传入则在插入首个key-value时从默认hash函数中寻找。

```go
type hashMap struct {
	arr   *vector.Vector   //第一层的vector
	hash  algorithm.Hasher //hash函数
	size  uint64           //当前存储数量
	cap   uint64           //vector的容量
	mutex sync.Mutex       //并发控制锁
}
```

​		索引结构体，存储key-value结构。

```go
type node struct {
	value interface{} //节点中存储的元素
	num   int         //该元素数量
	depth int         //该节点的深度
	left  *node       //左节点指针
	right *node       //右节点指针
}
```

#### 接口

```go
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
```

#### New

​		新建一个hashMap哈希映射容器并返回，初始vector长度为16，若有传入的hash函数,则将传入的第一个hash函数设为该hash映射的hash函数。

```go
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
```

#### Iterator

​		以hashMap哈希映射做接收者，将该hashMap中所有保存的索引指针的value放入迭代器中。

```go
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
```

#### Size

​		以hashMap哈希映射做接收者，返回该容器当前含有元素的数量，如果容器为nil返回0。

```go
func (hm *hashMap) Size() (num uint64) {
	if hm == nil {
		return 0
	}
	return hm.size
}
```

#### Cap

​		以hashMap哈希映射做接收者，返回该容器当前容量，如果容器为nil返回0。

```go
func (hm *hashMap) Cap() (num uint64) {
	if hm == nil {
		return 0
	}
	return hm.cap
}
```

#### Clear

​		以hashMap哈希映射做接收者，将该容器中所承载的元素清空，将该容器的size置0,容量置为16,vector重建并扩容到16。

```go
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
```

#### Empty

​		以hashMap哈希映射做接收者，判断该哈希映射中是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true。

```go
func (hm *hashMap) Empty() (b bool) {
	if hm == nil {
		return false
	}
	return hm.size > 0
}
```

#### Insert

​		以hashMap哈希映射做接收者，向哈希映射入元素e,若已存在相同key则进行覆盖,覆盖仍然视为插入成功，若不存在相同key则直接插入即可。

```go
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
```

##### expend

​		以hashMap哈希映射做接收者，对原vector进行扩容，将所有的key-value取出,让vector自行扩容并清空原有结点，扩容后将所有的key-value重新插入vector中。

```go
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
```

#### Erase

​		以hashMap哈希映射做接收者，从hashMap中删除以key为索引的value，若存在则删除,否则直接结束,删除成功后size-1，删除后可能出现size<0.75/2*cap且大于16,此时要缩容。

```go
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
```

##### shrink

​		以hashMap哈希映射做接收者，对原vector进行缩容，将所有的key-value取出,让vector自行缩容并清空所有结点，当vector容量与缩容开始时不同时则视为缩容结束，随容后将所有的key-value重新插入vector中。

```go
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
```

#### GetKeys

​		以hashMap哈希映射做接收者，返回该hashMap中所有的key。

```go
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
```

#### Get

​		以hashMap哈希映射做接收者，以key寻找到对应的value并返回。

```go
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
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/hashMap"
)

func main() {
	m:=hashMap.New()
	for i:=1;i<=17;i++{
		m.Insert(i,i)
	}
	fmt.Println("size=",m.Size())
	keys:=m.GetKeys()
	fmt.Println(keys)
	for i:=0;i< len(keys);i++{
		fmt.Println(m.Get(keys[i]))
	}
	for i:=m.Iterator().Begin();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	for i:=0;i< len(keys);i++{
		m.Erase(keys[i])
	}
}

```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> size= 17
> [1 8 16 2 14 3 4 9 12 5 15 17 6 10 13 7 11]
> 1
> 8
> 16
> 2
> 14
> 3
> 4
> 9
> 12
> 5
> 15
> 17
> 6
> 10
> 13
> 7
> 11
> 1
> 8
> 16
> 2
> 14
> 3
> 4
> 9
> 12
> 5
> 15
> 17
> 6
> 10
> 13
> 7
> 11

#### 时间开销

注：插入的是随机数，即可能出现重复插入覆盖的问题，所以数值会≤max

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/hashMap"
	"math/rand"
	"time"
)

func main() {
	max := 3000000
	num := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tm := time.Now()
	m := make(map[int]bool)
	for i := 0; i < max; i++ {
		num=r.Intn(4294967295)
		m[num] = true
	}
	fmt.Println("map消耗时间:", time.Since(tm))
	thm := time.Now()
	hm := hashMap.New()
	for i := 0; i < max; i++ {
		num=r.Intn(4294967295)
		hm.Insert(uint64(num), true)
	}
	fmt.Println("hashmap添加消耗时间:", time.Since(thm))
	thmd := time.Now()
	keys:=hm.GetKeys()
	for i := 0; i < len(keys); i++ {
		hm.Get(keys[i])
	}
	fmt.Println("size=",hm.Size())
	fmt.Println("hashmap遍历消耗时间:", time.Since(thmd))
}
```

> map消耗时间: 484.6992ms
> hashmap添加消耗时间: 4.6206379s
> size= 2999043
> hashmap遍历消耗时间: 2.2020792s
