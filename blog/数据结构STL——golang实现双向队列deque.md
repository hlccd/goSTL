github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		双向队列（deque）是一个封装了动态大小数组的顺序容器。它同其他数据结构一样都可以承载任意类型的元素，但相比于队列来说，它不同的点主要在于，队列只能首部出尾部入，而双向队列可以在首部和尾部都进行出入操作。

​		对于deque的实现来说，虽然它也是一个线性容器，但考虑到它将会在首部和尾部同时实现增删操作，如果仍然使用动态数组的形式去实现，虽然也可以做出来，但在空间扩容缩容的时候，将会进行大量的数据复制的过程，而这一过程其实也是可以省去的。方法就是**使用链表的形式对它进行存储**，这样对新元素的增删就可以不需要考虑中间仍然存在的元素的复制过程，从而免除复制导致的时间浪费。		

### 原理

​		对于deque来说，考虑到完全使用链表的形式：即一个结点承载一个元素，除了需要较多的空间分配操作的时间开销外，还需要更多的指针去记录该结点的前后结点的位置，即：既要更多的空间分配导致的时间开销，也有记录前后节点指针的空间开销。对于纯粹的链表实现来说，也并不是特别良好的策略。

​		于是，可以考虑将**链表和数组两者结合起来**，不采用动态数组的形式，而是将一个固定容量的数组放入一个链表结点中，这样即利用了链表可以减少复制数据从而减少时间开销的优势，也可以减少记录太多指针导致的空间浪费，同时由于一次分配固定量的空间大小，也可以留作一定的冗余量，更进一步的减少分配空间导致的时间开销。

​		因此，它的扩容和缩容策略相对于队列容器来说有着较大的不同：

#### 添加策略

​		考虑到上文所说，deque的扩容方案其实就只有一种了，即新增链表的首节点，将之前结点不能承载的元素放入新结点中，同时将新结点和链表建立连接。由新的结点称为首节点/尾结点。

​		对于deque添加元素的过程中，只会出现以下三种情况：

1. deque中无结点：新增一个结点同时作为首尾结点即可，并将元素放入其中，具体放入首部还是尾部根据添加方向来确定；
2. 首/尾结点无空间容纳：新增一个结点，将新结点设为首/尾结点，并将元素放弃新结点中
3. 首/尾结点仍有空间存放：直接放入即可

#### 删除策略

​		同上文一样，由于采用链表结合固定数组的形式进行存储，所以它的添加和删除策略是十分类似的，即在删除的时候只需要删除一个结点即可，同时需要将**被删除的结点和链表之间的联系全部断开**，这样才能让被删除的结点被回收。

​		对于deque删除结点的过程中，只会出现以下四种情况：

1. 链表不存在：无需删除操作
2. 对应结点仍有承载元素：将结点begin/end进行移动即可
3. 对应节点无承载元素了：释放该结点，同时解除该结点与链表的连接
4. 所有元素都被删除完成：销毁链表

### 实现

**注：整体实现过程中将deque和链表的node进行了分离**

​			deque双向队列结构体，包含链表的头尾节点指针，当删除节点时通过头尾节点指针进入链表进行删除，当一个节点全部被删除后则释放该节点,同时首尾节点做相应调整，当添加节点时若未占满节点空间时移动下标并做覆盖即可，当添加节点时空间已使用完毕时,根据添加位置新建一个新节点补充上去。

```go
type deque struct {
	first *node      //链表首节点指针
	last  *node      //链表尾节点指针
	size  uint64     //当前存储的元素个数
	mutex sync.Mutex //并发控制锁
}
```

​		deque双向队列中链表的node节点结构体，包含一个2^10空间的固定数组用以承载元素，使用begin和end两个下标用以表示新增的元素的下标,由于begin可能出现-1所以不选用uint16，pre和next是该节点的前后两个节点，用以保证链表整体是相连的。

```go
type node struct {
	data  [1024]interface{} //用于承载元素的股东数组
	begin int16             //该结点在前方添加结点的下标
	end   int16             //该结点在后方添加结点的下标
	pre   *node             //该结点的前一个结点
	next  *node             //该节点的后一个结点
}
```

#### 接口

```go
type dequer interface {
	Iterator() (i *Iterator.Iterator) //返回包含双向队列中所有元素的迭代器
	Size() (size uint64)              //返回该双向队列中元素的使用空间大小
	Clear()                           //清空该双向队列
	Empty() (b bool)                  //判断该双向队列是否为空
	PushFront(e interface{})          //将元素e添加到该双向队列的首部
	PushBack(e interface{})           //将元素e添加到该双向队列的尾部
	PopFront()                        //将该双向队列首元素弹出
	PopBack()                         //将该双向队列首元素弹出
	Front() (e interface{})           //获取该双向队列首部元素
	Back() (e interface{})            //获取该双向队列尾部元素
}
```

```go
type noder interface {
   nextNode() (m *node)                   //返回下一个结点
   preNode() (m *node)                    //返回上一个结点
   value() (es []interface{})             //返回该结点所承载的所有元素
   pushFront(e interface{}) (first *node) //在该结点头部添加一个元素,并返回新首结点
   pushBack(e interface{}) (last *node)   //在该结点尾部添加一个元素,并返回新尾结点
   popFront() (first *node)               //弹出首元素并返回首结点
   popBack() (last *node)                 //弹出尾元素并返回尾结点
   front() (e interface{})                //返回首元素
   back() (e interface{})                 //返回尾元素
}
```

#### New

​		创建一个deque容器并初始化，同时返回其指针。

```go
func New() (q *queue) {
	return &queue{
		data:  make([]interface{}, 1, 1),
		begin: 0,
		end:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}
```

​		创建一个首结点并初始化，同时返回其指针。

```go
func createFirst() (n *node) {
	return &node{
		data:  [1024]interface{}{},
		begin: 1023,
		end:   1024,
		pre:   nil,
		next:  nil,
	}
}
```

​		创建一个尾结点并初始化，同时返回其指针。

```go
func createLast() (n *node) {
   return &node{
      data:  [1024]interface{}{},
      begin: -1,
      end:   0,
      pre:   nil,
      next:  nil,
   }
}
```

#### Iterator

​		将双向队列中承载元素传入迭代器中，通过遍历链表中所有结点进行获取，不清除双向队列中的冗余量。返回迭代器指针，用于遍历双向队列。

```go
func (d *deque) Iterator() (i *Iterator.Iterator) {
	if d == nil {
		d = New()
	}
	tmp := make([]interface{}, 0, d.size)
	//遍历链表的所有节点,将其中承载的元素全部复制出来
	for m := d.first; m != nil; m = m.nextNode() {
		tmp = append(tmp, m.value()...)
	}
	return Iterator.New(&tmp)
}
```

​		返回当前结点的下一个结点。

```go
func (n *node) nextNode() (m *node) {
   if n == nil {
      return nil
   }
   return n.next
}
```

​		返回当前结点的上一个结点。

```go
func (n *node) preNode() (m *node) {
   if n == nil {
      return nil
   }
   return n.pre
}
```

​		返回当前结点所承载的所有元素。

```go
func (n *node) value() (es []interface{}) {
   es = make([]interface{}, 0, 0)
   if n == nil {
      return es
   }
   if n.begin > n.end {
      return es
   }
   es = n.data[n.begin+1 : n.end]
   return es
}
```

#### Size

​		返回deque中当前所包含的元素个数，由于其数值必然是非负整数，所以选用了**uint64**。

```go
func (d *deque) Size() (size uint64) {
	if d == nil {
		d = New()
	}
	return d.size
}
```

#### Clear

​		清空了deque中所承载的所有元素，同时销毁链表的首尾结点以实现销毁链表。

```go
func (d *deque) Clear() {
	if d == nil {
		d = New()
		return
	}
	d.mutex.Lock()
	d.first = nil
	d.last = nil
	d.size = 0
	d.mutex.Unlock()
}
```

#### Empty

​		判断deque中是否为空，通过Size()是否为0进行判断。

```go
func (d *deque) Empty() (b bool) {
	if d == nil {
		d = New()
	}
	return d.Size() == 0
}
```

#### PushFront

​		以deque为接收者，向其首部添加一个元素。

​		通过链表的首结点进行实现，当链表不存在时则创建一个首结点并以此作为链表首尾结点。

```go
func (d *deque) PushFront(e interface{}) {
	if d == nil {
		d = New()
	}
	d.mutex.Lock()
	d.size++
	//通过首节点进行添加
	if d.first == nil {
		d.first = createFirst()
		d.last = d.first
	}
	d.first = d.first.pushFront(e)
	d.mutex.Unlock()
}
```

​		以node结点做接收者，向该节点前方添加元素e，当该结点空间已经使用完毕后,新建一个结点并将新结点设为首结点，将插入元素放入新结点并返回新结点作为新的首结点，否则插入当前结点并返回当前结点,首结点不变。

```go
func (n *node) pushFront(e interface{}) (first *node) {
   if n == nil {
      return n
   }
   if n.begin >= 0 {
      //该结点仍有空间可用于承载元素
      n.data[n.begin] = e
      n.begin--
      return n
   }
   //该结点无空间承载,创建新的首结点用于存放
   m := createFirst()
   m.data[m.begin] = e
   m.next = n
   n.pre = m
   m.begin--
   return m
}
```

#### PushBack

​		以deque为接收者，向其尾部添加一个元素。

​		通过链表的尾结点进行实现，当链表不存在时则创建一个尾结点并以此作为链表首尾结点。

```go
func (d *deque) PushBack(e interface{}) {
	if d == nil {
		d = New()
	}
	d.mutex.Lock()
	d.size++
	//通过尾节点进行添加
	if d.last == nil {
		d.last = createLast()
		d.first = d.last
	}
	d.last = d.last.pushBack(e)
	d.mutex.Unlock()
}
```

​		以node结点做接收者，向该节点后方添加元素e，当该结点空间已经使用完毕后,新建一个结点并将新结点设为尾结点，将插入元素放入新结点并返回新结点作为新的尾结点，否则插入当前结点并返回当前结点,尾结点不变。

```go
func (n *node) pushBack(e interface{}) (last *node) {
	if n == nil {
		return n
	}
	if n.end < int16(len(n.data)) {
		//该结点仍有空间可用于承载元素
		n.data[n.end] = e
		n.end++
		return n
	}
	//该结点无空间承载,创建新的尾结点用于存放
	m := createLast()
	m.data[m.end] = e
	m.pre = n
	n.next = m
	m.end++
	return m
}
```

#### PopFront

​		以deque双向队列容器做接收者，利用首节点进行弹出元素,可能存在首节点全部释放要进行首节点后移的情况，当元素全部删除后,释放全部空间,将首尾节点都设为nil。

```go
func (d *deque) PopFront() {
	if d == nil {
		d = New()
	}
	if d.size == 0 {
		return
	}
	d.mutex.Lock()
	//利用首节点删除首元素
	//返回新的首节点
	d.first = d.first.popFront()
	d.size--
	if d.size == 0 {
		//全部删除完成,释放空间,并将首尾节点设为nil
		d.first = nil
		d.last = nil
	}
	d.mutex.Unlock()
}
```

​		 以node结点做接收者，利用首节点进行弹出元素,可能存在首节点全部释放要进行首节点后移的情况， 当发生首结点后移后将会返回新首结点,否则返回当前结点。

```go
func (n *node) popFront() (first *node) {
   if n == nil {
      return nil
   }
   if n.begin < int16(len(n.data))-2 {
      //该结点仍有承载元素
      n.begin++
      n.data[n.begin] = nil
      return n
   }
   if n.next != nil {
      //清除该结点下一节点的前结点指针
      n.next.pre = nil
   }
   return n.next
}
```

#### PopBack

​		以deque双向队列容器做接收者，利用尾节点进行弹出元素,可能存在尾节点全部释放要进行尾节点前移的情况，当元素全部删除后,释放全部空间,将首尾节点都设为nil。

```go
func (d *deque) PopBack() {
	if d == nil {
		d = New()
	}
	if d.size == 0 {
		return
	}
	d.mutex.Lock()
	//利用尾节点删除首元素
	//返回新的尾节点
	d.last = d.last.popBack()
	d.size--
	if d.size == 0 {
		//全部删除完成,释放空间,并将首尾节点设为nil
		d.first = nil
		d.last = nil
	}
	d.mutex.Unlock()
}
```

​		以node结点做接收者，利用尾节点进行弹出元素,可能存在尾节点全部释放要进行尾节点前移的情况，当发生尾结点前移后将会返回新尾结点,否则返回当前结点。

```go
func (n *node) popBack() (last *node) {
	if n == nil {
		return nil
	}
	if n.end > 1 {
		//该结点仍有承载元素
		n.end--
		n.data[n.end] = nil
		return n
	}
	if n.pre != nil {
		//清除该结点上一节点的后结点指针
		n.pre.next = nil
	}
	return n.pre
}
```

#### Front

​		以deque双向队列容器做接收者，返回该容器的第一个元素，利用首节点进行寻找，若该容器当前为空,则返回nil。

```go
func (d *deque) Front() (e interface{}) {
	if d == nil {
		d = New()
	}
	return d.first.front()
}
```

​		以node结点做接收者，返回该结点的第一个元素,利用首节点和begin进行查找，若该结点为nil,则返回nil。

```go
func (n *node) front() (e interface{}) {
   if n == nil {
      return nil
   }
   return n.data[n.begin+1]
}
```

#### Back

​		以deque双向队列容器做接收者，返回该容器的最后一个元素,利用尾节点进行寻找，若该容器当前为空,则返回nil。

```go
func (d *deque) Back() (e interface{}) {
	if d == nil {
		d = New()
	}
	return d.last.back()
}
```

​		以node结点做接收者，返回该结点的最后一个元素,利用尾节点和end进行查找，若该结点为nil,则返回nil。

```go
func (n *node) back() (e interface{}) {
	if n == nil {
		return nil
	}
	return n.data[n.end-1]
}
```

#### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/deque"
	"sync"
)

func main() {
	d:=deque.New()
	wg:=sync.WaitGroup{}
	for i := 0; i >=-4; i-- {
		wg.Add(1)
		go func(num int) {
			d.PushFront(num)
			wg.Done()
		}(i)
	}
	for i := 0; i <= 4; i++ {
		wg.Add(1)
		go func(num int) {
			d.PushBack(num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("使用迭代器遍历:")
	for i:=d.Iterator();i.HasNext();i.Next(){
		fmt.Println(i.Value())
	}
	len:=d.Size()
	for i := uint64(0); i < len; i += 2 {
		fmt.Println("当前首元素:",d.Front())
		fmt.Println("当前尾元素:",d.Back())
		d.PopFront()
		d.PopBack()
	}
}
```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 使用迭代器遍历:
> 0
> -1
> -4
> -3
> -2
> 4
> 0
> 1
> 2
> 3
> 当前首元素: 0
> 当前尾元素: 3
> 当前首元素: -1
> 当前尾元素: 2
> 当前首元素: -4
> 当前尾元素: 1
> 当前首元素: -3
> 当前尾元素: 0
> 当前首元素: -2
> 当前尾元素: 4
