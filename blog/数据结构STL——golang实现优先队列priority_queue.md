github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		优先队列（priority_queue）它虽然名字上是被称之为队列，单它底层其实是以堆的方式实现的，而堆这个数据结构，它是通过建立一棵**完全二叉树**来进行实现的。它在逻辑上并非是一个线性结构，但由于**二叉树可以用数组表示**的特性，**本次实现采用数组的形式实现**，后续会再使用完全二叉搜实现一次。

​		堆或者说优先队列的主要特点为：它的**父结点必然是小于或等于左右子结点**（该文章叙述中使用**小顶堆**，大顶堆大小情况相反），同时，考虑到用数组来模拟完全二叉搜的情况，它的所有结点必然在数组上是连续的，不会存在上一个结点和下一个结点间存在一段间隔的情况。

### 原理

##### 数组模拟完全二叉树

​		对于使用数组来模拟完全二叉树的情况来说，首先，完全二叉树除开底部的一层可能存在不满的情况外，其他上层必然是满的，同时，对于完全二叉搜的任意一个结点来说，假设该结点的下标为**p（从0开始）**，使用数组来模拟的情况下，其左子结点的下标则必然是**2p+1**，而其右子结点的下标则必然是**2p+2**，当然，也可能出现不存在右子结点或不存在左子结点的情况。

​		对于数组来说，可以将数组根据2^n来进行划分，n从0开始，即0位为第一层，1~2为第二层，3~6为第三层，7~17为第四层······以此可有
$$
2^{(n-1)}-1——2^{n}-2
$$
为第n层，同时，由于任意结点的左右子结点必然在其下一层，所以通过计算也可得知其左右子结点的下标也符合上文所述。当然，如果给定一个结点下标为**p**（除首结点，毕竟首节点无父结点），则它父结点的下标必然为**p/2**

##### 优先队列实现

​		对于一个优先队列来说，它是建立来完全二叉搜的基础上，同时其任意结点的值必然不大于其左右子结点，所以可知道，当插入节点时，可以利用尾部结点和父结点进行比较交换来实现，对于删除节点时，可以利用顶部结点与左右子结点的值进行比较进行下降即可。

#### 添加策略

​		当插入一个结点的时候，可以先将其放入二叉树的末尾，然后根据插入的值和它父结点进行比较，如果插入值小于父节点时则交换两结点值，随后在用交换后的结点与其父结点比较，重复该过程直到抵达顶点或不满足条件即可。当然，添加的过程中，由于底层使用动态数组来存储，所以可能出现扩容的情况，由于其情况和vector类似，则不在赘述，仅在下方列出即可。

##### 添加步骤

1. 先判断动态数组中是否有冗余空间，有则直接放入尾部，否则利用扩容策略进行扩容，扩容后再放入尾部即可
2. 结点放入尾部后，通过比较该结点与父结点的值进行交换
3. 满足交换条件，重复回到2过程
4. 到达顶部或不满足交换条件，添加结束，插入完成

##### 扩容策略

1. 固定扩容：固定的去增加一定量的空间，该方案时候在数组较大时添加空间，可以避免直接翻倍导致的**冗余过量**问题，到由于增加量是固定的，如果需要一次扩容很多量的话就会比较缓慢。
2. 翻倍扩容：将原有容量直接翻倍，该方案适合在数组不太大的适合添加空间，可以提高扩容量，增加扩容效率，但数组太大时使用的话会导致一次性增加太多空间，进而造成空间的浪费。

#### 删除策略

​		不同于添加时直接利用尾结点的情况，由于删除的是首结点，而同时删除需要减少一个空间，所以可以考虑将首位结点进行交换，或者其实直接**用尾结点覆写首节点**即可，这样就实现了首节点的删除，同时减小一位删去移到首位的尾结点实现实际的删除。同时，在删除后，需要判断实际使用空间是否满足缩容策略的情况，如果满足则需要使用缩容策略进行缩容，缩容策略同vector，已在下方列出。

​		上一过程完成后，首节点以及被删除了，但移到首节点的新首节点可能并不满足优先队列的父结点必然小于或等于左右子结点的情况，所以要通过比较父结点和其左右子结点来进行下将操作，即当存在一个在数组范围内的结点且大于父结点时则交换两结点，然后递归该过程，直到触底或不满足条件。

​		比较过程中，先比较左结点与父结点的情况，然后再比较右结点的情况，找到最小的一侧进行下降即可。

##### 删除步骤

1. 将尾结点覆盖首结点
2. 删除首结点，进行缩容判断，并可能利用缩容策略进行缩容
3. 以首节点为父结点，对其左右结点进行下降判断：左右节点是否存在，且存在的左右节点是否存在小于父结点的情况，
4. 比对左右两侧，找到满足小于父结点且最小的一侧进行下降即可，同时返回3过程进行递归下降
5. 触底或不满足条件时下降结束，删除完成。

##### 缩容策略

1. 固定缩容：释放一个固定的空间，该方案适合在当前数组较大的时候进行，可以减缓需要缩小的量，当空间较大的时候如果采取折半缩减的话;将可能在绝大多数时间内都不会被采用，从而造成空间浪费。
2. 折半缩容：将当前容量进行折半，该方案适合数组不太大的适合进行缩容，可以更快的缩小容量，但对于较大的空间来说并不会频繁使用到。

### 实现

​		priority_queue优先队列集合结构体，包含动态数组和比较器,同时包含其实际使用长度和实际占用空间容量，该数据结构可以存储多个相同的元素,并不会产生冲突，增删节点后会使用比较器保持该动态数组的相对有序性。

```go
type priority_queue struct {
	data  []interface{}         //动态数组
	len   uint64                //实际使用长度
	cap   uint64                //实际占用的空间的容量
	cmp   comparator.Comparator //该优先队列的比较器
	mutex sync.Mutex            //并发控制锁
}
```

#### 接口

```go
type priority_queueer interface {
	Size() (num uint64)   //返回该容器存储的元素数量
	Clear()               //清空该容器
	Empty() (b bool)      //判断该容器是否为空
	Push(e interface{})   //将元素e插入该容器
	Pop()                 //弹出顶部元素
	Top() (e interface{}) //返回顶部元素
}
```

#### New

​		新建一个priority_queue优先队列容器并返回，初始priority_queue的切片数组为空，如果有传入比较器,则将传入的第一个比较器设为可重复集合默认比较器，如果不传入比较器,在后续的增删过程中将会去寻找默认比较器。

```go
func New(cmps ...comparator.Comparator) (pq *priority_queue) {
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	//比较器为nil时后续的增删将会去寻找默认比较器
	return &priority_queue{
		data:  make([]interface{}, 1, 1),
		len:   0,
		cap:   1,
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}
```

#### Size

​		以priority_queue容器做接收者，返回该容器当前含有元素的数量，由于其数值必然是非负整数，所以选用了**uint64**。

```go
func (pq *priority_queue) Size() (num uint64) {
	if pq == nil {
		pq = New()
	}
	return pq.len
}
```

#### Clear

​		以priority_queue容器做接收者，将该容器中所承载的元素清空。

```go
func (pq *priority_queue) Clear() {
	if pq == nil {
		pq = New()
	}
	pq.mutex.Lock()
	//清空已分配的空间
	pq.data = make([]interface{}, 1, 1)
	pq.len = 0
	pq.cap = 1
	pq.mutex.Unlock()
}
```

#### Empty

​		以priority_queue容器做接收者，判断该priority_queue容器是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true，该判断过程通过含有元素个数进行判断。

```go
func (pq *priority_queue) Empty() bool {
	if pq == nil {
		pq = New()
	}
	return pq.len == 0
}
```

#### Push

​		以priority_queue容器做接收者，在该优先队列中插入元素e,利用比较器和交换使得优先队列保持相对有序状态，插入时,首先将该元素放入末尾,然后通过比较其逻辑上的父结点选择是否上移，扩容策略同vector,先进行翻倍扩容,在进行固定扩容,界限为2^16。

```go
func (pq *priority_queue) Push(e interface{}) {
	if pq == nil {
		pq = New()
	}
	pq.mutex.Lock()
	//判断是否存在比较器,不存在则寻找默认比较器,若仍不存在则直接结束
	if pq.cmp == nil {
		pq.cmp = comparator.GetCmp(e)
	}
	if pq.cmp == nil {
		pq.mutex.Unlock()
		return
	}
	//先判断是否需要扩容,同时使用和vector相同的扩容策略
	//即先翻倍扩容再固定扩容,随后在末尾插入元素e
	if pq.len < pq.cap {
		//还有冗余,直接添加
		pq.data[pq.len] = e
	} else {
		//冗余不足,需要扩容
		if pq.cap <= 65536 {
			//容量翻倍
			if pq.cap == 0 {
				pq.cap = 1
			}
			pq.cap *= 2
		} else {
			//容量增加2^16
			pq.cap += 65536
		}
		//复制扩容前的元素
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
		pq.data[pq.len] = e
	}
	pq.len++
	//到此时,元素以插入到末尾处,同时插入位的元素的下标为pq.len-1,随后将对该位置的元素进行上升
	//即通过比较它逻辑上的父结点进行上升
	pq.up(pq.len - 1)
	pq.mutex.Unlock()
}
```

##### up

​		以priority_queue容器做接收者，用于递归判断任意子结点和其父结点之间的关系，满足上升条件则递归上升，从而保证父节点必然都大于或都小于子节点。

```go
func (pq *priority_queue) up(p uint64) {
   if p == 0 {
      //以及上升到顶部,直接结束即可
      return
   }
   if pq.cmp(pq.data[(p-1)/2], pq.data[p]) > 0 {
      //判断该结点和其父结点的关系
      //满足给定的比较函数的关系则先交换该结点和父结点的数值,随后继续上升即可
      pq.data[p], pq.data[(p-1)/2] = pq.data[(p-1)/2], pq.data[p]
      pq.up((p - 1) / 2)
   }
}
```

#### Pop

​		以priority_queue容器做接收者，在该优先队列中删除顶部元素,利用比较器和交换使得优先队列保持相对有序状态，删除时首先将首结点移到最后一位进行交换,随后删除最后一位即可,然后对首节点进行下降即可，缩容时同vector一样,先进行固定缩容在进行折半缩容,界限为2^16。

```go
func (pq *priority_queue) Pop() {
	if pq == nil {
		pq = New()
	}
	if pq.Empty() {
		return
	}
	pq.mutex.Lock()
	//将最后一位移到首位,随后删除最后一位,即删除了首位,同时判断是否需要缩容
	pq.data[0] = pq.data[pq.len-1]
	pq.len--
	//缩容判断,缩容策略同vector,即先固定缩容在折半缩容
	if pq.cap-pq.len >= 65536 {
		//容量和实际使用差值超过2^16时,容量直接减去2^16
		pq.cap -= 65536
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
	} else if pq.len*2 < pq.cap {
		//实际使用长度是容量的一半时,进行折半缩容
		pq.cap /= 2
		tmp := make([]interface{}, pq.cap, pq.cap)
		copy(tmp, pq.data)
		pq.data = tmp
	}
	//判断是否为空,为空则直接结束
	if pq.Empty() {
		pq.mutex.Unlock()
		return
	}
	//对首位进行下降操作,即对比其逻辑上的左右结点判断是否应该下降,再递归该过程即可
	pq.down(0)
	pq.mutex.Unlock()
}
```

##### down

​		以priority_queue容器做接收者，判断待下沉节点与其左右子节点的大小关系以确定是否进行递归上升，从而保证父节点必然都大于或都小于子节点。

```go
func (pq *priority_queue) down(p uint64) {
	q := p
	//先判断其左结点是否在范围内,然后在判断左结点是否满足下降条件
	if 2*p+1 <= pq.len-1 && pq.cmp(pq.data[p], pq.data[2*p+1]) > 0 {
		q = 2*p + 1
	}
	//在判断右结点是否在范围内,同时若判断右节点是否满足下降条件
	if 2*p+2 <= pq.len-1 && pq.cmp(pq.data[q], pq.data[2*p+2]) > 0 {
		q = 2*p + 2
	}
	//根据上面两次判断,从最小一侧进行下降
	if p != q {
		//进行交互,递归下降
		pq.data[p], pq.data[q] = pq.data[q], pq.data[p]
		pq.down(q)
	}
}
```

#### Top

​		以priority_queue容器做接收者，返回该优先队列容器的顶部元素，如果容器不存在或容器为空,返回nil。

```go
func (pq *priority_queue) Top() (e interface{}) {
	if pq == nil {
		pq = New()
	}
	if pq.Empty() {
		return nil
	}
	pq.mutex.Lock()
	e = pq.data[0]
	pq.mutex.Unlock()
	return e
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/priority_queue"
	"sync"
)

func main() {
	pq := priority_queue.New()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			pq.Push(num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("遍历所有元素同时弹出：")
	size:=pq.Size()
	for i := uint64(0); i < size; i++ {
		fmt.Println(pq.Top())
		pq.Pop()
	}
}

```

注：虽然添加过程是随机的，但由于其本身是相对有序的，所以不论怎么添加都是一个输出结果

> 遍历所有元素同时弹出：
> 0
> 1
> 2
> 3
> 4
> 5
> 6
> 7
> 8
> 9
