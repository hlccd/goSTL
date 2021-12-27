github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		队列（queue）是一个封装了动态大小数组的顺序容器。除了可以包含任意类型的元素外，更主要是的它满足**FIFO**的先进先出模式，对于一些排队问题可以考虑使用队列来存储。

​		对于queue的实现,由于它也是一个线性容器，底层依然可以考虑使用动态数组来实现，但它和vector仍有一定的不同，vector的冗余量主要是在尾部，毕竟vector要实现随机读取的话中间和首部不能有空余量，而对于queue来说，它的添加只在尾部，而首部仅仅只做删除，所以除了在尾部留有一定的空间做添加之外，也可以在首部删除后留有不多的余量以**避免多次分配空间**。

### 原理

​		对于一个queue来说，它可以**在尾部添加元素**，**在首部删除元素**。而同vector一样，如果每次增删都要重新分配空间，将会极大的降低效率，所以可以考虑对前后都留有一定的冗余，以此来减少分配空间和复制的次数，从而减少时间开销。

​		对于首部的冗余来说，它的存在主要是为了减少删除后就重新分配空间并复制，同时对于尾部添加时，如果尾部冗余不足可以将使用的元素整体前移到首部，即将首部冗余挪到尾部去，从而减少了空间分配的次数。

​		对于尾部的冗余来说，它存在的目的和vector类似，仅仅是为了后续元素时更快而已。

​		和vector一样，都是通过设置一定量的冗余来换取操作次数的减少，从而提高效率减少时间损耗。

#### 扩容策略

​		对于queue的动态数组扩容来说，由于其添加是线性且连续的，即每一次只会增加一个，并且每一次都只会在它的尾部，不像vector一样可能在数组的任意位置进行添加，同时由于首部也有冗余，所以它的扩容也需要考虑首部的情况。

​		queue的扩容主要由三种方案：

1. 利用首部：**严格来说其实并没有进行扩容**，而是将承载的元素前移到首部，即将首部的冗余平移到最后进行利用。
2. 固定扩容：固定的去增加一定量的空间，该方案时候在数组较大时添加空间，可以避免直接翻倍导致的**冗余过量**问题，到由于增加量是固定的，如果需要一次扩容很多量的话就会比较缓慢。（同vector）
3. 翻倍扩容：将原有容量直接翻倍，该方案适合在数组不太大的适合添加空间，可以提高扩容量，增加扩容效率，但数组太大时使用的话会导致一次性增加太多空间，进而造成空间的浪费。（同vector）

#### 缩容策略

​		对于queue的数组来说，它的缩容不同于扩容，由于删除元素只会在首部进行，所以缩容其实也只会在首部进行，考虑到首部并不会进行添加，所以也不需要冗余太多的量进行减缓，即可以将设定上限减少一些，如(2^10)：

1. 固定缩容：释放一个固定的空间，该方案适合在当前数组较大的时候进行，可以减缓需要缩小的量，当首部冗余超过上限时进行缩容，**一次性全部释放即可**。
2. 折半缩容：当**首部冗余超过了实际承载的元素的数量**时，需要对首部进行缩容，同前者一样，不需要考虑首部添加的问题，即采用对首部冗余全部释放的方式。

### 实现

​		queue底层同样使用动态数组实现，同时，由于首部只做删除，尾部只做添加，而首尾两侧都有一定的冗余，所以需要对两侧都进行记录，也可以根据尾部下标减首部下标得出实际承载元素的量，与此同时，需要引入cap以记录实际分配的空间大小，也可以根据cap和end计算出尾部冗余量。同时，为了解决在高并发情况下的数据不一致问题，引入了并发控制锁。

```go
type queue struct {
	data  []interface{} //泛型切片
	begin uint64        //首节点下标
	end   uint64        //尾节点下标
	cap   uint64        //容量
	mutex sync.Mutex    //并发控制锁
}
```

#### 接口

```go
type queuer interface {
	Iterator() (i *Iterator.Iterator) //返回包含队列中所有元素的迭代器
	Size() (num uint64)               //返回该队列中元素的使用空间大小
	Clear()                           //清空该队列
	Empty() (b bool)                  //判断该队列是否为空
	Push(e interface{})               //将元素e添加到该队列末尾
	Pop() (e interface{})             //将该队列首元素弹出并返回
	Front() (e interface{})           //获取该队列首元素
	Back() (e interface{})            //获取该队列尾元素
}
```

#### New

​		创建一个queue容器并初始化，同时返回其指针。

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

#### Iterator

​		将队列中承载元素传入迭代器中，不清除队列中的冗余量。返回迭代器指针，用于遍历队列。

```go
func (q *queue) Iterator() (i *Iterator.Iterator) {
   if q == nil {
      q=New()
   }
   q.mutex.Lock()
   tmp:=make([]interface{},q.end-q.begin,q.end-q.begin)
   copy(tmp, q.data[q.begin:q.end])
   i = Iterator.New(&tmp)
   q.mutex.Unlock()
   return i
}
```

#### Size

​		返回queue中当前所包含的元素个数，由于其数值必然是非负整数，所以选用了**uint64**。

```go
func (q *queue) Size() (num uint64) {
	if q == nil {
		q = New()
	}
	return q.end - q.begin
}
```

#### Clear

​		清空了queue中所承载的所有元素。

```go
func (q *queue) Clear() {
	if q == nil {
		q = New()
	}
	q.mutex.Lock()
	q.data = make([]interface{}, 1, 1)
	q.begin = 0
	q.end = 0
	q.cap = 1
	q.mutex.Unlock()
}
```

#### Empty

​		判断queue中是否为空，通过Size()是否为0进行判断。

```go
func (q *queue) Empty() (b bool) {
	if q == nil {
		q = New()
	}
	return q.Size() <= 0
}
```

#### Push

​		以queue为接受器，向queue尾部添加一个元素，添加元素会出现两种情况，第三种是还有冗余量，此时**直接覆盖**以len为下标指向的位置即可，另一种情况是没有冗余量了，需要对动态数组进行**扩容**，此时就需要利用扩容策略。另一种是尾部没有冗余量，但首部仍有冗余量，此时可以将承载元素前移，把首部冗余量”借“过来使用。

​		扩容策略上文已做描述，可返回参考，该实现过程种将两者进行了结合使用，可参考下方注释。

​		固定扩容值设为2^16，翻倍扩容上限也为2^16。

```go
func (q *queue) Push(e interface{}) {
	if q == nil {
		q = New()
	}
	q.mutex.Lock()
	if q.end < q.cap {
		//不需要扩容
		q.data[q.end] = e
	} else {
		//需要扩容
		if q.begin > 0 {
			//首部有冗余,整体前移
			for i := uint64(0); i < q.end-q.begin; i++ {
				q.data[i] = q.data[i+q.begin]
			}
			q.end -= q.begin
			q.begin = 0
		} else {
			//冗余不足,需要扩容
			if q.cap <= 65536 {
				//容量翻倍
				if q.cap == 0 {
					q.cap = 1
				}
				q.cap *= 2
			} else {
				//容量增加2^16
				q.cap += 2 ^ 16
			}
			//复制扩容前的元素
			tmp := make([]interface{}, q.cap, q.cap)
			copy(tmp, q.data)
			q.data = tmp
		}
		q.data[q.end] = e
	}
	q.end++
	q.mutex.Unlock()
}
```

#### Pop

​		以queue队列容器做接收者，弹出容器首部元素,同时begin++即可，若容器为空,则不进行弹出，当弹出元素后,可能进行缩容，由于首部不会进行添加，所以不需要太多的冗余，即将首部冗余上限设为2^10，固定缩容和折半缩容都参考此值执行，具体缩容策略可参考上文介绍。

```go
func (q *queue) Pop() (e interface{}) {
	if q == nil {
		q = New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	q.mutex.Lock()
	e = q.data[q.begin]
	q.begin++
	if q.begin >= 1024 || q.begin*2>q.end {
		//首部冗余超过2^10或首部冗余超过实际使用
		q.cap -= q.begin
		q.end -= q.begin
		tmp := make([]interface{}, q.cap, q.cap)
		copy(tmp, q.data[q.begin:])
		q.data = tmp
		q.begin=0
	}
	q.mutex.Unlock()
	return e
}
```

#### Front

​		以queue为接受器，返回vector所承载的元素中位于**首部**的元素，如果queue为nil或者元素数组为nil或为空，则返回nil。考虑到仅仅只是读取元素，故不对该过程进行加锁操作。

```go
func (q *queue) Front() (e interface{}) {
	if q == nil {
		q=New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	return q.data[q.begin]
}
```

#### Back

​		以queue为接受器，返回queue所承载的元素中位于**尾部**的元素，如果queue为nil或者元素数组为nil或为空，则返回nil。考虑到仅仅只是读取元素，故不对该过程进行加锁操作。

```go
func (q *queue) Back() (e interface{}) {
	if q == nil {
		q=New()
		return nil
	}
	if q.Empty() {
		q.Clear()
		return nil
	}
	return q.data[q.end-1]
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/queue"
	"sync"
)

func main() {
	q := queue.New()
	wg := sync.WaitGroup{}
	//随机插入队列中
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(num int) {
			fmt.Println(num)
			q.Push(num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("输出首部:", q.Front())
	fmt.Println("输出尾部:", q.Back())
	fmt.Println("弹出并输出前4个:")
	for i := uint64(0); i < q.Size()-1; i++ {
		fmt.Println(q.Pop())
	}
	//在尾部再添加4个,从10开始以做区分
	for i := 10; i < 14; i++ {
		q.Push(i)
	}
	fmt.Println("从头输出全部:")
	for ;!q.Empty();{
		fmt.Println(q.Pop())
	}
}

```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 0
> 7
> 1
> 2
> 3
> 4
> 5
> 6
> 输出首部: 0
> 输出尾部: 6
> 弹出并输出前4个:
> 0
> 7
> 1
> 2
> 从头输出全部:
> 3
> 4
> 5
> 6
> 10
> 11
> 12
> 13
