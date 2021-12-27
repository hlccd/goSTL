github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		二叉搜索树（Binary Search Tree）不同于之前使用的线性结构，它是一种通过离散的多个点以指针的形式连接起来的**树形结构**。

​		二叉树由一个根节点和根节点下属的多层次的子结点构成，任意一个结点最多只能拥有两个子结点，即左右子结点。基于此种特性，在实现二叉搜索树时，可以仅持有根节点，然后通过根节点去递归访问其子结点以实现寻找到所有结点的。

​		对于二叉搜索树，由于其中缀表达式是有序的，即从小到大（**本文章描述中普遍采用此方案**），所以对于任意一个结点来说，其左子树必然是小于它的，其右子树必然是大于它的。

### 原理

​		对于一个二叉搜索树而言，由于其任意结点和其左右子树之间的关系来看，在插入结点时可以通过遍历找到对应的位置，然后插入新结点即可，当然，对于允许重复的情况来说，可以让结点存储的num+1即可，这样就不需要太多的空间去存储，以达到节约空间的目的。

​		而对于删除的情况来说，也可以通过根节点去找到要删除的结点的位置，然后将该结点删除，随后，可以自行选择左子树的结点或右子树的结点替代该结点，对于用于替代的结点，需要去寻找到它的**前缀节点或后继节点**，前缀节点指该子树中最大的比它小的结点，即左子树的最右结点，后继节点指该子树中最小的比它大的结点，即右子树的最左结点。

​		当要查找到对应的位置的情况时，通过根结点和给定的值，如果小于则进入左子树递归查找，如果大于进入右子树递归查找。

#### 添加策略

​		对于二叉搜索树的添加结点的情况来说，需要考虑的情况主要是该结点是否存在于二叉搜中，如果存在只需要结点num+1即可，否则找到的位置必然是叶子节点，此时只需要插入一个新结点，让其父结点指向它即可。

​		二叉搜索树添加结点的步骤：

1. 通过先以根节点为父结点，再以递归结点为父结点，比对待插入值和父结点之间的大小，小于递归左子树，大于递归右子树，找到对应插入的位置
2. 找到位置后，如果该元素已存在，则结点num+1即可，否则新建结点，同时建立父结点与新结点之间的关系

#### 删除策略

​		对于二叉搜索树的删除来说，首先需要先找到要删除的结点，如果没找到则直接结束即可，当找到后，判断该结点是否有子树，如果没有，直接删除即可，如果有左子树，则寻找其前缀节点，即左子树的最右结点，如果只有右子树，则寻找其后缀节点，即右子树的最左结点，将寻找到的结点替换为当前父结点，然后删除被替换的结点即可。

​		二叉搜索树删除结点的步骤：

1. 通过递归比对遍历寻找到要删除的结点
2. 如果该结点不存在则直接结束，如果有重复则num-1即可结束
3. 如果该结点无左右子树，删除该结点即可结束
4. 如果有左子树，寻找到前缀结点替换该结点，随后删除前缀节点即可结束
5. 如果只有右子树，寻找到后缀结点替换该结点，随后删除后缀节点即可结束
6. 删除完成

### 实现

​		bsTree二叉搜索树结构体，该实例存储二叉树的根节点，同时保存该二叉树已经存储了多少个元素，二叉树中排序使用的比较器在创建时传入,若不传入则在插入首个节点时从默认比较器中寻找，创建时传入是否允许该二叉树出现重复值,如果不允许则进行覆盖,允许则对节点数目增加即可。

```go
type bsTree struct {
	root    *node                 //根节点指针
	size    uint64                //存储元素数量
	cmp     comparator.Comparator //比较器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}
```

​		node树节点结构体，该节点是二叉搜索树的树节点，若该二叉搜索树允许重复则对节点num+1即可,否则对value进行覆盖，二叉搜索树节点不做平衡。

```go
type node struct {
	value interface{} //节点中存储的元素
	num   uint64      //该元素数量
	left  *node       //左节点指针
	right *node       //右节点指针
}
```

#### 接口

```go
type bsTreeer interface {
	Iterator() (i *Iterator.Iterator) //返回包含该二叉树的所有元素,重复则返回多个
	Size() (num uint64)               //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Insert(e interface{})             //向二叉树中插入元素e
	Erase(e interface{})              //从二叉树中删除元素e
	Count(e interface{}) (num uint64) //从二叉树中寻找元素e并返回其个数
}
```

#### New

​		新建一个bsTree二叉搜索树容器并返回，初始根节点为nil，传入该二叉树是否为可重复属性,如果为true则保存重复值,否则对原有相等元素进行覆盖，若有传入的比较器,则将传入的第一个比较器设为该二叉树的比较器。

```go
func New(isMulti bool, Cmp ...comparator.Comparator) (bs *bsTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(Cmp) == 0 {
		cmp = nil
	} else {
		cmp = Cmp[0]
	}
	return &bsTree{
		root:    nil,
		size:    0,
		cmp:     cmp,
		isMulti: isMulti,
		mutex:   sync.Mutex{},
	}
}
```

​		新建一个二叉搜索树节点并返回，将传入的元素e作为该节点的承载元素，该节点的num默认为1,左右子节点设为nil。

```go
func newNode(e interface{}) (n *node) {
   return &node{
      value: e,
      num:   1,
      left:  nil,
      right: nil,
   }
}
```

#### Iterator

​		以bsTree二叉搜索树做接收者，将该二叉树中所有保存的元素将从根节点开始以中缀序列的形式放入迭代器中，若允许重复存储则对于重复元素进行多次放入。

```go
func (bs *bsTree) Iterator() (i *Iterator.Iterator) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	es := bs.root.inOrder()
	i = Iterator.New(&es)
	bs.mutex.Unlock()
	return i
}
```

​		以node二叉搜索树节点做接收者，以中缀序列返回节点集合，若允许重复存储则对于重复元素进行多次放入。

```go
func (n *node) inOrder() (es []interface{}) {
   if n == nil {
      return es
   }
   if n.left != nil {
      es = append(es, n.left.inOrder()...)
   }
   for i := uint64(0); i < n.num; i++ {
      es = append(es, n.value)
   }
   if n.right != nil {
      es = append(es, n.right.inOrder()...)
   }
   return es
}
```

#### Size

​		以bsTree二叉搜索树做接收者，返回该容器当前含有元素的数量，如果容器为nil则创建一个并返回其承载的元素个数。

```go
func (bs *bsTree) Size() (num uint64) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	return bs.size
}
```

#### Clear

​		以bsTree二叉搜索树做接收者，将该容器中所承载的元素清空，将该容器的size置0。

```go
func (bs *bsTree) Clear() {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	bs.root = nil
	bs.size = 0
	bs.mutex.Unlock()
}
```

#### Empty

​		以bsTree二叉搜索树做接收者，判断该二叉搜索树是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true。

```go
func (bs *bsTree) Empty() (b bool) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	return bs.size == 0
}
```

#### Insert

​		以bsTree二叉搜索树做接收者，向二叉树插入元素e,若不允许重复则对相等元素进行覆盖，如果二叉树为空则之间用根节点承载元素e,否则以根节点开始进行查找，不做平衡。

```go
func (bs *bsTree) Insert(e interface{}) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	bs.mutex.Lock()
	if bs.Empty() {
		//二叉树为空,用根节点承载元素e
		if bs.cmp == nil {
			bs.cmp = comparator.GetCmp(e)
		}
		if bs.cmp == nil {
			bs.mutex.Unlock()
			return
		}
		bs.root = newNode(e)
		bs.size++
		bs.mutex.Unlock()
		return
	}
	//二叉树不为空,从根节点开始查找添加元素e
	if bs.root.insert(e, bs.isMulti, bs.cmp) {
		bs.size++
	}
	bs.mutex.Unlock()
}
```

​		以node二叉搜索树节点做接收者，从n节点中插入元素e，如果n节点中承载元素与e不同则根据大小从左右子树插入该元素，如果n节点与该元素相等,且允许重复值,则将num+1否则对value进行覆盖，插入成功返回true,插入失败或不允许重复插入返回false。

```go
func (n *node) insert(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
   if n == nil {
      return false
   }
   //n中承载元素小于e,从右子树继续插入
   if cmp(n.value, e) < 0 {
      if n.right == nil {
         //右子树为nil,直接插入右子树即可
         n.right = newNode(e)
         return true
      } else {
         return n.right.insert(e, isMulti, cmp)
      }
   }
   //n中承载元素大于e,从左子树继续插入
   if cmp(n.value, e) > 0 {
      if n.left == nil {
         //左子树为nil,直接插入左子树即可
         n.left = newNode(e)
         return true
      } else {
         return n.left.insert(e, isMulti, cmp)
      }
   }
   //n中承载元素等于e
   if isMulti {
      //允许重复
      n.num++
      return true
   }
   //不允许重复,直接进行覆盖
   n.value = e
   return false
}
```

#### Erase

​		以bsTree二叉搜索树做接收者，从搜素二叉树中删除元素e，若允许重复记录则对承载元素e的节点中数量记录减一即可，若不允许重复记录则删除该节点同时将前缀节点或后继节点更换过来以保证二叉树的不发送断裂，如果该二叉树仅持有一个元素且根节点等价于待删除元素,则将二叉树根节点置为nil。

```go
func (bs *bsTree) Erase(e interface{}) {
	if bs == nil {
		//创建一个允许插入重复值的二叉搜
		bs = New(true)
	}
	if bs.size == 0 {
		return
	}
	bs.mutex.Lock()
	if bs.size == 1 && bs.cmp(bs.root.value, e) == 0 {
		//二叉树仅持有一个元素且根节点等价于待删除元素,将二叉树根节点置为nil
		bs.root = nil
		bs.size = 0
		bs.mutex.Unlock()
		return
	}
	//从根节点开始删除元素e
	//如果删除成功则将size-1
	if bs.root.delete(e, bs.isMulti, bs.cmp) {
		bs.size--
	}
	bs.mutex.Unlock()
}
```

​		以node二叉搜索树节点做接收者，从n节点中删除元素e，如果n节点中承载元素与e不同则根据大小从左右子树删除该元素，如果n节点与该元素相等,且允许重复值,则将num-1否则直接删除该元素，删除时先寻找该元素的前缀节点,若不存在则寻找其后继节点进行替换，替换后删除该节点。

```go
func (n *node) delete(e interface{}, isMulti bool, cmp comparator.Comparator) (b bool) {
   if n == nil {
      return false
   }
   //n中承载元素小于e,从右子树继续删除
   if cmp(n.value, e) < 0 {
      if n.right == nil {
         //右子树为nil,删除终止
         return false
      }
      if cmp(e, n.right.value) == 0 && (!isMulti || n.right.num == 1) {
         if n.right.left == nil && n.right.right == nil {
            //右子树可直接删除
            n.right = nil
            return true
         }
      }
      //从右子树继续删除
      return n.right.delete(e, isMulti, cmp)
   }
   //n中承载元素大于e,从左子树继续删除
   if cmp(n.value, e) > 0 {
      if n.left == nil {
         //左子树为nil,删除终止
         return false
      }
      if cmp(e, n.left.value) == 0 && (!isMulti || n.left.num == 1) {
         if n.left.left == nil && n.left.right == nil {
            //左子树可直接删除
            n.left = nil
            return true
         }
      }
      //从左子树继续删除
      return n.left.delete(e, isMulti, cmp)
   }
   //n中承载元素等于e
   if (*n).num > 1 && isMulti {
      //允许重复且个数超过1,则减少num即可
      (*n).num--
      return true
   }
   if n.left == nil && n.right == nil {
      //该节点无前缀节点和后继节点,删除即可
      *(&n) = nil
      return true
   }
   if n.left != nil {
      //该节点有前缀节点,找到前缀节点进行交换并删除即可
      ln := n.left
      if ln.right == nil {
         n.value = ln.value
         n.num = ln.num
         n.left = ln.left
      } else {
         for ln.right.right != nil {
            ln = ln.right
         }
         n.value = ln.right.value
         n.num = ln.right.num
         ln.right = ln.right.left
      }
   } else if (*n).right != nil {
      //该节点有后继节点,找到后继节点进行交换并删除即可
      tn := n.right
      if tn.left == nil {
         n.value = tn.value
         n.num = tn.num
         n.right = tn.right
      } else {
         for tn.left.left != nil {
            tn = tn.left
         }
         n.value = tn.left.value
         n.num = tn.left.num
         tn.left = tn.left.right
      }
      return true
   }
   return true
}
```

#### Count

​		以bsTree二叉搜索树做接收者，从搜素二叉树中查找元素e的个数，如果找到则返回该二叉树中和元素e相同元素的个数，如果不允许重复则最多返回1，如果未找到则返回0。

```go
func (bs *bsTree) Count(e interface{}) (num uint64) {
	if bs == nil {
		//二叉树不存在,返回0
		return 0
	}
	if bs.Empty() {
		//二叉树为空,返回0
		return 0
	}
	bs.mutex.Lock()
	//从根节点开始查找并返回查找结果
	num = bs.root.search(e, bs.isMulti, bs.cmp)
	bs.mutex.Unlock()
	return num
}
```

​		以node二叉搜索树节点做接收者，从n节点中查找元素e并返回存储的个数，如果n节点中承载元素与e不同则根据大小从左右子树查找该元素，如果n节点与该元素相等,则直接返回其个数。

```go
func (n *node) search(e interface{}, isMulti bool, cmp comparator.Comparator) (num uint64) {
	if n == nil {
		return 0
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.search(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.search(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.num
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/bsTree"
	"sync"
)

func main() {
	bs := bsTree.New(true)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		bs.Insert(i)
		go func() {
			bs.Insert(i)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("遍历输出所有插入的元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	fmt.Println("删除一次二叉搜索树中存在的元素,存在重复的将会被剩下")
	for i := 0; i < 10; i++ {
		bs.Erase(i)
	}
	fmt.Println("输出剩余的重复元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
}

```

注：由于过程中的增删过程是并发执行的，所以其结果和下方示例并不完全相同

> 遍历输出所有插入的元素
> 0
> 1
> 2
> 3
> 3
> 4
> 5
> 6
> 6
> 6
> 7
> 7
> 7
> 7
> 7
> 8
> 9
> 10
> 10
> 10
> 删除一次二叉搜索树中存在的元素,存在重复的将会被剩下
> 输出剩余的重复元素
> 3
> 6
> 6
> 7
> 7
> 7
> 7
> 10
> 10
> 10
> 

