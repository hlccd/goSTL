github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		栈（stack）是一个**线性容器**，但不同于其他容器的特点在于，栈是仅仅支持从顶部插入和从顶部删除的操作，即**单向删除和添加**。

​		对于stack的实现，考虑到它是一个线性容器，并且其中的元素仅可通过顶部添加和删除即单项增删，所以可以考虑其**底层使用动态数组**的形式实现，考虑到动态数组需要自行对分配空间进行操作，同时它也类似于vector进行单项操作的情况，所以它的扩容缩容策略同vector一样，但不同点在于它只能使用单侧增删，所以不需要考虑从中间插入的一些情况。

### 原理

​		对于一个栈来说，可以将其看作是一个**木桶**，你要在这个容器里放入一些元素，就只能从已有元素的顶部放入，而不是插入的原有的元素之间去，当你要删除元素的时候，也只能从顶部一个一个删除。

​		它是可以动态的对元素进行添加和修改的，而如果说每一次增删元素都使得数组恰好可以容纳所有元素的话，那它在以后的每一次增删的过程中都需要去**重新分配空间**，并将原有的元素复制到新数组中，如果按照这样来操作的话，每次删减都会需要花费大量的时间去进行，将会极大的降低其性能。

​		为了提高效率，可以选择牺牲一定的空间作为冗余量，即**时空权衡**，通过每次扩容时多分配一定的空间，这样在后续添加时就可以减少扩容次数，而在删除元素时，如果冗余量不多时，可以选择不进行缩容，仅将顶部指针下移一位即可，，这样可以极大减少的由于容量不足或超出导致的扩容缩容造成的时间开销。

#### 扩容策略

​		对于动态数组扩容来说，由于其添加是线性且连续的，虽然stack一次只会添加一个，但为了减少频繁分配空间造成的时间浪费，选择多保留一些冗余空间以方便后续增删操作。

​		stack的扩容主要由两种方案：

**同vector的扩容策略**

1. 固定扩容：固定的去增加一定量的空间，该方案时候在数组较大时添加空间，可以避免直接翻倍导致的**冗余过量**问题，到由于增加量是固定的，如果需要一次扩容很多量的话就会比较缓慢。
2. 翻倍扩容：将原有容量直接翻倍，该方案适合在数组不太大的适合添加空间，可以提高扩容量，增加扩容效率，但数组太大时使用的话会导致一次性增加太多空间，进而造成空间的浪费。

#### 缩容策略

​		对于动态数组的缩容来说，同扩容类似，由于元素的减少只会是线性且连续的，即每一次只会减少一个，不会存在由于一个元素被删去导致其他已使用的空间也需要被释放。所以stack的缩容方案和stack的扩容方案十分类似，即是它的逆过程：

**同vector的缩容策略**

1. 固定缩容：释放一个固定的空间，该方案适合在当前数组较大的时候进行，可以减缓需要缩小的量，当空间较大的时候如果采取折半缩减的话;将可能在绝大多数时间内都不会被采用，从而造成空间浪费。
2. 折半缩容：将当前容量进行折半，该方案适合数组不太大的适合进行缩容，可以更快的缩小容量，但对于较大的空间来说并不会频繁使用到。

### 实现

​		stack栈结构体，包含动态数组和该数组的顶部指针，顶部指针指向实际顶部元素的下一位置，当删除节点时仅仅需要下移顶部指针一位即可，当新增结点时优先利用冗余空间，当冗余空间不足时先倍增空间至2^16，超过后每次增加2^16的空间，删除结点后如果冗余超过2^16,则释放掉，删除后若冗余量超过使用量，也释放掉冗余空间。

```go
type stack struct {
	data  []interface{} //用于存储元素的动态数组
	top   uint64        //顶部指针
	cap   uint64        //动态数组的实际空间
	mutex sync.Mutex    //并发控制锁
}
```

#### 接口

```go
type stacker interface {
	Iterator() (i *Iterator.Iterator) //返回一个包含栈中所有元素的迭代器
	Size() (num uint64)               //返回该栈中元素的使用空间大小
	Clear()                           //清空该栈容器
	Empty() (b bool)                  //判断该栈容器是否为空
	Push(e interface{})               //将元素e添加到栈顶
	Pop()                             //弹出栈顶元素
	Top() (e interface{})             //返回栈顶元素
}
```

#### New

​		新建一个stack栈容器并返回，初始stack的动态数组容量为1，初始stack的顶部指针置0，容量置1。

```go
func New() (s *stack) {
	return &stack{
		data:  make([]interface{}, 1, 1),
		top:   0,
		cap:   1,
		mutex: sync.Mutex{},
	}
}
```

#### Iterator

​		以stack栈容器做接收者，将stack栈容器中不使用空间释放掉，返回一个包含容器中所有使用元素的迭代器。

```go
func (s *stack) Iterator() (i *Iterator.Iterator) {
	if s == nil {
		s = New()
	}
	s.mutex.Lock()
	if s.data == nil {
		//data不存在,新建一个
		s.data = make([]interface{}, 1, 1)
		s.top = 0
		s.cap = 1
	} else if s.top < s.cap {
		//释放未使用的空间
		tmp := make([]interface{}, s.top, s.top)
		copy(tmp, s.data)
		s.data = tmp
	}
	//创建迭代器
	i = Iterator.New(&s.data)
	s.mutex.Unlock()
	return i
}
```

#### Size

​		以stack栈容器做接收者吗，返回该容器当前含有元素的数量。

```go
func (s *stack) Size() (num uint64) {
	if s == nil {
		s = New()
	}
	return s.top
}
```

#### Clear

​		以stack栈容器做接收者，将该容器中所承载的元素清空，将该容器的尾指针置0。

```go
func (s *stack) Clear() {
	if s == nil {
		s = New()
	}
	s.mutex.Lock()
	s.data = make([]interface{}, 0, 0)
	s.top = 0
	s.cap = 1
	s.mutex.Unlock()
}
```

#### Empty

​		以stack栈容器做接收者，判断该stack栈容器是否含有元素，如果含有元素则不为空,返回false， 如果不含有元素则说明为空,返回true，如果容器不存在,返回true，该判断过程通过顶部指针数值进行判断，当顶部指针数值为0时说明不含有元素，当顶部指针数值大于0时说明含有元素。

```go
func (s *stack) Empty() (b bool) {
	if s == nil {
		return true
	}
	return s.Size() == 0
}
```

#### Push

​		以stack栈容器做接收者，在容器顶部插入元素，若存储冗余空间，则在顶部指针位插入元素，随后上移顶部指针，否则进行扩容，扩容后获得冗余空间重复上一步即可。

​		固定扩容值设为2^16，翻倍扩容上限也为2^16。

```go
func (s *stack) Push(e interface{}) {
	if s == nil {
		s = New()
	}
	s.mutex.Lock()
	if s.top < s.cap {
		//还有冗余,直接添加
		s.data[s.top] = e
	} else {
		//冗余不足,需要扩容
		if s.cap <= 65536 {
			//容量翻倍
			if s.cap == 0 {
				s.cap = 1
			}
			s.cap *= 2
		} else {
			//容量增加2^16
			s.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
		s.data[s.top] = e
	}
	s.top++
	s.mutex.Unlock()
}
```

#### Pop

​		以stack栈容器做接收者，弹出容器顶部元素,同时顶部指针下移一位，当顶部指针小于容器切片实际使用空间的一半时,重新分配空间释放未使用部分，若容器为空,则不进行弹出。

```go
func (s *stack) Pop() {
	if s == nil {
		s = New()
		return
	}
	if s.Empty() {
		return
	}
	s.mutex.Lock()
	s.top--
	if s.cap-s.top >= 65536 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		s.cap -= 65536
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	} else if s.top*2 < s.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		s.cap /= 2
		tmp := make([]interface{}, s.cap, s.cap)
		copy(tmp, s.data)
		s.data = tmp
	}
	s.mutex.Unlock()
}
```

#### Top

​		以stack栈容器做接收者，返回该容器的顶部元素，若该容器当前为空,则返回nil。

```go
func (s *stack) Top() (e interface{}) {
	if s == nil {
		return nil
	}
	if s.Empty() {
		return nil
	}
	s.mutex.Lock()
	e = s.data[s.top-1]
	s.mutex.Unlock()
	return e
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/stack"
	"sync"
)

func main() {
	s := stack.New()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			s.Push(num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("使用迭代器遍历全部：")
	for i := s.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("使用size变删除边遍历：")
	size:=s.Size()
	for i := uint64(0); i < size; i++ {
		fmt.Println(s.Top())
		s.Pop()
	}
}
```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 使用迭代器遍历全部：
> 9
> 3
> 4
> 5
> 6
> 7
> 8
> 0
> 1
> 2
> 使用size变删除边遍历：
> 2
> 1
> 0
> 8
> 7
> 6
> 5
> 4
> 3
> 9

