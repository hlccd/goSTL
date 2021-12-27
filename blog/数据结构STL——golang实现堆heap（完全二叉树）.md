github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		本次采用**完全二叉树Complete Binary Tree**的形式实现堆。

​		堆（heap）是一类特殊的数据结构的统称，堆通常是一个可以被看做一棵树的数组对象。堆总是满足下列性质：

- 堆中某个结点的值总是不大于或不小于其父结点的值；
- 堆总是一棵完全二叉树。

​		堆的主要特点为：它的**父结点必然是小于或等于左右子结点**（该文章叙述中使用**小顶堆**，大顶堆大小情况相反）。

​		在本次实践过程中使用**完全二叉树**实现。

### 原理

##### 完全二叉树实现堆

​		对于一个完全二叉树来说，首先，完全二叉树除开底部的一层可能存在不满的情况外，其他上层必然是满的。对于完全二叉树的一个节点来说，它的左节点和右节点不使用下标来进行取得，而是使用其指针进行寻找。

​		同时，在实现中，需要使用**size**来记录该堆中累计存储了多少元素。

​		在插入和删除过程中，同数组实现类似，插入需要在最后一个节点插入，随后上升即可，删除需要将首节点替换为尾节点，随后下降首届点。而这两个过程中，核心点就在于如何**查找尾节点**。

##### 查找尾节点

​		在查找尾节点时，可以参考数组实现过程中的下标，同时由于只需要查找尾节点，而尾节点的下标，恰好是**size-1**的值。

​		所以，我们可以利用堆当前存储的元素数量去寻找最后一个节点和其父节点。

​		同数组类似，“下标”从0开始，观察前一组插入的节点：（size+1）且为二进制

- 第一个节点->10

- 第二个节点->11

- 第三个节点->100

- 第四个节点->101

- 第五个节点->110

- 第六个节点->111

- 第七个节点->1000

- 第八个节点->1001

- 第九个节点->1010

- ······

  

  同时观察其在二叉树的位置可知，其可以根据其下标转化为二进制的值及1010101串来找到插入节点的父节点，例：

  - 第11个节点：**1100**该节点从根节点获取的路径为：**右左左**
  - 第15给节点：**10000**该节点从根节点获取的路径为：**左左左左**
  - 第18给节点：**10011**该节点从根节点获取的路径为：**左左右右**
  - ···

观察可知，对于添加节点来说，可以通过其下标即当前的size+1转化为二进制，随后删掉其第一个bit，在根据后续的bit位，将1看作转入右节点，将0转入左节点，即可寻找到对应的节点，而若要找到其父节点只需要减少一步跳转即可。

​		而对于删除节点来说，则不需要先减少size，直接找到末尾节点即可，方法同上。

#### 添加策略

​		在解决了查找尾节点和其父节点的问题之后，对于添加策略其实和数组形式实现是基本一样的。

​		即：当插入一个结点的时候，可以先将其放入完全二叉搜的末尾，然后根据插入的值和它父结点进行比较，如果插入值小于父节点时则交换两结点值，随后在用交换后的结点与其父结点比较，重复该过程直到抵达顶点或不满足条件即可。

##### 添加步骤

1. 通过查找尾节点的算法找到对应尾节点和其父节点，利用父节点插入
2. 结点放入尾部后，通过比较该结点与父结点的值进行交换
3. 满足交换条件，重复回到2过程
4. 到达顶部或不满足交换条件，添加结束，插入完成

#### 删除策略

​		完全二叉树实现堆的方式的删除策略类似于数组形式实现堆的删除策略，区别只是在需要用特殊的查找尾节点及其父节点的方法（上文已做介绍）。

​		即：不同于添加时直接利用尾结点的情况，由于删除的是首结点，而同时删除需要减少一个空间，所以可以考虑将首位结点进行交换，或者其实直接**用尾结点覆写首节点**即可，这样就实现了首节点的删除，同时利用之前查找到的尾节点的父节点，将其对应侧的指针设为nil，即实现了尾节点的删除。

​		上一过程完成后，首节点以及被删除了，但移到首节点的新首节点可能并不满足优先队列的父结点必然小于或等于左右子结点的情况，所以要通过比较父结点和其左右子结点来进行下将操作，即当存在一个在数组范围内的结点且大于父结点时则交换两结点，然后递归该过程，直到触底或不满足条件。

​		比较过程中，先比较左结点与父结点的情况，然后再比较右结点的情况，找到最小的一侧进行下降即可。

##### 删除步骤

1. 通过查找尾节点的方法找到尾节点的父节点，并用尾节点覆盖首节点。
2. 删除尾节点
3. 以首节点为父结点，对其左右结点进行下降判断：左右节点是否存在，且存在的左右节点是否存在小于父结点的情况，
4. 比对左右两侧，找到满足小于父结点且最小的一侧进行下降即可，同时返回3过程进行递归下降
5. 触底或不满足条件时下降结束，删除完成。

### 实现

​		cbTree完全二叉树树结构体，该实例存储二叉树的根节点，同时保存该二叉树已经存储了多少个元素，二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找。

```go
type cbTree struct {
	root  *node                 //根节点指针
	size  uint64                //存储元素数量
	cmp   comparator.Comparator //比较器
	mutex sync.Mutex            //并发控制锁
}
```

​		node树节点结构体，该节点是完全二叉树的树节点，该节点除了保存承载元素外,还将保存父节点、左右子节点的指针。

```go
type node struct {
   value  interface{} //节点中存储的元素
   parent *node       //父节点指针
   left   *node       //左节点指针
   right  *node       //右节点指针
}
```

#### 接口

```go
type cbTreer interface {
	Iterator() (i *Iterator.Iterator) //返回包含该二叉树的所有元素
	Size() (num uint64)               //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Push(e interface{})               //向二叉树中插入元素e
	Pop()                             //从二叉树中弹出顶部元素
	Top() (e interface{})             //返回该二叉树的顶部元素
}
```

#### New

​		新建一个cbTree完全二叉树容器并返回，初始根节点为nil，若有传入的比较器,则将传入的第一个比较器设为该二叉树的比较器。

```go
func New(Cmp ...comparator.Comparator) (cb *cbTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(Cmp) > 0 {
		cmp = Cmp[0]
	}
	return &cbTree{
		root:  nil,
		size:  0,
		cmp:   cmp,
		mutex: sync.Mutex{},
	}
}
```

##### newNode

​		新建一个完全二叉树节点并返回，将传入的元素e作为该节点的承载元素，将传入的parent节点作为其父节点,左右节点设为nil。

```go
func newNode(parent *node, e interface{}) (n *node) {
   return &node{
      value:  e,
      parent: parent,
      left:   nil,
      right:  nil,
   }
}
```

#### Iterator

​		以cbTree完全二叉树做接收者，将该二叉树中所有保存的元素将从根节点开始以前缀序列的形式放入迭代器中。

```go
func (cb *cbTree) Iterator() (i *Iterator.Iterator) {
   if cb == nil {
      cb = New()
   }
   cb.mutex.Lock()
   es := cb.root.frontOrder()
   i = Iterator.New(&es)
   cb.mutex.Unlock()
   return i
}
```

##### frontOrder

​		以node节点做接收者，以前缀序列返回节点集合。

```go
func (n *node) frontOrder() (es []interface{}) {
   if n == nil {
      return es
   }
   es = append(es, n.value)
   if n.left != nil {
      es = append(es, n.left.frontOrder()...)
   }
   if n.right != nil {
      es = append(es, n.right.frontOrder()...)
   }
   return es
}
```

#### Size

​		以cbTree完全二叉树做接收者，返回该容器当前含有元素的数量，如果容器为nil返回0。

```go
func (cb *cbTree) Size() (num uint64) {
	if cb == nil {
		cb = New()
	}
	return cb.size
}
```

#### Clear

​		以cbTree完全二叉树做接收者，将该容器中所承载的元素清空，将该容器的size置0。

```go
func (cb *cbTree) Clear() {
	if cb == nil {
		cb = New()
	}
	cb.mutex.Lock()
	cb.root = nil
	cb.size = 0
	cb.mutex.Unlock()
}
```

#### Empty

​		以cbTree完全二叉树做接收者，判断该完全二叉树树是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true。

```go
func (cb *cbTree) Empty() (b bool) {
	if cb == nil {
		cb = New()
	}
	return cb.size == 0
}
```

#### lastParent

​		以node节点做接收者，根据传入数值通过转化为二进制的方式模拟查找最后一个父节点，由于查找父节点的路径等同于转化为二进制后除开首位的中间值,故该方案是可行的。

```go
func (n *node) lastParent(num uint64) (ans *node) {
   if num > 3 {
      //去掉末尾的二进制值
      arr := make([]byte, 0, 64)
      ans = n
      for num > 0 {
         //转化为二进制
         arr = append(arr, byte(num%2))
         num /= 2
      }
      //去掉首位的二进制值
      for i := len(arr) - 2; i >= 1; i-- {
         if arr[i] == 1 {
            ans = ans.right
         } else {
            ans = ans.left
         }
      }
      return ans
   }
   return n
}
```

#### Push

​		以cbTree完全二叉树做接收者，向二叉树插入元素e,将其放入完全二叉树的最后一个位置,随后进行元素上升，如果二叉树本身为空,则直接将根节点设为插入节点元素即可。

```go
func (cb *cbTree) Push(e interface{}) {
	if cb == nil {
		cb=New()
	}
	cb.mutex.Lock()
	if cb.Empty() {
		if cb.cmp == nil {
			cb.cmp = comparator.GetCmp(e)
		}
		if cb.cmp == nil {
			cb.mutex.Unlock()
			return
		}
		cb.root = newNode(nil, e)
		cb.size++
	} else {
		cb.size++
		cb.root.insert(cb.size, e, cb.cmp)
	}
	cb.mutex.Unlock()
}
```

##### insert

​		以node节点做接收者，从该节点插入元素e,并根据传入的num寻找最后一个父节点用于插入最后一位值，随后对插入值进行上升处理。

```go
func (n *node) insert(num uint64, e interface{}, cmp comparator.Comparator) {
   if n == nil {
      return
   }
   //寻找最后一个父节点
   n = n.lastParent(num)
   //将元素插入最后一个节点
   if num%2 == 0 {
      n.left = newNode(n, e)
      n = n.left
   } else {
      n.right = newNode(n, e)
      n = n.right
   }
   //对插入的节点进行上升
   n.up(cmp)
}
```

##### up

​		以node节点做接收者，对该节点进行上升，当该节点存在且父节点存在时,若该节点小于夫节点，则在交换两个节点值后继续上升即可。

```go
func (n *node) up(cmp comparator.Comparator) {
   if n == nil {
      return
   }
   if n.parent == nil {
      return
   }
   //该节点和父节点都存在
   if cmp(n.parent.value, n.value) > 0 {
      //该节点值小于父节点值,交换两节点值,继续上升
      n.parent.value, n.value = n.value, n.parent.value
      n.parent.up(cmp)
   }
}
```

#### Pop

​		以cbTree完全二叉树做接收者，从完全二叉树中删除顶部元素e，将该顶部元素于最后一个元素进行交换，随后删除最后一个元素，再将顶部元素进行下沉处理即可。

```go
func (cb *cbTree) Pop() {
	if cb == nil {
		return
	}
	if cb.Empty() {
		return
	}
	cb.mutex.Lock()
	if cb.size == 1 {
		//该二叉树仅剩根节点,直接删除即可
		cb.root = nil
	} else {
		//该二叉树删除根节点后还有其他节点可生为跟节点
		cb.root.delete(cb.size, cb.cmp)
	}
	cb.size--
	cb.mutex.Unlock()
}
```

##### delete

​		以node节点做接收者，从删除该,并根据传入的num寻找最后一个父节点用于替换删除，随后对替换后的值进行下沉处理即可

```go
func (n *node) delete(num uint64, cmp comparator.Comparator) {
	if n == nil {
		return
	}
	//寻找最后一个父节点
	ln := n.lastParent(num)
	if num%2 == 0 {
		n.value = ln.left.value
		ln.left = nil
	} else {
		n.value = ln.right.value
		ln.right = nil
	}
	//对交换后的节点进行下沉
	n.down(cmp)
}
```

##### down

​		以node节点做接收者，对该节点进行下沉，当该存在右节点且小于自身元素时,与右节点进行交换并继续下沉，否则当该存在左节点且小于自身元素时,与左节点进行交换并继续下沉，当左右节点都不存在或都大于自身时下沉停止。

```go
func (n *node) down(cmp comparator.Comparator) {
	if n == nil {
		return
	}
	if n.right != nil && cmp(n.left.value, n.right.value) >= 0 {
		n.right.value, n.value = n.value, n.right.value
		n.right.down(cmp)
		return
	}
	if n.left != nil && cmp(n.value, n.left.value) >= 0 {
		n.left.value, n.value = n.value, n.left.value
		n.left.down(cmp)
		return
	}
}
```

#### Top

​		以cbTree完全二叉树做接收者，返回该完全二叉树的顶部元素，当该完全二叉树不存在或根节点不存在时返回nil。

```go
func (cb *cbTree) Top() (e interface{}) {
	if cb == nil {
		cb=New()
	}
	cb.mutex.Lock()
	e = cb.root.value
	cb.mutex.Unlock()
	return e
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/heap"
	"sync"
)

func main() {
	h := heap.New()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			h.Push(num)
		}(i)
	}
	fmt.Println("利用迭代器输出堆中存储的所有元素:")
	for i := h.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("依次输出顶部元素:")
	for !h.Empty() {
		fmt.Println(h.Top())
		h.Pop()
	}
}
```

注：虽然添加过程是随机的，但由于其本身是相对有序的，所以不论怎么添加都是一个输出结果

> 利用迭代器输出堆中存储的所有元素:
> 0
> 3
> 4
> 6
> 8
> 5
> 9
> 1
> 2
> 7
> 依次输出顶部元素:
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
