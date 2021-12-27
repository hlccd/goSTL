github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		链表（list），是一种**离散的线性结构**，它是又多个分布在不同物理空间的结点，通过指针链接建立逻辑连接而形成的线性结构。

​		由于它的一个个结点相互之间是分离开的，所以它增加和删除结点的过程就会变得十分简单，只需要找到对应节点并将其增加/删除即可，同时修改该结点前后结点的指针以保证整个链表不断开即可，对整个链表的大多数元素来说几乎没有影响，适用于频繁增删的情况，但它也会频繁的分配空间。

### 原理

​		对于一个线性结构来说，主要是要建立其它的连接，相比较于用数组的形式实现来说，数组实现需要分配一段连续的空间用以存储数据，但当数组分配的空间不足时，则需要再重新分配空间并将原有的元素复制过去，考虑到复制数据的时间开销也是不容忽略的，所以需要找到另一个结构去尽量减少复制元素的时间开销。

​		对此，可以选择使用链表，它是**将多个分离的点通过指针连接起来**，对于增删元素都只需要修改对应的结点和它前后连接的结点即可，对于链表中存储的其他结点不需要做修改，即可以极大的省去频繁复制元素的时间开销，但相应的，由于每个结点是离散的，每次增删结点都需要分配空间，增加了分配空间的时间开销。

​		链表和数组两种数据结构各有千秋，当在需要承载的元素过多的时候，可以选择使用链表以减少频繁复制的时间开销，当元素不多时可以选择之间分配一个较大的数组去存储以减少分配空间的时间开销。

#### 添加策略

​		对于链表来说，由于它的结点之间的物理间隔并不是连续的，而是根据需要随机的分布在整个内存中的，同时，对于一个链表来说，核心的主要是它的首尾结点，可以通过首尾结点进行元素的添加，添加时需要先遍历到对应的位置再进行添加，可以根据位置和总承载元素的情况选择从前或从后进行遍历以减少时间开销。

​		list添加主要会有以下四种情况：

1. 链表不存在：创建新的结点设为头尾结点即可
2. 添加在首部：从头结点添加并修改链表的头结点即可
3. 添加在尾部：从尾结点添加并修改链表的尾结点即可
4. 添加的中部：从最近的首尾结点开始遍历找到对应位置后插入即可，同时要**修改前后结点之间的指针以指向自身**，保持链表的连续性。

#### 删除策略

​		同添加策略一样，链表的删除特性和添加是类似的，只需要删除一个结点即可，同时修改前后连接结点的指针以保证链表的连续性，当删除的是头尾结点时在修改记录的头尾结点即可。当所有元素删除完毕后，需要**销毁链表**。

​		list删除主要会有以下五种情况：

1. 链表不存在：结束
2. 删除首结点：将首节点设为首节点的下一结点，同时销毁首结点
3. 删除尾结点：将尾节点设为尾节点的上一结点，同时销毁尾结点
4. 删除中间结点：从最近的首尾结点开始遍历找到对应位置后删除即可，同时**将前后结点连接起来**，以保证链表的连续性。
5. 删除后链表中不存在元素：销毁链表

### 实现

​		list链表结构体，包含链表的头尾节点指针，当增删结点时只需要找到对应位置进行操作即可，当一个节点进行增删时需要同步修改其临接结点的前后指针，结构体中记录整个链表的首尾指针,同时记录其当前已承载的元素，使用并发控制锁以保证数据一致性。

```go
type list struct {
	first *node      //链表首节点指针
	last  *node      //链表尾节点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}
```

​		链表的node节点结构体，pre和next是该节点的前后两个节点的指针，用以保证链表整体是相连的。

```go
type node struct {
   data interface{} //结点所承载的元素
   pre  *node       //前结点指针
   next *node       //后结点指针
}
```

#### 接口

```go
type lister interface {
	Iterator() (i *Iterator.Iterator)                              //创建一个包含链表中所有元素的迭代器并返回其指针
	Sort(Cmp ...comparator.Comparator)                             //将链表中所承载的所有元素进行排序
	Size() (size uint64)                                           //返回链表所承载的元素个数
	Clear()                                                        //清空该链表
	Empty() (b bool)                                               //判断该链表是否位空
	Insert(idx uint64, e interface{})                              //向链表的idx位(下标从0开始)插入元素组e
	Erase(idx uint64)                                              //删除第idx位的元素(下标从0开始)
	Get(idx uint64) (e interface{})                                //获得下标为idx的元素
	Set(idx uint64, e interface{})                                 //在下标为idx的位置上放置元素e
	IndexOf(e interface{}, Equ ...comparator.Equaler) (idx uint64) //返回和元素e相同的第一个下标
	SubList(begin, num uint64) (newList *list)                     //从begin开始复制最多num个元素以形成新的链表
}
```

```go
type noder interface {
	preNode() (m *node)     //返回前结点指针
	nextNode() (m *node)    //返回后结点指针
	insertPre(pre *node)    //在该结点前插入结点并建立连接
	insertNext(next *node)  //在该结点后插入结点并建立连接
	erase()                 //删除该结点,并使该结点前后两结点建立连接
	value() (e interface{}) //返回该结点所承载的元素
	setValue(e interface{}) //修改该结点承载元素为e
}
```

#### New

​		新建一个list链表容器并返回，初始链表首尾节点为nil，初始size为0。

```go
func New() (l *list) {
	return &list{
		first: nil,
		last:  nil,
		size:  0,
		mutex: sync.Mutex{},
	}
}
```

​		新建一个结点并返回其指针， 初始首结点的前后结点指针都为nil。

```go
func newNode(e interface{}) (n *node) {
   return &node{
      data: e,
      pre:  nil,
      next: nil,
   }
}
```

#### Iterator

​		以list链表容器做接收者，将list链表容器中所承载的元素放入迭代器中。

```go
func (l *list) Iterator() (i *Iterator.Iterator) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//将所有元素复制出来放入迭代器中
	tmp := make([]interface{}, l.size, l.size)
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	i = Iterator.New(&tmp)
	l.mutex.Unlock()
	return i
}
```

​		以node结点做接收者，返回该结点的前结点。

```go
func (n *node) preNode() (pre *node) {
   if n == nil {
      return
   }
   return n.pre
}
```

​		以node结点做接收者，返回该结点的后结点。

```go
func (n *node) nextNode() (next *node) {
   if n == nil {
      return
   }
   return n.next
}
```

#### Sort

​		以list链表容器做接收者，将list链表容器中所承载的元素利用比较器进行排序，可以自行传入比较函数,否则将调用默认比较函数。

```go
func (l *list) Sort(Cmp ...comparator.Comparator) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//将所有元素复制出来用于排序
	tmp := make([]interface{}, l.size, l.size)
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	if len(Cmp) > 0 {
		comparator.Sort(&tmp, Cmp[0])
	} else {
		comparator.Sort(&tmp)
	}
	//将排序结果再放入链表中
	for n, idx := l.first, uint64(0); n != nil && idx < l.size; n, idx = n.nextNode(), idx+1 {
		n.setValue(tmp[idx])
	}
	l.mutex.Unlock()
}
```

#### Size

​		以list链表容器做接收者，返回该容器当前含有元素的数量。

```go
func (l *list) Size() (size uint64) {
	if l == nil {
		l = New()
	}
	return l.size
}
```

#### Clear

​		以list链表容器做接收者，将该容器中所承载的元素清空，将该容器的首尾指针均置nil,将size重置为0。

```go
func (l *list) Clear() {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	//销毁链表
	l.first = nil
	l.last = nil
	l.size = 0
	l.mutex.Unlock()
}
```

#### Empty

​		以list链表容器做接收者，判断该list链表容器是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true，该判断过程通过size进行判断,为0则为true,否则为false。

```go
func (l *list) Empty() (b bool) {
	if l == nil {
		l = New()
	}
	return l.size == 0
}
```

#### Insert

​		以list链表容器做接收者，通过链表的首尾结点进行元素插入，插入的元素可以有很多个，通过判断idx。

```go
func (l *list) Insert(idx uint64, e interface{}) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	n := newNode(e)
	if l.size == 0 {
		//链表中原本无元素,新建链表
		l.first = n
		l.last = n
	} else {
		//链表中存在元素
		if idx == 0 {
			//插入头节点
			n.insertNext(l.first)
			l.first = n
		} else if idx >= l.size {
			//插入尾节点
			l.last.insertNext(n)
			l.last = n
		} else {
			//插入中间节点
			//根据插入的位置选择从前或从后寻找
			if idx < l.size/2 {
				//从首节点开始遍历寻找
				m := l.first
				for i := uint64(0); i < idx-1; i++ {
					m = m.nextNode()
				}
				m.insertNext(n)
			} else {
				//从尾节点开始遍历寻找
				m := l.last
				for i := l.size - 1; i > idx; i-- {
					m = m.preNode()
				}
				m.insertPre(n)
			}
		}
	}
	l.size++
	l.mutex.Unlock()
}
```

​		以node结点做接收者，对该结点插入前结点，并建立前结点和该结点之间的连接。

```go
func (n *node) insertPre(pre *node) {
   if n == nil || pre == nil {
      return
   }
   pre.next = n
   pre.pre = n.pre
   if n.pre != nil {
      n.pre.next = pre
   }
   n.pre = pre
}
```

​		以node结点做接收者，对该结点插入后结点，并建立后结点和该结点之间的连接。

```go
func (n *node) insertNext(next *node) {
   if n == nil || next == nil {
      return
   }
   next.pre = n
   next.next = n.next
   if n.next != nil {
      n.next.pre = next
   }
   n.next = next
}
```

#### Erase

​		以list链表容器做接收者，先判断是否为首尾结点,如果是首尾结点,在删除后将设置新的首尾结点，当链表所承载的元素全部删除后则销毁链表，删除时通过idx与总元素数量选择从前或从后进行遍历以找到对应位置，删除后,将该位置的前后结点连接起来,以保证链表不断裂。

```go
func (l *list) Erase(idx uint64) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	if l.size > 0 && idx < l.size {
		//链表中存在元素,且要删除的点在范围内
		if idx == 0 {
			//删除头节点
			l.first = l.first.next
		} else if idx == l.size-1 {
			//删除尾节点
			l.last = l.last.pre
		} else {
			//删除中间节点
			//根据删除的位置选择从前或从后寻找
			if idx < l.size/2 {
				//从首节点开始遍历寻找
				m := l.first
				for i := uint64(0); i < idx; i++ {
					m = m.nextNode()
				}
				m.erase()
			} else {
				//从尾节点开始遍历寻找
				m := l.last
				for i := l.size - 1; i > idx; i-- {
					m = m.preNode()
				}
				m.erase()
			}
		}
		l.size--
		if l.size == 0 {
			//所有节点都被删除,销毁链表
			l.first = nil
			l.last = nil
		}
	}
	l.mutex.Unlock()
}
```

​		以node结点做接收者，销毁该结点，同时建立该节点前后节点之间的连接。

```go
func (n *node) erase() {
   if n == nil {
      return
   }
   if n.pre == nil && n.next == nil {
      return
   } else if n.pre == nil {
      n.next.pre = nil
   } else if n.next == nil {
      n.pre.next = nil
   } else {
      n.pre.next = n.next
      n.next.pre = n.pre
   }
   n = nil
}
```

#### Get

​		以list链表容器做接收者，获取第idx位结点所承载的元素,若不在链表范围内则返回nil

```go
func (l *list) Get(idx uint64) (e interface{}) {
	if l == nil {
		l = New()
	}
	if idx >= l.size {
		return nil
	}
	l.mutex.Lock()
	if idx < l.size/2 {
		//从首节点开始遍历寻找
		m := l.first
		for i := uint64(0); i < idx; i++ {
			m = m.nextNode()
		}
		e = m.value()
	} else {
		//从尾节点开始遍历寻找
		m := l.last
		for i := l.size - 1; i > idx; i-- {
			m = m.preNode()
		}
		e = m.value()
	}
	l.mutex.Unlock()
	return e
}
```

​		以node结点做接收者，返回该结点所要承载的元素。

```go
func (n *node) value() (e interface{}) {
   if n == nil {
      return nil
   }
   return n.data
}
```

#### Set

​		以list链表容器做接收者，修改第idx为结点所承载的元素,超出范围则不修改。

```go
func (l *list) Set(idx uint64, e interface{}) {
	if l == nil {
		l = New()
	}
	if idx >= l.size {
		return
	}
	l.mutex.Lock()
	if idx < l.size/2 {
		//从首节点开始遍历寻找
		m := l.first
		for i := uint64(0); i < idx; i++ {
			m = m.nextNode()
		}
		m.setValue(e)
	} else {
		//从尾节点开始遍历寻找
		m := l.last
		for i := l.size - 1; i > idx; i-- {
			m = m.preNode()
		}
		m.setValue(e)
	}
	l.mutex.Unlock()
}
```

​		以node结点做接收者，对该结点设置其承载的元素。

```go
func (n *node) setValue(e interface{}) {
   if n == nil {
      return
   }
   n.data = e
}
```

#### IndexOf

​		以list链表容器做接收者，返回与e相同的元素的首个位置，可以自行传入用于判断相等的相等器进行处理，遍历从头至尾,如果不存在则返回l.size。

```go
func (l *list) IndexOf(e interface{}, Equ ...comparator.Equaler) (idx uint64) {
	if l == nil {
		l = New()
	}
	l.mutex.Lock()
	var equ comparator.Equaler
	if len(Equ) > 0 {
		equ = Equ[0]
	} else {
		equ = comparator.GetEqual()
	}
	n := l.first
	//从头寻找直到找到相等的两个元素即可返回
	for idx = 0; idx < l.size && n != nil; idx++ {
		if equ(n.value(), e) {
			break
		}
		n = n.nextNode()
	}
	l.mutex.Unlock()
	return idx
}
```

#### SubList

​		以list链表容器做接收者，以begin为起点(包含),最多复制num个元素进入新链表，并返回**新链表指针**。

```go
func (l *list) SubList(begin, num uint64) (newList *list) {
	if l == nil {
		l = New()
	}
	newList = New()
	l.mutex.Lock()
	if begin < l.size {
		//起点在范围内,可以复制
		n := l.first
		for i := uint64(0); i < begin; i++ {
			n = n.nextNode()
		}
		m := newNode(n.value())
		newList.first = m
		newList.size++
		for i := uint64(0); i < num-1 && i+begin < l.size-1; i++ {
			n = n.nextNode()
			m.insertNext(newNode(n.value()))
			m = m.nextNode()
			newList.size++
		}
		newList.last = m
	}
	l.mutex.Unlock()
	return newList
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/list"
	"sync"
)

func main() {
	l := list.New()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx uint64) {
			l.Insert(idx, idx)
			wg.Done()
		}(uint64(i))
	}
	wg.Wait()
	fmt.Println("使用迭代器访问全部元素；")
	for i := l.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	l.Set(5, "测试")
	fmt.Println("输出刚设定的测试元素的位置：", l.IndexOf("测试"))
	fmt.Println("从测试位生产新list，长度上限为10，并从头部输出：")
	newList := l.SubList(l.IndexOf("测试"), 10)
	len := newList.Size()
	fmt.Println("新链表的长度：", len)
	for i := uint64(0); i < len; i++ {
		fmt.Println(newList.Get(0))
		newList.Erase(0)
	}
	fmt.Println("从结尾向首部输出原链表：")
	for i := l.Size(); i > 0; i-- {
		fmt.Println(l.Get(i - 1))
		l.Erase(i - 1)
	}
}

```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 使用迭代器访问全部元素；
> 0
> 1
> 2
> 3
> 4
> 5
> 6
> 7
> 9
> 8
> 输出刚设定的测试元素的位置： 5
> 从测试位生产新list，长度上限为10，并从头部输出：
> 新链表的长度： 5
> 测试
> 6
> 7
> 9
> 8
> 从结尾向首部输出原链表：
> 8
> 9
> 7
> 6
> 测试
> 4
> 3
> 2
> 1
> 0
