github仓库存储地址：https://github.com/hlccd/goSTL

### 概述

​		自平衡二叉查找树（Self-Balancing Binary Search Tree），简称为平衡二叉树，一般以其发明者的名称缩写命名为avl树。

​		对于一颗平衡二叉树来说，一方面它需要**满足二叉搜索树的性质**，即父结点大于左结点小于右结点，另一方面，**该树中每个结点的左右子结点的高度差不能超过1**，即其平衡因子最大为1，当插入或删除结点后导致平衡因子超过1时，则需要通过旋转的方式对其进行调节。

### 原理

​		对于平衡二叉搜来说，它的搜索同二叉搜索树一样，只需要定向遍历树结点即可，而在增加和删除的过程中，第一步可以先按照二叉搜索树的方式进行增加和删除，即直接插入增加、删除前缀or后继结点即可，增删完成后，则需要对增删结点的父结点进行平衡因子的判断，即重新计算其左右子树的高度差，当期高度差超过1时则需要进行旋转，旋转策略如下：

1. 左旋：将该结点向其左结点方向旋转，使得右结点放在原结点的位置，右节点的左子树结点即为该结点。
2. 右旋：将该结点向其右结点方向旋转，使得左结点放在原结点的位置，左节点的右子树结点即为该结点。
3. 右左旋：将原节点的右节点先进行右旋再将原节点进行左旋
4. 左右旋：将原节点的左结点先进行左旋再将原结点进行右旋

#### 情况分析

​		以插入结点为例，删除节点可视为插入结点的逆运算，即**删除左结点等于插入右结点，删除右结点等于插入左结点**。

##### 单左旋转

​		按照{4，5，6}的顺序插入树内，在插入6时会出现4结点的不平衡情况，即单侧不平衡，对于此类情况，只需要将4结点进行左旋即可，即将5结点放到4结点的位置，4结点设为5节点的左结点。

##### 单右旋转

​		在刚刚的基础上，插入{3，2}结点，在插入2结点的时候会出现4结点的不平衡，并且也是单侧不平衡，对于此类情况，只需要将4结点进行右旋即可，即将3结点放到4结点的位置上，4结点设为3结点的右结点。

##### 右左旋转

​		在此基础上继续插入{8，7}结点，在插入7结点的时候会出现6结点的不平衡，同时，由于8结点存在左结点，故只能先对8结点进行右旋，使得7结点处于8结点的位置，8结点设为7结点的右结点，然后就变成了单左旋转的情况，再进行一次对6结点的单左旋转即可。

##### 左右旋转

​		在此基础上继续插入{0，1}结点，在插入1结点的时候会出现2结点的不平衡，考虑到0结点存在右结点，故需要先对0结点进行一次左旋，将1结点放到0结点的位置，0结点设为1结点的左结点，然后就变成了单右旋转的情况，再对2结点进行一次单右旋转即可。

### 实现

​		avlTree平衡二叉树结构体，该实例存储平衡二叉树的根节点，同时保存该二叉树已经存储了多少个元素，二叉树中排序使用的比较器在创建时传入，若不传入则在插入首个节点时从默认比较器中寻找，创建时传入是否允许该二叉树出现重复值,如果不允许则进行覆盖,允许则对节点数目增加即可。

```go
type avlTree struct {
	root    *node                 //根节点指针
	size    int                   //存储元素数量
	cmp     comparator.Comparator //比较器
	isMulti bool                  //是否允许重复
	mutex   sync.Mutex            //并发控制锁
}
```

​		node树节点结构体，该节点是平衡二叉树的树节点，若该平衡二叉树允许重复则对节点num+1即可,否则对value进行覆盖，平衡二叉树节点当左右子节点深度差超过1时进行左右旋转以实现平衡。

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
type avlTreer interface {
	Iterator() (i *Iterator.Iterator) //返回包含该二叉树的所有元素,重复则返回多个
	Size() (num int)                  //返回该二叉树中保存的元素个数
	Clear()                           //清空该二叉树
	Empty() (b bool)                  //判断该二叉树是否为空
	Insert(e interface{})             //向二叉树中插入元素e
	Erase(e interface{})              //从二叉树中删除元素e
	Count(e interface{}) (num int)    //从二叉树中寻找元素e并返回其个数
}
```

#### New

​		新建一个avlTree平衡二叉树容器并返回，初始根节点为nil，传入该二叉树是否为可重复属性,如果为true则保存重复值,否则对原有相等元素进行覆盖，若有传入的比较器,则将传入的第一个比较器设为该二叉树的比较器。

```go
func New(isMulti bool, cmps ...comparator.Comparator) (avl *avlTree) {
	//判断是否有传入比较器,若有则设为该二叉树默认比较器
	var cmp comparator.Comparator
	if len(cmps) == 0 {
		cmp = nil
	} else {
		cmp = cmps[0]
	}
	return &avlTree{
		root:    nil,
		size:    0,
		cmp:     cmp,
		isMulti: isMulti,
	}
}
```

​		新建一个平衡二叉树节点并返回，将传入的元素e作为该节点的承载元素，该节点的num和depth默认为1,左右子节点设为nil。

```go
func newNode(e interface{}) (n *node) {
	return &node{
		value: e,
		num:   1,
		depth: 1,
		left:  nil,
		right: nil,
	}
}
```

#### Iterator

​		以avlTree平衡二叉树做接收者，将该二叉树中所有保存的元素将从根节点开始以中缀序列的形式放入迭代器中，若允许重复存储则对于重复元素进行多次放入。

```go
func (avl *avlTree) Iterator() (i *Iterator.Iterator) {
	if avl == nil {
		return nil
	}
	avl.mutex.Lock()
	es:=avl.root.inOrder()
	i = Iterator.New(&es)
	avl.mutex.Unlock()
	return i
}
```

​		以node平衡二叉树节点做接收者，以中缀序列返回节点集合，若允许重复存储则对于重复元素进行多次放入。

```go
func (n *node) inOrder() (es []interface{}) {
	if n == nil {
		return es
	}
	if n.left != nil {
		es = append(es, n.left.inOrder()...)
	}
	for i := 0; i < n.num; i++ {
		es = append(es, n.value)
	}
	if n.right != nil {
		es = append(es, n.right.inOrder()...)
	}
	return es
}
```

#### Size

​		以avlTree平衡二叉树做接收者，返回该容器当前含有元素的数量，如果容器为nil返回0。

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

​		以avlTree平衡二叉树做接收者，将该容器中所承载的元素清空，将该容器的size置0。

```go
func (avl *avlTree) Clear() {
	if avl == nil {
		return
	}
	avl.mutex.Lock()
	avl.root = nil
	avl.size = 0
	avl.mutex.Unlock()
}
```

#### Empty

​		以avlTree平衡二叉树做接收者，判断该二叉搜索树是否含有元素，如果含有元素则不为空,返回false，如果不含有元素则说明为空,返回true，如果容器不存在,返回true。

```go
func (avl *avlTree) Empty() (b bool) {
	if avl == nil {
		return true
	}
	if avl.size > 0 {
		return false
	}
	return true
}
```

#### 旋转

##### getDepth

​		以node平衡二叉树节点做接收者，返回该节点的深度,节点不存在返回0。

```go
func (n *node) getDepth() (depth int) {
   if n == nil {
      return 0
   }
   return n.depth
}
```

##### max

​		返回a和b中较大的值

```go
func max(a, b int) (m int) {
	if a > b {
		return a
	} else {
		return b
	}
}
```

##### leftRotate

​		以node平衡二叉树节点做接收者，将该节点向左节点方向转动,使右节点作为原来节点,并返回右节点，同时将右节点的左节点设为原节点的右节点。

```go
func (n *node) leftRotate() (m *node) {
   //左旋转
   headNode := n.right
   n.right = headNode.left
   headNode.left = n
   //更新结点高度
   n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
   headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
   return headNode
}
```

##### rightRotate

​		以node平衡二叉树节点做接收者，将该节点向右节点方向转动,使左节点作为原来节点,并返回左节点，同时将左节点的右节点设为原节点的左节点。

```go
func (n *node) rightRotate() (m *node) {
   //右旋转
   headNode := n.left
   n.left = headNode.right
   headNode.right = n
   //更新结点高度
   n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
   headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
   return headNode
}
```

##### rightLeftRotate

​		以node平衡二叉树节点做接收者，将原节点的右节点先进行右旋再将原节点进行左旋。

```go
func (n *node) rightLeftRotate() (m *node) {
   //右旋转,之后左旋转
   //以失衡点右结点先右旋转
   sonHeadNode := n.right.rightRotate()
   n.right = sonHeadNode
   //再以失衡点左旋转
   return n.leftRotate()
}
```

##### leftRightRotate

​		以node平衡二叉树节点做接收者，将原节点的左节点先进行左旋再将原节点进行右旋。

```go
func (n *node) leftRightRotate() (m *node) {
   //左旋转,之后右旋转
   //以失衡点左结点先左旋转
   sonHeadNode := n.left.leftRotate()
   n.left = sonHeadNode
   //再以失衡点左旋转
   return n.rightRotate()
}
```

##### getMin

​		以node平衡二叉树节点做接收者，返回n节点的最小元素的承载的元素和数量，即该节点父节点的后缀节点的元素和数量。

```go
func (n *node) getMin() (e interface{}, num int) {
   if n == nil {
      return nil, 0
   }
   if n.left == nil {
      return n.value, n.num
   } else {
      return n.left.getMin()
   }
}
```

##### adjust

​		以node平衡二叉树节点做接收者，对n节点进行旋转以保持节点左右子树平衡。

```go
func (n *node) adjust() (m *node) {
   if n.right.getDepth()-n.left.getDepth() >= 2 {
      //右子树高于左子树且高度差超过2,此时应当对n进行左旋
      if n.right.right.getDepth() > n.right.left.getDepth() {
         //由于右右子树高度超过右左子树,故可以直接左旋
         n = n.leftRotate()
      } else {
         //由于右右子树不高度超过右左子树
         //所以应该先右旋右子树使得右子树高度不超过左子树
         //随后n节点左旋
         n = n.rightLeftRotate()
      }
   } else if n.left.getDepth()-n.right.getDepth() >= 2 {
      //左子树高于右子树且高度差超过2,此时应当对n进行右旋
      if n.left.left.getDepth() > n.left.right.getDepth() {
         //由于左左子树高度超过左右子树,故可以直接右旋
         n = n.rightRotate()
      } else {
         //由于左左子树高度不超过左右子树
         //所以应该先左旋左子树使得左子树高度不超过右子树
         //随后n节点右旋
         n = n.leftRightRotate()
      }
   }
   return n
}
```

#### Insert

​		以avlTree平衡二叉树做接收者，向二叉树插入元素e,若不允许重复则对相等元素进行覆盖，如果二叉树为空则之间用根节点承载元素e,否则以根节点开始进行查找，当节点左右子树高度差超过1时将进行旋转以保持平衡。

```go
func (avl *avlTree) Insert(e interface{}) {
	if avl == nil {
		return
	}
	avl.mutex.Lock()
	if avl.Empty() {
		if avl.cmp == nil {
			avl.cmp = comparator.GetCmp(e)
		}
		if avl.cmp == nil {
			return
		}
		//二叉树为空,用根节点承载元素e
		avl.root = newNode(e)
		avl.size = 1
		avl.mutex.Unlock()
		return
	}
	//从根节点进行插入,并返回节点,同时返回是否插入成功
	var b bool
	avl.root, b = avl.root.insert(e, avl.isMulti, avl.cmp)
	if b {
		//插入成功,数量+1
		avl.size++
	}
	avl.mutex.Unlock()
}
```

​		以node平衡二叉树节点做接收者，从n节点中插入元素e，如果n节点中承载元素与e不同则根据大小从左右子树插入该元素，如果n节点与该元素相等,且允许重复值,则将num+1否则对value进行覆盖，插入成功返回true,插入失败或不允许重复插入返回false。

```go
func (n *node) insert(e interface{}, isMulti bool, cmp comparator.Comparator) (m *node, b bool) {
	//节点不存在,应该创建并插入二叉树中
	if n == nil {
		return newNode(e), true
	}
	if cmp(e, n.value) < 0 {
		//从左子树继续插入
		n.left, b = n.left.insert(e, isMulti, cmp)
		if b {
			//插入成功,对该节点进行平衡
			n = n.adjust()
		}
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		return n, b
	}
	if cmp(e, n.value) > 0 {
		//从右子树继续插入
		n.right, b = n.right.insert(e, isMulti, cmp)
		if b {
			//插入成功,对该节点进行平衡
			n = n.adjust()
		}
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		return n, b
	}
	//该节点元素与待插入元素相同
	if isMulti {
		//允许重复,数目+1
		n.num++
		return n, true
	}
	//不允许重复,对值进行覆盖
	n.value = e
	return n, false
}
```

#### Erase

​		以avlTree平衡二叉树做接收者，从平衡二叉树中删除元素e，若允许重复记录则对承载元素e的节点中数量记录减一即可，若不允许重复记录则删除该节点同时将前缀节点或后继节点更换过来以保证二叉树的不发送断裂，如果该二叉树仅持有一个元素且根节点等价于待删除元素,则将二叉树根节点置为nil。

```go
func (avl *avlTree) Erase(e interface{}) {
	if avl == nil {
		return
	}
	if avl.Empty() {
		return
	}
	avl.mutex.Lock()
	if avl.size == 1 && avl.cmp(avl.root.value, e) == 0 {
		//二叉树仅持有一个元素且根节点等价于待删除元素,将二叉树根节点置为nil
		avl.root = nil
		avl.size = 0
		avl.mutex.Unlock()
		return
	}
	//从根节点进行插入,并返回节点,同时返回是否删除成功
	var b bool
	avl.root, b = avl.root.erase(e, avl.cmp)
	if b {
		avl.size--
	}
	avl.mutex.Unlock()
}
```

​		以node平衡二叉树节点做接收者，从n节点中删除元素e，如果n节点中承载元素与e不同则根据大小从左右子树删除该元素，如果n节点与该元素相等,且允许重复值,则将num-1否则直接删除该元素，删除时先寻找该元素的前缀节点,若不存在则寻找其后继节点进行替换，替换后删除该节点。

```go
func (n *node) erase(e interface{}, cmp comparator.Comparator) (m *node, b bool) {
	if n == nil {
		//待删除值不存在,删除失败
		return n, false
	}
	if cmp(e, n.value) < 0 {
		//从左子树继续删除
		n.left, b = n.left.erase(e, cmp)
	} else if cmp(e, n.value) > 0 {
		//从右子树继续删除
		n.right, b = n.right.erase(e, cmp)
	} else if cmp(e, n.value) == 0 {
		//存在相同值,从该节点删除
		b = true
		if n.num > 1 {
			//有重复值,节点无需删除,直接-1即可
			n.num--
		} else {
			//该节点需要被删除
			if n.left != nil && n.right != nil {
				//找到该节点后继节点进行交换删除
				n.value, n.num = n.right.getMin()
				//从右节点继续删除,同时可以保证删除的节点必然无左节点
				n.right, b = n.right.erase(n.value, cmp)
			} else if n.left != nil {
				n = n.left
			} else {
				n = n.right
			}
		}
	}
	//当n节点仍然存在时,对其进行调整
	if n != nil {
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		n = n.adjust()
	}
	return n, b
}
```

#### Count

​		以avlTree平衡二叉树做接收者，从搜素二叉树中查找元素e的个数，如果找到则返回该二叉树中和元素e相同元素的个数，如果不允许重复则最多返回1，如果未找到则返回0。

```go
func (avl *avlTree) Count(e interface{}) (num int) {
	if avl == nil {
		//二叉树为空,返回0
		return 0
	}
	if avl.Empty() {
		return 0
	}
	avl.mutex.Lock()
	num = avl.root.count(e, avl.isMulti, avl.cmp)
	avl.mutex.Unlock()
	return num
}
```

​		以node二叉搜索树节点做接收者，从n节点中查找元素e并返回存储的个数，如果n节点中承载元素与e不同则根据大小从左右子树查找该元素，如果n节点与该元素相等,则直接返回其个数。

```go
func (n *node) count(e interface{}, isMulti bool, cmp comparator.Comparator) (num int) {
	if n == nil {
		return 0
	}
	//n中承载元素小于e,从右子树继续查找并返回结果
	if cmp(n.value, e) < 0 {
		return n.right.count(e, isMulti, cmp)
	}
	//n中承载元素大于e,从左子树继续查找并返回结果
	if cmp(n.value, e) > 0 {
		return n.left.count(e, isMulti, cmp)
	}
	//n中承载元素等于e,直接返回结果
	return n.num
}
```

#### Find

​		以avlTree平衡二叉树做接收者，从搜素二叉树中查找以元素e为索引信息的全部信息，如果找到则返回该二叉树中和索引元素e相同的元素的全部信息，如果未找到则返回nil。

```go
func (avl *avlTree) Find(e interface{}) (ans interface{}) {
   if avl == nil {
      //二叉树为空,返回0
      return 0
   }
   if avl.Empty() {
      return 0
   }
   avl.mutex.Lock()
   ans = avl.root.find(e, avl.isMulti, avl.cmp)
   avl.mutex.Unlock()
   return ans
}
```

​		以node二叉搜索树节点做接收者，从n节点中查找元素e并返回其存储的全部信息，如果n节点中承载元素与e不同则根据大小从左右子树查找该元素，如果n节点与该元素相等,则直接返回其信息即可，若未找到该索引信息则返回nil。

```go
func (n *node) find(e interface{}, isMulti bool, cmp comparator.Comparator) (ans interface{}) {
   if n == nil {
      return nil
   }
   //n中承载元素小于e,从右子树继续查找并返回结果
   if cmp(n.value, e) < 0 {
      return n.right.count(e, isMulti, cmp)
   }
   //n中承载元素大于e,从左子树继续查找并返回结果
   if cmp(n.value, e) > 0 {
      return n.left.count(e, isMulti, cmp)
   }
   //n中承载元素等于e,直接返回结果
   return n.value
}
```

### 使用示例

```go
package main

import (
	"fmt"
	"github.com/hlccd/goSTL/data_structure/avlTree"
	"sync"
)

func main() {
	bs := avlTree.New(true)
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
	fmt.Println("删除一次平衡二叉树中存在的元素,存在重复的将会被剩下")
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
> 4
> 5
> 6
> 7
> 7
> 8
> 8
> 8
> 8
> 9
> 9
> 9
> 9
> 9
> 10
> 10
> 删除一次平衡二叉树中存在的元素,存在重复的将会被剩下
> 输出剩余的重复元素
> 7
> 8
> 8
> 8
> 9
> 9
> 9
> 9
> 10
> 10

#### 时间开销

```go
package main

import (
   "fmt"
   "github.com/hlccd/goSTL/data_structure/avlTree"
   "time"
)

func main() {
   max:=10000000
   tv := time.Now()
   v := make([]interface{},max,max)
   num:=0
   for i := 0; i < max; i++ {
      if i%2==1{
         num=max-i
      }else{
         num=i
      }
      v[num]=num
   }
   fmt.Println("插入切片的消耗时间:",time.Since(tv))
   tt := time.Now()
   t := avlTree.New(false)
   for i := 0; i < max; i++ {
      if i%2==1{
         num=max-i
      }else{
         num=i
      }
      t.Insert(num)
   }
   fmt.Println("插入平衡二叉树的消耗时间:",time.Since(tt))
}
```

> 插入切片的消耗时间: 195.5018ms
> 插入平衡二叉树的消耗时间: 5.6738274s

