github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		环（ring），是一种**离散的环状线性结构**，即**首尾连接的线性结构**，它是又多个分布在不同物理空间的结点，通过指针链接建立逻辑连接而形成的线性结构。

​		但不同的地方在于环本身没有首尾结点之分，甚至说，它没有首尾结点，它可以看作是**将链表的首尾结点连接起来**，即任何一个结点都可以看作是首节点也可以看作是尾结点，可以通过对环中任意一个结点向前或向后遍历得到全部结点。

​		同链表一样，它的所有结点之间都是相互分离的，基本不存在在物理上临接的可能（不绝对），所以它要增加或删除结点的过程也是十分简单的，在添加时建立待添加结点和前后节点的连接即可，删除时建立前后节点的连接即可。适用于频繁删减的情况，或者循环遍历的情况。

### 原理

​		对于一个环状线性结构来说，主要是要建立其它的连接，相比较于用数组的形式实现来说，数组实现需要分配一段连续的空间用以存储数据，但当数组分配的空间不足时，则需要再重新分配空间并将原有的元素复制过去，考虑到复制数据的时间开销也是不容忽略的，所以需要找到另一个结构去尽量减少复制元素的时间开销。

​		同时，对于使用数组去实现环状结构来说，它需要通过逻辑上去建立首尾结点的连接，而不是将这两个点连起来，直观上并不是十分友好，特别是在扩容缩容的时候需要做的事情太多了，时间开销太大了，所以本次**实现并不采用数组形式实现**，但**数组形式仍然是可行的**。

#### 添加策略

​		对于环来说，由于它的结点之间的物理间隔并不是连续的，而是根据需要随机的分布在整个内存中的，同时，对于一个环来说，**环并不存在一个核心结点**，它是任意一个结点都可以当作核心结点去使用。

​		对于此种特性，ring在添加结点时只会出现两种情况：

1. 环不存在：创造一个自环的结点并持有即可，自环指自己前后节点均指向自己
2. 环存在：在当前持有的结点后面插入一个结点即可，同时建立新结点和前后节点之间的关系，此时持有结点不变

#### 删除策略

​		同添加策略类似，考虑到环本身无核心的特性，只需要判断销毁该结点后是否仍有结点即可，有则更换持有结点，否则销毁环即可。

​		ring删除主要会有以下三种情况：

1. 环不存在：结束
2. 环仅有一个结点：销毁环
3. 环有多个结点：销毁当前持有结点，并将持有结点切换为原持有节点的**下一个结点**，同时通过在新持有结点前插入原持有结点的前结点实现连接。

### 实现

​		ring环结构体，包含环的头尾节点指针，当增删结点时只需要移动到对应位置进行操作即可，当一个节点进行增删时需要同步修改其临接结点的前后指针，结构体中记录该环中当前所持有的结点的指针即可，同时记录该环中存在多少元素即size，使用并发控制锁以保证数据一致性。

```go
type ring struct {
	now   *node      //环当前持有的结点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}
```

​		环的node节点结构体，pre和next是该节点的前后两个节点的指针，用以保证环整体是相连的。

```go
type node struct {
	data interface{} //结点所承载的元素
	pre  *node       //前结点指针
	next *node       //后结点指针
}
```

#### 接口

```go
type ringer interface {
	Iterator() (i *Iterator.Iterator) //创建一个包含环中所有元素的迭代器并返回其指针
	Size() (size uint64)              //返回环所承载的元素个数
	Clear()                           //清空该环
	Empty() (b bool)                  //判断该环是否位空
	Insert(e interface{})             //向环当前位置后方插入元素e
	Erase()                           //删除当前结点并持有下一结点
	Value() (e interface{})           //返回当前持有结点的元素
	Set(e interface{})                //在当前结点设置其承载的元素为e
	Next()                            //持有下一节点
	Pre()                             //持有上一结点
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

​		新建一个ring环容器并返回，初始持有的结点不存在,即为nil，初始size为0。

```go
func New() (r *ring) {
	return &ring{
		now:   nil,
		size:  0,
		mutex: sync.Mutex{},
	}
}
```

​		新建一个自环结点并返回其指针，初始首结点的前后结点指针都为自身。

```go
func newNode(e interface{}) (n *node) {
	n = &node{
		data: e,
		pre:  nil,
		next: nil,
	}
	n.pre = n
	n.next = n
	return n
}
```

#### Iterator

​		以ring环容器做接收者，将ring环容器中所承载的元素放入迭代器中，从该结点开始向后遍历获取全部承载的元素。

```go
func (r *ring) Iterator() (i *Iterator.Iterator) {
	if r == nil {
		r = New()
	}
	r.mutex.Lock()
	//将所有元素复制出来放入迭代器中
	tmp := make([]interface{}, r.size, r.size)
	//从当前结点开始向后遍历
	for n, idx := r.now, uint64(0); n != nil && idx < r.size; n, idx = n.nextNode(), idx+1 {
		tmp[idx] = n.value()
	}
	i = Iterator.New(&tmp)
	r.mutex.Unlock()
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

#### Size

​		以ring环容器做接收者，返回该容器当前含有元素的数量。

```go
func (r *ring) Size() (size uint64) {
	if r == nil {
		r = New()
	}
	return r.size
}
```

#### Clear

​		以ring环容器做接收者，将该容器中所承载的元素清空，将该容器的当前持有的结点置为nil,长度初始为0。

```go
func (r *ring) Clear() {
	if r == nil {
		r = New()
	}
	r.mutex.Lock()
	//销毁环
	r.now = nil
	r.size = 0
	r.mutex.Unlock()
}
```

#### Empty

​		以ring环容器做接收者，判断该ring环容器是否含有元素，该判断过程通过size进行判断,size为0则为true,否则为false。

```go
func (r *ring) Empty() (b bool) {
	if r == nil {
		r = New()
	}
	return r.size == 0
}
```

#### Insert

​		以ring环容器做接收者，通过环中当前持有的结点进行添加，如果环为建立,则新建一个自环结点设为环，存在持有的结点,则在其后方添加即可。

```go
func (r *ring) Insert(e interface{}) {
	if r == nil {
		r = New()
	}
	r.mutex.Lock()
	//新建自环结点
	n := newNode(e)
	if r.size == 0 {
		//原本无环,设为新环
		r.now = n
	} else {
		//持有结点,在后方插入
		r.now.insertNext(n)
	}
	r.size++
	r.mutex.Unlock()
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

​		以ring环容器做接收者，先判断是否仅持有一个结点，若仅有一个结点,则直接销毁环，否则将当前持有结点设为下一节点,并前插原持有结点的前结点即可。

```go
func (r *ring) Erase() {
	if r == nil {
		r = New()
	}
	if r.size == 0 {
		return
	}
	r.mutex.Lock()
	//删除开始
	if r.size == 1 {
		//环内仅有一个结点,销毁环即可
		r.now = nil
	} else {
		//环内还有其他结点,将持有结点后移一位
		//后移后将当前结点前插原持有结点的前结点
		r.now = r.now.nextNode()
		r.now.insertPre(r.now.preNode().preNode())
	}
	r.size--
	r.mutex.Unlock()
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

#### Value

​		以ring环容器做接收者，获取环中当前持有节点所承载的元素，若环中持有的结点不存在,直接返回nil。

```go
func (r *ring) Value() (e interface{}) {
	if r == nil {
		r = New()
	}
	if r.now == nil {
		//无持有结点,直接返回nil
		return nil
	}
	return r.now.value()
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

​		以ring环容器做接收者，修改当前持有结点所承载的元素，若未持有结点,直接结束即可。

```go
func (r *ring) Set(e interface{}) {
	if r == nil {
		r = New()
	}
	if r.now == nil {
		return
	}
	r.mutex.Lock()
	r.now.setValue(e)
	r.mutex.Unlock()
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

#### Next

​		以ring环容器做接收者，将当前持有的结点后移一位，若当前无持有结点,则直接结束。

```go
func (r *ring) Next() {
	if r == nil {
		r = New()
	}
	if r.now == nil {
		return
	}
	r.mutex.Lock()
	r.now = r.now.nextNode()
	r.mutex.Unlock()
}
```

#### Pre

​		以ring环容器做接收者，将当前持有的结点前移一位，若当前无持有结点,则直接结束。

```go
func (r *ring) Pre() {
	if r == nil {
		r = New()
	}
	if r.size == 0 {
		return
	}
	r.mutex.Lock()
	r.now = r.now.preNode()
	r.mutex.Unlock()
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/ring"
	"sync"
)

func main() {
	r:=ring.New()
	wg:=sync.WaitGroup{}
	for i:=0;i<10;i++{
		wg.Add(1)
		go func(num int) {
			r.Insert(num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i:=uint64(0);i<r.Size();i++{
		fmt.Println(r.Value())
		r.Pre()
	}
	fmt.Println("ring当前持有结点的元素:",r.Value())
	r.Set(-1)
	fmt.Println("ring修改有持有的结点元素:",r.Value())
	fmt.Println("删除后")
	r.Erase()
	for i:=r.Iterator();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
}
```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 2
> 0
> 1
> 5
> 3
> 4
> 7
> 6
> 8
> 9
> ring当前持有结点的元素: 2
> ring修改有持有的结点元素: -1
> 删除后
> 9
> 8
> 6
> 7
> 4
> 3
> 5
> 1
> 0
