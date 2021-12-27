github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		LRU-Least Recently Used，最近最少使用链结构，相对于仅考虑时间因素的FIFO和仅考虑访问频率的LFU，LRU算法可以认为是相对平衡的一种淘汰算法，LRU认为，如果数据最近被访问过，那么将来被访问的概率也会更高，LRU 算法的实现非常简单，维护一个队列，如果某条记录被访问了，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。

​		它适用于作为一个cache缓存，去保存带有索引的数据，对于其要存储的数据值来说，通过实现一个Len()的函数的泛型来表示。

### 原理

​		本次实现过程中，采用双向链表的形式做为用于存储值的队列，同时利用map建立key和存储值的链表结点之间的联系。

​		该方案可以加快从队列中找到对应的值并做修改，同时也便于将最近访问的value移到队首。

​		本次实现不做定时淘汰，若要实现定时淘汰可自行创建一个协程同时添加一个map去保存其上次访问时间去完成定时淘汰即可。

​		对于LRU的修改策略来说：

- 初次添加/访问，直接放到队首
- 添加后若容量超过上限值则移除最远访问的value，即队尾元素，移除过程直至使用容量低于上限值时才结束
- 再次访问时从map中找到刚刚访问的结点随后将其移动到队首
- 删除时直接利用map找到对应的结点进行删除即可

### 实现

​		LRU链结构体，包含了该LRU结构中能承载的byte数上限和当前已承载的数量，以链表的形式存储承载元素，使用map建立key索引和链表结点之间的联系,链表结点中存放其value，onRemove函数是用于在删除时的执行,可用于将数据持久化，使用并发控制所保证线程安全。

```go
type LRU struct {
	maxBytes int64                         //所能承载的最大byte数量
	nowBytes int64                         //当前承载的byte数
	ll       *list.List                    //用于存储的链表
	cache    map[string]*list.Element      //链表元素与key的映射表
	onRemove func(key string, value Value) //删除元素时的执行函数
	mutex    sync.Mutex                    //并发控制锁
}
```

​		索引结构体，保存了一个待存储值的索引和值，value值为一个interface{},需要实现它的长度函数即Len()。

```go
type indexes struct {
	key   string //索引
	value Value  //存储值
}
```

​		存储值的函数。

```go
type Value interface {
   Len() int //用于计算该存储值所占用的长度
}
```

#### 接口

```go
type lruer interface {
	Size() (num int64)                     //返回lru中当前存放的byte数
	Cap() (num int64)                      //返回lru能存放的byte树的最大值
	Clear()                                //清空lru,将其中存储的所有元素都释放
	Empty() (b bool)                       //判断该lru中是否存储了元素
	Insert(key string, value Value)        //向lru中插入以key为索引的value
	Erase(key string)                      //从lru中删除以key为索引的值
	Get(key string) (value Value, ok bool) //从lru中获取以key为索引的value和是否获取成功?
}
```

#### New

​		以LRU链容器做接收者，返回该LRU链当前所存储的byte数。

```go
func (l *LRU) Size() (num int64) {
	if l == nil {
		return 0
	}
	return l.nowBytes
}
```

#### Cap

​		以LRU链容器做接收者，返回该LRU链能存储的最大byte数。

```go
func (l *LRU) Cap() (num int64) {
	if l == nil {
		return 0
	}
	return l.maxBytes
}
```

#### Clear

​		以LRU链容器做接收者，将LRU链中的所有承载的元素清除，同时将map清空,将当前存储的byte数清零。

```go
func (l *LRU) Clear() {
	if l == nil {
		return
	}
	l.mutex.Lock()
	l.ll = list.New()
	l.cache = make(map[string]*list.Element)
	l.nowBytes = 0
	l.mutex.Unlock()
}
```

#### Empty

​		以LRU链容器做接收者，判断该LRU链中是否存储了元素，没存储则返回true,LRU不存在也返回true，否则返回false。

```go
func (l *LRU) Empty() (b bool) {
	if l == nil {
		return true
	}
	return l.nowBytes <= 0
}
```

#### Insert

​		以LRU链容器做接收者，向该LRU中插入以key为索引的value，若已经存在则将其放到队尾，若空间充足且不存在则直接插入队尾，若空间不足则淘汰队首元素再将其插入队尾，插入完成后将当前已存的byte数增加。

```go
func (l *LRU) Insert(key string, value Value) {
	l.mutex.Lock()
	//利用map从已存的元素中寻找
	if ele, ok := l.cache[key]; ok {
		//该key已存在,直接替换即可
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*indexes)
		//此处是一个替换,即将cache中的value替换为新的value,同时根据实际存储量修改其当前存储的实际大小
		l.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//该key不存在,需要进行插入
		ele := l.ll.PushFront(&indexes{key, value})
		l.cache[key] = ele
		//此处是一个增加操作,即原本不存在,所以直接插入即可,同时在当前数值范围内增加对应的占用空间
		l.nowBytes += int64(len(key)) + int64(value.Len())
	}
	//添加完成后根据删除掉尾部一部分数据以保证实际使用空间小于设定的空间上限
	for l.maxBytes != 0 && l.maxBytes < l.nowBytes {
		//删除队尾元素
		//删除后执行创建时传入的删除执行函数
		ele := l.ll.Back()
		if ele != nil {
			l.ll.Remove(ele)
			kv := ele.Value.(*indexes)
			//删除最末尾的同时将其占用空间减去
			delete(l.cache, kv.key)
			l.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
			if l.onRemove != nil {
				//删除后的回调函数,可用于持久化该部分数据
				l.onRemove(kv.key, kv.value)
			}
		}
	}
	l.mutex.Unlock()
}
```

#### Erase

​		以LRU链容器做接收者，从LRU链中删除以key为索引的value，如果不存在则直接结束，删除完成后将其占用的byte数减去，删除后执行函数回调函数,一般用于持久化。

```go
func (l *LRU) Erase(key string) {
	l.mutex.Lock()
	if ele, ok := l.cache[key]; ok {
		l.ll.Remove(ele)
		kv := ele.Value.(*indexes)
		//删除的同时将其占用空间减去
		delete(l.cache, kv.key)
		l.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if l.onRemove != nil {
			//删除后的回调函数,可用于持久化该部分数据
			l.onRemove(kv.key, kv.value)
		}
	}
	l.mutex.Unlock()
}
```

#### Get

​		以LRU链容器做接收者，从LRU链中寻找以key为索引的value，找到对应的value后将其移到首部,随后返回其value，如果没找到以key为索引的value,直接结束即可。

```go
func (l *LRU) Get(key string) (value Value, ok bool) {
	l.mutex.Lock()
	if ele, ok := l.cache[key]; ok {
		//找到了value,将其移到链表首部
		l.ll.MoveToFront(ele)
		kv := ele.Value.(*indexes)
		l.mutex.Unlock()
		return kv.value, true
	}
	l.mutex.Unlock()
	return nil, false
}
```

### 使用示例

```go
package main

import (
	"fmt"
	lru2 "github.com/hlccd/goSTL/data_structure/lru"
)

type Value struct {
	Bytes []byte
}

func (v Value) Len() int {
	return len(v.Bytes)
}
func main() {
	space :=0
	lru := lru2.New(2<<10, nil)
	for i := 0; i < 10000; i++ {
		v := Value{Bytes: []byte(string(i))}
		lru.Insert(string(i), v)
		space +=v.Len()
	}
	fmt.Println("应该占用空间:",space)
	fmt.Println("LRU中存放的byte数量:",lru.Size())
	fmt.Println("LRU的byte数量上限:",lru.Cap())
	lru.Erase(string(9999))
	fmt.Println("删除后的LRU中存放的byte数量:",lru.Size())
	fmt.Println("从LRU中找9998")
	if v,ok:=lru.Get(string(9998));ok{
		fmt.Println(string(v.(Value).Bytes))
		fmt.Println(string(9998))
	}
	fmt.Println("从LRU中找1")
	if v,ok:=lru.Get(string(1));ok{
		fmt.Println(string(v.(Value).Bytes))
		fmt.Println(string(1))
	}
}
```

> 应该占用空间: 27824
> LRU中存放的byte数量: 2046
> LRU的byte数量上限: 2048
> 删除后的LRU中存放的byte数量: 2040
> 从LRU中找9998
> ✎
> ✎
> 从LRU中找1
> 
> 
